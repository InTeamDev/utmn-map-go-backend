package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	routeentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MapService interface {
	CreateObject(ctx context.Context, input mapentities.CreateObjectInput) (mapentities.Object, error)
	UpdateObject(ctx context.Context, input mapentities.UpdateObjectInput) (mapentities.Object, error)
	CreateBuilding(ctx context.Context, input mapentities.CreateBuildingInput) (mapentities.Building, error)
	DeleteBuilding(ctx context.Context, id uuid.UUID) error
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
		api.POST("/floors/:floor_id/objects", p.CreateObjectHandler)
		api.PATCH("/objects/:object_id", p.UpdateObjectHandler)
		api.POST("/buildings", p.CreateBuildingHandler)
		api.POST("/route/intersections", p.AddIntersection)
		api.POST("/route/connections", p.AddConnection)
		api.DELETE("/buildings/:building_id", p.DeleteBuildingHandler)
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

	createInput := mapentities.CreateObjectInput{
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
		switch {
		case errors.Is(err, mapentities.ErrFloorNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, mapentities.ErrObjectTypeNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, mapentities.ErrInvalidDimensions):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, mapentities.ErrPositionConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			log.Printf("Internal server error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": mapentities.ErrInternalServer.Error()})
		}
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

	updatedObj := mapentities.UpdateObjectInput{
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
		updatedObj.ObjectType = mapentities.ObjectType(*input.ObjectType)
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

	buildingInput := mapentities.CreateBuildingInput{
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
