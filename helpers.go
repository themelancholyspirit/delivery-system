package main

import (
	"database/sql"
	"math"
)

const EarthRadius = 6371 // Earth's radius in kilometers

// haversine calculates the great-circle distance between two points on Earth in meters
func haversine(lat1, lon1, lat2, lon2 float64) int {
	// Convert latitude and longitude from degrees to radians
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	distanceKm := EarthRadius * c

	// Convert to meters and return as an integer
	return int(math.Round(distanceKm * 1000))

}

func setupTestDatabase() (*sql.DB, error) {
	// Create a new in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Create the necessary tables for testing
	_, err = db.Exec(`CREATE TABLE orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		distance INTEGER,
		status TEXT
	);`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
