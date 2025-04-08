package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	mapentites "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
)

const defaultPageLimit = 30

type MapService interface {
	GetBuildings(ctx context.Context) ([]mapentites.Building, error)
	GetFloors(ctx context.Context, buildID uuid.UUID) ([]mapentites.Floor, error)
	GetObjectCategories(ctx context.Context, buildID uuid.UUID) ([]mapentites.ObjectType, error)
	GetObjectsByBuilding(ctx context.Context, buildID uuid.UUID) ([]mapentites.Object, error)
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
		// GET получить все категории объектов корпуса (для фильтрации объектов на карте)
		api.GET("/buildings/:build_id/categories", p.GetObjectCategories)
		// GET поиск объектов
		api.GET("/buildings/:build_id/search", p.SearchHandler)
		// ADMIN API
		api.PATCH("/objects/:id", p.UpdateObjectHandler)
	}
}

func (p *PublicAPI) GetObjectsByBuildingHandler(c *gin.Context) {
	buildIDStr := c.Param("build_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid build_id"})
		return
	}

	objects, err := p.mapService.GetObjectsByBuilding(c.Request.Context(), buildID)
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
	buildIDStr := c.Param("build_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid build_id"})
		return
	}

	categories, err := p.mapService.GetObjectCategories(c.Request.Context(), buildID)
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

func (p *PublicAPI) UpdateObjectHandler(c *gin.Context) {
	// objectIDStr := c.Param("id")
	// objectID, err := uuid.Parse(objectIDStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id"})
	// 	return
	// }

	// var input struct {
	// 	Name        *string `json:"name"`
	// 	Alias       *string `json:"alias"`
	// 	Description *string `json:"description"`
	// 	ObjectType  *string `json:"object_type"`
	// }

	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
	// 	return
	// }

	// updatedObj := mapentites.UpdateObjectInput{
	// 	ID: objectID,
	// }

	// if input.Name != nil {
	// 	updatedObj.Name = *input.Name
	// }
	// if input.Alias != nil {
	// 	updatedObj.Alias = *input.Alias
	// }
	// if input.Description != nil {
	// 	updatedObj.Description = *input.Description
	// }
	// if input.ObjectType != nil {
	// 	updatedObj.ObjectType = mapentites.ObjectType(*input.ObjectType)
	// }

	// result, err := p.mapService.UpdateObject(c.Request.Context(), updatedObj)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"result": "ok"})
}
