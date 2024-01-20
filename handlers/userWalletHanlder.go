package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"assesment.sqlc.dev/app/models"
	"assesment.sqlc.dev/app/postgres"
	"assesment.sqlc.dev/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type WalletHanlder struct {
	conn *pgx.Conn
}

func CreateWalletHanlder(conn *pgx.Conn) *WalletHanlder {
	return &WalletHanlder{conn: conn}
}

func (w *WalletHanlder) Create(c *gin.Context) {
	queries := postgres.New(w.conn)
	data := &models.UserWallet{}
	if err := c.ShouldBindJSON(data); err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusNoContent, gin.H{"error": "Body can not be empty"})
		return
	}

	if data.UserID < 1 || data.Amount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "one or both of user_id and amount are invalid"})
		return
	}
	fmt.Println(data.UserID)
	_, err := queries.GetUserByID(context.Background(), int64(data.UserID))
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User with this id does not exsists"})
		return
	}
	_, err = queries.GetUserWallet(context.Background(), pgtype.Int4{Int32: data.UserID, Valid: true})
	if err == nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusConflict, gin.H{"error": "User wallet already exsists"})
		return
	}
	insertedWallet, err := queries.CreateUserWallet(context.Background(), postgres.CreateUserWalletParams{
		UserID: pgtype.Int4{Int32: data.UserID, Valid: true},
		Amount: pgtype.Float8{Float64: data.Amount, Valid: true},
	})
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wallet of this user may already exsists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "wallet created", "result": insertedWallet})
}

func (w *WalletHanlder) GetWallet(c *gin.Context) {
	queries := postgres.New(w.conn)
	idStr := c.Query("user_id")
	user_id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	foundWallet, err := queries.GetUserWallet(context.Background(), pgtype.Int4{Int32: int32(user_id), Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet does not exsists"})
		return
	}
	result := utils.ParseUserWalletData(foundWallet)
	c.JSON(http.StatusOK, gin.H{"message": "Wallet Found Sucessfully", "result": result})
}

func (w *WalletHanlder) Withdraw(c *gin.Context) {
	tx, err := w.conn.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while starting a transaction"})
		return
	}

	// Defer a function to handle rollback if necessary
	defer func() {
		if err != nil {
			// Rollback only if there was an error
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				fmt.Printf("Error rolling back transaction: %v\n", rollbackErr)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error rolling back transaction"})
				return
			}
		}
	}()

	queries := postgres.New(w.conn)
	data := &models.UserTransaction{}
	if err := c.ShouldBindJSON(data); err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body can not be empty"})
		return
	}

	if data.TransactionAmount < 1 || data.UserWalletID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction_amount and user_wallet_id are required and should be greater than 0"})
		return
	}

	userWallet, err := queries.GetUserWalletByID(context.Background(), int64(data.UserWalletID))
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet with this id does not exist"})
		return
	}

	if userWallet.Amount.Float64 < data.TransactionAmount {
		c.JSON(http.StatusNotFound, gin.H{"error": "Low Balance to withdraw"})
		return
	}

	err = utils.UpdateWalletAmount(data.UserWalletID, (userWallet.Amount.Float64 - data.TransactionAmount), queries)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error occurred while withdrawal"})
		return
	}

	_, err = queries.CreateUserTransaction(context.Background(), postgres.CreateUserTransactionParams{
		UserWalletID:      pgtype.Int4{Int32: data.UserWalletID, Valid: true},
		TransactionAmount: pgtype.Float8{Float64: data.TransactionAmount, Valid: true},
		Action:            "Withdraw",
	})

	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Error while performing transaction"})
		return
	}

	userWallet, err = queries.GetUserWalletByID(context.Background(), int64(data.UserWalletID))

	// Commit the transaction explicitly
	if err := tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Error while committing transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdraw Successful", "result": userWallet})
}

func (w *WalletHanlder) Deposit(c *gin.Context) {
	tx, err := w.conn.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while starting a transaction"})
		return
	}

	// Defer a function to handle rollback if necessary
	defer func() {
		if err != nil {
			// Rollback only if there was an error
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				fmt.Printf("Error rolling back transaction: %v\n", rollbackErr)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error rolling back transaction"})
				return
			}
		}
	}()
	queries := postgres.New(w.conn)
	data := &models.UserTransaction{}
	if err := c.ShouldBindJSON(data); err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusNoContent, gin.H{"error": "Body can not be empty"})
		return
	}
	if data.TransactionAmount < 1 || data.UserWalletID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction_amount and user_wallet_id are required and should be greater then 0"})
		return
	}
	userWallet, err := queries.GetUserWalletByID(context.Background(), int64(data.UserWalletID))
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet with this id does not exsists"})
		return
	}
	if utils.UpdateWalletAmount(data.UserWalletID, (userWallet.Amount.Float64+data.TransactionAmount), queries) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while desposit"})
		return
	}
	_, err = queries.CreateUserTransaction(context.Background(), postgres.CreateUserTransactionParams{
		UserWalletID:      pgtype.Int4{Int32: data.UserWalletID, Valid: true},
		TransactionAmount: pgtype.Float8{Float64: data.TransactionAmount, Valid: true},
		Action:            "Deposit",
	})

	if err != nil {
		fmt.Printf("%v", err)
		utils.UpdateWalletAmount(data.UserWalletID, (userWallet.Amount.Float64 - data.TransactionAmount), queries)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while performing transaction"})
		return
	}
	userWallet, err = queries.GetUserWalletByID(context.Background(), int64(data.UserWalletID))
	// Commit the transaction explicitly
	if err := tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Error while committing transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Amount Deposited", "result": userWallet})
}

func (w *WalletHanlder) GetUserTransactions(c *gin.Context) {
	queries := postgres.New(w.conn)
	idStr := c.Query("user_wallet_id")
	user_wallet_id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	foundUserTransactions, err := queries.GetUserWalletTransactions(context.Background(), pgtype.Int4{Int32: int32(user_wallet_id), Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error while fetching transactions"})
		return
	}
	result, err := utils.ParseUserTransactionData(foundUserTransactions)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No transactions found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transactions fetched", "result": result})
}
