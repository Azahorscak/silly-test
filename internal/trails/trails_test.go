package trails

import (
	"context"
	"math"
	"testing"

	"github.com/user/hikefinder/internal/models"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lng1     float64
		lat2     float64
		lng2     float64
		wantDist float64
		tolerance float64
	}{
		{
			name:     "same point",
			lat1:     40.0, lng1: -105.0,
			lat2:     40.0, lng2: -105.0,
			wantDist: 0, tolerance: 0.001,
		},
		{
			name:     "Denver to Boulder (~25 miles)",
			lat1:     39.7392, lng1: -104.9903,
			lat2:     40.0150, lng2: -105.2705,
			wantDist: 25.0, tolerance: 3.0,
		},
		{
			name:     "New York to Los Angeles (~2450 miles)",
			lat1:     40.7128, lng1: -74.0060,
			lat2:     34.0522, lng2: -118.2437,
			wantDist: 2450.0, tolerance: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HaversineDistance(tt.lat1, tt.lng1, tt.lat2, tt.lng2)
			if math.Abs(got-tt.wantDist) > tt.tolerance {
				t.Errorf("HaversineDistance() = %f, want ~%f (tolerance %f)", got, tt.wantDist, tt.tolerance)
			}
		})
	}
}

func TestFindNearby(t *testing.T) {
	ctx := context.Background()

	// Search near Yosemite - should find Yosemite hikes
	loc := models.Location{Lat: 37.7490, Lng: -119.5885, Name: "Yosemite"}
	hikes, err := FindNearby(ctx, loc, 20)
	if err != nil {
		t.Fatalf("FindNearby() error: %v", err)
	}
	if len(hikes) == 0 {
		t.Fatal("FindNearby() returned no hikes near Yosemite")
	}
	for _, h := range hikes {
		if h.Distance > 4.0 {
			t.Errorf("hike %q has distance %.1f miles, want <= 4.0", h.Name, h.Distance)
		}
	}
}

func TestFindNearbyNoResults(t *testing.T) {
	ctx := context.Background()

	// Middle of the ocean
	loc := models.Location{Lat: 0, Lng: 0, Name: "Nowhere"}
	hikes, err := FindNearby(ctx, loc, 10)
	if err != nil {
		t.Fatalf("FindNearby() error: %v", err)
	}
	if len(hikes) != 0 {
		t.Errorf("FindNearby() returned %d hikes, want 0", len(hikes))
	}
}

func TestFindNearbySorted(t *testing.T) {
	ctx := context.Background()

	loc := models.Location{Lat: 37.7490, Lng: -119.5885, Name: "Yosemite"}
	hikes, err := FindNearby(ctx, loc, 50)
	if err != nil {
		t.Fatalf("FindNearby() error: %v", err)
	}
	if len(hikes) < 2 {
		t.Skip("not enough hikes to test sorting")
	}

	// Verify sorted by distance from user
	for i := 1; i < len(hikes); i++ {
		d1 := HaversineDistance(loc.Lat, loc.Lng, hikes[i-1].Lat, hikes[i-1].Lng)
		d2 := HaversineDistance(loc.Lat, loc.Lng, hikes[i].Lat, hikes[i].Lng)
		if d1 > d2 {
			t.Errorf("hikes not sorted: %q (%.1f mi) before %q (%.1f mi)", hikes[i-1].Name, d1, hikes[i].Name, d2)
		}
	}
}
