package routes

import (
	"assesment.sqlc.dev/app/handlers"
	"assesment.sqlc.dev/app/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RegisterUserRoutes(router *gin.Engine, conn *pgx.Conn) {
	server := router.Group("/users")
	{
		UserHanlder := handlers.CreateUserHanlder(conn)
		server.POST("/create", UserHanlder.CreateUser)
		server.POST("/sign-in", UserHanlder.SignIn)
		server.GET("/get", UserHanlder.GetUser)
		server.GET("/get-all-users", middleware.AuthMiddleware(), UserHanlder.GetAllUser)
		server.POST("/upload-profile", middleware.AuthMiddleware(), UserHanlder.UploadProfilePicture)
		server.POST("/send-otp", UserHanlder.SendOtp)
		server.POST("/verify-otp", UserHanlder.VerifyOtp)
	}
}
