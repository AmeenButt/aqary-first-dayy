package main

import (
	"context"
	"net/http"
	"os"

	"assesment.sqlc.dev/app/config"
	"assesment.sqlc.dev/app/postgres"
	"assesment.sqlc.dev/app/routes"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

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

func main() {
	ctx := context.Background()

	// Loading ENV file
	godotenv.Load()

	// Gin object creation
	server := gin.Default()

	// initialing database
	conn := config.Config(ctx)

	// SQLC STORE WHICH IN THIS CASE IS QUERIES
	store := postgres.NewStore(conn)

	// REGISTER THE SERVER AND APIS
	MyNewServer(server, conn, &ctx, store)

	// Starting server
	server.Run(os.Getenv("PORT"))

	// Close DB connection after job is done
	defer conn.Close(ctx)
}
