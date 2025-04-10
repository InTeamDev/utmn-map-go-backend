package entities

import "github.com/google/uuid"

type GetObjectsRequest struct {
	BuildID uuid.UUID
	FloorID uuid.UUID
}

// Тип для возвращаемой информации по объектам здания, сгруппированным по этажам.
type GetObjectsResponse struct {
	Building Building        `json:"building"` // Общая информация о здании.
	Floors   []FloorWithData `json:"floors"`   // Данные по каждому этажу.
}

// Данные по этажу, включая информацию об этаже, объекты и фон.
type FloorWithData struct {
	Floor      Floor                    `json:"floor"`
	Objects    []Object                 `json:"objects"`
	Background []FloorBackgroundElement `json:"background"`
}

type UpdateObjectInput struct {
	ID          uuid.UUID
	Name        string
	Alias       string
	Description string
	ObjectType  ObjectType
}
