package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	mapentites "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	routeentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MapService interface {
	UpdateObject(ctx context.Context, input mapentites.UpdateObjectInput) (mapentites.Object, error)
	CreateBuilding(ctx context.Context, input mapentites.CreateBuildingInput) (mapentites.Building, error)
	DeleteBuilding(ctx context.Context, id uuid.UUID) error
	UpdateBuilding(ctx context.Context, id uuid.UUID, input mapentites.UpdateBuildingInput) (mapentites.Building, error)
}

type RouteService interface {
	// GetRoute строит маршрут между точками
	// (первая точка - начальная, промежуточные, последняя - конечная).
	// Точки - ID Объектов.
	// BuildRoute(ctx context.Context, start uuid.UUID, end uuid.UUID, waypoints []uuid.UUID) ([]entities.Edge, error)
	// Admin. AddIntersection добавляет новый узел в граф.
	AddIntersection(ctx context.Context, x, y float64) (uuid.UUID, error)
	// Admin. AddConnection добавляет новое ребро в граф.
	AddConnection(ctx context.Context, fromID, toID uuid.UUID, weight float64) (routeentities.Edge, error)
	// Admin. DeleteNode удаляет узел из графа.
	// DeleteNode(ctx context.Context, id uuid.UUID) error
}

type AdminAPI struct {
	mapService   MapService
	routeService RouteService
}

func NewAdminAPI(mapService MapService, routeService RouteService) *AdminAPI {
	return &AdminAPI{mapService: mapService, routeService: routeService}
}

func (p *AdminAPI) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.PATCH("/objects/:object_id", p.UpdateObjectHandler)
		api.POST("/buildings", p.CreateBuildingHandler)
		api.POST("/route/intersections", p.AddIntersection)
		api.POST("/route/connections", p.AddConnection)
		api.DELETE("/buildings/:building_id", p.DeleteBuildingHandler)
		api.PATCH("/buildings/:building_id", p.UpdateBuilding)
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
	if input.X != nil {
		updatedObj.X = *input.X
	}
	if input.Y != nil {
		updatedObj.Y = *input.Y
	}
	if input.Width != nil {
		updatedObj.Width = *input.Width
	}
	if input.Height != nil {
		updatedObj.Height = *input.Height
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

func (p *AdminAPI) DeleteBuildingHandler(c *gin.Context) {
	idParam := c.Param("building_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	err = p.mapService.DeleteBuilding(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": "building not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (p *AdminAPI) CreateBuildingHandler(c *gin.Context) {
	var input struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	buildingInput := mapentites.CreateBuildingInput{
		Name:    input.Name,
		Address: input.Address,
	}
	building, err := p.mapService.CreateBuilding(c.Request.Context(), buildingInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, building)
}

func (p *AdminAPI) AddIntersection(c *gin.Context) {
	var input struct {
		X *float64 `json:"x"`
		Y *float64 `json:"y"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if input.X == nil || input.Y == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "x and y must be non-zero"})
		return
	}

	nodeID, err := p.routeService.AddIntersection(c.Request.Context(), *input.X, *input.Y)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, routeentities.AddIntersectionResponse{
		ID: nodeID,
	})
}

func (p *AdminAPI) AddConnection(c *gin.Context) {
	var input struct {
		FromID uuid.UUID `json:"from_id"`
		ToID   uuid.UUID `json:"to_id"`
		Weight float64   `json:"weight"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	result, err := p.routeService.AddConnection(c.Request.Context(), input.FromID, input.ToID, input.Weight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (p *AdminAPI) UpdateBuilding(c *gin.Context) {
	idParam := c.Param("building_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}
	var input mapentites.UpdateBuildingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	building, err := p.mapService.UpdateBuilding(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, building)
}
