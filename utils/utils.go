package utils

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"assesment.sqlc.dev/app/models"
	"assesment.sqlc.dev/app/postgres"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"gopkg.in/gomail.v2"
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
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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

	return nil, fmt.Errorf("token is not valid")
}

func ParseUserTransactionData(foundUserTransactions []postgres.GetUserWalletTransactionsRow) ([]models.UserTransaction, error) {
	var result []models.UserTransaction
	if len(foundUserTransactions) < 1 {
		return nil, fmt.Errorf("no rows to parse")
	}
	for i := 0; i < len(foundUserTransactions); i++ {
		var transaction models.UserTransaction
		transaction.ID = int32(foundUserTransactions[i].ID)
		transaction.TransactionAmount = foundUserTransactions[i].TransactionAmount.Float64
		transaction.CreatedAt = foundUserTransactions[i].CreatedAt.Time.String()
		transaction.UpdatedAt = foundUserTransactions[i].UpdatedAt.Time.String()
		transaction.UserWalletID = foundUserTransactions[i].UserWalletID.Int32
		transaction.UserWalletData.ID = int32(foundUserTransactions[i].ID_2)
		transaction.UserWalletData.Amount = foundUserTransactions[i].Amount.Float64
		transaction.UserWalletData.CreatedAt = foundUserTransactions[i].CreatedAt_2.Time.String()
		transaction.UserWalletData.UpdatedAt = foundUserTransactions[i].UpdatedAt_2.Time.String()
		transaction.UserWalletData.User.ID = foundUserTransactions[i].ID_3
		transaction.UserWalletData.User.Email = foundUserTransactions[i].Email.String
		transaction.UserWalletData.User.Name = foundUserTransactions[i].Name
		transaction.UserWalletData.User.Password = foundUserTransactions[i].Password.String
		transaction.UserWalletData.User.CreatedAt = foundUserTransactions[i].CreatedAt_3.Time.String()
		transaction.UserWalletData.User.UpdatedAt = foundUserTransactions[i].UpdatedAt_3.Time.String()

		result = append(result, transaction)
	}
	return result, nil
}

func ParseUserWalletData(foundWallet postgres.GetUserWalletRow) interface{} {
	var result models.UserWallet
	result.ID = int32(foundWallet.ID)
	result.Amount = foundWallet.Amount.Float64
	result.UserID = foundWallet.UserID.Int32
	result.User.ID = foundWallet.ID_2
	result.User.Name = foundWallet.Name
	result.User.Email = foundWallet.Email.String
	result.User.Password = foundWallet.Password.String
	return result
}

func ParsePropertyData(property postgres.GetPropertyByIDRow) models.Property {
	var result models.Property
	result.ID = property.ID
	result.SizeInSqFeet = int64(property.Sizeinsqfeet.Int32)
	result.Demand = property.Demand.String
	result.Status = property.Status.String
	result.Location = property.Location.String
	result.UserId = int64(property.UserID.Int32)
	result.Images = property.Images
	result.UpdatedAt = property.UpdatedAt.Time.String()
	result.CreatedAt = property.CreatedAt.Time.String()
	result.User.Email = property.Email.String
	result.User.ID = int64(property.UserID.Int32)
	result.User.Name = property.Name
	result.User.Password = property.Password.String
	result.User.UpdatedAt = property.UpdatedAt_2.Time.String()
	result.User.CreatedAt = property.CreatedAt_2.Time.String()
	return result
}

func ParsePropertyDataArray(property []postgres.GetPropertyByUserIDRow) []models.Property {
	var finalResult []models.Property
	for i := 0; i < len(property); i++ {
		var result models.Property
		result.ID = property[i].ID
		result.SizeInSqFeet = int64(property[i].Sizeinsqfeet.Int32)
		result.Demand = property[i].Demand.String
		result.Status = property[i].Status.String
		result.Location = property[i].Location.String
		result.Images = property[i].Images
		result.UserId = int64(property[i].UserID.Int32)
		result.UpdatedAt = property[i].UpdatedAt.Time.String()
		result.CreatedAt = property[i].CreatedAt.Time.String()
		result.User.Email = property[i].Email.String
		result.User.ID = int64(property[i].UserID.Int32)
		result.User.Name = property[i].Name
		result.User.Password = property[i].Password.String
		result.User.UpdatedAt = property[i].UpdatedAt_2.Time.String()
		result.User.CreatedAt = property[i].CreatedAt_2.Time.String()
		finalResult = append(finalResult, result)
	}
	return finalResult
}

func GenerateRandomCode() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(999999-100000+1) + 100000
}

func GetErrorMessage(err error) string {
	if err == nil {
		return "No error"
	}
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return "Record does not exsists"
	case errors.Is(err, pgx.ErrTxClosed):
		return "Transaction already done"
	default:
		return fmt.Sprintf("Internal server error: %v", err)
	}
}

func SendMail(email, subject, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_USERNAME"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	// Set up the email server connection information
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"))

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	} else {
		return nil
	}
}

func GetBindErrorMessage(err error) string {
	inputString := err.Error()
	errorString := ""
	lines := strings.Split(inputString, "\n")
	errors := []models.CustomError{}
	// Split the field name by period to get the field name
	errors = findError(lines, err, errors)
	errorString = convertErrorArrayToString(errors, errorString)
	return errorString
}

func findError(lines []string, err error, errors []models.CustomError) []models.CustomError {
	for _, line := range lines {
		var errorFor string
		var errorMessage string
		if strings.Contains(err.Error(), "unmarshal") {
			keyParts := strings.Split(line, ".")
			words := strings.Fields(keyParts[1])
			errorFor = words[0]
			errorMessage = "Invalid data type"
			errors = append(errors, models.CustomError{
				ErrorFor:     errorFor,
				ErrorMessage: errorMessage,
			})
		} else {
			keyParts := strings.Split(line, "'")

			if len(keyParts) >= 2 {
				fieldWithName := keyParts[1]

				fieldParts := strings.Split(fieldWithName, ".")
				if len(fieldParts) >= 2 {
					errorFor = fieldParts[1]

					if strings.Contains(line, "required") {
						errorMessage = "Field is required"
					} else {
						errorMessage = "Incorrect type of field"
					}

					errors = append(errors, models.CustomError{
						ErrorFor:     errorFor,
						ErrorMessage: errorMessage,
					})
				} else {
					fmt.Println("Field name not found in the input string.")
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println("Field name not found in the input string.")
				fmt.Println(err.Error())
			}
		}
	}
	return errors
}

func convertErrorArrayToString(errors []models.CustomError, errorString string) string {
	foundRequired := 0
	for _, customError := range errors {
		if strings.Contains(customError.ErrorMessage, "required") {
			foundRequired = foundRequired + 1
			errorString = errorString + customError.ErrorFor + ", "
		}
	}
	if foundRequired > 0 {
		if foundRequired == 1 {
			errorString = errorString + "is required"
		} else {
			errorString = errorString + "are required"
		}
	}
	foundInvalid := false
	for index, customError := range errors {
		if strings.Contains(customError.ErrorMessage, "Invalid") {
			foundInvalid = true
			if index == 0 && foundRequired > 0 {
				errorString = " and " + errorString + customError.ErrorFor + ","
			} else {
				errorString = errorString + customError.ErrorFor + ","
			}
		}
	}
	if foundInvalid {
		errorString = errorString + " has invalid datatype"
	}
	return errorString
}
