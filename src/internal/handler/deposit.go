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
	"jzio/src/internal/lib"
)

type Deposit struct {
	id                  int64
	GoldTenthMilliGrams uint32 `json:"gold_tenth_milli_grams"`
	Cents               uint32 `json:"cents"`
	UserId              int64  `json:"user_id"`
	CreatedAt           int    `json:"created_at"`
	DepositPublicId     string `json:"deposit_public_id"`
	UserPublicId        string `json:"user_public_id"`
}

// deposit (AMOUNT USD)
// Tranaction
// Insert USD to Vault
// Update Available to Trade USD
// Fix the above concepts with ACH
// {
//        "gold_tenth_milli_grams": 0,
//        "cents": 1,
//        "user_public_id": "rails_jz_test_with_vault",
//        "deposit_public_id": "dep_012312DJKSFDSCJ"
// }
func PostDeposit(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := Deposit{}
		json.NewDecoder(r.Body).Decode(&d)

		ctx := context.Background()
		tx, err := db.SQL.BeginTx(ctx, nil)

		if err != nil {
			log.Fatal(err)

		}

		_, err = validate(ctx, &d)
		if err != nil {
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		_, err = saveDeposit(ctx, tx, &d)
		if err != nil {
			tx.Rollback()
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		var depositType lib.TransactionType
		if d.GoldTenthMilliGrams > 0 {
			depositType = lib.DepositGold
		} else {
			depositType = lib.DepositUSD
		}

		transaction := &Transaction{
			UserPublicId:        d.UserPublicId,
			DepositPublicId:     d.DepositPublicId,
			TransactionType:     depositType,
			Cents:               d.Cents,
			Side:                lib.Sell,
			GoldTenthMilliGrams: d.GoldTenthMilliGrams,
		}

		_, err = transaction.save(ctx, tx)
		if err != nil {
			tx.Rollback()
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		vault, err := loadVaultByPublicId(ctx, tx, d.UserPublicId)
		if vault.UserId == 0 {
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", "user_public_id Not found")})
			return
		}
		if err != nil {
			tx.Rollback()
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		vault.deposit(ctx, &d)
		err = vault.update(ctx, tx)
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

func saveDeposit(ctx context.Context, tx *sql.Tx, d *Deposit) (code int, err error) {
	code, err = validate(ctx, d)
	if err != nil {
		return code, err
	}

	query := "INSERT INTO deposits SET gold_tenth_milli_grams=?, cents=?, user_id=?, deposit_public_id=?, user_public_id=?, created_at=UNIX_TIMESTAMP()"
	s, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 500, err
	}
	result, err := s.ExecContext(
		ctx,
		d.GoldTenthMilliGrams,
		d.Cents,
		d.UserId,
		d.DepositPublicId,
		d.UserPublicId,
	)

	defer s.Close()

	if err != nil {
		return 500, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 500, err
	}

	if rows != 1 {
		return 500, errors.New("expected single row affected. multiple rows affected")
	}
	return 201, err
}

func validate(ctx context.Context, d *Deposit) (int, error) {
	if d.GoldTenthMilliGrams == 0 && d.Cents == 0 {
		return 400, errors.New("set gold_tenth_milli_grams or cents. A new deposit can not be empty.")
	}

	if d.UserPublicId == "" {
		return 400, errors.New("user_public_id must be set")
	}

	if d.DepositPublicId == "" {
		return 400, errors.New("depsoit_public_id must be set")
	}
	return 0, nil
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(response)
}
