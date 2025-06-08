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
)

type MapService interface {
	CreateBuilding(ctx context.Context, input mapentities.CreateBuildingInput) (mapentities.Building, error)
	DeleteBuilding(ctx context.Context, id uuid.UUID) error
	UpdateBuilding(
		ctx context.Context,
		id uuid.UUID,
		input mapentities.UpdateBuildingInput,
	) (mapentities.Building, error)
	GetBuildings(ctx context.Context) ([]mapentities.Building, error)

	CreateObject(
		ctx context.Context,
		floorID uuid.UUID,
		input mapentities.CreateObjectInput,
	) (mapentities.Object, error)
	UpdateObject(ctx context.Context, id uuid.UUID, input mapentities.UpdateObjectInput) (mapentities.Object, error)
	DeleteObject(ctx context.Context, objectID uuid.UUID) error
	GetObjectCategories(ctx context.Context) ([]mapentities.ObjectTypeInfo, error)
	GetObjectsResponse(ctx context.Context, buildingID uuid.UUID) (mapentities.GetObjectsResponse, error)

	GetDoor(
		ctx context.Context,
		buildingID uuid.UUID,
		doorID uuid.UUID,
	) (mapentities.Door, error)

	CreateFloor(ctx context.Context, buildingID uuid.UUID, floor mapentities.Floor) error
	CreateDoor(ctx context.Context, objectID uuid.UUID, door mapentities.Door) error
	GetPolygonByID(ctx context.Context, id uuid.UUID) (mapentities.FloorPolygon, error)

	CreatePolygon(ctx context.Context, polygon mapentities.Polygon) (mapentities.Polygon, error)
	CreatePolygonPoint(
		ctx context.Context,
		polygonID uuid.UUID,
		order int32,
		x, y float64,
	) (mapentities.PolygonPoint, error)
	DeletePolygonPoints(ctx context.Context, request mapentities.DeletePolygonPointsRequest) error
}

type RouteService interface {
	// GetRoute строит маршрут между точками
	// (первая точка - начальная, промежуточные, последняя - конечная).
	// Точки - ID Объектов.
	// BuildRoute(ctx context.Context, start uuid.UUID, end uuid.UUID, waypoints []uuid.UUID) ([]entities.Connection,
	// error)
	// Admin. AddIntersection добавляет новый узел в граф.
	AddIntersection(ctx context.Context, req routeentities.AddIntersectionRequest) (routeentities.Node, error)
	GetIntersections(ctx context.Context, buildingID uuid.UUID) ([]routeentities.Node, error)
	// Admin. AddConnection добавляет новое ребро в граф.
	AddConnection(ctx context.Context, fromID, toID uuid.UUID, weight float64) (routeentities.Connection, error)
	// Admin. DeleteNode удаляет узел из графа.
	// DeleteNode(ctx context.Context, id uuid.UUID) error
	GetConnections(ctx context.Context, buildingID uuid.UUID) ([]routeentities.Connection, error)
	DeleteIntersection(ctx context.Context, buildingID, intersectionID uuid.UUID) error
}

type AdminAPI struct {
	mapService   MapService
	routeService RouteService
}

func NewAdminAPI(mapService MapService, routeService RouteService) *AdminAPI {
	return &AdminAPI{mapService: mapService, routeService: routeService}
}

