package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

type ResponseBody struct {
	Result struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"result"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Token   string `json:"json"`
}
type CreateWalletResponse struct {
	Result struct {
		ID     int32   `json:"id"`
		UserID int32   `json:"user_id"`
		Amount float64 `json:"amount"`
	} `json:"result"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Token   string `json:"json"`
}

var TestEmail = "test@gmail.com"
var TestUserID int
var TestUserWalletID int

func CreateUserTest(rr *httptest.ResponseRecorder, t *testing.T) {
	var response ResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, rr.Code)
	if rr.Code == http.StatusOK {
		t.Log("User Created")
	} else if rr.Code == http.StatusNoContent {
		t.Fatal("Body Not recieved")
	} else if rr.Code == http.StatusBadRequest {
		t.Fatal("All Required variables destructured")
	} else if rr.Code == http.StatusConflict {
		t.Fatal("Email Already exsists")
	} else if rr.Code == http.StatusInternalServerError {
		t.Fatal("Error in saving using in database")
	} else {
		t.Fatal(response.Error)
	}
}

func SiginInUserTest(rr *httptest.ResponseRecorder, t *testing.T) {
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	_, tokenPresent := response["jwt-token"].(string)
	assert.True(t, tokenPresent, "Token is not present in the response")
	if rr.Code == http.StatusOK {
		t.Log("User Signed in")
	} else if rr.Code == http.StatusNoContent {
		t.Fatal("Body Not recieved")
	} else if rr.Code == http.StatusBadRequest {
		t.Fatal("All Required variables destructured")
	} else if rr.Code == http.StatusNotFound {
		t.Fatal("Email does not exsists")
	} else if rr.Code == http.StatusUnauthorized {
		t.Fatal("Password validation failed")
	} else {
		t.Fatal("Server crashing")
	}
}

func GetUserTest(rr *httptest.ResponseRecorder, t *testing.T) {
	if rr.Code == http.StatusOK {
		t.Log("User Data recieved")
	} else if rr.Code == http.StatusBadRequest {
		t.Fatal("ID not sent")
	} else if rr.Code == http.StatusNotFound {
		t.Fatal("User not found")
	} else {
		t.Fatal("Server crashing")
	}
}

func GetAllUserTest(rr *httptest.ResponseRecorder, t *testing.T) {
	if rr.Code == http.StatusOK {
		t.Log("All Users Fetched")
	} else {
		t.Fatal("Server Error")
	}
}

func CreateWalletTest(rr *httptest.ResponseRecorder, t *testing.T) {
	if rr.Code == http.StatusOK {
		t.Log("User Wallet Created")
	} else if rr.Code == http.StatusNoContent {
		t.Fatal("Body Not recieved")
	} else if rr.Code == http.StatusBadRequest {
		t.Fatal("All Required variables are not destructured")
	} else if rr.Code == http.StatusNotFound {
		t.Fatal("User does not exsists")
	} else if rr.Code == http.StatusConflict {
		t.Fatal("Wallet Already Exsists")
	} else if rr.Code == http.StatusInternalServerError {
		t.Fatal("Error while adding wallet in db")
	} else {
		t.Fatal("Server crashing")
	}
}

func DepositTest(rr *httptest.ResponseRecorder, t *testing.T) {
	if rr.Code == http.StatusOK {
		t.Log("Deposit Sucessfull")
	} else if rr.Code == http.StatusNoContent {
		t.Fatal("Body Not recieved")
	} else if rr.Code == http.StatusBadRequest {
		t.Fatal("All Required variables are not destructured")
	} else if rr.Code == http.StatusNotFound {
		t.Fatal("Wallet does not exsists")
	} else if rr.Code == http.StatusInternalServerError {
		t.Fatal("Error while despoiting wallet in db")
	} else {
		t.Fatal("Server crashing")
	}
}

func WithdrawTest(rr *httptest.ResponseRecorder, t *testing.T) {
	if rr.Code == http.StatusOK {
		t.Log("Withdrawal Sucessfull")
	} else if rr.Code == http.StatusNoContent {
		t.Fatal("Body Not recieved")
	} else if rr.Code == http.StatusBadRequest {
		t.Fatal("All Required variables are not destructured")
	} else if rr.Code == http.StatusNotFound {
		t.Fatal("Wallet does not exsists")
	} else if rr.Code == http.StatusInternalServerError {
		t.Fatal("Error while withdrawal wallet in db")
	} else {
		t.Fatal("Server crashing")
	}
}

func GetUserTransactionsTest(rr *httptest.ResponseRecorder, t *testing.T) {
	if rr.Code == http.StatusOK {
		t.Log("Transactions fetched")
	} else if rr.Code == http.StatusBadRequest {
		t.Fatal("ID not sent")
	} else if rr.Code == http.StatusNotFound {
		t.Fatal("Transaction not found")
	} else {
		t.Fatal("Server crashing")
	}
}

func GetUserWalletTest(rr *httptest.ResponseRecorder, t *testing.T) {
	if rr.Code == http.StatusOK {
		t.Log("Wallet Fetched")
	} else if rr.Code == http.StatusBadRequest {
		t.Fatal("ID not sent")
	} else if rr.Code == http.StatusNotFound {
		t.Fatal("wallet not found")
	} else {
		t.Fatal("Server crashing")
	}
}

func CleanUp(conn *pgx.Conn) {
	fmt.Println(TestUserWalletID)
	conn.Exec(context.Background(), "DELETE FROM user_transactions WHERE user_wallet_id = $1", TestUserWalletID)
	conn.Exec(context.Background(), "DELETE FROM user_wallet WHERE user_id = $1", TestUserID)
	conn.Exec(context.Background(), "DELETE FROM users WHERE email = $1", TestEmail)
}

func CreatePayLoads(createUserPayload, signInUserPayload *string) {
	*createUserPayload = `{"name": "test", "email": "test@gmail.com", "password": "123456"}`
	*signInUserPayload = `{"email": "test@gmail.com", "password": "123456"}`
}

func CreateRequests(t *testing.T) (*http.Request, *http.Request, *http.Request, *http.Request, *http.Request, *http.Request, error) {
	var createUserPayload, signInUserPayload string
	CreatePayLoads(&createUserPayload, &signInUserPayload)
	createUserReq, err := http.NewRequest("POST", "/create", bytes.NewBufferString(createUserPayload))
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	signInUserReq, err := http.NewRequest("POST", "/sign-in", bytes.NewBufferString(signInUserPayload))
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	getUserReq, err := http.NewRequest("GET", "/get", nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	getAllUserReq, err := http.NewRequest("GET", "/get-all-users", nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	getUserTransactionsReq, err := http.NewRequest("GET", "/wallet/get-user-Transactions", nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	getUserWalletReq, err := http.NewRequest("GET", "/wallet/get", nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	return createUserReq, signInUserReq, getUserReq, getAllUserReq, getUserTransactionsReq, getUserWalletReq, nil
}

// TEST FUNCTIONS
func GetUserWallet(response *ResponseBody, getUserWalletReq *http.Request, router *gin.Engine, t *testing.T) {
	getUserWalletRecorder := httptest.NewRecorder()
	TestUserID = int(response.Result.ID)
	guwr := getUserWalletReq.URL.Query()
	guwr.Add("user_id", strconv.FormatInt(int64(response.Result.ID), 10))
	getUserWalletReq.URL.RawQuery = guwr.Encode()
	getUserWalletReq.Header.Set("Authorization", response.Token)
	router.ServeHTTP(getUserWalletRecorder, getUserWalletReq)
	GetUserTransactionsTest(getUserWalletRecorder, t)
}

func GetUserTransactions(response *ResponseBody, getUserTransactionsReq *http.Request, router *gin.Engine, t *testing.T, createWalletResponse *CreateWalletResponse) {
	getUserTransactionRecorder := httptest.NewRecorder()
	TestUserID = int(response.Result.ID)
	wq := getUserTransactionsReq.URL.Query()
	wq.Add("user_wallet_id", strconv.FormatInt(int64(createWalletResponse.Result.ID), 10))
	getUserTransactionsReq.URL.RawQuery = wq.Encode()
	getUserTransactionsReq.Header.Set("Authorization", response.Token)
	router.ServeHTTP(getUserTransactionRecorder, getUserTransactionsReq)
	GetUserTransactionsTest(getUserTransactionRecorder, t)
}

func WithdrawFromWallet(response *ResponseBody, router *gin.Engine, t *testing.T, createWalletResponse *CreateWalletResponse) {
	widthdrawUserWalletRecorder := httptest.NewRecorder()
	widthdrawlPayload := map[string]interface{}{
		"user_wallet_id": createWalletResponse.Result.ID,
		"amount":         500,
	}
	widthdrawlWalletPayload, err := json.Marshal(widthdrawlPayload)
	if err != nil {
		t.Fatal(err)
	}

	withdrawUserWalletReq, err := http.NewRequest("PUT", "/wallet/withdraw", bytes.NewReader(widthdrawlWalletPayload))
	if err != nil {
		t.Fatal(err)
	}
	withdrawUserWalletReq.Header.Set("Authorization", response.Token)
	router.ServeHTTP(widthdrawUserWalletRecorder, withdrawUserWalletReq)
	WithdrawTest(widthdrawUserWalletRecorder, t)
}

func DepositInWallet(response *ResponseBody, router *gin.Engine, t *testing.T, createWalletResponse *CreateWalletResponse) {
	depositUserWalletRecorder := httptest.NewRecorder()
	TestUserWalletID = int(createWalletResponse.Result.ID)
	payload := map[string]interface{}{
		"user_wallet_id": createWalletResponse.Result.ID,
		"amount":         500,
	}
	depositWalletPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	depositUserWalletReq, err := http.NewRequest("PUT", "/wallet/deposit", bytes.NewReader(depositWalletPayload))
	if err != nil {
		t.Fatal(err)
	}
	depositUserWalletReq.Header.Set("Authorization", response.Token)
	router.ServeHTTP(depositUserWalletRecorder, depositUserWalletReq)
	DepositTest(depositUserWalletRecorder, t)
}

func CreateWallet(response *ResponseBody, router *gin.Engine, createUserWalletRecorder *httptest.ResponseRecorder, t *testing.T, walletResponse *ResponseBody) {
	createWalletpayload := map[string]interface{}{
		"user_id": walletResponse.Result.ID,
		"amount":  0,
	}
	createUserWalletPayload, err := json.Marshal(createWalletpayload)
	if err != nil {
		t.Fatal(err)
	}
	createUserWalletReq, err := http.NewRequest("POST", "/wallet/create", bytes.NewReader(createUserWalletPayload))
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(createUserWalletRecorder, createUserWalletReq)
	CreateWalletTest(createUserWalletRecorder, t)
}

func GetUser(response *ResponseBody, getUserReq *http.Request, router *gin.Engine, getUserRecorder *httptest.ResponseRecorder) {
	TestUserID = int(response.Result.ID)
	q := getUserReq.URL.Query()
	q.Add("id", strconv.FormatInt(response.Result.ID, 10))
	getUserReq.URL.RawQuery = q.Encode()
	router.ServeHTTP(getUserRecorder, getUserReq)
}
