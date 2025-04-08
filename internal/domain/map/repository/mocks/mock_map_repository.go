// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/service (interfaces: MapRepository)
//
// Generated by this command:
//
//	mockgen -destination=mocks/mock_map_repository.go -package=mocks github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/service MapRepository
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockMapRepository is a mock of MapRepository interface.
type MockMapRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMapRepositoryMockRecorder
	isgomock struct{}
}

// MockMapRepositoryMockRecorder is the mock recorder for MockMapRepository.
type MockMapRepositoryMockRecorder struct {
	mock *MockMapRepository
}

// NewMockMapRepository creates a new mock instance.
func NewMockMapRepository(ctrl *gomock.Controller) *MockMapRepository {
	mock := &MockMapRepository{ctrl: ctrl}
	mock.recorder = &MockMapRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMapRepository) EXPECT() *MockMapRepositoryMockRecorder {
	return m.recorder
}

// GetBuildings mocks base method.
func (m *MockMapRepository) GetBuildings(ctx context.Context) ([]entities.Building, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBuildings", ctx)
	ret0, _ := ret[0].([]entities.Building)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBuildings indicates an expected call of GetBuildings.
func (mr *MockMapRepositoryMockRecorder) GetBuildings(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBuildings", reflect.TypeOf((*MockMapRepository)(nil).GetBuildings), ctx)
}

// GetFloors mocks base method.
func (m *MockMapRepository) GetFloors(ctx context.Context, buildID uuid.UUID) ([]entities.Floor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFloors", ctx, buildID)
	ret0, _ := ret[0].([]entities.Floor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFloors indicates an expected call of GetFloors.
func (mr *MockMapRepositoryMockRecorder) GetFloors(ctx, buildID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFloors", reflect.TypeOf((*MockMapRepository)(nil).GetFloors), ctx, buildID)
}

// GetObjectTypes mocks base method.
func (m *MockMapRepository) GetObjectTypes(ctx context.Context) ([]entities.ObjectType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectTypes", ctx)
	ret0, _ := ret[0].([]entities.ObjectType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectTypes indicates an expected call of GetObjectTypes.
func (mr *MockMapRepositoryMockRecorder) GetObjectTypes(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectTypes", reflect.TypeOf((*MockMapRepository)(nil).GetObjectTypes), ctx)
}

// GetObjectsByBuilding mocks base method.
func (m *MockMapRepository) GetObjectsByBuilding(ctx context.Context, buildID uuid.UUID) ([]entities.Object, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectsByBuilding", ctx, buildID)
	ret0, _ := ret[0].([]entities.Object)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectsByBuilding indicates an expected call of GetObjectsByBuilding.
func (mr *MockMapRepositoryMockRecorder) GetObjectsByBuilding(ctx, buildID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectsByBuilding", reflect.TypeOf((*MockMapRepository)(nil).GetObjectsByBuilding), ctx, buildID)
}

// UpdateObject mocks base method.
func (m *MockMapRepository) UpdateObject(ctx context.Context, input entities.UpdateObjectInput) (entities.Object, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateObject", ctx, input)
	ret0, _ := ret[0].(entities.Object)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateObject indicates an expected call of UpdateObject.
func (mr *MockMapRepositoryMockRecorder) UpdateObject(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateObject", reflect.TypeOf((*MockMapRepository)(nil).UpdateObject), ctx, input)
}
