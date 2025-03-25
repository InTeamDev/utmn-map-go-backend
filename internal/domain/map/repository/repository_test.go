package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/mocks"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/sqlc"
)

type RepositoryTestSuite struct {
	suite.Suite
	db            *sql.DB
	mockDB        sqlmock.Sqlmock
	ctrl          *gomock.Controller
	mockConverter *mocks.MockMapConverter
	repo          *repository.Map
}

func (s *RepositoryTestSuite) SetupTest() {
	var err error
	s.db, s.mockDB, err = sqlmock.New()
	s.Require().NoError(err, "failed to create sqlmock")
	s.ctrl = gomock.NewController(s.T())
	s.mockConverter = mocks.NewMockMapConverter(s.ctrl)
	s.repo = repository.NewMap(s.db, s.mockConverter)
}

func (s *RepositoryTestSuite) TearDownTest() {
	s.Require().NoError(s.mockDB.ExpectationsWereMet())
	s.ctrl.Finish()
	s.db.Close()
}

//------------------------------------------------------------------------------
// GetObjects tests
//------------------------------------------------------------------------------

// TestGetObjects_Success проверяет happy path: объекты и двери найдены.
func (s *RepositoryTestSuite) TestGetObjects_Success() {
	testBuildingID := uuid.New()
	testFloorID := uuid.New()
	testObjectID := uuid.New()
	testDoorID := uuid.New()

	// Ожидаемый запрос для объектов.
	objQueryRegex := regexp.MustCompile(
		`(?s).*FROM objects o.*WHERE b\.id = \$1::uuid AND f\.id = \$2::uuid.*`)
	objRows := sqlmock.NewRows([]string{
		"id", "name", "alias", "description", "x", "y", "width", "height",
		"object_type", "floor_id", "floor_name", "building_id", "building_name",
	}).AddRow(
		testObjectID, "object name", "object alias",
		"object description", 1, 2, 3, 4, "object type", testFloorID,
		"floor name", testBuildingID, "building name",
	)
	s.mockDB.ExpectQuery(objQueryRegex.String()).
		WithArgs(testBuildingID, testFloorID).
		WillReturnRows(objRows)

	// Ожидаемый запрос для дверей.
	doorQueryRegex := regexp.MustCompile(
		`(?s).*FROM doors d.*JOIN object_doors od ON d\.id = od\.door_id.*` +
			`WHERE od\.object_id = ANY\(\$1::uuid\[\]\).*`)
	doorRows := sqlmock.NewRows([]string{"id", "x", "y", "width", "height", "object_id"}).
		AddRow(testDoorID, 0, 0, 0, 0, testObjectID)
	s.mockDB.ExpectQuery(doorQueryRegex.String()).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(doorRows)

	// Формируем ожидаемые структуры, получаемые из sqlmock.
	expectedObjRows := []sqlc.GetObjectsByBuildingAndFloorRow{
		{
			ID:           testObjectID,
			Name:         "object name",
			Alias:        "object alias",
			Description:  sql.NullString{String: "object description", Valid: true},
			X:            1,
			Y:            2,
			Width:        3,
			Height:       4,
			ObjectType:   "object type",
			FloorID:      testFloorID,
			FloorName:    "floor name",
			BuildingID:   testBuildingID,
			BuildingName: "building name",
		},
	}
	expectedDoorRows := []sqlc.GetDoorsByObjectIDsRow{
		{
			ID:       testDoorID,
			X:        0,
			Y:        0,
			Width:    0,
			Height:   0,
			ObjectID: testObjectID,
		},
	}
	expectedDoors := map[uuid.UUID][]entities.Door{
		testObjectID: {{ID: testDoorID}},
	}
	s.mockConverter.EXPECT().
		DoorsSqlcToEntity(gomock.Eq(expectedDoorRows)).
		Return(expectedDoors)
	expectedObjects := []entities.Object{
		{ID: testObjectID, Doors: expectedDoors[testObjectID]},
	}
	s.mockConverter.EXPECT().
		ObjectsSqlcToEntity(gomock.Eq(expectedObjRows), gomock.Eq(expectedDoors)).
		Return(expectedObjects)

	objects, err := s.repo.GetObjects(context.Background(), entities.GetObjectsRequest{
		BuildID: testBuildingID,
		FloorID: testFloorID,
	})
	s.Require().NoError(err)
	s.Require().Len(objects, 1)
	s.Equal(testObjectID, objects[0].ID)
}

