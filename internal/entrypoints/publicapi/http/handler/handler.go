package handler

import (
	"context"
	"net/http"
	"strconv"

	searhentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MapService interface {
	GetObjects(ctx context.Context, req entities.GetObjectsRequest) ([]entities.Object, error)
	GetBuildings(ctx context.Context) ([]entities.Building, error)
	GetFloors(ctx context.Context, buildID uuid.UUID) ([]entities.Floor, error)
}
type SearchService interface {
	Search(ctx context.Context,
		query string,
		userFloor uuid.UUID,
		UserContext *searhentities.UserContext,
	) ([]searhentities.SearchResult, error)
}

type PublicAPI struct {
	mapService    MapService
	searchService SearchService
}

func NewPublicAPI(mapService MapService, searchService SearchService) *PublicAPI {
	return &PublicAPI{mapService: mapService, searchService: searchService}
}

func (p *PublicAPI) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/buildings", p.GetBuildingsHandler)
		api.GET("/buildings/:build_id/floors", p.GetFloorsHandler)
		api.GET("/buildings/:build_id/floors/:floor_id/objects", p.GetObjectsHandler)
		api.GET("/search", p.SearchHandler)
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

func (p *PublicAPI) SearchHandler(c *gin.Context) {
	query := c.Query("q")
	userFloor := c.Query("floor")
	floorID, err := uuid.Parse(userFloor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid floor_id"})
		return
	}

	var userContext *searhentities.UserContext
	if latStr := c.Query("lat"); latStr != "" {
		if lonStr := c.Query("lon"); lonStr != "" {
			lat, err := strconv.ParseFloat(latStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat"})
				return
			}
			lon, err := strconv.ParseFloat(lonStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lon"})
				return
			}
			userContext = &searhentities.UserContext{
				Location: &searhentities.Location{
					X: lat,
					Y: lon,
				},
			}
		}
	}

	results, err := p.searchService.Search(c.Request.Context(), query, floorID, userContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
