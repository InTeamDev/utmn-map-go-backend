package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	mapentites "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
)

type MapService interface {
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
		api.PATCH("/objects/:object_id", p.UpdateObjectHandler)
	}
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
