package balance

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"jzio/src/internal/driver"
	"jzio/src/internal/handler"

	"github.com/thedevsaddam/gojsonq"
)

func FetchAndSaveLastGoldPrice(ctx context.Context) (err error) {
	price := handler.Price{}
	dollars, err := FetchGoldPrice(ctx)
	price.Cents = uint32(dollars * 100)
	now := time.Now()
	price.DateString = now.Format("2006-01-02")
	fmt.Printf("price %+v \n", price)

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	connection, err := driver.ConnectSQL(dbHost, dbPort, "root", dbPass, dbName)
	if err != nil {
		os.Exit(-1)
	}

	tx, err := connection.SQL.BeginTx(ctx, nil)
	price.Save(ctx, tx)
	err = tx.Commit()
	return err
}

// External API
// Returns dollars
func FetchGoldPrice(ctx context.Context) (dollars float64, err error) {
	dollars, err = getJson("https://forex-data-feed.swissquote.com/public-quotes/bboquotes/instrument/XAU/USD")
	if err != nil {
		return 0, err
	}
	return dollars, err
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string) (float64, error) {
	r, err := myClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	res := gojsonq.New().Reader(r.Body).Find("[0].spreadProfilePrices.[0].ask")

	return res.(float64), err
}
