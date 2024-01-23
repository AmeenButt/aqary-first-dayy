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

type PropertiesHanlder struct {
	conn    *pgx.Conn
	ctx     *context.Context
	queries postgres.Store
}

func GetPropertiesHandlers(conn *pgx.Conn, ctx *context.Context, store postgres.Store) *PropertiesHanlder {
	return &PropertiesHanlder{conn: conn, ctx: ctx, queries: store}
}

func (p *PropertiesHanlder) Add(c *gin.Context) {
	data := &models.CreateProperty{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetBindErrorMessage(err)})
		return
	}
	property, err := p.queries.InsertProperty(*p.ctx, postgres.InsertPropertyParams{
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

func (p *PropertiesHanlder) Update(c *gin.Context) {
	data := &models.UpdateProperty{}
	var updatedProperty interface{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetBindErrorMessage(err)})
		return
	}
	_, err := p.queries.GetPropertyByID(*p.ctx, data.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetBindErrorMessage(err)})
		return
	}
	if data.Demand != "" {
		updatedProperty, err = p.queries.UpdatePropertyDemand(*p.ctx, postgres.UpdatePropertyDemandParams{
			ID:     data.ID,
			Demand: pgtype.Text{String: data.Demand, Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
			return
		}
	}
	if data.Location != "" {
		updatedProperty, err = p.queries.UpdatePropertyLocation(*p.ctx, postgres.UpdatePropertyLocationParams{
			ID:       data.ID,
			Location: pgtype.Text{String: data.Location, Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
			return
		}
	}
	if data.SizeInSqFeet != 0 {
		updatedProperty, err = p.queries.UpdatePropertySize(*p.ctx, postgres.UpdatePropertySizeParams{
			ID:           data.ID,
			Sizeinsqfeet: pgtype.Int4{Int32: int32(data.SizeInSqFeet), Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
			return
		}
	}
	if data.SizeInSqFeet != 0 {
		updatedProperty, err = p.queries.UpdatePropertySize(*p.ctx, postgres.UpdatePropertySizeParams{
			ID:           data.ID,
			Sizeinsqfeet: pgtype.Int4{Int32: int32(data.SizeInSqFeet), Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
			return
		}
	}
	if data.Images != nil {
		updatedProperty, err = p.queries.UpdatePropertyImages(*p.ctx, postgres.UpdatePropertyImagesParams{
			ID:     data.ID,
			Images: data.Images,
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
			return
		}
	}
	if data.Status != "" {
		updatedProperty, err = p.queries.UpdateStatus(*p.ctx, postgres.UpdateStatusParams{
			ID:     data.ID,
			Status: pgtype.Text{String: data.Status, Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": utils.GetErrorMessage(err)})
			return
		}
	}

	// updatedProperty, err := p.queries.UpdateProperty(*p.ctx, postgres.UpdatePropertyParams{
	// 	ID:           data.ID,
	// 	Sizeinsqfeet: pgtype.Int4{Int32: int32(data.SizeInSqFeet), Valid: true},
	// 	Location:     pgtype.Text{String: data.Location, Valid: true},
	// 	Status:       pgtype.Text{String: data.Status, Valid: true},
	// 	Images:       data.Images,
	// 	Demand:       pgtype.Text{String: data.Demand, Valid: true},
	// })
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property Updated", "result": updatedProperty})
}

func (p *PropertiesHanlder) GetByID(c *gin.Context) {
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
	foundProperty, err := p.queries.GetPropertyByID(*p.ctx, int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	result := utils.ParsePropertyData(foundProperty)
	c.JSON(http.StatusOK, gin.H{"message": "Property Fetched", "result": result})
}

func (p *PropertiesHanlder) GetByUserID(c *gin.Context) {
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
	foundProperties, err := p.queries.GetPropertyByUserID(*p.ctx, pgtype.Int4{Int32: int32(user_id), Valid: true})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	result := utils.ParsePropertyDataArray(foundProperties)
	c.JSON(http.StatusOK, gin.H{"message": "PropertiesHanlder Fetched", "result": result})
}

func (p *PropertiesHanlder) DeleteProperty(c *gin.Context) {
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
	deletedProperty, err := p.queries.DeleteProperty(*p.ctx, int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "property with this id does not exsists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Property Deleted", "result": deletedProperty})
}
