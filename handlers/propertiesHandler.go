package handlers

import (
	"context"
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
	data := &models.InputProperty{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, size_in_feet, location and demand are required"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": utils.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property Added", "result": property})
}

func (p *Properties) Update(c *gin.Context) {
	queries := postgres.New(p.conn)
	data := &models.InputProperty{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusNoContent, gin.H{"error": "No content found"})
		return
	}
	_, err := queries.GetPropertyByID(context.Background(), data.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property Updated", "result": updatedProperty})
}

func (p *Properties) UpdateStatus(c *gin.Context) {
	queries := postgres.New(p.conn)
	data := &models.Property{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no content found"})
		return
	}
	updatedData, err := queries.UpdateStatus(context.Background(), postgres.UpdateStatusParams{
		ID:     data.ID,
		Status: pgtype.Text{String: data.Status, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": utils.GetErrorMessage(err)})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetErrorMessage(err)})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetErrorMessage(err)})
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
