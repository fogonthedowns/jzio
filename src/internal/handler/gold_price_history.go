package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"jzio/src/internal/driver"
)

type Price struct {
	DateString string `json:"date_string"`
	PriceDate  uint32 `json:"price_date"`
	Cents      uint32 `json:"cents"`
}

func GetGoldPriceHistory(db *driver.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
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
		prices, err := loadGoldPriceHistory(ctx, tx, int(start), int(end))

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		respondwithJSON(w, http.StatusCreated, prices)
	}
}

// lastGoldPrice() Loads the last gold price.
// Returns cents
func LastGoldPrice(ctx context.Context, tx *sql.Tx) (uint32, error) {
	query := "select cents from gold_price_history order by id desc limit 1"

	rows, err := tx.QueryContext(
		ctx,
		query,
	)

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	price := Price{}
	for rows.Next() {
		switch err := rows.Scan(&price.Cents); err {
		case sql.ErrNoRows:
		case nil:
			return price.Cents, nil
		default:
			panic(err)
		}
	}
	return price.Cents, nil
}

// loadGoldPriceHistory() Loads a range of prices between two timestamps, start and end
func loadGoldPriceHistory(ctx context.Context, tx *sql.Tx, start int, end int) (map[string]uint32, error) {
	query := "select date_string, cents from gold_price_history where price_date>? and price_date<? and fast_search=true"
	rows, err := tx.QueryContext(
		ctx,
		query,
		start,
		end,
	)

	fmt.Printf("query %v start %v end %v\n", query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prices := make(map[string]uint32, 0)
	for rows.Next() {
		price := Price{}
		switch err := rows.Scan(&price.DateString, &price.Cents); err {
		case sql.ErrNoRows:
		case nil:
			prices[price.DateString] = (price.Cents / 100)
		default:
			panic(err)
		}
		fmt.Println(price)
	}

	return prices, nil
}

func (p *Price) Save(ctx context.Context, tx *sql.Tx) (priceId int64, err error) {
	fastSearch := false

	i, err := strconv.ParseInt(strings.Split(p.DateString, "-")[2], 10, 64)
	if err != nil {
		panic(err)
	}

	if (i % 5) == 0 {
		fastSearch = true
	}

	query := "INSERT INTO gold_price_history SET date_string=?, cents=?, fast_search=?, price_date=UNIX_TIMESTAMP()"
	s, err := tx.PrepareContext(ctx, query)

	result, err := s.ExecContext(
		ctx,
		p.DateString,
		p.Cents,
		fastSearch,
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

	priceId, err = result.LastInsertId()

	return priceId, err
}
