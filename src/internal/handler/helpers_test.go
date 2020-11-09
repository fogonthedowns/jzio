package handler

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"jzio/src/internal/driver"

	"github.com/stretchr/testify/require"
)

// PostRequest(t, url, nil, requestBody)
func PostRequest(t *testing.T, url string, h http.Header, requestBody string) (*http.Response, string) {
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if h != nil {
		req.Header = h
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if t != nil {
		require.NoError(t, err)
	} else {
		if err != nil {
			return nil, ""
		}
	}

	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)

	if t != nil {
		require.NoError(t, err)
	}
	return res, string(responseBody)
}

// createTestUser() is used for test tear down/clean up
func createTestUser(db *driver.DB) {
	sqlStatement := `
INSERT INTO users SET user_public_id=?, created_at=UNIX_TIMESTAMP();`
	_, err := db.SQL.Exec(sqlStatement, "test_user")
	if err != nil {
		panic(err)
	}

}

// createTestVault() is used for test tear down/clean up
func createTestVault(db *driver.DB) {
	sqlStatement := `
INSERT INTO vaults SET user_public_id=?, cents_available=99999,gold_tenth_milli_grams_available=10001, user_id=(SELECT id from users order by id desc limit 1), vault_public_id=?, updated_at=UNIX_TIMESTAMP();`
	_, err := db.SQL.Exec(sqlStatement, "test_user", "test_vault")

	if err != nil {
		panic(err)
	}
}

// deleteTestTransactions() is used for test tear down/clean up
func deleteTestTransactions(db *driver.DB) {
	sqlStatement := `
DELETE FROM transactions
WHERE user_public_id = ?;`
	_, err := db.SQL.Exec(sqlStatement, "test_user")
	if err != nil {
		panic(err)
	}
}

// deleteTestUser() is used for test tear down/clean up
func deleteTestUser(db *driver.DB) {
	sqlStatement := `
DELETE FROM users
WHERE user_public_id = ?;`
	_, err := db.SQL.Exec(sqlStatement, "test_user")
	if err != nil {
		panic(err)
	}
}

// deleteTestUser() is used for test tear down/clean up
func deleteTestVault(db *driver.DB) {
	sqlStatement := `
DELETE FROM vaults
WHERE user_public_id = ?;`
	_, err := db.SQL.Exec(sqlStatement, "test_user")
	if err != nil {
		panic(err)
	}
}
