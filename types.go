package main

import (
	"errors"
	"strconv"
)

type Coordinates struct {
	Origin      [2]string `json:"origin"`
	Destination [2]string `json:"destination"`
}

type Order struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Distance int    `json:"distance"`
	Status   string `json:"status"`
	Version  int    `json:"version"` // Optimistic locking field

}

func validateCoordinate(coord [2]string) error {
	lat, err := strconv.ParseFloat(coord[0], 64)
	if err != nil || lat < -90 || lat > 90 {
		return errors.New("invalid latitude value")
	}

	lon, err := strconv.ParseFloat(coord[1], 64)
	if err != nil || lon < -180 || lon > 180 {
		return errors.New("invalid longitude value")
	}

	return nil
}

// ValidateCoordinates checks if both origin and destination are valid coordinates
func ValidateCoordinates(c Coordinates) error {
	// Validate origin coordinates
	if err := validateCoordinate(c.Origin); err != nil {
		return errors.New("invalid origin coordinates: " + err.Error())
	}

	// Validate destination coordinates
	if err := validateCoordinate(c.Destination); err != nil {
		return errors.New("invalid destination coordinates: " + err.Error())
	}

	return nil
}
