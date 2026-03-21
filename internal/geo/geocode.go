package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/user/hikefinder/internal/models"
)

var nominatimURL = "https://nominatim.openstreetmap.org/search"

func setNominatimURL(url string) {
	nominatimURL = url
}

type nominatimResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

// Geocode converts a location string to coordinates using the Nominatim API.
func Geocode(ctx context.Context, location string) (models.Location, error) {
	u, err := url.Parse(nominatimURL)
	if err != nil {
		return models.Location{}, fmt.Errorf("parsing URL: %w", err)
	}

	q := u.Query()
	q.Set("q", location)
	q.Set("format", "json")
	q.Set("limit", "1")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return models.Location{}, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "hikefinder/1.0 (local hike finder CLI)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.Location{}, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Location{}, fmt.Errorf("geocoding returned status %d", resp.StatusCode)
	}

	var results []nominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return models.Location{}, fmt.Errorf("decoding geocode response: %w", err)
	}

	if len(results) == 0 {
		return models.Location{}, fmt.Errorf("no results found for location %q", location)
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return models.Location{}, fmt.Errorf("parsing latitude: %w", err)
	}

	lng, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return models.Location{}, fmt.Errorf("parsing longitude: %w", err)
	}

	return models.Location{
		Lat:  lat,
		Lng:  lng,
		Name: results[0].DisplayName,
	}, nil
}
