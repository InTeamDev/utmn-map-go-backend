package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	routeentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MapService interface {
	GetObjectByID(ctx context.Context, objectID uuid.UUID) (mapentities.Object, error)
	CreateObject(
		ctx context.Context,
		floorID uuid.UUID,
		input mapentities.CreateObjectInput,
	) (mapentities.Object, error)
	UpdateObject(ctx context.Context, id uuid.UUID, input mapentities.UpdateObjectInput) (mapentities.Object, error)
	DeleteObject(ctx context.Context, objectID uuid.UUID) error
	CreateBuilding(ctx context.Context, input mapentities.CreateBuildingInput) (mapentities.Building, error)
	DeleteBuilding(ctx context.Context, id uuid.UUID) error
	UpdateBuilding(
		ctx context.Context,
		id uuid.UUID,
		input mapentities.UpdateBuildingInput,
	) (mapentities.Building, error)
	CreatePolygon(ctx context.Context, floorID uuid.UUID, label string, zIndex int32) (mapentities.Polygon, error)
	CreatePolygonPoint(
		ctx context.Context,
		polygonID uuid.UUID,
		order int32,
		x, y float64,
	) (mapentities.PolygonPoint, error)
	DeletePolygonPoint(ctx context.Context, id int) error
	DeletePolygonPoints(ctx context.Context, ids []int) error
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
		api.GET("/buildings/:building_id/floors/:floor_id/objects/:object_id", p.GetObjectByIDHandler)
		api.POST("/buildings/:building_id/floors/:floor_id/objects", p.CreateObjectHandler)
		api.PATCH("/objects/:object_id", p.UpdateObjectHandler)
		api.POST("/buildings", p.CreateBuildingHandler)
		api.POST("/route/intersections", p.AddIntersection)
		api.POST("/route/connections", p.AddConnection)
		api.DELETE("/buildings/:building_id", p.DeleteBuildingHandler)
		api.DELETE("/buildings/:building_id/objects/:object_id", p.DeleteObjectHandler)
		api.PATCH("/buildings/:building_id", p.UpdateBuilding)
		api.POST("/floors/:floor_id/poligons", p.CreatePolygonHandler)
		api.POST("/floors/:floor_id/poligons/:p_id/points", p.CreatePolygonPointsHandler)
		api.DELETE("/floors/:floor_id/poligons/:polygon_id/points/:point_id", p.DeletePolygonPointHandler)
		api.DELETE("/floors/:floor_id/poligons/:polygon_id/points", p.DeletePolygonPointsHandler)
	}
}

func (p *AdminAPI) GetObjectByIDHandler(c *gin.Context) {
	objectID, err := uuid.Parse(c.Param("object_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object_id"})
		return
	}

	result, err := p.mapService.GetObjectByID(c.Request.Context(), objectID)
	if err != nil {
		if errors.Is(err, mapentities.ErrObjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (p *AdminAPI) CreateObjectHandler(c *gin.Context) {
	floorID, err := uuid.Parse(c.Param("floor_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid floor_id"})
		return
	}

	var input mapentities.CreateObjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := p.mapService.CreateObject(c.Request.Context(), floorID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (p *AdminAPI) DeleteObjectHandler(c *gin.Context) {
	objectIDStr := c.Param("object_id")
	objectID, err := uuid.Parse(objectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object_id"})
		return
	}

	err = p.mapService.DeleteObject(c.Request.Context(), objectID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": "object not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (p *AdminAPI) UpdateObjectHandler(c *gin.Context) {
	objectIDStr := c.Param("object_id")
	objectID, err := uuid.Parse(objectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object_id"})
		return
	}

	var input *mapentities.UpdateObjectInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	updatedObj := mapentities.UpdateObjectInput{
		Name:         input.Name,
		Alias:        input.Alias,
		Description:  input.Description,
		X:            input.X,
		Y:            input.Y,
		Width:        input.Width,
		Height:       input.Height,
		ObjectTypeID: input.ObjectTypeID,
	}

	result, err := p.mapService.UpdateObject(c.Request.Context(), objectID, updatedObj)
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

func (p *AdminAPI) UpdateBuilding(c *gin.Context) {
	idParam := c.Param("building_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}
	var input mapentities.UpdateBuildingInput
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

type CreatePolygonRequest struct {
	Label  string `json:"label"   binding:"required"`
	ZIndex int32  `json:"z_index"`
}

func (p *AdminAPI) CreatePolygonHandler(c *gin.Context) {
	floorID, err := uuid.Parse(c.Param("floor_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid floor_id"})
		return
	}

	var req CreatePolygonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	polygon, err := p.mapService.CreatePolygon(c.Request.Context(), floorID, req.Label, req.ZIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, polygon)
}

func (p *AdminAPI) CreatePolygonPointsHandler(c *gin.Context) {
	polygonID, err := uuid.Parse(c.Param("p_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid polygon_id"})
		return
	}

	var points []mapentities.PolygonPointRequest
	if err := c.ShouldBindJSON(&points); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result []mapentities.PolygonPoint
	for _, pt := range points {
		res, err := p.mapService.CreatePolygonPoint(c.Request.Context(), polygonID, pt.PointOrder, pt.X, pt.Y)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result = append(result, res)
	}

	c.JSON(http.StatusCreated, result)
}

func (p *AdminAPI) DeletePolygonPointHandler(c *gin.Context) {
	pointID, err := strconv.Atoi(c.Param("point_
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point_id"})
		return
	}

	if err := p.mapService.DeletePolygonPoint(c.Request.Context(), pointID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
func (p *AdminAPI) DeletePolygonPointsHandler(c *gin.Context) {
	var req mapentities.DeletePointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := p.mapService.DeletePolygonPoints(c.Request.Context(), req.Points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
