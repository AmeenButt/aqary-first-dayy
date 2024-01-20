package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"assesment.sqlc.dev/app/models"
	"assesment.sqlc.dev/app/postgres"
	"assesment.sqlc.dev/app/utils"
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
	data := &models.Property{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusNoContent, gin.H{"error": "No content found"})
		return
	}
	if data.SizeInSqFeet == 0 || data.Location == "" || data.Demand == "" || data.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No content found"})
		return
	}
	property, err := queries.InsertProperty(context.Background(), postgres.InsertPropertyParams{
		Sizeinsqfeet: pgtype.Int4{Int32: int32(data.SizeInSqFeet), Valid: true},
		Location:     pgtype.Text{String: data.Location, Valid: true},
		Demand:       pgtype.Text{String: data.Demand, Valid: true},
		UserID:       pgtype.Int4{Int32: int32(data.UserId), Valid: true},
		Images:       data.Images,
	})
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while adding property"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property Added", "result": property})
}

func (p *Properties) Update(c *gin.Context) {
	queries := postgres.New(p.conn)
	data := &models.Property{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusNoContent, gin.H{"error": "No content found"})
		return
	}
	_, err := queries.GetPropertyByID(context.Background(), data.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property with this id does not exsists"})
		return
	}
	updatedProperty, err := queries.UpdateProperty(context.Background(), postgres.UpdatePropertyParams{
		ID:           data.ID,
		Sizeinsqfeet: pgtype.Int4{Int32: int32(data.SizeInSqFeet), Valid: true},
		Location:     pgtype.Text{String: data.Location, Valid: true},
		Status:       pgtype.Text{String: data.Status, Valid: true},
		Images:       data.Images,
		Demand:       pgtype.Text{String: data.Demand, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating properties"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property Updated", "result": updatedProperty})
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
		Status: pgtype.Text{String: data.Status, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no property found with this id"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property status updated", "result": updatedData})
}

func (p *Properties) GetByID(c *gin.Context) {
	queries := postgres.New(p.conn)
	s_id := c.Query("id")
	id, err := strconv.Atoi(s_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can not get query parameter id"})
		return
	}
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required and can not be 0"})
		return
	}
	foundProperty, err := queries.GetPropertyByID(context.Background(), int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "property with this id does not exsists"})
		return
	}
	result := utils.ParsePropertyData(foundProperty)
	c.JSON(http.StatusOK, gin.H{"message": "Property Fetched", "result": result})
}

func (p *Properties) GetByUserID(c *gin.Context) {
	queries := postgres.New(p.conn)
	s_id := c.Query("user_id")
	user_id, err := strconv.Atoi(s_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can not get query parameter user_id"})
		return
	}
	if user_id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required and can not be 0"})
		return
	}
	foundProperties, err := queries.GetPropertyByUserID(context.Background(), pgtype.Int4{Int32: int32(user_id), Valid: true})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "property with this user_id does not exsists"})
		return
	}
	result := utils.ParsePropertyDataArray(foundProperties)
	c.JSON(http.StatusOK, gin.H{"message": "Properties Fetched", "result": result})
}

func (p *Properties) DeleteProperty(c *gin.Context) {
	queries := postgres.New(p.conn)
	s_id := c.Query("id")
	id, err := strconv.Atoi(s_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can not get query parameter id"})
		return
	}
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required and can not be 0"})
		return
	}
	deletedProperty, err := queries.DeleteProperty(context.Background(), int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "property with this id does not exsists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property Deleted", "result": deletedProperty})
}
