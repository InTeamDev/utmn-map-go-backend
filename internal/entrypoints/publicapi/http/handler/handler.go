package handler

import (
	"context"
	"net/http"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/gin-gonic/gin"
)

type MapService interface {
	GetObjects(ctx context.Context) []entities.Object
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
	objects := p.mapService.GetObjects(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"objects": objects})
}
