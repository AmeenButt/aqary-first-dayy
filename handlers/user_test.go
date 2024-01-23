package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"assesment.sqlc.dev/app/config"
	"assesment.sqlc.dev/app/postgres"
	mock_postgres "assesment.sqlc.dev/app/postgres/mock"
	"assesment.sqlc.dev/app/routes"
	"assesment.sqlc.dev/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name     string
	Method   string
	Endpoint string
	Payload  string
	Header   map[string]string
	Expect   func(*gomock.Controller, *mock_postgres.MockStore)
	Verify   func(*testing.T, *httptest.ResponseRecorder)
}

func MyNewServer(server *gin.Engine, conn *pgx.Conn, ctx *context.Context, store postgres.Store) {
	// Path to uploads folder
	uploadsFolderPath := "uploads"

	// Serve the uploads folder statically
	server.Static("/uploads", uploadsFolderPath)

	// Default route
	server.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"Message": "Server is running"}) })

	// User Routes
	routes.RegisterUserRoutes(server, conn, ctx, store)

	// User wallet Routes
	routes.RegisterUserWalletRoutes(server, conn, ctx, store)

	// Property Route
	routes.RegisterPropertiesRoutes(server, conn, ctx, store)
}

func runTest(t *testing.T, tc testCase, ctrl *gomock.Controller, store *mock_postgres.MockStore, server *gin.Engine) {

	req, err := http.NewRequest(tc.Method, tc.Endpoint, bytes.NewBufferString(tc.Payload))
	if err != nil {
		t.Fatal("error creating request:", err)
	}
	tc.Expect(ctrl, store)
	req.Header.Set("Authorization", tc.Header["Authorization"])
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)

	tc.Verify(t, recorder)
}

