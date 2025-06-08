package validate

import (
	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
)

func CreateDoorValidation(door mapentities.Door, object mapentities.Object) error {
	if door.Width < 0 || door.Height < 0 {
		return mapentities.ErrInvalidInput
	}
	if door.X < 0 || door.Y < 0 {
		return mapentities.ErrInvalidInput
	}
	switch {
	case (door.X < object.X-10 || object.X+object.Width+10 < door.X) && (object.Y-10 > door.Y || door.Y < object.Y+10):
		return mapentities.ErrInvalidCoordinates
	case (door.X < object.X-10 || object.X+object.Width+10 < door.X) && (object.Y+object.Height-10 > door.Y || door.Y > object.Y+object.Height+10):
		return mapentities.ErrInvalidCoordinates
	case (door.X < object.X-10 || door.X > object.X+10) && (door.Y < object.Y-10 || door.Y > object.Y+object.Height+10):
		return mapentities.ErrInvalidCoordinates
	case (door.X < object.X+object.Width-10 || door.X > object.X+object.Width+10) && (door.Y < object.Y-10 || door.Y > object.Y+object.Height+10):
		return mapentities.ErrInvalidCoordinates
	}

	return nil
}