func (p *AdminAPI) RegisterRoutes(router *gin.Engine, m ...gin.HandlerFunc) {
	api := router.Group("/api", m...)
	{
		// buildings
		api.POST("/buildings", p.CreateBuildingHandler)
		api.PATCH("/buildings/:building_id", p.UpdateBuilding)
		api.DELETE("/buildings/:building_id", p.DeleteBuildingHandler)
		// TODO: floors post, patch and delete
		// objects
		api.POST("/buildings/:building_id/floors/:floor_id/objects", p.CreateObjectHandler)
		api.PATCH("/buildings/:building_id/floors/:floor_id/objects/:object_id", p.UpdateObjectHandler)
		api.DELETE("/buildings/:building_id/floors/:floor_id/objects/:object_id", p.DeleteObjectHandler)
		// TODO: doors post, patch and delete
		api.GET("/buildings/:building_id/doors/:door_id", p.GetDoorHandler)
		// route
		api.POST("/buildings/:building_id/route/intersections", p.AddIntersection)
		api.POST("/buildings/:building_id/route/connections", p.AddConnection)
		api.DELETE("/buildings/:building_id/intersections/:intersection_id", p.DeleteIntersectionHandler)
		// polygons
		api.POST("/buildings/:building_id/floors/:floor_id/poligons", p.CreatePolygonHandler)
		api.POST("/buildings/:building_id/floors/:floor_id/poligons/:p_id/points", p.CreatePolygonPointsHandler)
		api.DELETE("/buildings/:building_id/floors/:floor_id/poligons:poligon_id/points", p.DeletePolygonPointsHandler)
		// sync
		api.POST("/sync", p.SyncDatabaseHandler)
		api.GET("/sync", p.GetDatabaseHandler)
		api.GET("/floors/:floor_id/poligons/:poligon_id", p.GetPolygonByIDHandler)
	}
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

	c.JSON(http.StatusCreated, gin.H{"object": result})
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

	c.JSON(http.StatusOK, gin.H{"object": result})
}

func (p *AdminAPI) GetDoorHandler(c *gin.Context) {
	buildingID, err := uuid.Parse(c.Param("building_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	doorID, err := uuid.Parse(c.Param("door_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid door_id"})
		return
	}

	door, err := p.mapService.GetDoor(c.Request.Context(), buildingID, doorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "door not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, door)
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

	c.JSON(http.StatusOK, gin.H{"building": building})
}

func (p *AdminAPI) AddIntersection(c *gin.Context) {
	var input routeentities.AddIntersectionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	node, err := p.routeService.AddIntersection(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"node": node})
}

func (p *AdminAPI) DeleteIntersectionHandler(c *gin.Context) {
	buildingID, err := uuid.Parse(c.Param("building_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid building_id"})
		return
	}

	intersectionID, err := uuid.Parse(c.Param("intersection_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid intersection_id"})
		return
	}

	err = p.routeService.DeleteIntersection(c.Request.Context(), buildingID, intersectionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "intersection not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
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

	c.JSON(http.StatusCreated, gin.H{"connection": result})
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
	c.JSON(http.StatusOK, gin.H{"building": building})
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

	polygon, err := p.mapService.CreatePolygon(
		c.Request.Context(),
		mapentities.Polygon{FloorID: floorID, Label: req.Label, ZIndex: req.ZIndex},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"polygon": polygon})
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

	c.JSON(http.StatusCreated, gin.H{"points": result})
}

