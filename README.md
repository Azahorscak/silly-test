# Hike Finder

A CLI tool that finds short hikes (4 miles or less) near any location. It geocodes your location, searches a curated dataset of 75+ trails across US national parks and popular hiking areas, and suggests a random hike.

## Build

```bash
go build -o hikefinder ./cmd/hikefinder
```

## Usage

```bash
# Find a random hike near a location
./hikefinder -location "Boulder, CO"

# Search within a specific radius (default: 50 miles)
./hikefinder -location "Moab, UT" -radius 30

# List all matching hikes
./hikefinder -location "Asheville, NC" -all
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-location` | (required) | Location to search near (e.g. "Denver, CO") |
| `-radius` | 50 | Search radius in miles |
| `-all` | false | List all matching hikes instead of picking one |

## Example Output

```
Geocoding "Boulder, CO"...
Found: Boulder, Boulder County, Colorado, USA (40.0150, -105.2705)

Hike Found!

  Name:       Chautauqua: Royal Arch Trail
  Distance:   3.4 miles
  Difficulty: Hard
  Location:   39.9990, -105.2830
  About:      Challenging hike to a natural rock arch with views of Boulder and the Flatirons.
```

## Testing

```bash
go test ./...
```

## How It Works

1. Geocodes the input location using the Nominatim (OpenStreetMap) API
2. Searches an embedded dataset of curated short hikes using the Haversine formula
3. Filters to hikes within the search radius that are 4 miles or less
4. Randomly selects one hike (or lists all with `-all`)

No third-party dependencies — uses only the Go standard library.
