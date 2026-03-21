package picker

import (
	"testing"

	"github.com/user/hikefinder/internal/models"
)

func TestPickEmptySlice(t *testing.T) {
	_, err := Pick([]models.Hike{})
	if err == nil {
		t.Fatal("Pick(empty slice) should return error")
	}
}

func TestPickPreservesAllFields(t *testing.T) {
	hikes := []models.Hike{
		{
			Name:        "Test Trail",
			Distance:    2.5,
			Lat:         39.7392,
			Lng:         -104.9903,
			Description: "A scenic test trail",
			Difficulty:  "Moderate",
			URL:         "https://example.com",
		},
	}

	h, err := Pick(hikes)
	if err != nil {
		t.Fatalf("Pick() error: %v", err)
	}

	if h.Name != "Test Trail" {
		t.Errorf("Name = %q, want %q", h.Name, "Test Trail")
	}
	if h.Distance != 2.5 {
		t.Errorf("Distance = %f, want 2.5", h.Distance)
	}
	if h.Lat != 39.7392 {
		t.Errorf("Lat = %f, want 39.7392", h.Lat)
	}
	if h.Lng != -104.9903 {
		t.Errorf("Lng = %f, want -104.9903", h.Lng)
	}
	if h.Description != "A scenic test trail" {
		t.Errorf("Description = %q, want %q", h.Description, "A scenic test trail")
	}
	if h.Difficulty != "Moderate" {
		t.Errorf("Difficulty = %q, want %q", h.Difficulty, "Moderate")
	}
	if h.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", h.URL, "https://example.com")
	}
}

func TestPickDistribution(t *testing.T) {
	hikes := []models.Hike{
		{Name: "A", Distance: 1.0},
		{Name: "B", Distance: 2.0},
		{Name: "C", Distance: 3.0},
	}

	counts := make(map[string]int)
	iterations := 300

	for i := 0; i < iterations; i++ {
		h, err := Pick(hikes)
		if err != nil {
			t.Fatalf("Pick() error: %v", err)
		}
		counts[h.Name]++
	}

	// Each hike should be picked at least once in 300 iterations
	for _, name := range []string{"A", "B", "C"} {
		if counts[name] == 0 {
			t.Errorf("hike %q was never picked in %d iterations", name, iterations)
		}
	}
}
