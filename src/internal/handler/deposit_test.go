package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessfullCentsDeposit(t *testing.T) {
	requestBody := `{
		"gold_tenth_milli_grams": 0,
        	"cents": 1,
        	"user_public_id": "test_user",
        	"deposit_public_id": "dep_012312DJKSFDSCJ"
	}`

	res, bodyString := PostRequest(t, "http://localhost:3001/deposit", nil, requestBody)
	require.Contains(t, bodyString, `{"success":true`)
	require.Equal(t, "201 Created", res.Status)
}

func TestSuccessfullGoldDeposit(t *testing.T) {
	requestBody := `{
                "gold_tenth_milli_grams": 101,
                "cents": 0,
                "user_public_id": "test_user",
                "deposit_public_id": "dep_012312DJKSFDSCJ"
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/deposit", nil, requestBody)
	require.Contains(t, bodyString, `{"success":true`)
	require.Equal(t, "201 Created", res.Status)
}

func Test400MissingUser(t *testing.T) {
	requestBody := `{
                "gold_tenth_milli_grams": 101,
                "cents": 0,
                "deposit_public_id": "dep_012312DJKSFDSCJ"
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/deposit", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"user_public_id must be set"`)
	require.Equal(t, "400 Bad Request", res.Status)
}

func Test400EmptyDeposit(t *testing.T) {
	requestBody := `{
                "gold_tenth_milli_grams": 0,
                "cents": 0,
		"user_public_id": "test_user",
                "deposit_public_id": "dep_012312DJKSFDSCJ"
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/deposit", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"set gold_tenth_milli_grams or cents. A new deposit can not be empty."`)
	require.Equal(t, "400 Bad Request", res.Status)
}

func Test500BadUser(t *testing.T) {
	requestBody := `{
                "gold_tenth_milli_grams": 101,
                "cents": 0,
		"user_public_id": "999991232132",
                "deposit_public_id": "dep_012312DJKSFDSCJ"
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/deposit", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"user_public_id Not found"`)
	require.Equal(t, "500 Internal Server Error", res.Status)
}
