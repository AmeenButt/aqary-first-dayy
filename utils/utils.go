package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"assesment.sqlc.dev/app/postgres"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v5/pgtype"
)

func UpdateWalletAmount(ID int32, Amount float64, queries *postgres.Queries) error {
	deductedAmount := queries.UpdateUserWalletAmount(context.Background(), postgres.UpdateUserWalletAmountParams{
		ID:     int64(ID),
		Amount: pgtype.Float8{Float64: Amount, Valid: true},
	})
	return deductedAmount
}
func GenerateToken(userID int64) (string, error) {
	// Set the secret key for signing the token
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	// Create a new token object with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// Set the secret key for validating the token
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Verify token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		println(claims)
		return claims, nil
	}

	return nil, fmt.Errorf("Token is not valid")
}
