package geo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGeocodeHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	origURL := nominatimURL
	setNominatimURL(server.URL)
	defer setNominatimURL(origURL)

	_, err := Geocode(context.Background(), "Denver, CO")
	if err == nil {
		t.Fatal("expected error for HTTP 500 response")
	}
}

func TestGeocodeInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()

	origURL := nominatimURL
	setNominatimURL(server.URL)
	defer setNominatimURL(origURL)

	_, err := Geocode(context.Background(), "Denver, CO")
	if err == nil {
		t.Fatal("expected error for invalid JSON response")
	}
}

func TestGeocodeInvalidLatitude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		results := []nominatimResult{
			{Lat: "not-a-number", Lon: "-104.9903", DisplayName: "Denver"},
		}
		json.NewEncoder(w).Encode(results)
	}))
	defer server.Close()

	origURL := nominatimURL
	setNominatimURL(server.URL)
	defer setNominatimURL(origURL)

	_, err := Geocode(context.Background(), "Denver, CO")
	if err == nil {
		t.Fatal("expected error for invalid latitude")
	}
}

func TestGeocodeInvalidLongitude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		results := []nominatimResult{
			{Lat: "39.7392", Lon: "not-a-number", DisplayName: "Denver"},
		}
		json.NewEncoder(w).Encode(results)
	}))
	defer server.Close()

	origURL := nominatimURL
	setNominatimURL(server.URL)
	defer setNominatimURL(origURL)

	_, err := Geocode(context.Background(), "Denver, CO")
	if err == nil {
		t.Fatal("expected error for invalid longitude")
	}
}

func TestGeocodeContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		results := []nominatimResult{
			{Lat: "39.7392", Lon: "-104.9903", DisplayName: "Denver"},
		}
		json.NewEncoder(w).Encode(results)
	}))
	defer server.Close()

	origURL := nominatimURL
	setNominatimURL(server.URL)
	defer setNominatimURL(origURL)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := Geocode(ctx, "Denver, CO")
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestGeocodeQueryParams(t *testing.T) {
	var capturedQuery string
	var capturedFormat string
	var capturedLimit string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query().Get("q")
		capturedFormat = r.URL.Query().Get("format")
		capturedLimit = r.URL.Query().Get("limit")

		results := []nominatimResult{
			{Lat: "39.7392", Lon: "-104.9903", DisplayName: "Denver, CO, USA"},
		}
		json.NewEncoder(w).Encode(results)
	}))
	defer server.Close()

	origURL := nominatimURL
	setNominatimURL(server.URL)
	defer setNominatimURL(origURL)

	_, err := Geocode(context.Background(), "Denver, CO")
	if err != nil {
		t.Fatalf("Geocode() error: %v", err)
	}

	if capturedQuery != "Denver, CO" {
		t.Errorf("q param = %q, want %q", capturedQuery, "Denver, CO")
	}
	if capturedFormat != "json" {
		t.Errorf("format param = %q, want %q", capturedFormat, "json")
	}
	if capturedLimit != "1" {
		t.Errorf("limit param = %q, want %q", capturedLimit, "1")
	}
}
