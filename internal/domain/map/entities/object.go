package entities

import "github.com/google/uuid"

type ObjectType string

const (
	// Аудитория.
	ObjectTypeCabinet ObjectType = "cabinet"
	// Кафедра.
	ObjectTypeDepartment ObjectType = "department"
	// Мужской Туалет.
	ObjectTypeManToilet ObjectType = "man-toilet"
	// Женский Туалет.
	ObjectTypeWomanToilet ObjectType = "woman-toilet"
	// Лестница.
	ObjectTypeStair ObjectType = "stair"
	// Гардероб.
	ObjectTypeWardrobe ObjectType = "wardrobe"
	// Gym.
	ObjectTypeGym ObjectType = "gym"
	// Кафе.
	ObjectTypeCafe ObjectType = "cafe"
	// Столовая.
	ObjectTypeCanteen ObjectType = "canteen"
	// Зона отдыха.
	ObjectTypeChillZone ObjectType = "chill-zone"
)

type Object struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Alias        string    `json:"alias"`
	Description  string    `json:"description"`
	X            float64   `json:"x"`
	Y            float64   `json:"y"`
	Width        float64   `json:"width"`
	Height       float64   `json:"height"`
	ObjectTypeID int32     `json:"object_type_id"`
	Doors        []Door    `json:"doors"`
	Floor        Floor     `json:"floor"`
}

type ObjectTypeInfo struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}
