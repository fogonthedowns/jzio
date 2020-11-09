package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"jzio/src/internal/driver"

	"github.com/google/uuid"
)

type User struct {
	active       bool
	UserPublicId string `json:"user_public_id"`
	CreatedAt    int32  `json:"created_at"`
	Id           int64  `json"id"`
	Vault        *Vault
}

func PostUser(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := User{}
		json.NewDecoder(r.Body).Decode(&user)

		ctx := context.Background()
		tx, err := db.SQL.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal(err)

		}

		err = validateUser(ctx, &user)
		if err != nil {
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		id, err := saveUser(ctx, tx, &user)
		if err != nil {
			tx.Rollback()
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		err = createVault(ctx, tx, user, id, user.UserPublicId)
		if err != nil {
			tx.Rollback()
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		respondwithJSON(w, http.StatusCreated, map[string]bool{"success": true})
	}
}

// GetUserVault()
// GET balance?user=UserPublicId
func GetUserVault(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := User{}

		val := r.FormValue("user")
		user.UserPublicId = val
		ctx := context.Background()
		tx, err := db.SQL.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal(err)

		}

		err = validateUser(ctx, &user)
		if err != nil {
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		vault, err := loadVaultByPublicId(ctx, tx, user.UserPublicId)
		if vault.UserPublicId == "" {
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", "user_public_id Not found")})
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		respondwithJSON(w, http.StatusCreated, vault)
	}
}

func validateUser(ctx context.Context, u *User) error {
	if u.UserPublicId == "" {
		return errors.New("user_public_id must not be empty.")
	}

	return nil
}

func saveUser(ctx context.Context, tx *sql.Tx, u *User) (user_id int64, err error) {
	query := "INSERT INTO users SET user_public_id=?, created_at=UNIX_TIMESTAMP()"
	s, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	result, err := s.ExecContext(
		ctx,
		u.UserPublicId,
	)

	defer s.Close()

	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rows != 1 {
		return 0, errors.New("expected single new user record. multiple rows affected")
	}

	user_id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return user_id, err
}

func generatePublicId(prefix string) string {
	id := uuid.New()
	return fmt.Sprintf("%s_%s", prefix, id.String())
}

// loadAllUsers() loads all users
// TODO In the future add an offset
func LoadAllUsers(ctx context.Context, tx *sql.Tx) ([]*User, error) {
	query := "SELECT u.id, u.user_public_id, v.cents, v.gold_tenth_milli_grams FROM `users` u JOIN `vaults` v ON u.id = v.user_id WHERE u.active=true;"
	rows, err := tx.QueryContext(
		ctx,
		query,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := User{}
		vault := Vault{}
		switch err := rows.Scan(&user.Id, &user.UserPublicId, &vault.Cents, &vault.GoldTenthMilliGrams); err {
		case sql.ErrNoRows:
		case nil:
			user.Vault = &vault
			users = append(users, &user)
		default:
			panic(err)
		}
	}

	return users, nil
}