func TestHandlers(t *testing.T) {
	testCases := []testCase{
		{
			Name:     "CreateUser",
			Method:   "POST",
			Endpoint: "/users/create",
			Payload:  `{"name": "test", "email": "test234234@gmail.com", "password": "123456"}`,
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), pgtype.Text{String: "test234234@gmail.com", Valid: true}).Return(postgres.User{}, errors.New("user not found"))
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any())
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Method:   "GET",
			Endpoint: "/users/get?id=1",
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserByID(gomock.Any(), gomock.Any())
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
				// Add additional verification logic if needed
			},
		},
		{
			Method:   "POST",
			Endpoint: "/users/sign-in",
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(postgres.User{Password: pgtype.Text{String: "$2a$10$vt4xv8SE.B6O80BVzIupvus/y2Q7IEZLBz79sa2/MnJXeW23wiG6q"}}, nil)
				// Add expectations for other methods if needed
			},
			Payload: `{"email": "test234234@gmail.com", "password": "123456"}`,
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
				// Add additional verification logic if needed
			},
		},
		{
			Method:   "GET",
			Endpoint: "/users/get-all-users",
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().ListUsers(gomock.Any())
				// Add expectations for other methods if needed
			},
			Header: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYwOTcyODEsInVzZXJJRCI6Mn0.0-xjEz6Eu5f0W_RVpFGLIR0OZWxNLwe7aZxpWGgjU4g",
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
				// Add additional verification logic if needed
			},
		},
		{
			Method:   "POST",
			Endpoint: "/users/send-otp?id=1",
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(postgres.User{Email: pgtype.Text{String: "test@gmail.com"}}, nil)
				store.EXPECT().UpdateOTP(gomock.Any(), gomock.Any()).Return(postgres.User{Email: pgtype.Text{String: "test@gmail.com"}}, nil)
				// Add expectations for other methods if needed
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
				// Add additional verification logic if needed
			},
		},
		{
			Method:   "POST",
			Endpoint: "/users/verify-otp?id=1&otp=123456",
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(postgres.User{Otp: pgtype.Int4{Int32: 123456, Valid: true}}, nil)
				store.EXPECT().UpdateOTP(gomock.Any(), gomock.Any()).Return(postgres.User{Email: pgtype.Text{String: "test@gmail.com"}}, nil)
				// Add expectations for other methods if needed
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
				// Add additional verification logic if needed
			},
		},
		{
			Name:     "CreateWallet",
			Method:   "POST",
			Endpoint: "/wallet/create",
			Payload:  `{"user_id": 1}`,
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(postgres.User{}, nil)
				store.EXPECT().GetUserWallet(gomock.Any(), gomock.Any()).Return(postgres.GetUserWalletRow{}, errors.New("wallet not found"))
				store.EXPECT().CreateUserWallet(gomock.Any(), gomock.Any()).Return(postgres.UserWallet{}, nil)
			},
			Header: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYwOTcyODEsInVzZXJJRCI6Mn0.0-xjEz6Eu5f0W_RVpFGLIR0OZWxNLwe7aZxpWGgjU4g",
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Name:     "GetWallet",
			Method:   "GET",
			Endpoint: "/wallet/get?user_id=1",
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserWallet(gomock.Any(), gomock.Any()).Return(postgres.GetUserWalletRow{}, nil)
			},
			Header: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYwOTcyODEsInVzZXJJRCI6Mn0.0-xjEz6Eu5f0W_RVpFGLIR0OZWxNLwe7aZxpWGgjU4g",
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Name:     "Withdraw",
			Method:   "PUT",
			Endpoint: "/wallet/withdraw",
			Payload:  `{"user_wallet_id": 1, "amount":200.0}`,
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserWalletByID(gomock.Any(), gomock.Any()).Return(postgres.GetUserWalletByIDRow{
					Amount: pgtype.Float8{Float64: 300.0},
				}, nil)
				store.EXPECT().GetUserWalletByID(gomock.Any(), gomock.Any()).Return(postgres.GetUserWalletByIDRow{
					Amount: pgtype.Float8{Float64: 100.0},
				}, nil)
				store.EXPECT().UpdateUserWalletAmount(gomock.Any(), gomock.Any()).Return(nil)
				store.EXPECT().CreateUserTransaction(gomock.Any(), gomock.Any()).Return(postgres.UserTransaction{}, nil)
			},
			Header: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYwOTcyODEsInVzZXJJRCI6Mn0.0-xjEz6Eu5f0W_RVpFGLIR0OZWxNLwe7aZxpWGgjU4g",
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Name:     "Deposit",
			Method:   "PUT",
			Endpoint: "/wallet/deposit",
			Payload:  `{"user_wallet_id": 1, "amount":200.0}`,
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserWalletByID(gomock.Any(), gomock.Any()).Return(postgres.GetUserWalletByIDRow{
					Amount: pgtype.Float8{Float64: 300.0},
				}, nil)
				store.EXPECT().GetUserWalletByID(gomock.Any(), gomock.Any()).Return(postgres.GetUserWalletByIDRow{
					Amount: pgtype.Float8{Float64: 100.0},
				}, nil)
				store.EXPECT().UpdateUserWalletAmount(gomock.Any(), gomock.Any()).Return(nil)
				store.EXPECT().CreateUserTransaction(gomock.Any(), gomock.Any()).Return(postgres.UserTransaction{}, nil)
			},
			Header: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYwOTcyODEsInVzZXJJRCI6Mn0.0-xjEz6Eu5f0W_RVpFGLIR0OZWxNLwe7aZxpWGgjU4g",
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response utils.ResponseBody
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Name:     "GetUserTransactions",
			Method:   "GET",
			Endpoint: "/wallet/get-user-Transactions?user_wallet_id=1",
			Expect: func(ctrl *gomock.Controller, store *mock_postgres.MockStore) {
				store.EXPECT().GetUserWalletTransactions(gomock.Any(), gomock.Any()).Return([]postgres.GetUserWalletTransactionsRow{
					postgres.GetUserWalletTransactionsRow{
						Amount:            pgtype.Float8{Float64: 100.0},
						Action:            "Despot",
						UserWalletID:      pgtype.Int4{Int32: 1},
						TransactionAmount: pgtype.Float8{Float64: 100.0},
						CreatedAt:         pgtype.Timestamp{Time: time.Now()},
						UpdatedAt:         pgtype.Timestamp{Time: time.Now()},

						// ... other fields
					},
					postgres.GetUserWalletTransactionsRow{
						Amount:            pgtype.Float8{Float64: 100.0},
						Action:            "Despot",
						UserWalletID:      pgtype.Int4{Int32: 1},
						TransactionAmount: pgtype.Float8{Float64: 100.0},
						CreatedAt:         pgtype.Timestamp{Time: time.Now()},
						UpdatedAt:         pgtype.Timestamp{Time: time.Now()},

						// ... other fields
					},
				}, nil)
			},
			Header: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYwOTcyODEsInVzZXJJRCI6Mn0.0-xjEz6Eu5f0W_RVpFGLIR0OZWxNLwe7aZxpWGgjU4g",
			},
			Verify: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// var response []postgres.GetUserWalletTransactionsRow
				// err := json.Unmarshal(recorder.Body.Bytes(), &response)
				// if err != nil {
				// 	t.Fatal(err)
				// }
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	godotenv.Load("../.env")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	conn := config.Config(ctx)
	store := mock_postgres.NewMockStore(ctrl)
	server := gin.Default()

	MyNewServer(server, conn, &ctx, store)
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			runTest(t, tc, ctrl, store, server)
		})
	}
}

