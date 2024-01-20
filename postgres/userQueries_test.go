package postgres

import (
	"context"
	"testing"

	"assesment.sqlc.dev/app/config"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

var u_id int
var u_w_id int
var u_t_id int
var email = "test1@gmail.com"

// USER QUERIES TEST
func TestCreateUser(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	userData, err := queries.CreateUser(context.Background(), CreateUserParams{
		Name:     "Test",
		Email:    pgtype.Text{String: email, Valid: true},
		Password: pgtype.Text{String: "123456", Valid: true},
	})
	if err != nil {
		t.Fatal(err)
		t.Fatal("error in creating user")
	}
	u_id = int(userData.ID)
	assert.Equal(t, userData.Name, "Test")
	assert.Equal(t, userData.Email.String, email)
	assert.NotEmpty(t, userData.Password)
}

func TestGetUserByID(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	userData, err := queries.GetUserByID(context.Background(), int64(u_id))
	if err != nil {
		t.Fatal("error in fetching user user")
	}
	assert.Equal(t, userData.Name, "Test")
	assert.Equal(t, userData.Email.String, email)
	assert.NotEmpty(t, userData.Password)
}

func TestListUsers(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	_, err := queries.ListUsers(context.Background())
	if err != nil {
		t.Fatal("error in fetching all users")
	} else {
		t.Log("All Users Fetched")
	}
}

func TestUpdateUser(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	_, err := queries.UpdateUser(context.Background(), UpdateUserParams{
		ID:   int64(u_id),
		Name: "test_update",
	})
	if err != nil {
		t.Fatal("error in updating user user")
	} else {
		t.Log("User Updated")
	}
}

// USER WALLET QUERIES TEST
func TestCreateUserWallet(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	userWallet, err := queries.CreateUserWallet(context.Background(), CreateUserWalletParams{
		UserID: pgtype.Int4{Int32: int32(u_id), Valid: true},
		Amount: pgtype.Float8{Float64: 0.0, Valid: true},
	})
	if err != nil {
		t.Fatal("error in creating wallet")
	}
	u_w_id = int(userWallet.ID)
	assert.Equal(t, userWallet.UserID.Int32, int32(u_id))
	assert.Equal(t, userWallet.Amount.Float64, 0.0)
	assert.NotEmpty(t, userWallet.ID)
}

func TestGetUserWallet(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	userWalletData, err := queries.GetUserWallet(context.Background(), pgtype.Int4{Int32: int32(u_id), Valid: true})
	if err != nil {
		t.Fatal(err)
		t.Fatal("error in getting wallet")
	}
	u_w_id = int(userWalletData.ID)
	assert.Equal(t, userWalletData.UserID.Int32, int32(u_id))
	assert.Equal(t, userWalletData.Amount.Float64, 0.0)
	assert.NotEmpty(t, userWalletData.ID)
}

func TestGetUserWalletByID(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	userWalletData, err := queries.GetUserWalletByID(context.Background(), int64(u_w_id))
	if err != nil {
		t.Fatal(err)
		t.Fatal("error in getting wallet")
	}
	assert.Equal(t, userWalletData.UserID.Int32, int32(u_id))
	assert.Equal(t, userWalletData.Amount.Float64, 0.0)
	assert.Equal(t, userWalletData.ID, int64(u_w_id))
}

func TestListUserWallets(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	_, err := queries.ListUserWallets(context.Background())
	if err != nil {
		t.Fatal("error in fetching user user")
	} else {
		t.Log("User Created")
	}
}

func TestUpdateUserWalletAmount(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	err := queries.UpdateUserWalletAmount(context.Background(), UpdateUserWalletAmountParams{
		ID:     int64(u_w_id),
		Amount: pgtype.Float8{Float64: 500.0, Valid: true},
	})
	if err != nil {
		t.Fatal("error in updating user user")
	} else {
		t.Log("User Created")
	}
}

// TRANSACTION QUERIES TEST
func TestCreateUserTransaction(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	userTransactionData, err := queries.CreateUserTransaction(context.Background(), CreateUserTransactionParams{
		UserWalletID:      pgtype.Int4{Int32: int32(u_w_id), Valid: true},
		TransactionAmount: pgtype.Float8{Float64: 500.0, Valid: true},
		Action:            "Deposit",
	})
	if err != nil {
		t.Fatal(err)
		t.Fatal("error in creating user")
	}
	u_t_id = int(userTransactionData.ID)
	assert.Equal(t, userTransactionData.UserWalletID.Int32, int32(u_w_id))
	assert.Equal(t, userTransactionData.TransactionAmount.Float64, 500.0)
	assert.NotEmpty(t, userTransactionData.ID)
}

func TestGetUserWalletTransactions(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	_, err := queries.GetUserWalletTransactions(context.Background(), pgtype.Int4{Int32: int32(u_w_id), Valid: true})
	if err != nil {
		t.Fatal("error in getting user transactions")
	} else {
		t.Log("User Transaction fetched")
	}
}

func TestListTransactions(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	_, err := queries.ListTransactions(context.Background())
	if err != nil {
		t.Fatal("error in fetching transactions")
	} else {
		t.Log("Transactions Created")
	}
}

// DELETE QUERIES
func TestDeleteUserTransaction(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	err := queries.DeleteUserTransaction(context.Background(), int64(u_t_id))
	if err != nil {
		t.Fatal("error in deleting user user")
	} else {
		t.Log("User Created")
	}
}

func TestDeleteUserWallet(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	err := queries.DeleteUserWallet(context.Background(), int64(u_w_id))
	if err != nil {
		t.Fatal("error in deleting user wallet")
	} else {
		t.Log("Wallet Deleted")
	}
}

func TestDeleteUser(t *testing.T) {
	conn := config.TestConfig(context.Background())
	defer conn.Close(context.Background())
	queries := New(conn)
	err := queries.DeleteUser(context.Background(), int64(u_id))
	if err != nil {
		t.Fatal("error in deleting user user")
	} else {
		t.Log("User Deleted")
	}
}
