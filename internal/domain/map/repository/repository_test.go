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

type RepositoryNewTestSuite struct {
	suite.Suite
	db            *sql.DB
	mockDB        sqlmock.Sqlmock
	ctrl          *gomock.Controller
	mockConverter *mocks.MockMapConverter
	repo          *repository.Map
}

func (s *RepositoryNewTestSuite) SetupTest() {
	var err error
	s.db, s.mockDB, err = sqlmock.New()
	s.Require().NoError(err, "failed to create sqlmock")
	s.ctrl = gomock.NewController(s.T())
	s.mockConverter = mocks.NewMockMapConverter(s.ctrl)
	s.repo = repository.NewMap(s.db, s.mockConverter)
}

func (s *RepositoryNewTestSuite) TearDownTest() {
	s.Require().NoError(s.mockDB.ExpectationsWereMet())
	s.ctrl.Finish()
	s.db.Close()
}

//------------------------------------------------------------------------------
// GetBuildings tests
//------------------------------------------------------------------------------

func (s *RepositoryNewTestSuite) TestGetBuildings_Success() {
	testBuildingID := uuid.New()

	buildQueryRegex := regexp.MustCompile(`(?s).*FROM buildings b.*`)
	rows := sqlmock.NewRows([]string{"id", "name", "address"}).
		AddRow(testBuildingID, "BuildingName", "Address")
	s.mockDB.ExpectQuery(buildQueryRegex.String()).WillReturnRows(rows)

	expectedSQLCBuildings := []sqlc.Building{
		{ID: testBuildingID, Name: "BuildingName", Address: "Address"},
	}
	expectedEntitiesBuildings := []entities.Building{
		{ID: testBuildingID},
	}

	s.mockConverter.EXPECT().
		BuildingsSqlcToEntity(gomock.Eq(expectedSQLCBuildings)).
		Return(expectedEntitiesBuildings)

	buildings, err := s.repo.GetBuildings(context.Background())
	s.Require().NoError(err)
	s.Require().Len(buildings, 1)
	s.Equal(testBuildingID, buildings[0].ID)
}

func (s *RepositoryNewTestSuite) TestGetBuildings_Error() {
	buildQueryRegex := regexp.MustCompile(`(?s).*FROM buildings b.*`)
	s.mockDB.ExpectQuery(buildQueryRegex.String()).WillReturnError(errors.New("db error"))

	buildings, err := s.repo.GetBuildings(context.Background())
	s.Require().Error(err)
	s.Nil(buildings)
	s.Contains(err.Error(), "get buildings")
}

//------------------------------------------------------------------------------
// GetFloors tests
//------------------------------------------------------------------------------

func (s *RepositoryNewTestSuite) TestGetFloors_Success() {
	testBuildingID := uuid.New()
	testFloorID := uuid.New()

	floorQueryRegex := regexp.MustCompile(`(?s).*FROM floors f.*WHERE f\.building_id = \$1.*`)
	rows := sqlmock.NewRows([]string{"id", "name", "alias", "building_id"}).
		AddRow(testFloorID, "FloorName", "Alias", testBuildingID)
	s.mockDB.ExpectQuery(floorQueryRegex.String()).
		WithArgs(testBuildingID).
		WillReturnRows(rows)

	expectedSQLCFloors := []sqlc.Floor{
		{ID: testFloorID, Name: "FloorName", Alias: "Alias", BuildingID: testBuildingID},
	}
	expectedEntitiesFloors := []entities.Floor{
		{ID: testFloorID},
	}

	s.mockConverter.EXPECT().
		FloorsSqlcToEntity(gomock.Eq(expectedSQLCFloors)).
		Return(expectedEntitiesFloors)

	floors, err := s.repo.GetFloors(context.Background(), testBuildingID)
	s.Require().NoError(err)
	s.Require().Len(floors, 1)
	s.Equal(testFloorID, floors[0].ID)
}

func (s *RepositoryNewTestSuite) TestGetFloors_Error() {
	testBuildingID := uuid.New()
	floorQueryRegex := regexp.MustCompile(`(?s).*FROM floors f.*WHERE f\.building_id = \$1.*`)
	s.mockDB.ExpectQuery(floorQueryRegex.String()).
		WithArgs(testBuildingID).
		WillReturnError(errors.New("db error"))

	floors, err := s.repo.GetFloors(context.Background(), testBuildingID)
	s.Require().Error(err)
	s.Nil(floors)
	s.Contains(err.Error(), "get floors")
}

func TestRepositoryNewTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryNewTestSuite))
}
