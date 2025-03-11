package handler

import (
	"context"
	"net/http"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MapService interface {
	GetObjects(ctx context.Context, req entities.GetObjectsRequest) ([]entities.Object, error)
	GetBuildings(ctx context.Context) ([]entities.Building, error)
	GetFloors(ctx context.Context, buildID uuid.UUID) ([]entities.Floor, error)
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
		api.GET("/buildings", p.GetBuildingsHandler)
		api.GET("/buildings/:build_id/floors", p.GetFloorsHandler)
		api.GET("/buildings/:build_id/floors/:floor_id/objects", p.GetObjectsHandler)
	}
}

func (p *PublicAPI) GetObjectsHandler(c *gin.Context) {
	buildIDStr := c.Param("build_id")
	floorIDStr := c.Param("floor_id")

	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid build_id"})
		return
	}

	floorID, err := uuid.Parse(floorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid floor_id"})
		return
	}

	request := entities.GetObjectsRequest{
		BuildID: buildID,
		FloorID: floorID,
	}

	objects, err := p.mapService.GetObjects(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"objects": objects})
}

func (p *PublicAPI) GetBuildingsHandler(c *gin.Context) {
	buildings, err := p.mapService.GetBuildings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"buildings": buildings})
}

func (p *PublicAPI) GetFloorsHandler(c *gin.Context) {
	buildIDStr := c.Param("build_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid build_id"})
		return
	}

	floors, err := p.mapService.GetFloors(c.Request.Context(), buildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"floors": floors})
}
