package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"assesment.sqlc.dev/app/models"
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

func ParseUserTransactionData(foundUserTransactions []postgres.GetUserWalletTransactionsRow) ([]models.UserTransaction, error) {
	var result []models.UserTransaction
	if len(foundUserTransactions) < 1 {
		return nil, fmt.Errorf("No rows to parse")
	}
	println(len(foundUserTransactions))
	for i := 0; i < len(foundUserTransactions); i++ {
		var transaction models.UserTransaction
		transaction.ID = int32(foundUserTransactions[i].ID)
		transaction.TransactionAmount = foundUserTransactions[i].TransactionAmount.Float64
		transaction.UserWalletID = foundUserTransactions[i].UserWalletID.Int32
		transaction.UserWalletData.ID = int32(foundUserTransactions[i].ID_2)
		transaction.UserWalletData.Amount = foundUserTransactions[i].Amount.Float64
		transaction.UserWalletData.User.ID = foundUserTransactions[i].ID_3
		transaction.UserWalletData.User.Email = foundUserTransactions[i].Email.String
		transaction.UserWalletData.User.Name = foundUserTransactions[i].Name
		transaction.UserWalletData.User.Password = foundUserTransactions[i].Password.String

		result = append(result, transaction)
	}
	return result, nil
}
func ParseUserWalletData(foundWallet postgres.GetUserWalletRow) interface{} {
	var result models.UserWallet
	result.ID = foundWallet.UserID.Int32
	result.Amount = foundWallet.Amount.Float64
	result.UserID = foundWallet.UserID.Int32
	result.User.ID = foundWallet.ID_2
	result.User.Name = foundWallet.Name
	result.User.Email = foundWallet.Email.String
	result.User.Password = foundWallet.Password.String
	return result
}
