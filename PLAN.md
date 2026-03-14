# Hike Finder - Go Application Plan

## Overview
A CLI application written in Go that finds hikes of 4 miles or less near a user-provided location. The app geocodes the location, queries a trail/hiking API, filters results to ≤4 miles, and presents a randomly selected hike.

---

## Architecture

```
cmd/hikefinder/main.go       # Entry point, CLI flag parsing
internal/geo/geocode.go       # Location string → lat/lng conversion
internal/trails/trails.go     # Trail API client + filtering
internal/picker/picker.go     # Random hike selection logic
internal/models/models.go     # Shared data types (Hike, Location)
go.mod / go.sum               # Module definition
```

---

## Step-by-Step Implementation Plan

### Step 1: Project Initialization
- Run `go mod init github.com/user/hikefinder`
- Create the directory structure: `cmd/hikefinder/`, `internal/geo/`, `internal/trails/`, `internal/picker/`, `internal/models/`

### Step 2: Define Data Models (`internal/models/models.go`)
- **`Location`** struct: `Lat float64`, `Lng float64`, `Name string`
- **`Hike`** struct: `Name string`, `Distance float64` (miles), `Lat float64`, `Lng float64`, `Description string`, `Difficulty string`, `URL string`

### Step 3: Geocoding (`internal/geo/geocode.go`)
- Use the **Nominatim (OpenStreetMap)** geocoding API — free, no API key required
- Endpoint: `https://nominatim.openstreetmap.org/search?q=<location>&format=json&limit=1`
- Function: `Geocode(ctx context.Context, location string) (models.Location, error)`
- Set a proper `User-Agent` header (Nominatim requires this)
- Parse JSON response, extract `lat` and `lon`, return as `models.Location`
- Return a clear error if no results are found

### Step 4: Trail Data Source (`internal/trails/trails.go`)
- **Primary approach: Embedded curated dataset**
  - Embed a JSON file (`internal/trails/data/hikes.json`) containing ~50-100 popular short hikes across major US regions using Go's `embed` package
  - Each entry has: name, lat, lng, distance (miles), difficulty, description
  - This avoids API key requirements and rate limits
- **Query function**: `FindNearby(ctx context.Context, loc models.Location, radiusMiles float64) ([]models.Hike, error)`
  - Calculate distance from the user's location to each hike using the **Haversine formula**
  - Filter to hikes within the given radius AND ≤ 4 miles in length
  - Sort results by distance from user

### Step 5: Haversine Distance Calculation (in `internal/trails/trails.go` or `internal/geo/distance.go`)
- Implement `HaversineDistance(lat1, lng1, lat2, lng2 float64) float64` returning miles
- Standard formula using Earth's radius (3959 miles)

### Step 6: Hike Picker (`internal/picker/picker.go`)
- Function: `Pick(hikes []models.Hike) (models.Hike, error)`
- If no hikes available, return an error with a helpful message
- If hikes are available, select one at random using `math/rand`
- Optionally weight selection toward closer hikes

### Step 7: CLI Entry Point (`cmd/hikefinder/main.go`)
- Parse CLI flags:
  - `-location` (required): free-text location string, e.g. `"Denver, CO"`
  - `-radius` (optional, default `50`): search radius in miles
  - `-all` (optional): list all matching hikes instead of picking one
- Flow:
  1. Validate input — ensure location is provided
  2. Call `geo.Geocode()` to convert location to coordinates
  3. Call `trails.FindNearby()` with coordinates and radius
  4. Call `picker.Pick()` to select a hike (or list all if `-all` flag set)
  5. Print the result in a formatted, readable way

### Step 8: Output Formatting
- Display selected hike with:
  ```
  🥾 Hike Found!

  Name:       Bear Lake Loop
  Distance:   3.2 miles
  Difficulty: Easy
  Location:   40.3128° N, 105.6458° W
  About:      A scenic loop around Bear Lake in RMNP...
  ```
- If `-all` flag, display a numbered table of all matching hikes

### Step 9: Curated Hike Dataset (`internal/trails/data/hikes.json`)
- Create a JSON file with ~50-100 short hikes (≤4 miles) across popular hiking regions:
  - Colorado (RMNP, Garden of the Gods, etc.)
  - California (Yosemite, Joshua Tree, etc.)
  - Utah (Zion, Arches, Bryce Canyon, etc.)
  - Pacific Northwest (Olympic, Rainier, etc.)
  - Northeast (Acadia, White Mountains, etc.)
  - Southeast (Great Smoky Mountains, Shenandoah, etc.)
- Each entry includes: name, latitude, longitude, distance_miles, difficulty, description

### Step 10: Tests
- `internal/geo/geocode_test.go` — test geocoding with mock HTTP server
- `internal/trails/trails_test.go` — test Haversine formula, test filtering logic
- `internal/picker/picker_test.go` — test random selection, test empty input error
- `cmd/hikefinder/main_test.go` — integration test of CLI flags

### Step 11: Documentation
- Add a `README.md` with:
  - What the app does
  - How to build: `go build -o hikefinder ./cmd/hikefinder`
  - How to run: `./hikefinder -location "Boulder, CO"`
  - Available flags
  - Example output

---

## Dependencies
- **Go standard library only** — no third-party dependencies needed
  - `net/http` for geocoding API calls
  - `encoding/json` for JSON parsing
  - `embed` for embedding the hike dataset
  - `flag` for CLI argument parsing
  - `math` for Haversine calculation
  - `math/rand` for random selection
  - `fmt`, `os`, `context`, `strings`

## Key Design Decisions
1. **Embedded dataset over external API** — avoids API key setup, works offline after geocoding, always available
2. **Nominatim for geocoding** — free, no signup, reliable for location-to-coordinates
3. **No third-party deps** — keeps the project simple and easy to build
4. **4-mile filter is hard-coded but configurable** — the dataset only contains ≤4 mile hikes, and the filter enforces it programmatically as well
5. **Random selection** — keeps it fun and surprising; `-all` flag available for full control

## Build & Run
```bash
go build -o hikefinder ./cmd/hikefinder
./hikefinder -location "Portland, OR"
./hikefinder -location "Moab, UT" -radius 30
./hikefinder -location "Asheville, NC" -all
```
