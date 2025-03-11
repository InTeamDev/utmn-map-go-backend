package handler

import (
	"context"
	"net/http"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/gin-gonic/gin"
)

type MapService interface {
	GetObjects(ctx context.Context, req entities.GetObjectsRequest) ([]entities.Object, error)
}

type PublicAPI struct {
	mapService MapService
}

func NewPublicAPI(mapService MapService) *PublicAPI {
	return &PublicAPI{mapService: mapService}
}

func (p *PublicAPI) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/objects", p.GetObjectsHandler)
	}
}

func (p *PublicAPI) GetObjectsHandler(c *gin.Context) {
	objects, err := p.mapService.GetObjects(c.Request.Context(), entities.GetObjectsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"objects": objects})
}
