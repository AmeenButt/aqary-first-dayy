package routes

import (
	"context"

	"assesment.sqlc.dev/app/handlers"
	"assesment.sqlc.dev/app/middleware"
	"assesment.sqlc.dev/app/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RegisterUserRoutes(router *gin.Engine, conn *pgx.Conn, ctx *context.Context, store postgres.Store) {
	server := router.Group("/users")
	{
		UserHanlder := handlers.CreateUserHanlder(conn, ctx, store)
		server.POST("/create", UserHanlder.CreateUser)
		server.POST("/sign-in", UserHanlder.SignIn)
		server.GET("/get", UserHanlder.GetUser)
		server.GET("/get-all-users", middleware.AuthMiddleware(), UserHanlder.GetAllUser)
		server.POST("/upload-profile", middleware.AuthMiddleware(), UserHanlder.UploadProfilePicture)
		server.POST("/send-otp", UserHanlder.SendOtp)
		server.POST("/verify-otp", UserHanlder.VerifyOtp)
	}
}
