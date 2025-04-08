// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository (interfaces: MapConverter)
//
// Generated by this command:
//
//	mockgen -destination=mocks/mock_map_converter.go -package=mocks github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository MapConverter
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	entities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	sqlc "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/sqlc"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockMapConverter is a mock of MapConverter interface.
type MockMapConverter struct {
	ctrl     *gomock.Controller
	recorder *MockMapConverterMockRecorder
	isgomock struct{}
}

// MockMapConverterMockRecorder is the mock recorder for MockMapConverter.
type MockMapConverterMockRecorder struct {
	mock *MockMapConverter
}

// NewMockMapConverter creates a new mock instance.
func NewMockMapConverter(ctrl *gomock.Controller) *MockMapConverter {
	mock := &MockMapConverter{ctrl: ctrl}
	mock.recorder = &MockMapConverterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMapConverter) EXPECT() *MockMapConverterMockRecorder {
	return m.recorder
}

// BuildingsSqlcToEntity mocks base method.
func (m *MockMapConverter) BuildingsSqlcToEntity(buildings []sqlc.Building) []entities.Building {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildingsSqlcToEntity", buildings)
	ret0, _ := ret[0].([]entities.Building)
	return ret0
}

// BuildingsSqlcToEntity indicates an expected call of BuildingsSqlcToEntity.
func (mr *MockMapConverterMockRecorder) BuildingsSqlcToEntity(buildings any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildingsSqlcToEntity", reflect.TypeOf((*MockMapConverter)(nil).BuildingsSqlcToEntity), buildings)
}

// DoorsSqlcToEntityMap mocks base method.
func (m *MockMapConverter) DoorsSqlcToEntityMap(doors []sqlc.GetDoorsByObjectIDsRow) map[uuid.UUID][]entities.Door {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoorsSqlcToEntityMap", doors)
	ret0, _ := ret[0].(map[uuid.UUID][]entities.Door)
	return ret0
}

// DoorsSqlcToEntityMap indicates an expected call of DoorsSqlcToEntityMap.
func (mr *MockMapConverterMockRecorder) DoorsSqlcToEntityMap(doors any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoorsSqlcToEntityMap", reflect.TypeOf((*MockMapConverter)(nil).DoorsSqlcToEntityMap), doors)
}

// FloorsSqlcToEntity mocks base method.
func (m *MockMapConverter) FloorsSqlcToEntity(floors []sqlc.Floor) []entities.Floor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FloorsSqlcToEntity", floors)
	ret0, _ := ret[0].([]entities.Floor)
	return ret0
}

// FloorsSqlcToEntity indicates an expected call of FloorsSqlcToEntity.
func (mr *MockMapConverterMockRecorder) FloorsSqlcToEntity(floors any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FloorsSqlcToEntity", reflect.TypeOf((*MockMapConverter)(nil).FloorsSqlcToEntity), floors)
}

// ObjectSqlcToEntity mocks base method.
func (m *MockMapConverter) ObjectSqlcToEntity(object sqlc.GetObjectsByBuildingRow, doors []entities.Door) entities.Object {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ObjectSqlcToEntity", object, doors)
	ret0, _ := ret[0].(entities.Object)
	return ret0
}

// ObjectSqlcToEntity indicates an expected call of ObjectSqlcToEntity.
func (mr *MockMapConverterMockRecorder) ObjectSqlcToEntity(object, doors any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObjectSqlcToEntity", reflect.TypeOf((*MockMapConverter)(nil).ObjectSqlcToEntity), object, doors)
}

// ObjectTypesSqlcToEntity mocks base method.
func (m *MockMapConverter) ObjectTypesSqlcToEntity(objectTypes []sqlc.ObjectType) []entities.ObjectType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ObjectTypesSqlcToEntity", objectTypes)
	ret0, _ := ret[0].([]entities.ObjectType)
	return ret0
}

// ObjectTypesSqlcToEntity indicates an expected call of ObjectTypesSqlcToEntity.
func (mr *MockMapConverterMockRecorder) ObjectTypesSqlcToEntity(objectTypes any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObjectTypesSqlcToEntity", reflect.TypeOf((*MockMapConverter)(nil).ObjectTypesSqlcToEntity), objectTypes)
}

// ObjectsSqlcToEntityByBuilding mocks base method.
func (m *MockMapConverter) ObjectsSqlcToEntityByBuilding(objects []sqlc.GetObjectsByBuildingRow, doors map[uuid.UUID][]entities.Door) []entities.Object {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ObjectsSqlcToEntityByBuilding", objects, doors)
	ret0, _ := ret[0].([]entities.Object)
	return ret0
}

// ObjectsSqlcToEntityByBuilding indicates an expected call of ObjectsSqlcToEntityByBuilding.
func (mr *MockMapConverterMockRecorder) ObjectsSqlcToEntityByBuilding(objects, doors any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObjectsSqlcToEntityByBuilding", reflect.TypeOf((*MockMapConverter)(nil).ObjectsSqlcToEntityByBuilding), objects, doors)
}
