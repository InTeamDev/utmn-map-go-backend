package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	mapentites "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
)

type MapService interface {
	CreateObject(ctx context.Context, input mapentites.CreateObjectInput) (mapentites.Object, error)
	UpdateObject(ctx context.Context, input mapentites.UpdateObjectInput) (mapentites.Object, error)
}

type AdminAPI struct {
	mapService MapService
}

func NewAdminAPI(mapService MapService) *AdminAPI {
	return &AdminAPI{mapService: mapService}
}

func (p *AdminAPI) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/floors/:floor_id/objects", p.CreateObjectHandler)
		api.PATCH("/objects/:object_id", p.UpdateObjectHandler)
	}
}

func (p *AdminAPI) CreateObjectHandler(c *gin.Context) {
	floorIDStr := c.Param("floor_id")
	floorID, err := uuid.Parse(floorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid floor_id"})
		return
	}

	var input struct {
		Name         string  `json:"name" binding:"required,max=255"`
		Alias        string  `json:"alias" binding:"required,max=255"`
		Description  string  `json:"description" binding:"max=255"`
		X            float64 `json:"x" binding:"required"`
		Y            float64 `json:"y" binding:"required"`
		Width        float64 `json:"width" binding:"required,gte=1"`
		Height       float64 `json:"height" binding:"required,gte=1"`
		ObjectTypeID int32   `json:"object_type_id" binding:"required"`
		FloorID      int32   `json:"floor_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createInput := mapentites.CreateObjectInput{
		FloorID:      floorID,
		Name:         input.Name,
		Alias:        input.Alias,
		Description:  input.Description,
		X:            input.X,
		Y:            input.Y,
		Width:        input.Width,
		Height:       input.Height,
		ObjectTypeID: input.ObjectTypeID,
	}

	result, err := p.mapService.CreateObject(c.Request.Context(), createInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (p *AdminAPI) UpdateObjectHandler(c *gin.Context) {
	objectIDStr := c.Param("object_id")
	objectID, err := uuid.Parse(objectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object_id"})
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Alias       *string `json:"alias"`
		Description *string `json:"description"`
		ObjectType  *string `json:"object_type"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	updatedObj := mapentites.UpdateObjectInput{
		ID: objectID,
	}

	if input.Name != nil {
		updatedObj.Name = *input.Name
	}
	if input.Alias != nil {
		updatedObj.Alias = *input.Alias
	}
	if input.Description != nil {
		updatedObj.Description = *input.Description
	}
	if input.ObjectType != nil {
		updatedObj.ObjectType = mapentites.ObjectType(*input.ObjectType)
	}

	result, err := p.mapService.UpdateObject(c.Request.Context(), updatedObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
