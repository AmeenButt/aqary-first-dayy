package main

import (
	"context"
	"net/http"
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

	// Path to uploads folder
	uploadsFolderPath := "uploads"

	// Serve the uploads folder statically
	server.Static("/uploads", uploadsFolderPath)

	// Default route
	server.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"Message": "Server is running"}) })

	// User Routes
	routes.RegisterUserRoutes(server, conn)

	// User wallet Routes
	routes.RegisterUserWalletRoutes(server, conn)

	// Property Route
	routes.RegisterPropertiesRoutes(server, conn)

	// Starting server
	server.Run(os.Getenv("PORT"))

	// Close DB connection after job is done
	defer conn.Close(ctx)
}
