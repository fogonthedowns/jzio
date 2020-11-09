package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessfullCreateUser(t *testing.T) {
	requestBody := `{
                "user_public_id": "test_user"
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/users", nil, requestBody)
	require.Contains(t, bodyString, `{"success":true`)
	require.Equal(t, "201 Created", res.Status)
}

func Test400MissingUserPublicId(t *testing.T) {
	requestBody := `{
        }`

	res, bodyString := PostRequest(t, "http://localhost:3001/users", nil, requestBody)
	require.Contains(t, bodyString, `{"error":"user_public_id must not be empty."`)
	require.Equal(t, "400 Bad Request", res.Status)
}
