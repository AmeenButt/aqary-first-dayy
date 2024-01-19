package handlers

import (
	"context"
	"log"
	"net/http"

	"assesment.sqlc.dev/app/models"
	"assesment.sqlc.dev/app/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Properties struct {
	conn *pgx.Conn
}

func GetPropertiesHandlers(conn *pgx.Conn) *Properties {
	return &Properties{conn: conn}
}

func (p *Properties) Add(c *gin.Context) {
	queries := postgres.New(p.conn)
	log.Fatal(queries)
	data := &models.Property{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusNoContent, gin.H{"error": "No content found"})
		return
	}
	property, err := queries.InsertProperty(context.Background(), postgres.InsertPropertyParams{
		Sizeinsqfeet: pgtype.Int4{Int32: int32(data.SizeInSqFeet), Valid: true},
		Location:     pgtype.Text{String: data.Location, Valid: true},
		Status:       pgtype.Text{String: data.Status, Valid: true},
		Demand:       pgtype.Text{String: data.Demand, Valid: true},
		UserID:       pgtype.Int4{Int32: int32(data.UserId), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while adding property"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property Added", "result": property})
}

func (p *Properties) Update(c *gin.Context) {

}

func (p *Properties) UpdateStatus(c *gin.Context) {
	queries := postgres.New(p.conn)
	data := &models.Property{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "no content found"})
		return
	}
	updatedData, err := queries.UpdateStatus(context.Background(), postgres.UpdateStatusParams{
		ID:     data.ID,
		Status: pgtype.Text{String: data.Status},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no property found with this id"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property status updated", "result": updatedData})
}
func (p *Properties) GetByID(c *gin.Context) {

}
func (p *Properties) GetByUserID(c *gin.Context) {

}
func (p *Properties) DeleteProperty(c *gin.Context) {

}
