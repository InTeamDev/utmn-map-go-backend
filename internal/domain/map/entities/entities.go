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

// модель создания (структура данных и валидация полей)
type CreateObjectRequest struct {
	Name         string  `json:"name" binding:"required,max=255"`
	Alias        string  `json:"alias" binding:"required,max=255"`
	Description  string  `json:"description" binding:"max=255"`
	X            float64 `json:"x" binding:"required"`
	Y            float64 `json:"y" binding:"required"`
	Width        float64 `json:"width" binding:"required,gte=1"`
	Height       float64 `json:"height" binding:"required,gte=1"`
	ObjectTypeID int32   `json:"object_type_id" binding:"required"`
}