// func TestCreateUser(t *testing.T) {
// 	godotenv.Load("../.env")
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	ctx := context.Background()
// 	conn := config.Config(ctx)
// 	store := mock_postgres.NewMockStore(ctrl)
// 	server := gin.Default()
// 	// t.Run("Get User", func(t *testing.T) {
// 	MyNewServer(server, conn, &ctx, store)
// 	store.EXPECT().GetUserByEmail(gomock.Any(), pgtype.Text{String: "test234234@gmail.com", Valid: true}).Return(postgres.User{}, errors.New("user not found"))
// 	store.EXPECT().CreateUser(gomock.Any(), gomock.Any())
// 	createUserPayload := `{"name": "test", "email": "test234234@gmail.com", "password": "123456"}`
// 	createUserReq, err := http.NewRequest("POST", "/users/create", bytes.NewBufferString(createUserPayload))
// 	if err != nil {
// 		t.Fatal("error creating create user request")
// 	}
// 	createUserRecorder := httptest.NewRecorder()
// 	server.ServeHTTP(createUserRecorder, createUserReq)
// 	// utils.CreateUserTest(createUserRecorder, t)
// 	var response utils.ResponseBody
// 	err = json.Unmarshal(createUserRecorder.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, http.StatusOK, createUserRecorder.Code)
// 	// })
// }
// func TestGetUser(t *testing.T) {
// 	godotenv.Load("../.env")
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	ctx := context.Background()
// 	conn := config.Config(ctx)
// 	store := mock_postgres.NewMockStore(ctrl)
// 	server := gin.Default()
// 	// t.Run("Get User", func(t *testing.T) {
// 	MyNewServer(server, conn, &ctx, store)
// 	createUserReq, err := http.NewRequest("GET", "/users/get?id=1", nil)
// 	if err != nil {
// 		t.Fatal("error creating create user request")
// 	}
// 	store.EXPECT().GetUserByID(gomock.Any(), gomock.Any())
// 	createUserRecorder := httptest.NewRecorder()
// 	server.ServeHTTP(createUserRecorder, createUserReq)
// 	// utils.CreateUserTest(createUserRecorder, t)

