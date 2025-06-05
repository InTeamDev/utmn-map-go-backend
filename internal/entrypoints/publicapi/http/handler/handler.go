package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	routeentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
)

const defaultPageLimit = 30

type MapService interface {
	GetBuildings(ctx context.Context) ([]mapentities.Building, error)
	GetBuildingByID(ctx context.Context, id uuid.UUID) (mapentities.Building, error)
	GetFloors(ctx context.Context, buildID uuid.UUID) ([]mapentities.Floor, error)
	GetObjectByID(ctx context.Context, objectID uuid.UUID) (mapentities.Object, error)
	GetObjectCategories(ctx context.Context) ([]mapentities.ObjectTypeInfo, error)
	GetObjectsByBuilding(ctx context.Context, buildID uuid.UUID) ([]mapentities.Object, error)
	GetObjectsResponse(ctx context.Context, buildingID uuid.UUID) (mapentities.GetObjectsResponse, error)
	GetDoors(ctx context.Context, buildID uuid.UUID) ([]mapentities.GetDoorsResponse, error)
}

type RouteService interface {
	BuildRoute(
		ctx context.Context,
		buildingID uuid.UUID,
		input routeentities.BuildRouteRequest,
	) ([]routeentities.Connection, error)
	GetIntersections(ctx context.Context, buildingID uuid.UUID) ([]routeentities.Node, error)
	GetDoors(ctx context.Context, buildingID uuid.UUID) ([]routeentities.Node, error)
	GetConnections(ctx context.Context, buildingID uuid.UUID) ([]routeentities.Connection, error)
}

type SearchService interface {
	Search(ctx context.Context, req searchentities.SearchRequest) ([]searchentities.SearchResult, error)
}

type PublicAPI struct {
	mapService    MapService
	searchService SearchService
	routeService  RouteService
}

func NewPublicAPI(mapService MapService, searchService SearchService, routeService RouteService) *PublicAPI {
	return &PublicAPI{mapService: mapService, searchService: searchService, routeService: routeService}
}

func (p *PublicAPI) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// buildings
		api.GET("/buildings", p.GetBuildingsHandler)
		api.GET("/buildings/:building_id", p.GetBuildingByIDHandler)
		// floors
		api.GET("/buildings/:building_id/floors", p.GetFloorsHandler)
		// objects
		api.GET("/buildings/:building_id/objects", p.GetObjectsByBuildingHandler)
		api.GET("/buildings/:building_id/objects/:object_id", p.GetObjectByIDHandler)
		// doors
		api.GET("/buildings/:building_id/doors", p.GetDoorsHandler)
		// route
		api.GET("/buildings/:building_id/intersections", p.GetIntersectionsHandler)
		api.GET("/buildings/:building_id/connections", p.GetConnectionsHandler)
		api.GET("/buildings/:building_id/graph/nodes", p.GetNodesHandler)
		api.POST("/buildings/:building_id/route", p.BuildRouteHandler)
		// polygons
		// search
		api.GET("/buildings/:building_id/search", p.SearchHandler)
		// categories
		api.GET("/categories", p.GetObjectCategories)
	}
}

func (p *PublicAPI) GetBuildingsHandler(c *gin.Context) {
	buildings, err := p.mapService.GetBuildings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"buildings": buildings})
}

func (p *PublicAPI) GetBuildingByIDHandler(c *gin.Context) {
	idParam := c.Param("building_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
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

func (p *PublicAPI) GetFloorsHandler(c *gin.Context) {
	buildIDStr := c.Param("building_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	floors, err := p.mapService.GetFloors(c.Request.Context(), buildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"floors": floors})
}

func (p *PublicAPI) GetObjectsByBuildingHandler(c *gin.Context) {
	buildIDStr := c.Param("building_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	objects, err := p.mapService.GetObjectsResponse(c.Request.Context(), buildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"objects": objects})
}

func (p *PublicAPI) GetObjectByIDHandler(c *gin.Context) {
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

func (p *PublicAPI) GetDoorsHandler(c *gin.Context) {
	buildIDStr := c.Param("building_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	doors, err := p.mapService.GetDoors(c.Request.Context(), buildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"doors": doors})
}

func (p *PublicAPI) GetIntersectionsHandler(c *gin.Context) {
	buildIDStr := c.Param("building_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	intersections, err := p.routeService.GetIntersections(c.Request.Context(), buildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"intersections": intersections})
}

func (p *PublicAPI) GetConnectionsHandler(c *gin.Context) {
	buildingID, err := uuid.Parse(c.Param("building_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	connections, err := p.routeService.GetConnections(c.Request.Context(), buildingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"connections": connections})
}

func (p *PublicAPI) GetNodesHandler(c *gin.Context) {
	buildIDStr := c.Param("building_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}
	intersections, err := p.routeService.GetIntersections(c.Request.Context(), buildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	doors, err := p.routeService.GetDoors(c.Request.Context(), buildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	nodes := make([]routeentities.Node, 0, len(intersections)+len(doors))
	nodes = append(nodes, doors...)
	nodes = append(nodes, intersections...)
	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

func (p *PublicAPI) BuildRouteHandler(c *gin.Context) {
	buildIDStr := c.Param("building_id")
	buildingID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	var input routeentities.BuildRouteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	edges, err := p.routeService.BuildRoute(c.Request.Context(), buildingID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"edges": edges})
}

func (p *PublicAPI) SearchHandler(c *gin.Context) {
	buildIDStr := c.Param("building_id")
	buildID, err := uuid.Parse(buildIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	userCategories := c.QueryArray("category")
	if len(userCategories) == 0 {
		userCategories = []string{}
	}

	categories := make([]mapentities.ObjectType, len(userCategories))
	for i, category := range userCategories {
		categories[i] = mapentities.ObjectType(category)
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

func (p *PublicAPI) GetObjectCategories(c *gin.Context) {
	categories, err := p.mapService.GetObjectCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
