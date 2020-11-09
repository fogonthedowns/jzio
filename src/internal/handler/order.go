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

type Order struct {
	id                  int64
	UserPublicId        string `json:"user_public_id"`
	OrderPublicId       string
	createdAt           int
	updatedAt           int
	active              bool
	Status              string
	Cents               uint32        `json:"cents"`
	Side                lib.TradeSide `json:"side"`
	GoldTenthMilliGrams uint32        `json:"gold_tenth_milli_grams"`
}

func PostOrder(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		order := Order{}
		json.NewDecoder(r.Body).Decode(&order)

		ctx := context.Background()
		tx, err := db.SQL.BeginTx(ctx, nil)

		if err != nil {
			log.Fatal(err)

		}

		fmt.Printf("order %+v \n", order)

		if order.UserPublicId == "" {
			tx.Rollback()
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v", "user_public_id must not be blank")})
			return
		}

		// defined in deposit.go
		vault, err := loadVaultByPublicId(ctx, tx, order.UserPublicId)
		if vault.UserId == 0 {
			tx.Rollback()
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		if err != nil {
			tx.Rollback()
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		err = order.validate(ctx, vault)
		if err != nil {
			tx.Rollback()
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		order.OrderPublicId, err = order.save(ctx, tx)
		if err != nil {
			tx.Rollback()
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		err = vault.decrementVaultWith(ctx, &order)
		if err != nil {
			tx.Rollback()
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

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

		respondwithJSON(w, http.StatusCreated, map[string]string{"order_public_id": order.OrderPublicId})
	}
}

func (o *Order) validate(ctx context.Context, v *Vault) error {

	fmt.Printf("order is here: %+v \n", o)
	if o.Cents == 0 {
		return errors.New("cents must not be empty.")
	}

	if o.Side == "" {
		return errors.New("side must not be empty.")
	}

	if o.GoldTenthMilliGrams == 0 {
		return errors.New("gold_tenth_milli_grams must not be empty.")
	}

	if v.GoldTenthMilliGramsAvailable < o.GoldTenthMilliGrams {
		return errors.New("not enough available gold to trade.")
	}

	if o.Side == "buy" {
		if v.CentsAvailable < o.Cents {
			return errors.New("not enough USD Available to trade.")
		}
	}

	return nil
}

func (o *Order) save(ctx context.Context, tx *sql.Tx) (orderPublicId string, err error) {
	query := "INSERT INTO orders SET user_public_id=?, order_public_id=?, cents=?, side=?, gold_tenth_milli_grams=?, created_at=UNIX_TIMESTAMP(), updated_at=UNIX_TIMESTAMP()"
	s, err := tx.PrepareContext(ctx, query)

	orderPublicId = generatePublicId("ord")

	if err != nil {
		return "", err
	}
	result, err := s.ExecContext(
		ctx,
		o.UserPublicId,
		orderPublicId,
		o.Cents,
		o.Side,
		o.GoldTenthMilliGrams,
	)

	defer s.Close()

	if err != nil {
		return "", err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rows != 1 {
		return "", errors.New("expected single row affected. multiple rows affected")
	}

	// order_id, err = result.LastInsertId()

	return orderPublicId, err
}

// LoadOpenOrders() Loads a single vault record by user_id
func LoadOpenOrders(ctx context.Context, tx *sql.Tx) ([]*Order, error) {
	query := "SELECT id, gold_tenth_milli_grams, cents, side, order_public_id FROM orders WHERE status='open'"

	rows, err := tx.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*Order, 0)
	for rows.Next() {
		o := Order{}
		switch err := rows.Scan(&o.id, &o.GoldTenthMilliGrams, &o.Cents, &o.Side, &o.OrderPublicId); err {
		case sql.ErrNoRows:
		case nil:
			orders = append(orders, &o)
		default:
			panic(err)
		}
	}

	return orders, nil
}
