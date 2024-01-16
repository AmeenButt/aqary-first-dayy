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
		userHanlder := handlers.CreateUserHanlder(conn)
		server.POST("/create", userHanlder.CreateUser)
		server.POST("/sign-in", userHanlder.SignIn)
		server.GET("/get", userHanlder.GetUser)
		server.GET("/get-all-users", middleware.AuthMiddleware(), userHanlder.GetAllUser)
	}
}
