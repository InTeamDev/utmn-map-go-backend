package entities

import "errors"

var (
	ErrInvalidIDFormat = errors.New("invalid ID format")
	ErrInvalidRequest  = errors.New("invalid request format")
	ErrInternalServer  = errors.New("internal server error")

	ErrObjectNotFound     = errors.New("object not found")
	ErrFloorNotFound      = errors.New("floor not found")
	ErrObjectTypeNotFound = errors.New("object type not found")
	ErrDoorNotFound       = errors.New("door not found")
	ErrInvalidInput       = errors.New("invalid input")

	ErrInvalidDimensions  = errors.New("invalid object dimensions")
	ErrInvalidCoordinates = errors.New("invalid coordinates")
	ErrPositionConflict   = errors.New("object position conflict")
)
