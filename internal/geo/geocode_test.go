package geo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeocode(t *testing.T) {
	// Mock Nominatim server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			t.Error("expected non-empty q parameter")
		}

		// Check User-Agent
		if r.Header.Get("User-Agent") == "" {
			t.Error("expected User-Agent header")
		}

		results := []nominatimResult{
			{Lat: "39.7392", Lon: "-104.9903", DisplayName: "Denver, CO, USA"},
		}
		json.NewEncoder(w).Encode(results)
	}))
	defer server.Close()

	// Override the URL for testing
	origURL := nominatimURL
	setNominatimURL(server.URL)
	defer setNominatimURL(origURL)

	loc, err := Geocode(context.Background(), "Denver, CO")
	if err != nil {
		t.Fatalf("Geocode() error: %v", err)
	}

	if loc.Lat < 39.7 || loc.Lat > 39.8 {
		t.Errorf("Lat = %f, want ~39.74", loc.Lat)
	}
	if loc.Lng > -104.9 || loc.Lng < -105.0 {
		t.Errorf("Lng = %f, want ~-104.99", loc.Lng)
	}
	if loc.Name == "" {
		t.Error("Name should not be empty")
	}
}

func TestGeocodeNoResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]nominatimResult{})
	}))
	defer server.Close()

	origURL := nominatimURL
	setNominatimURL(server.URL)
	defer setNominatimURL(origURL)

	_, err := Geocode(context.Background(), "xyznonexistent")
	if err == nil {
		t.Fatal("expected error for no results")
	}
}
