package handlers

import (
	"context"
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
	ctx  *context.Context
}

func CreateWalletHanlder(conn *pgx.Conn, ctx *context.Context) *WalletHanlder {
	return &WalletHanlder{conn: conn, ctx: ctx}
}

func (w *WalletHanlder) Create(c *gin.Context) {
	queries := postgres.New(w.conn)
	data := &models.CreateUserWallet{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetBindErrorMessage(err)})
		return
	}
	_, err := queries.GetUserByID(*w.ctx, int64(data.UserID))
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	_, err = queries.GetUserWallet(*w.ctx, pgtype.Int4{Int32: data.UserID, Valid: true})
	if err == nil {

		c.JSON(http.StatusConflict, gin.H{"error": "User wallet already exsists"})
		return
	}
	insertedWallet, err := queries.CreateUserWallet(*w.ctx, postgres.CreateUserWalletParams{
		UserID: pgtype.Int4{Int32: data.UserID, Valid: true},
		Amount: pgtype.Float8{Float64: data.Amount, Valid: true},
	})
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.GetErrorMessage(err)})
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
	foundWallet, err := queries.GetUserWallet(*w.ctx, pgtype.Int4{Int32: int32(user_id), Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	result := utils.ParseUserWalletData(foundWallet)
	c.JSON(http.StatusOK, gin.H{"message": "Wallet Found Sucessfully", "result": result})
}

func (w *WalletHanlder) Withdraw(c *gin.Context) {
	tx, err := w.conn.Begin(*w.ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while starting a transaction"})
		return
	}

	// Defer a function to handle rollback if necessary
	defer func() {
		if err != nil {
			// Rollback only if there was an error
			if rollbackErr := tx.Rollback(*w.ctx); rollbackErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error rolling back transaction"})
				return
			}
		}
	}()

	queries := postgres.New(w.conn)
	data := &models.InputUserTransaction{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetBindErrorMessage(err)})
		return
	}
	userWallet, err := queries.GetUserWalletByID(*w.ctx, int64(data.UserWalletID))
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}

	if userWallet.Amount.Float64 < data.TransactionAmount {
		c.JSON(http.StatusNotFound, gin.H{"error": "Low Balance to withdraw"})
		return
	}

	err = utils.UpdateWalletAmount(data.UserWalletID, (userWallet.Amount.Float64 - data.TransactionAmount), queries)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}

	_, err = queries.CreateUserTransaction(*w.ctx, postgres.CreateUserTransactionParams{
		UserWalletID:      pgtype.Int4{Int32: data.UserWalletID, Valid: true},
		TransactionAmount: pgtype.Float8{Float64: data.TransactionAmount, Valid: true},
		Action:            "Withdraw",
	})

	if err != nil {

		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}

	userWallet, err = queries.GetUserWalletByID(*w.ctx, int64(data.UserWalletID))

	// Commit the transaction explicitly
	if err := tx.Commit(*w.ctx); err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Error while committing transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdraw Successful", "result": userWallet})
}

func (w *WalletHanlder) Deposit(c *gin.Context) {
	tx, err := w.conn.Begin(*w.ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while starting a transaction"})
		return
	}

	// Defer a function to handle rollback if necessary
	defer func() {
		if err != nil {
			// Rollback only if there was an error
			if rollbackErr := tx.Rollback(*w.ctx); rollbackErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error rolling back transaction"})
				return
			}
		}
	}()
	queries := postgres.New(w.conn)
	data := &models.InputUserTransaction{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetBindErrorMessage(err)})
		return
	}
	userWallet, err := queries.GetUserWalletByID(*w.ctx, int64(data.UserWalletID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	if utils.UpdateWalletAmount(data.UserWalletID, (userWallet.Amount.Float64+data.TransactionAmount), queries) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while desposit"})
		return
	}
	_, err = queries.CreateUserTransaction(*w.ctx, postgres.CreateUserTransactionParams{
		UserWalletID:      pgtype.Int4{Int32: data.UserWalletID, Valid: true},
		TransactionAmount: pgtype.Float8{Float64: data.TransactionAmount, Valid: true},
		Action:            "Deposit",
	})

	if err != nil {
		utils.UpdateWalletAmount(data.UserWalletID, (userWallet.Amount.Float64 - data.TransactionAmount), queries)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while performing transaction"})
		return
	}
	userWallet, err = queries.GetUserWalletByID(*w.ctx, int64(data.UserWalletID))
	// Commit the transaction explicitly
	if err := tx.Commit(*w.ctx); err != nil {
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
	foundUserTransactions, err := queries.GetUserWalletTransactions(*w.ctx, pgtype.Int4{Int32: int32(user_wallet_id), Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	result, err := utils.ParseUserTransactionData(foundUserTransactions)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No transactions found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transactions fetched", "result": result})
}
