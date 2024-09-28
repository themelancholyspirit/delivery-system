package main

import (
	"math"
	"testing"
)

func TestHaversine(t *testing.T) {
	// Define test cases with city pairs in the same continent and expected distances

	tests := []struct {
		lat1, lon1, lat2, lon2 float64
		expected               int
		description            string
	}{
		// Europe
		{52.5200, 13.4050, 48.8566, 2.3522, 878000, "Berlin -> Paris"},    // Berlin -> Paris
		{51.5074, -0.1278, 52.2297, 21.0122, 1448000, "London -> Warsaw"}, // London -> Warsaw
		{41.9028, 12.4964, 40.4168, -3.7038, 1360000, "Rome -> Madrid"},   // Rome -> Madrid

		// North America
		{40.7128, -74.0060, 45.4215, -75.6972, 540000, "New York -> Ottawa"},         // New York -> Ottawa
		{34.0522, -118.2437, 36.1699, -115.1398, 367000, "Los Angeles -> Las Vegas"}, // Los Angeles -> Las Vegas
		{25.7617, -80.1918, 30.3322, -81.6557, 527000, "Miami -> Jacksonville"},      // Miami -> Jacksonville

		// Asia

		{39.9042, 116.4074, 31.2304, 121.4737, 1064000, "Beijing -> Shanghai"}, // Beijing -> Shanghai

	}

	// Iterate through the test cases
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Calculate the distance using the haversine formula
			result := haversine(tt.lat1, tt.lon1, tt.lat2, tt.lon2)

			// Tolerance ensures that the test still validates the logic of the haversine function while being lenient enough to allow for small deviations.

			tolerance := 5000

			// Allow a margin of error based on tolerance
			if math.Abs(float64(result)-float64(tt.expected)) > float64(tolerance) {
				t.Errorf("%s: Expected distance: %v, got: %v", tt.description, tt.expected, result)
			}
		})
	}
}
