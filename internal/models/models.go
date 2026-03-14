package models

// Location represents a geocoded position.
type Location struct {
	Lat  float64
	Lng  float64
	Name string
}

// Hike represents a hiking trail.
type Hike struct {
	Name        string  `json:"name"`
	Distance    float64 `json:"distance_miles"`
	Lat         float64 `json:"latitude"`
	Lng         float64 `json:"longitude"`
	Description string  `json:"description"`
	Difficulty  string  `json:"difficulty"`
	URL         string  `json:"url,omitempty"`
}
