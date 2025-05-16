package entities

import "errors"

var (
	ErrInvalidIDFormat = errors.New("invalid ID format")
	ErrInvalidRequest  = errors.New("invalid request format")
	ErrInternalServer  = errors.New("internal server error")

	ErrFloorNotFound      = errors.New("floor not found")
	ErrObjectTypeNotFound = errors.New("object type not found")
	ErrInvalidDimensions  = errors.New("invalid object dimensions")
	ErrPositionConflict   = errors.New("object position conflict")
)