// 	var response utils.ResponseBody
// 	err = json.Unmarshal(createUserRecorder.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, http.StatusOK, createUserRecorder.Code)
// 	// })
// }
// func TestSignIn(t *testing.T) {
// 	godotenv.Load("../.env")
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	ctx := context.Background()
// 	conn := config.Config(ctx)
// 	store := mock_postgres.NewMockStore(ctrl)
// 	server := gin.Default()
// 	// t.Run("Get User", func(t *testing.T) {
// 	MyNewServer(server, conn, &ctx, store)
// 	store.EXPECT().GetUserByEmail(gomock.Any(), pgtype.Text{String: "test234234@gmail.com", Valid: true}).Return(postgres.User{Password: pgtype.Text{String: "$2a$10$vt4xv8SE.B6O80BVzIupvus/y2Q7IEZLBz79sa2/MnJXeW23wiG6q"}}, nil)
// 	// store.EXPECT().CreateUser(gomock.Any(), gomock.Any())
// 	signInUserPayload := `{"email": "test234234@gmail.com", "password": "123456"}`
// 	signInUserReq, err := http.NewRequest("POST", "/users/sign-in", bytes.NewBufferString(signInUserPayload))
// 	if err != nil {
// 		t.Fatal("error creating create user request")
// 	}
// 	createUserRecorder := httptest.NewRecorder()
// 	server.ServeHTTP(createUserRecorder, signInUserReq)
// 	var response utils.ResponseBody
// 	err = json.Unmarshal(createUserRecorder.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, http.StatusOK, createUserRecorder.Code)
// 	// })
// }
// func TestGetAllUser(t *testing.T) {
// 	godotenv.Load("../.env")
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	ctx := context.Background()
// 	conn := config.Config(ctx)
// 	store := mock_postgres.NewMockStore(ctrl)
// 	server := gin.Default()
// 	// t.Run("Get User", func(t *testing.T) {
// 	MyNewServer(server, conn, &ctx, store)
// 	getAllUserReq, err := http.NewRequest("GET", "/users/get-all-users", nil)
// 	if err != nil {
// 		t.Fatal("error creating create user request")
// 	}
// 	getAllUserReq.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYwOTcyODEsInVzZXJJRCI6Mn0.0-xjEz6Eu5f0W_RVpFGLIR0OZWxNLwe7aZxpWGgjU4g")
// 	store.EXPECT().ListUsers(gomock.Any())
// 	createUserRecorder := httptest.NewRecorder()
// 	server.ServeHTTP(createUserRecorder, getAllUserReq)
// 	var response utils.ResponseBody
// 	err = json.Unmarshal(createUserRecorder.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, http.StatusOK, createUserRecorder.Code)
// 	// })
// }
// func TestSendOtp(t *testing.T) {
// 	godotenv.Load("../.env")
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	ctx := context.Background()
// 	conn := config.Config(ctx)
// 	store := mock_postgres.NewMockStore(ctrl)
// 	server := gin.Default()
// 	// t.Run("Get User", func(t *testing.T) {
// 	MyNewServer(server, conn, &ctx, store)
// 	createUserReq, err := http.NewRequest("POST", "/users/send-otp?id=1", nil)
// 	if err != nil {
// 		t.Fatal("error creating create user request")
// 	}
// 	store.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(postgres.User{Email: pgtype.Text{String: "test@gmail.com"}}, nil)
// 	store.EXPECT().UpdateOTP(gomock.Any(), gomock.Any()).Return(postgres.User{Email: pgtype.Text{String: "test@gmail.com"}}, nil)
// 	createUserRecorder := httptest.NewRecorder()
// 	server.ServeHTTP(createUserRecorder, createUserReq)
// 	var response utils.ResponseBody
// 	err = json.Unmarshal(createUserRecorder.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, http.StatusOK, createUserRecorder.Code)
// 	// })
// }
// func TestVerifyOtp(t *testing.T) {
// 	godotenv.Load("../.env")
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	ctx := context.Background()
// 	conn := config.Config(ctx)
// 	store := mock_postgres.NewMockStore(ctrl)
// 	server := gin.Default()
// 	// t.Run("Get User", func(t *testing.T) {
// 	MyNewServer(server, conn, &ctx, store)
// 	createUserReq, err := http.NewRequest("POST", "/users/verify-otp?id=1&otp=123456", nil)
// 	if err != nil {
// 		t.Fatal("error creating create user request")
// 	}
// 	store.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(postgres.User{Otp: pgtype.Int4{Int32: 123456, Valid: true}}, nil)
// 	store.EXPECT().UpdateOTP(gomock.Any(), gomock.Any()).Return(postgres.User{Email: pgtype.Text{String: "test@gmail.com"}}, nil)
// 	createUserRecorder := httptest.NewRecorder()
// 	server.ServeHTTP(createUserRecorder, createUserReq)
// 	var response utils.ResponseBody
// 	err = json.Unmarshal(createUserRecorder.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, http.StatusOK, createUserRecorder.Code)
// 	// })
// }
