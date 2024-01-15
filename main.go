package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"tutorial.sqlc.dev/app/config"
	"tutorial.sqlc.dev/app/handlers"
)

func main() {
	ctx := context.Background()

	// Loading ENV file
	godotenv.Load()

	// Gin object creation
	server := gin.Default()

	// initialing database
	conn := config.Config(ctx)

	// Creating user handler object
	userHanlder := handlers.CreateUserHanlder(conn)

	// User Routes
	server.POST("/user/create", userHanlder.CreateUser)
	server.GET("/user/get", userHanlder.GetUser)

	// Starting server
	server.Run(os.Getenv("PORT"))

	// Close DB connection after job is done
	defer conn.Close(ctx)
}