// TestGetObjects_EmptyObjects проверяет случай, когда запрос объектов возвращает пустой результат.
func (s *RepositoryTestSuite) TestGetObjects_EmptyObjects() {
	testBuildingID := uuid.New()
	testFloorID := uuid.New()

	objQueryRegex := regexp.MustCompile(
		`(?s).*FROM objects o.*WHERE b\.id = \$1::uuid AND f\.id = \$2::uuid.*`)
	emptyObjRows := sqlmock.NewRows([]string{
		"id", "name", "alias", "description", "x", "y", "width", "height",
		"object_type", "floor_id", "floor_name", "building_id", "building_name",
	})
	s.mockDB.ExpectQuery(objQueryRegex.String()).
		WithArgs(testBuildingID, testFloorID).
		WillReturnRows(emptyObjRows)

	// Даже если объектов нет, запрос дверей выполняется с пустым срезом.
	doorQueryRegex := regexp.MustCompile(
		`(?s).*FROM doors d.*JOIN object_doors od ON d\.id = od\.door_id.*` +
			`WHERE od\.object_id = ANY\(\$1::uuid\[\]\).*`)
	emptyDoorRows := sqlmock.NewRows([]string{"id", "x", "y", "width", "height", "object_id"})
	s.mockDB.ExpectQuery(doorQueryRegex.String()).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(emptyDoorRows)

	s.mockConverter.EXPECT().
		DoorsSqlcToEntity(gomock.Eq([]sqlc.GetDoorsByObjectIDsRow(nil))).
		Return(map[uuid.UUID][]entities.Door{})
	s.mockConverter.EXPECT().
		ObjectsSqlcToEntity(gomock.Eq([]sqlc.GetObjectsByBuildingAndFloorRow(nil)),
			gomock.Eq(map[uuid.UUID][]entities.Door{})).
		Return([]entities.Object{})

	objects, err := s.repo.GetObjects(context.Background(), entities.GetObjectsRequest{
		BuildID: testBuildingID,
		FloorID: testFloorID,
	})
	s.Require().NoError(err)
	s.Empty(objects)
}

