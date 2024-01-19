package main

import (
	"context"
	"os"

	"assesment.sqlc.dev/app/config"
	"assesment.sqlc.dev/app/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// Loading ENV file
	godotenv.Load()

	// Gin object creation
	server := gin.Default()

	// initialing database
	conn := config.Config(ctx)

	uploadsFolderPath := "uploads"

	// Serve the uploads folder statically
	server.Static("/uploads", uploadsFolderPath)

	// User Routes
	routes.RegisterUserRoutes(server, conn)

	// User wallet Routes
	routes.RegisterUserWalletRoutes(server, conn)

	// Starting server
	server.Run(os.Getenv("PORT"))

	// Close DB connection after job is done
	defer conn.Close(ctx)
}
