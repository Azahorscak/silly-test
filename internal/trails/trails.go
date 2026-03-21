package trails

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/user/hikefinder/internal/models"
)

//go:embed data/hikes.json
var hikesFS embed.FS

const maxHikeDistance = 4.0 // miles

// loadHikes reads the embedded hike dataset.
func loadHikes() ([]models.Hike, error) {
	data, err := hikesFS.ReadFile("data/hikes.json")
	if err != nil {
		return nil, fmt.Errorf("reading embedded hike data: %w", err)
	}

	var hikes []models.Hike
	if err := json.Unmarshal(data, &hikes); err != nil {
		return nil, fmt.Errorf("parsing hike data: %w", err)
	}
	return hikes, nil
}

// FindNearby returns hikes within radiusMiles of loc that are ≤4 miles long,
// sorted by distance from the user's location.
func FindNearby(_ context.Context, loc models.Location, radiusMiles float64) ([]models.Hike, error) {
	allHikes, err := loadHikes()
	if err != nil {
		return nil, err
	}

	type hikeWithDist struct {
		hike models.Hike
		dist float64
	}

	var nearby []hikeWithDist
	for _, h := range allHikes {
		if h.Distance > maxHikeDistance {
			continue
		}
		dist := HaversineDistance(loc.Lat, loc.Lng, h.Lat, h.Lng)
		if dist <= radiusMiles {
			nearby = append(nearby, hikeWithDist{hike: h, dist: dist})
		}
	}

	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].dist < nearby[j].dist
	})

	result := make([]models.Hike, len(nearby))
	for i, n := range nearby {
		result[i] = n.hike
	}
	return result, nil
}

// HaversineDistance returns the distance in miles between two lat/lng points.
func HaversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusMiles = 3959.0

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMiles * c
}
