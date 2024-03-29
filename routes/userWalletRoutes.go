package routes

import (
	"context"

	"assesment.sqlc.dev/app/handlers"
	"assesment.sqlc.dev/app/middleware"
	"assesment.sqlc.dev/app/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RegisterUserWalletRoutes(router *gin.Engine, conn *pgx.Conn, ctx *context.Context, store postgres.Store) {
	server := router.Group("/wallet")
	{
		walletHanlder := handlers.CreateWalletHanlder(conn, ctx, store)
		server.POST("/create", middleware.AuthMiddleware(), walletHanlder.Create)
		server.GET("/get", middleware.AuthMiddleware(), walletHanlder.GetWallet)
		server.PUT("/withdraw", middleware.AuthMiddleware(), walletHanlder.Withdraw)
		server.PUT("/deposit", middleware.AuthMiddleware(), walletHanlder.Deposit)
		server.GET("/get-user-Transactions", middleware.AuthMiddleware(), walletHanlder.GetUserTransactions)
	}
}
