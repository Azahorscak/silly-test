package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/user/hikefinder/internal/geo"
	"github.com/user/hikefinder/internal/models"
	"github.com/user/hikefinder/internal/picker"
	"github.com/user/hikefinder/internal/trails"
)

func main() {
	location := flag.String("location", "", "Location to search near (e.g. \"Denver, CO\")")
	radius := flag.Float64("radius", 50, "Search radius in miles")
	all := flag.Bool("all", false, "List all matching hikes instead of picking one")
	flag.Parse()

	if *location == "" {
		fmt.Fprintln(os.Stderr, "Error: -location is required")
		fmt.Fprintln(os.Stderr, "Usage: hikefinder -location \"Boulder, CO\" [-radius 50] [-all]")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("Geocoding %q...\n", *location)
	loc, err := geo.Geocode(ctx, *location)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found: %s (%.4f, %.4f)\n\n", loc.Name, loc.Lat, loc.Lng)

	hikes, err := trails.FindNearby(ctx, loc, *radius)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding hikes: %v\n", err)
		os.Exit(1)
	}

	if len(hikes) == 0 {
		fmt.Println("No hikes found within the search area. Try a larger radius or different location.")
		os.Exit(0)
	}

	if *all {
		printAllHikes(hikes)
	} else {
		h, err := picker.Pick(hikes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		printHike(h)
	}
}

func printHike(h models.Hike) {
	fmt.Println("Hike Found!")
	fmt.Println()
	fmt.Printf("  Name:       %s\n", h.Name)
	fmt.Printf("  Distance:   %.1f miles\n", h.Distance)
	fmt.Printf("  Difficulty: %s\n", h.Difficulty)
	fmt.Printf("  Location:   %.4f, %.4f\n", h.Lat, h.Lng)
	fmt.Printf("  About:      %s\n", h.Description)
}

func printAllHikes(hikes []models.Hike) {
	fmt.Printf("Found %d hikes:\n\n", len(hikes))
	fmt.Printf("  %-4s %-45s %-10s %s\n", "#", "Name", "Distance", "Difficulty")
	fmt.Printf("  %-4s %-45s %-10s %s\n", "---", "----", "--------", "----------")
	for i, h := range hikes {
		fmt.Printf("  %-4d %-45s %-10.1f %s\n", i+1, h.Name, h.Distance, h.Difficulty)
	}
}