// TestGetObjects_EmptyDoors проверяет случай, когда объекты найдены, но дверей нет.
func (s *RepositoryTestSuite) TestGetObjects_EmptyDoors() {
	testBuildingID := uuid.New()
	testFloorID := uuid.New()
	testObjectID := uuid.New()

	objQueryRegex := regexp.MustCompile(
		`(?s).*FROM objects o.*WHERE b\.id = \$1::uuid AND f\.id = \$2::uuid.*`)
	objRows := sqlmock.NewRows([]string{
		"id", "name", "alias", "description", "x", "y", "width", "height",
		"object_type", "floor_id", "floor_name", "building_id", "building_name",
	}).AddRow(
		testObjectID, "object name", "object alias",
		"object description", 1, 2, 3, 4, "object type", testFloorID,
		"floor name", testBuildingID, "building name",
	)
	s.mockDB.ExpectQuery(objQueryRegex.String()).
		WithArgs(testBuildingID, testFloorID).
		WillReturnRows(objRows)

	doorQueryRegex := regexp.MustCompile(
		`(?s).*FROM doors d.*JOIN object_doors od ON d\.id = od\.door_id.*` +
			`WHERE od\.object_id = ANY\(\$1::uuid\[\]\).*`)
	emptyDoorRows := sqlmock.NewRows([]string{"id", "x", "y", "width", "height", "object_id"})
	s.mockDB.ExpectQuery(doorQueryRegex.String()).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(emptyDoorRows)

	s.mockConverter.EXPECT().
		DoorsSqlcToEntity(gomock.Eq([]sqlc.GetDoorsByObjectIDsRow(nil))).
		Return(map[uuid.UUID][]entities.Door{})

	expectedObjRows := []sqlc.GetObjectsByBuildingAndFloorRow{
		{
			ID:           testObjectID,
			Name:         "object name",
			Alias:        "object alias",
			Description:  sql.NullString{String: "object description", Valid: true},
			X:            1,
			Y:            2,
			Width:        3,
			Height:       4,
			ObjectType:   "object type",
			FloorID:      testFloorID,
			FloorName:    "floor name",
			BuildingID:   testBuildingID,
			BuildingName: "building name",
		},
	}
	s.mockConverter.EXPECT().
		ObjectsSqlcToEntity(gomock.Eq(expectedObjRows), gomock.Eq(map[uuid.UUID][]entities.Door{})).
		Return([]entities.Object{
			{ID: testObjectID, Doors: nil},
		})

	objects, err := s.repo.GetObjects(context.Background(), entities.GetObjectsRequest{
		BuildID: testBuildingID,
		FloorID: testFloorID,
	})
	s.Require().NoError(err)
	s.Require().Len(objects, 1)
	s.Equal(testObjectID, objects[0].ID)
	s.Empty(objects[0].Doors)
}

// TestGetObjects_ErrorOnGetObjects проверяет ошибку при запросе объектов.
func (s *RepositoryTestSuite) TestGetObjects_ErrorOnGetObjects() {
	testBuildingID := uuid.New()
	testFloorID := uuid.New()

	objQueryRegex := regexp.MustCompile(
		`(?s).*FROM objects o.*WHERE b\.id = \$1::uuid AND f\.id = \$2::uuid.*`)
	s.mockDB.ExpectQuery(objQueryRegex.String()).
		WithArgs(testBuildingID, testFloorID).
		WillReturnError(errors.New("db error"))

	_, err := s.repo.GetObjects(context.Background(), entities.GetObjectsRequest{
		BuildID: testBuildingID,
		FloorID: testFloorID,
	})
	s.Require().Error(err)
	s.Contains(err.Error(), "get objects")
}

//------------------------------------------------------------------------------
// GetBuildings tests
//------------------------------------------------------------------------------

// TestGetBuildings_Success проверяет успешное выполнение.
func (s *RepositoryTestSuite) TestGetBuildings_Success() {
	testBuildingID := uuid.New()
	buildQueryRegex := regexp.MustCompile(`(?s).*FROM buildings b.*`)
	buildRows := sqlmock.NewRows([]string{"id", "name", "address"}).
		AddRow(testBuildingID, "BuildingName", "Address")
	s.mockDB.ExpectQuery(buildQueryRegex.String()).WillReturnRows(buildRows)

	expectedBuildingsSqlc := []sqlc.Building{
		{
			ID:      testBuildingID,
			Name:    "BuildingName",
			Address: "Address",
		},
	}
	s.mockConverter.EXPECT().
		BuildingsSqlcToEntity(gomock.Eq(expectedBuildingsSqlc)).
		Return([]entities.Building{{ID: testBuildingID}})

	buildings, err := s.repo.GetBuildings(context.Background())
	s.Require().NoError(err)
	s.Require().Len(buildings, 1)
	s.Equal(testBuildingID, buildings[0].ID)
}

// TestGetBuildings_Error проверяет ошибку при запросе зданий.
func (s *RepositoryTestSuite) TestGetBuildings_Error() {
	buildQueryRegex := regexp.MustCompile(`(?s).*FROM buildings b.*`)
	s.mockDB.ExpectQuery(buildQueryRegex.String()).WillReturnError(errors.New("db error"))

	buildings, err := s.repo.GetBuildings(context.Background())
	s.Require().Error(err)
	s.Nil(buildings)
	s.Contains(err.Error(), "get buildings")
}

