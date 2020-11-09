package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"jzio/src/internal/driver"
	"jzio/src/internal/lib"
)

type Transaction struct {
	Id                  int64
	UserId              int64               `json:"user_id"`
	UserPublicId        string              `json:"user_public_id"`
	DepositPublicId     string              `json:"deposit_public_id"`
	OrderPublicId       string              `json:"order_public_id"`
	TransactionType     lib.TransactionType `json:"transaction_type"`
	UnixTime            int                 `json:"unix_time"`
	DateString          string              `json:"date_string"`
	Cents               uint32              `json:"cents"`
	Side                lib.TradeSide       `json:"side"`
	GoldTenthMilliGrams uint32              `json:"gold_tenth_milli_grams"`
}

func GetTransactions(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		qStart := r.FormValue("start")
		qEnd := r.FormValue("end")
		uid := r.FormValue("user_public_id")

		start, err := strconv.ParseInt(qStart, 10, 32)

		if err != nil {
			log.Fatal(err)
		}

		end, err := strconv.ParseInt(qEnd, 10, 32)

		if err != nil {
			log.Fatal(err)
		}

		tx, err := db.SQL.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		// start and end Unixtime stamps
		transactions, err := loadTransactions(ctx, tx, uid, int(start), int(end))

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		respondwithJSON(w, http.StatusCreated, transactions)
	}
}

// loadTransactions() Loads a range of transactions between two timestamps, start and end
func loadTransactions(ctx context.Context, tx *sql.Tx, user_public_id string, start int, end int) ([]Transaction, error) {
	query := "select user_id, user_public_id, deposit_public_id, order_public_id, transaction_type, date_string, cents, side, gold_tenth_milli_grams from transactions where user_public_id=?"
	rows, err := tx.QueryContext(
		ctx,
		query,
		user_public_id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		t := Transaction{}
		switch err := rows.Scan(&t.UserId, &t.UserPublicId, &t.DepositPublicId, &t.OrderPublicId, &t.TransactionType, &t.DateString, &t.Cents, &t.Side, &t.GoldTenthMilliGrams); err {
		case sql.ErrNoRows:
		case nil:
			transactions = append(transactions, t)
		default:
			panic(err)
		}
	}

	return transactions, nil
}

func (t *Transaction) save(ctx context.Context, tx *sql.Tx) (id int64, err error) {
	query := "INSERT INTO `transactions` SET user_id=?, user_public_id=?, deposit_public_id=?, order_public_id=?, transaction_type=?, unix_time=UNIX_TIMESTAMP(), date_string=?, cents=?, side=?, gold_tenth_milli_grams=?"
	s, err := tx.PrepareContext(ctx, query)

	if err != nil {
		return 0, err
	}

	ts := time.Now().UTC()
	dateString := ts.Format("2006-1-2")
	result, err := s.ExecContext(
		ctx,
		t.UserId,
		t.UserPublicId,
		t.DepositPublicId,
		t.OrderPublicId,
		t.TransactionType,
		dateString,
		t.Cents,
		t.Side,
		t.GoldTenthMilliGrams,
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
		return 0, errors.New("expected single row affected. multiple rows affected")
	}

	id, err = result.LastInsertId()

	return id, err
}