func (p *AdminAPI) DeletePolygonPointsHandler(c *gin.Context) {
	var req mapentities.DeletePolygonPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := p.mapService.DeletePolygonPoints(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (p *AdminAPI) SyncDatabaseHandler(c *gin.Context) {
	var dataRequest mapentities.SyncAllDataRequest
	if err := c.ShouldBindJSON(&dataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := dataRequest.Data

	ctx := c.Request.Context()
	// iterate over buildings
	for _, b := range data.Buildings {
		if _, err := p.mapService.CreateBuilding(ctx, mapentities.CreateBuildingInput{ID: b.ID, Name: b.Name, Address: b.Address}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, f := range b.Floors {
			if err := p.mapService.CreateFloor(ctx, b.ID, mapentities.Floor{ID: f.ID, Name: f.Name, Alias: f.Alias}); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// objects
			for _, obj := range f.Objects {
				input := mapentities.CreateObjectInput{
					ID:           obj.ID,
					Name:         obj.Name,
					Alias:        obj.Alias,
					Description:  obj.Description,
					X:            obj.X,
					Y:            obj.Y,
					Width:        obj.Width,
					Height:       obj.Height,
					ObjectTypeID: obj.ObjectTypeID,
				}
				_, err := p.mapService.CreateObject(ctx, f.ID, input)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
			for _, d := range f.Doors {
				if err := p.mapService.CreateDoor(ctx, d.ObjectID, d); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
			// polygons
			for _, poly := range f.FloorPolygons {
				if _, err := p.mapService.CreatePolygon(ctx, mapentities.Polygon{ID: poly.ID, FloorID: f.ID, Label: poly.Label, ZIndex: poly.ZIndex}); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				for _, pt := range poly.Points {
					if _, err := p.mapService.CreatePolygonPoint(ctx, poly.ID, pt.Order, pt.X, pt.Y); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				}
			}
			// intersections
			for _, in := range f.Intersections {
				if _, err := p.routeService.AddIntersection(ctx, routeentities.AddIntersectionRequest{ID: in.ID, X: in.X, Y: in.Y, FloorID: f.ID}); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			for _, conn := range f.Connections {
				if _, err := p.routeService.AddConnection(ctx, conn.FromID, conn.ToID, conn.Weight); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}
	}

	c.Status(http.StatusCreated)
}

func (p *AdminAPI) GetDatabaseHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var result mapentities.SyncAllData

	types, err := p.mapService.GetObjectCategories(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, t := range types {
		result.ObjectTypes = append(result.ObjectTypes, mapentities.ObjectTypeInfo{
			ID:    t.ID,
			Name:  t.Name,
			Alias: t.Alias,
		})
	}

	buildings, err := p.mapService.GetBuildings(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, b := range buildings {
		bSync := mapentities.SyncBuildings{ID: b.ID, Name: b.Name, Address: b.Address}

		objResp, err := p.mapService.GetObjectsResponse(ctx, b.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, fl := range objResp.Floors {
			fSync := mapentities.SyncFloors{
				ID:         fl.Floor.ID,
				Name:       fl.Floor.Name,
				Alias:      fl.Floor.Alias,
				BuildingID: b.ID,
				Objects:    fl.Objects,
			}
			// gather doors per floor
			for _, obj := range fl.Objects {
				fSync.Doors = append(fSync.Doors, obj.Doors...)
			}
			// polygons
			for _, bg := range fl.Background {
				var pts []mapentities.PolygonPoint
				for _, bp := range bg.Points {
					pts = append(pts, mapentities.PolygonPoint{
						PolygonID: bg.ID,
						Order:     int32(bp.Order),
						X:         bp.X,
						Y:         bp.Y,
					})
				}
				fSync.FloorPolygons = append(fSync.FloorPolygons, mapentities.Polygon{
					ID:      bg.ID,
					FloorID: fl.Floor.ID,
					Label:   bg.Label,
					ZIndex:  int32(bg.ZIndex),
					Points:  pts,
				})
			}

			intersections, err := p.routeService.GetIntersections(ctx, b.ID)
			if err == nil {
				for _, inter := range intersections {
					if inter.FloorID == fl.Floor.ID {
						fSync.Intersections = append(fSync.Intersections, routeentities.Intersection{
							ID:      inter.ID,
							X:       inter.X,
							Y:       inter.Y,
							FloorID: inter.FloorID,
						})
					}
				}
			}
			connections, err := p.routeService.GetConnections(ctx, b.ID)
			if err == nil {
				fSync.Connections = connections
			}

			bSync.Floors = append(bSync.Floors, fSync)
		}
		result.Buildings = append(result.Buildings, bSync)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (p *AdminAPI) GetPolygonByIDHandler(c *gin.Context) {
	idParam := c.Param("poligon_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poligon_id"})
		return
	}

	polygon, err := p.mapService.GetPolygonByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "polygon not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, polygon)
}
