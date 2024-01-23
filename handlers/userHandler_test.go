package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"assesment.sqlc.dev/app/middleware"
	"assesment.sqlc.dev/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func TestUserHandler(t *testing.T) {
	// You might want to use a testing database or a mock for more extensive testing
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:admin@localhost:5432/test_sqlc_practice")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
	defer utils.CleanUp(conn)
	// Create a Gin router
	router := gin.New()
	
	// Create a request with JSON payload
	createUserReq, signInUserReq, getUserReq, getAllUserReq, getUserTransactionsReq, getUserWalletReq, err := utils.CreateRequests(t)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	createUserRecorder := httptest.NewRecorder()
	signInUserRecorder := httptest.NewRecorder()
	getUserRecorder := httptest.NewRecorder()
	getAllUserRecorder := httptest.NewRecorder()
	createUserWalletRecorder := httptest.NewRecorder()

	// Create an instance of CreateUserHandler with the temporary database connection
	userHandler := &UserHanlder{conn: conn}
	userWalletHandler := &WalletHanlder{conn: conn}

	// Set up the router with the CreateUserHandler methods
	Routes(router, userHandler, userWalletHandler)

	// CREATE USER API TEST
	router.ServeHTTP(createUserRecorder, createUserReq)
	utils.CreateUserTest(createUserRecorder, t)

	// SIGN IN API TESTING
	router.ServeHTTP(signInUserRecorder, signInUserReq)
	utils.SiginInUserTest(signInUserRecorder, t)

	// GET USER API TESTING
	var response utils.ResponseBody
	err = json.Unmarshal(signInUserRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	utils.GetUser(&response, getUserReq, router, getUserRecorder)
	utils.GetUserTest(getUserRecorder, t)

	// GET ALL USERS API TEST
	getAllUserReq.Header.Set("Authorization", response.Token)
	router.ServeHTTP(getAllUserRecorder, getAllUserReq)
	utils.GetAllUserTest(getUserRecorder, t)

	// CREATE WALLET API TESTING
	var walletResponse *utils.ResponseBody
	_ = json.Unmarshal(signInUserRecorder.Body.Bytes(), &walletResponse)
	utils.CreateWallet(&response, router, createUserWalletRecorder, t, walletResponse)

	// DEPOSIT IN USER WALLET API TEST
	var createWalletResponse *utils.CreateWalletResponse
	_ = json.Unmarshal(createUserWalletRecorder.Body.Bytes(), &createWalletResponse)
	utils.DepositInWallet(&response, router, t, createWalletResponse)

	// DEPOSIT IN USER WALLET API TEST
	utils.WithdrawFromWallet(&response, router, t, createWalletResponse)

	// GET USER TRANSACTIONS API TESTING
	utils.GetUserTransactions(&response, getUserTransactionsReq, router, t, createWalletResponse)

	// GET USER WALLET API TESTING
	utils.GetUserWallet(&response, getUserWalletReq, router, t)

}

func Routes(router *gin.Engine, userHandler *UserHanlder, userWalletHandler *WalletHanlder) {
	router.POST("/create", userHandler.CreateUser)
	router.POST("/sign-in", userHandler.SignIn)
	router.POST("/wallet/create", userWalletHandler.Create)
	router.PUT("/wallet/deposit", userWalletHandler.Deposit)
	router.PUT("/wallet/withdraw", userWalletHandler.Withdraw)
	router.GET("/get", userHandler.GetUser)
	router.GET("/get-all-users", middleware.AuthMiddleware(), userHandler.GetAllUser)
	router.GET("/wallet/get-user-Transactions", userWalletHandler.GetUserTransactions)
	router.GET("/wallet/get", userWalletHandler.GetWallet)
}
