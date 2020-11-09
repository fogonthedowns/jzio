package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessfullSellGold(t *testing.T) {
	requestBody := `{
                "user_id":1,
		"user_public_id":"test_user",
        	"cents":967,
        	"gold_tenth_milli_grams":100,
        	"side":"sell"
	}`

	res, bodyString := PostRequest(t, "http://localhost:3001/orders", nil, requestBody)
	require.Contains(t, bodyString, `{"success":true`)
	require.Equal(t, "201 Created", res.Status)
}

func TestSuccessfullBuyGold(t *testing.T) {
	requestBody := `{
                "user_id":1,
        "user_public_id":"test_user",
	"cents":967,
        	"gold_tenth_milli_grams":100,
        	"side":"buy"
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/orders", nil, requestBody)
	require.Contains(t, bodyString, `{"success":true`)
	require.Equal(t, "201 Created", res.Status)
}

func Test400OrdersMissingUserId(t *testing.T) {
	requestBody := `{
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/orders", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"user_public_id must not be blank"`)
	require.Equal(t, "400 Bad Request", res.Status)
}

func Test400OrdersMissingCents(t *testing.T) {
	requestBody := `{
"user_public_id":"test_user",
        	"user_id":1
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/orders", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"cents must not be empty."`)
	require.Equal(t, "400 Bad Request", res.Status)
}

func Test400OrdersMissingUserSide(t *testing.T) {
	requestBody := `{
"user_public_id":"test_user",
        	"user_id":1,
		"cents":52
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/orders", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"side must not be empty."`)
	require.Equal(t, "400 Bad Request", res.Status)
}

func Test400OrdersMissingGold(t *testing.T) {
	requestBody := `{
"user_public_id":"test_user",
                "user_id":1,
                "cents":52,
		"side":"buy"
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/orders", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"gold_tenth_milli_grams must not be empty."`)
	require.Equal(t, "400 Bad Request", res.Status)
}

func Test400OrdersInsufficientUSD(t *testing.T) {
	requestBody := `{
"user_public_id":"test_user",
                "user_id":1,
                "cents":52999999,
		"gold_tenth_milli_grams": 10000,
                "side":"buy"
        }`
	res, bodyString := PostRequest(t, "http://localhost:3001/orders", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"not enough USD Available to trade."`)
	require.Equal(t, "400 Bad Request", res.Status)
}

func Test400OrdersInsufficientGld(t *testing.T) {
	requestBody := `{
"user_public_id":"test_user",
                "user_id":1,
                "cents":52999999,
                "gold_tenth_milli_grams": 100000000,
                "side":"sell"
        }`
	res, bodyString := PostRequest(t, "http://localhost:3001/orders", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"not enough available gold to trade."`)
	require.Equal(t, "400 Bad Request", res.Status)
}
