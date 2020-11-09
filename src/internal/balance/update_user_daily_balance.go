package balance

import (
	"context"
	"os"

	"jzio/src/internal/driver"
	"jzio/src/internal/handler"
)

func SaveUserDailyBalance(ctx context.Context) (err error) {

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	connection, err := driver.ConnectSQL(dbHost, dbPort, "root", dbPass, dbName)
	if err != nil {
		os.Exit(-1)
	}
	// Load the most recent gold price
	tx, err := connection.SQL.BeginTx(ctx, nil)
	cents, err := handler.LastGoldPrice(ctx, tx)

	users, err := handler.LoadAllUsers(ctx, tx)
	if err != nil {
		return err
	}

	for _, user := range users {
		err = handler.CreateDailyBalance(ctx, tx, user, cents)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}
