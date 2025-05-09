package handler

import (
	"context"
	"net/http"

	mapentites "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		Name        *string  `json:"name"`
		Alias       *string  `json:"alias"`
		Description *string  `json:"description"`
		X           *float64 `json:"x"`
		Y           *float64 `json:"y"`
		Width       *float64 `json:"width"`
		Height      *float64 `json:"height"`
		ObjectType  *string  `json:"object_type"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	updatedObj := mapentites.UpdateObjectInput{
		ID:          objectID,
		Name:        input.Name,
		Alias:       input.Alias,
		Description: input.Description,
		X:           input.X,
		Y:           input.Y,
		Width:       input.Width,
		Height:      input.Height,
	}

	if input.ObjectType != nil {
		objType := mapentites.ObjectType(*input.ObjectType)
		updatedObj.ObjectType = &objType
	}

	result, err := p.mapService.UpdateObject(c.Request.Context(), updatedObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
