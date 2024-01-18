package utils_test

import (
	"context"
	"os"
	"testing"

	"assesment.sqlc.dev/app/postgres"
	"assesment.sqlc.dev/app/utils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestUpdateWalletAmount(t *testing.T) {
	// Create a new instance of postgres.Queries
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:admin@localhost:5432/test_sqlc_practice")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
	queries := postgres.New(conn)

	// Call the function under test with the actual implementation
	err = utils.UpdateWalletAmount(1, 50.0, queries)

	// Perform assertions or verifications as needed
	assert.Nil(t, err)
}
func TestGenerateToken(t *testing.T) {
	// Set the JWT_SECRET environment variable for the test
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() { os.Setenv("JWT_SECRET", originalSecret) }()
	os.Setenv("JWT_SECRET", "testsecret")

	// Call the function under test
	userID := int64(123)
	token, err := utils.GenerateToken(userID)

	// Perform assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// TODO: Add additional assertions as needed
}
func TestParseToken(t *testing.T) {
	// Set the JWT_SECRET environment variable for the test
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() { os.Setenv("JWT_SECRET", originalSecret) }()
	os.Setenv("JWT_SECRET", "testsecret")

	// Generate a test token
	userID := int64(123)
	testToken, err := utils.GenerateToken(userID)
	assert.NoError(t, err)

	// Call the function under test
	claims, err := utils.ParseToken(testToken)

	// Perform assertions
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Verify that the claims contain the expected userID
	userIDFromClaims, ok := claims["userID"].(float64)
	assert.True(t, ok)
	assert.Equal(t, float64(userID), userIDFromClaims)

	// TODO: Add additional assertions as needed
}