// TestGetBuildings_EmptyResult проверяет случай пустого результата.
func (s *RepositoryTestSuite) TestGetBuildings_EmptyResult() {
	buildQueryRegex := regexp.MustCompile(`(?s).*FROM buildings b.*`)
	emptyRows := sqlmock.NewRows([]string{"id", "name", "address"})
	s.mockDB.ExpectQuery(buildQueryRegex.String()).WillReturnRows(emptyRows)

	s.mockConverter.EXPECT().
		BuildingsSqlcToEntity(gomock.Eq([]sqlc.Building(nil))).
		Return([]entities.Building{})

	buildings, err := s.repo.GetBuildings(context.Background())
	s.Require().NoError(err)
	s.Empty(buildings)
}

//------------------------------------------------------------------------------
// GetFloors tests
//------------------------------------------------------------------------------

// TestGetFloors_Success проверяет успешное выполнение.
func (s *RepositoryTestSuite) TestGetFloors_Success() {
	testBuildingID := uuid.New()
	testFloorID := uuid.New()
	floorQueryRegex := regexp.MustCompile(
		`(?s).*FROM floors f.*WHERE f\.building_id = \$1::uuid.*`)
	floorRows := sqlmock.NewRows([]string{"id", "name", "alias", "building_id"}).
		AddRow(testFloorID, "FloorName", "Alias", testBuildingID)
	s.mockDB.ExpectQuery(floorQueryRegex.String()).
		WithArgs(testBuildingID).
		WillReturnRows(floorRows)

	expectedFloorsSqlc := []sqlc.Floor{
		{
			ID:         testFloorID,
			Name:       "FloorName",
			Alias:      "Alias",
			BuildingID: testBuildingID,
		},
	}
	s.mockConverter.EXPECT().
		FloorsSqlcToEntity(gomock.Eq(expectedFloorsSqlc)).
		Return([]entities.Floor{{ID: testFloorID}})

	floors, err := s.repo.GetFloors(context.Background(), testBuildingID)
	s.Require().NoError(err)
	s.Require().Len(floors, 1)
	s.Equal(testFloorID, floors[0].ID)
}

// TestGetFloors_Error проверяет ошибку при запросе этажей.
func (s *RepositoryTestSuite) TestGetFloors_Error() {
	testBuildingID := uuid.New()
	floorQueryRegex := regexp.MustCompile(
		`(?s).*FROM floors f.*WHERE f\.building_id = \$1::uuid.*`)
	s.mockDB.ExpectQuery(floorQueryRegex.String()).
		WithArgs(testBuildingID).
		WillReturnError(errors.New("db error"))

	floors, err := s.repo.GetFloors(context.Background(), testBuildingID)
	s.Require().Error(err)
	s.Nil(floors)
	s.Contains(err.Error(), "get floors")
}

// TestGetFloors_EmptyResult проверяет случай пустого результата.
func (s *RepositoryTestSuite) TestGetFloors_EmptyResult() {
	testBuildingID := uuid.New()
	floorQueryRegex := regexp.MustCompile(
		`(?s).*FROM floors f.*WHERE f\.building_id = \$1::uuid.*`)
	emptyRows := sqlmock.NewRows([]string{"id", "name", "alias", "building_id"})
	s.mockDB.ExpectQuery(floorQueryRegex.String()).
		WithArgs(testBuildingID).
		WillReturnRows(emptyRows)

	s.mockConverter.EXPECT().
		FloorsSqlcToEntity(gomock.Eq([]sqlc.Floor(nil))).
		Return([]entities.Floor{})

	floors, err := s.repo.GetFloors(context.Background(), testBuildingID)
	s.Require().NoError(err)
	s.Empty(floors)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
