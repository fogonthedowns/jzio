package handler

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"jzio/src/internal/driver"
	"jzio/src/internal/lib"
)

var db *driver.DB

func TestMain(m *testing.M) {

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	var err error
	db, err = driver.ConnectSQL(dbHost, dbPort, "root", dbPass, dbName)
	if err != nil {
		os.Exit(-1)
	}

	createTestUser(db)
	createTestVault(db)
	ec := m.Run()
	deleteTestUser(db)
	deleteTestTransactions(db)
	deleteTestVault(db)
	os.Exit(ec)
}

func TestSuccessTransactionSave(t *testing.T) {
	transaction := &Transaction{
		UserPublicId:        "test_user",
		OrderPublicId:       "test_order",
		TransactionType:     lib.SellGold,
		Cents:               172399,
		Side:                lib.Sell,
		GoldTenthMilliGrams: 311035,
	}

	ctx := context.Background()
	tx, err := db.SQL.BeginTx(ctx, nil)

	if err != nil {
		log.Fatal(err)

	}

	_, err = transaction.save(ctx, tx)
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	require.NoError(t, err)
}
