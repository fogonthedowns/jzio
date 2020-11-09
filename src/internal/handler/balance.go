package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"jzio/src/internal/driver"
)

type DailyBalance struct {
	id                  int64
	GoldTenthMilliGrams uint32 `json:"gold_tenth_milli_grams"`
	Cents               uint32 `json:"cents"`
	Balance             uint32 `json:"balance"`
	UserId              int64  `json:"user_id"`
	UserPublicId        string `json:"user_public_id"`
	DateString          string `json:"date_string"`
	UnixTime            int    `json:"time_stamp"`
}

// GetUserDailyBalanceHistory() GET returns a balance history
// request with a user_public_id
func GetUserDailyBalanceHistory(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userPublicId := r.FormValue("user_public_id")

		qStart := r.FormValue("start")
		qEnd := r.FormValue("end")

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
		prices, err := loadUserDailyBalanceHistory(ctx, tx, userPublicId, int(start), int(end))

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		respondwithJSON(w, http.StatusCreated, prices)
	}
}

// CreateDailyBalance() creates a daily balance record for a User
// requires *User and lastPriceCents
func CreateDailyBalance(ctx context.Context, tx *sql.Tx, u *User, lastPriceCents uint32) (err error) {
	query := "INSERT INTO daily_balance SET gold_tenth_milli_grams=?, cents=?, balance=?, user_id=?, user_public_id=?, date_string=?, unix_time=UNIX_TIMESTAMP()"
	s, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	// converts unix time to date string
	i := time.Now().Unix()
	tm := time.Unix(i, 0)

	mg := tenthMilliGramTomg(u.Vault.GoldTenthMilliGrams)
	troyOunce := mgToTroyOunce(mg)
	balance := uint32(float64(lastPriceCents)*troyOunce) + u.Vault.Cents
	result, err := s.ExecContext(
		ctx,
		u.Vault.GoldTenthMilliGrams,
		u.Vault.Cents,
		balance,
		u.Id,
		u.UserPublicId,
		tm.Format("2006-01-02"),
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

// tenthMilliGramTomg
// divide by 10
func tenthMilliGramTomg(g uint32) uint32 {
	return (g / 10)
}

// Approximate result, divide the mass value by 31103
func mgToTroyOunce(mg uint32) float64 {
	return float64(mg) / float64(31103.4768)
}

// loadUserDailyBalanceHistory() Queries daily_balance table by user_public_id
// and start and end unix time stamps
// returns a set of daily balances
func loadUserDailyBalanceHistory(ctx context.Context, tx *sql.Tx, userPublicId string, start int, end int) (map[string][]uint32, error) {
	query := "select date_string, gold_tenth_milli_grams, cents, balance from daily_balance where user_public_id=? and unix_time>? and unix_time<?"
	rows, err := tx.QueryContext(
		ctx,
		query,
		userPublicId,
		start,
		end,
	)

	fmt.Printf("%v  %v start%v end %v\n", query, userPublicId, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prices := make(map[string][]uint32, 0)
	for rows.Next() {
		b := DailyBalance{}
		switch err := rows.Scan(&b.DateString, &b.GoldTenthMilliGrams, &b.Balance, &b.Cents); err {
		case sql.ErrNoRows:
		case nil:
			result := []uint32{b.GoldTenthMilliGrams, b.Balance, b.Cents}
			prices[b.DateString] = result
		default:
			panic(err)
		}
		fmt.Println(b)
	}

	return prices, nil
}
