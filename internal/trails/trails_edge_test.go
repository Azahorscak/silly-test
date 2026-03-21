package trails

import (
	"context"
	"testing"

	"github.com/user/hikefinder/internal/models"
)

func TestLoadHikesDataIntegrity(t *testing.T) {
	hikes, err := loadHikes()
	if err != nil {
		t.Fatalf("loadHikes() error: %v", err)
	}

	if len(hikes) == 0 {
		t.Fatal("loadHikes() returned no hikes")
	}

	for i, h := range hikes {
		if h.Name == "" {
			t.Errorf("hike[%d]: name is empty", i)
		}
		if h.Distance <= 0 {
			t.Errorf("hike[%d] %q: distance %.1f should be > 0", i, h.Name, h.Distance)
		}
		if h.Distance > 4.0 {
			t.Errorf("hike[%d] %q: distance %.1f exceeds 4-mile max", i, h.Name, h.Distance)
		}
		if h.Lat < -90 || h.Lat > 90 {
			t.Errorf("hike[%d] %q: latitude %.4f out of range [-90, 90]", i, h.Name, h.Lat)
		}
		if h.Lng < -180 || h.Lng > 180 {
			t.Errorf("hike[%d] %q: longitude %.4f out of range [-180, 180]", i, h.Name, h.Lng)
		}
		if h.Difficulty == "" {
			t.Errorf("hike[%d] %q: difficulty is empty", i, h.Name)
		}
		if h.Description == "" {
			t.Errorf("hike[%d] %q: description is empty", i, h.Name)
		}
	}
}

func TestLoadHikesDifficultyValues(t *testing.T) {
	hikes, err := loadHikes()
	if err != nil {
		t.Fatalf("loadHikes() error: %v", err)
	}

	validDifficulties := map[string]bool{"Easy": true, "Moderate": true, "Hard": true}
	for _, h := range hikes {
		if !validDifficulties[h.Difficulty] {
			t.Errorf("hike %q has unexpected difficulty %q", h.Name, h.Difficulty)
		}
	}
}

func TestLoadHikesUniqueNames(t *testing.T) {
	hikes, err := loadHikes()
	if err != nil {
		t.Fatalf("loadHikes() error: %v", err)
	}

	seen := make(map[string]bool)
	for _, h := range hikes {
		if seen[h.Name] {
			t.Errorf("duplicate hike name: %q", h.Name)
		}
		seen[h.Name] = true
	}
}

func TestFindNearbyRadiusBoundary(t *testing.T) {
	ctx := context.Background()

	// Denver location
	loc := models.Location{Lat: 39.7392, Lng: -104.9903, Name: "Denver"}

	// Very small radius should return fewer results
	smallRadius, err := FindNearby(ctx, loc, 5)
	if err != nil {
		t.Fatalf("FindNearby() small radius error: %v", err)
	}

	// Larger radius should return at least as many results
	largeRadius, err := FindNearby(ctx, loc, 100)
	if err != nil {
		t.Fatalf("FindNearby() large radius error: %v", err)
	}

	if len(largeRadius) < len(smallRadius) {
		t.Errorf("larger radius (%d results) returned fewer hikes than smaller radius (%d results)",
			len(largeRadius), len(smallRadius))
	}
}

func TestFindNearbyZeroRadius(t *testing.T) {
	ctx := context.Background()
	loc := models.Location{Lat: 39.7392, Lng: -104.9903, Name: "Denver"}

	hikes, err := FindNearby(ctx, loc, 0)
	if err != nil {
		t.Fatalf("FindNearby() error: %v", err)
	}
	if len(hikes) != 0 {
		t.Errorf("FindNearby() with 0 radius returned %d hikes, want 0", len(hikes))
	}
}

func TestFindNearbyMultipleRegions(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name   string
		loc    models.Location
		radius float64
	}{
		{"Colorado", models.Location{Lat: 39.7392, Lng: -104.9903, Name: "Denver"}, 100},
		{"California", models.Location{Lat: 37.7490, Lng: -119.5885, Name: "Yosemite"}, 50},
		{"Utah", models.Location{Lat: 37.2982, Lng: -113.0263, Name: "Zion"}, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hikes, err := FindNearby(ctx, tt.loc, tt.radius)
			if err != nil {
				t.Fatalf("FindNearby() error: %v", err)
			}
			if len(hikes) == 0 {
				t.Errorf("expected hikes near %s, got none", tt.name)
			}
		})
	}
}

func TestHaversineDistanceSymmetric(t *testing.T) {
	lat1, lng1 := 39.7392, -104.9903
	lat2, lng2 := 40.0150, -105.2705

	d1 := HaversineDistance(lat1, lng1, lat2, lng2)
	d2 := HaversineDistance(lat2, lng2, lat1, lng1)

	if d1 != d2 {
		t.Errorf("HaversineDistance not symmetric: %f != %f", d1, d2)
	}
}

func TestHaversineDistanceNonNegative(t *testing.T) {
	points := [][2]float64{
		{0, 0}, {90, 0}, {-90, 0}, {0, 180}, {39.7392, -104.9903},
	}

	for _, p1 := range points {
		for _, p2 := range points {
			d := HaversineDistance(p1[0], p1[1], p2[0], p2[1])
			if d < 0 {
				t.Errorf("HaversineDistance(%.1f,%.1f -> %.1f,%.1f) = %f, want >= 0",
					p1[0], p1[1], p2[0], p2[1], d)
			}
		}
	}
}

func TestFindNearbyAllHikesWithin4Miles(t *testing.T) {
	ctx := context.Background()

	// Use a very large radius to get all hikes
	loc := models.Location{Lat: 39.0, Lng: -105.0, Name: "Center US"}
	hikes, err := FindNearby(ctx, loc, 10000)
	if err != nil {
		t.Fatalf("FindNearby() error: %v", err)
	}

	for _, h := range hikes {
		if h.Distance > 4.0 {
			t.Errorf("hike %q has distance %.1f miles, exceeds 4-mile max", h.Name, h.Distance)
		}
	}
}
