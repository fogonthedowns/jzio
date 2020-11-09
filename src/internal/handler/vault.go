package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"jzio/src/internal/driver"
	"jzio/src/internal/lib"
)

type Vault struct {
	id                           int64
	GoldTenthMilliGrams          uint32 `json:"gold_tenth_milli_grams"`
	GoldTenthMilliGramsAvailable uint32 `json:"gold_tenth_milli_grams_available"`
	Cents                        uint32 `json:"cents"`
	CentsAvailable               uint32 `json:"cents_available"`
	UserId                       int64  `json:"user_id"`
	UpdatedAt                    int    `json:"updated_at"`
	VaultPublicId                string `json:"vault_public_id"`
	UserPublicId                 string `json:"user_public_id"`
}

// GetVault() GET returns a User Vault
// request with a user_public_id
func GetVault(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userPublicId := r.FormValue("user_public_id")

		tx, err := db.SQL.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		// start and end Unixtime stamps
		vault, err := loadVaultByPublicId(ctx, tx, userPublicId)

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		respondwithJSON(w, http.StatusCreated, vault)
	}
}

// createVault() creates a new vault after a user is created
func createVault(ctx context.Context, tx *sql.Tx, u User, user_id int64, user_public_id string) (err error) {
	query := "INSERT INTO vaults SET user_id=?, user_public_id=?, vault_public_id=?, updated_at=UNIX_TIMESTAMP()"
	s, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	result, err := s.ExecContext(
		ctx,
		user_id,
		user_public_id,
		generatePublicId("vlt"),
	)

	defer s.Close()

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("expected single new vault record. multiple rows affected")
	}
	return err
}

// loadVault() Loads a single vault record by user_id
func loadVault(ctx context.Context, tx *sql.Tx, user_id int64) (*Vault, error) {
	query := "SELECT id, gold_tenth_milli_grams, gold_tenth_milli_grams_available, cents, cents_available, user_id, vault_public_id FROM vaults WHERE user_id=? LIMIT 1"

	rows, err := tx.QueryContext(
		ctx,
		query,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vault := Vault{}
	for rows.Next() {
		switch err := rows.Scan(&vault.id, &vault.GoldTenthMilliGrams, &vault.GoldTenthMilliGramsAvailable, &vault.Cents, &vault.CentsAvailable, &vault.UserId, &vault.VaultPublicId); err {
		case sql.ErrNoRows:
		case nil:
			return &vault, nil
		default:
			panic(err)
		}
	}

	return &vault, nil
}

// loadVault() Loads a single vault record by user_id
func loadVaultByPublicId(ctx context.Context, tx *sql.Tx, user_public_id string) (*Vault, error) {
	query := "SELECT id, gold_tenth_milli_grams, gold_tenth_milli_grams_available, cents, cents_available, user_id, vault_public_id, user_public_id FROM vaults WHERE user_public_id=? LIMIT 1"

	rows, err := tx.QueryContext(
		ctx,
		query,
		user_public_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vault := Vault{}
	for rows.Next() {
		switch err := rows.Scan(&vault.id, &vault.GoldTenthMilliGrams, &vault.GoldTenthMilliGramsAvailable, &vault.Cents, &vault.CentsAvailable, &vault.UserId, &vault.VaultPublicId, &vault.UserPublicId); err {
		case sql.ErrNoRows:
		case nil:
			return &vault, nil
		default:
			panic(err)
		}
	}

	return &vault, nil
}

// deposit() recieves a new Deposit and modifies
// the users vault, increasing the gold, USD,
// available gold and available gold
func (v *Vault) deposit(ctx context.Context, d *Deposit) {
	v.GoldTenthMilliGramsAvailable += d.GoldTenthMilliGrams
	v.GoldTenthMilliGrams += d.GoldTenthMilliGrams
	v.CentsAvailable += d.Cents
	v.Cents += d.Cents
}

// decrementVaultWith() recieves a new order and modifies
// the users Vault depending upon the order side.
// Buy reduces vault available USD, Sell reduces vault available Gold.
// TODO return errors if values are negative
func (v *Vault) decrementVaultWith(ctx context.Context, o *Order) error {
	switch o.Side {
	case lib.Buy:
		v.CentsAvailable -= o.Cents
	case lib.Sell:
		v.GoldTenthMilliGramsAvailable -= o.GoldTenthMilliGrams
	default:
		return errors.New("side not set.")
	}
	return nil
}

func (v *Vault) update(ctx context.Context, tx *sql.Tx) error {
	query := "UPDATE VAULTS SET gold_tenth_milli_grams_available=?, gold_tenth_milli_grams=?, cents_available=?, cents=?, vault_public_id=?, updated_at=UNIX_TIMESTAMP() where user_id=? LIMIT 1"
	s, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	result, err := s.ExecContext(
		ctx,
		v.GoldTenthMilliGramsAvailable,
		v.GoldTenthMilliGrams,
		v.CentsAvailable,
		v.Cents,
		v.VaultPublicId,
		v.UserId,
	)

	defer s.Close()

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("expected single row affected.")
	}

	return err
}
