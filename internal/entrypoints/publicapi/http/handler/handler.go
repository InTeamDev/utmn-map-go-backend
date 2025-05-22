package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	mapentites "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const defaultPageLimit = 30

type MapService interface {
	GetBuildings(ctx context.Context) ([]mapentites.Building, error)
	GetFloors(ctx context.Context, buildID uuid.UUID) ([]mapentites.Floor, error)
	GetObjectCategories(ctx context.Context) ([]mapentites.ObjectTypeInfo, error)
	GetObjectsByBuilding(ctx context.Context, buildID uuid.UUID) ([]mapentites.Object, error)
	GetObjectsResponse(ctx context.Context, buildingID uuid.UUID) (mapentites.GetObjectsResponse, error)
	GetBuildingByID(ctx context.Context, id uuid.UUID) (mapentites.Building, error)
}

type SearchService interface {
	Search(ctx context.Context, req searchentities.SearchRequest) ([]searchentities.SearchResult, error)
}

type PublicAPI struct {
	mapService    MapService
	searchService SearchService
}

func NewPublicAPI(mapService MapService, searchService SearchService) *PublicAPI {
	return &PublicAPI{mapService: mapService, searchService: searchService}
}

func (p *PublicAPI) RegisterRoutes(router *gin.Engine) {
	// GET получить информацию об объекте (для отображения информации об объекте)
	// GET получение маршрута

	api := router.Group("/api")
	{
		// GET получить все корпуса (просто список корпусов, id;name;address)
		api.GET("/buildings", p.GetBuildingsHandler)
		// GET получить все этажи корпуса
		api.GET("/buildings/:build_id/floors", p.GetFloorsHandler)
		// GET получить объекты этажа корпуса (для отрисовки объектов на карте)
		api.GET("/buildings/:build_id/objects", p.GetObjectsByBuildingHandler)
		// GET поиск объектов
		api.GET("/buildings/:build_id/search", p.SearchHandler)
		// GET получить все категории объектов корпуса (для фильтрации объектов на карте)
		api.GET("/categories", p.GetObjectCategories)
		api.GET("/buildings/:build_id", p.GetBuildingByIDHandler)
	}
}

func (p *PublicAPI) GetObjectsByBuildingHandler(c *gin.Context) {
	buildIDStr := c.Param("build_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid build_id"})
		return
	}

	objects, err := p.mapService.GetObjectsResponse(c.Request.Context(), buildID)
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

func (p *PublicAPI) GetObjectCategories(c *gin.Context) {
	categories, err := p.mapService.GetObjectCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func (p *PublicAPI) SearchHandler(c *gin.Context) {
	buildIDStr := c.Param("build_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	userCategories := c.QueryArray("category")
	if len(userCategories) == 0 {
		userCategories = []string{}
	}

	categories := make([]mapentites.ObjectType, len(userCategories))
	for i, category := range userCategories {
		categories[i] = mapentites.ObjectType(category)
	}

	query := c.Query("query")

	results, err := p.searchService.Search(c.Request.Context(), searchentities.SearchRequest{
		Query:      query,
		Limit:      defaultPageLimit,
		Offset:     0,
		BuildingID: buildID,
		Categories: categories,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (p *PublicAPI) GetBuildingByIDHandler(c *gin.Context) {
	idParam := c.Param("build_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid build_id"})
		return
	}

	result, err := p.mapService.GetBuildingByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "build not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
