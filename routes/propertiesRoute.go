package routes

import (
	"assesment.sqlc.dev/app/handlers"
	"assesment.sqlc.dev/app/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RegisterPropertiesRoutes(router *gin.Engine, conn *pgx.Conn) {
	server := router.Group("/property")
	{
		propertyHandlers := handlers.GetPropertiesHandlers(conn)
		server.POST("/add", middleware.AuthMiddleware(), propertyHandlers.Add)
		server.PUT("/update", middleware.AuthMiddleware(), propertyHandlers.Update)
		server.PUT("/update-status", middleware.AuthMiddleware(), propertyHandlers.UpdateStatus)
		server.GET("/get-by-id", propertyHandlers.GetByID)
		server.GET("/get-by-user-id", propertyHandlers.GetByUserID)
		server.DELETE("/delete", middleware.AuthMiddleware(), propertyHandlers.DeleteProperty)
	}
}