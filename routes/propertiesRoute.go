package routes

import (
	"assesment.sqlc.dev/app/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func PropertiesRoutes(router *gin.Engine, conn *pgx.Conn) {
	server := router.Group("/users")
	{
		propertyHandler := handlers.GetPropertiesHandlers(conn)
		server.POST("/add", propertyHandler.Add)
		server.PUT("/update", propertyHandler.Update)
		server.PUT("/update-status", propertyHandler.UpdateStatus)
		server.GET("/get-by-id", propertyHandler.GetByID)
		server.GET("/get-by-user-id", propertyHandler.GetByUserID)
		server.DELETE("/delete", propertyHandler.DeleteProperty)
	}
}
