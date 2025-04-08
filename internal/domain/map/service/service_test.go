package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/mocks"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/service"
)

type MapServiceTestSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	mockRepo   *mocks.MockMapRepository
	mapService *service.Map
}

func (s *MapServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mocks.NewMockMapRepository(s.ctrl)
	s.mapService = service.NewMap(s.mockRepo)
}

func (s *MapServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

// TestGetBuildings_Success проверяет успешное получение зданий.
func (s *MapServiceTestSuite) TestGetBuildings_Success() {
	expectedBuildings := []entities.Building{
		{ID: uuid.New()},
	}
	s.mockRepo.EXPECT().
		GetBuildings(gomock.Any()).
		Return(expectedBuildings, nil)

	buildings, err := s.mapService.GetBuildings(context.Background())
	s.Require().NoError(err)
	s.Equal(expectedBuildings, buildings)
}

// TestGetBuildings_Error проверяет, что ошибка репозитория корректно оборачивается.
func (s *MapServiceTestSuite) TestGetBuildings_Error() {
	repoErr := errors.New("repository error")
	s.mockRepo.EXPECT().
		GetBuildings(gomock.Any()).
		Return(nil, repoErr)

	buildings, err := s.mapService.GetBuildings(context.Background())
	s.Require().Error(err)
	s.Nil(buildings)
	s.Contains(err.Error(), "get buildings")
	s.Contains(err.Error(), "repository error")
}

// TestGetFloors_Success проверяет успешное получение этажей.
func (s *MapServiceTestSuite) TestGetFloors_Success() {
	buildID := uuid.New()
	expectedFloors := []entities.Floor{
		{ID: uuid.New()},
	}
	s.mockRepo.EXPECT().
		GetFloors(gomock.Any(), buildID).
		Return(expectedFloors, nil)

	floors, err := s.mapService.GetFloors(context.Background(), buildID)
	s.Require().NoError(err)
	s.Equal(expectedFloors, floors)
}

// TestGetFloors_Error проверяет, что ошибка репозитория корректно оборачивается.
func (s *MapServiceTestSuite) TestGetFloors_Error() {
	buildID := uuid.New()
	repoErr := errors.New("repository error")
	s.mockRepo.EXPECT().
		GetFloors(gomock.Any(), buildID).
		Return(nil, repoErr)

	floors, err := s.mapService.GetFloors(context.Background(), buildID)
	s.Require().Error(err)
	s.Nil(floors)
	s.Contains(err.Error(), "get floors")
	s.Contains(err.Error(), "repository error")
}

func TestMapServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MapServiceTestSuite))
}
