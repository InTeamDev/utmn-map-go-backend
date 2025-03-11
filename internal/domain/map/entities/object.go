package entities

import "github.com/google/uuid"

type ObjectType string

const (
	// Аудитория
	ObjectTypeCabinet ObjectType = "cabinet"
	// Кафедра
	ObjectTypeDepartment ObjectType = "department"
	// Мужской Туалет
	ObjectTypeManToilet ObjectType = "man-toilet"
	// Женский Туалет
	ObjectTypeWomanToilet ObjectType = "woman-toilet"
	// Лестница
	ObjectTypeStair ObjectType = "stair"
	// Гардероб
	ObjectTypeWardrobe ObjectType = "wardrobe"
	// Gym
	ObjectTypeGym ObjectType = "gym"
)

type Object struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Alias       string     `json:"alias"`
	Description string     `json:"description"`
	X           float64    `json:"x"`
	Y           float64    `json:"y"`
	Width       float64    `json:"width"`
	Height      float64    `json:"height"`
	ObjectType  ObjectType `json:"object_type"`
	Doors       []Door     `json:"doors"`
}
