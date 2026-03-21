package picker

import (
	"testing"

	"github.com/user/hikefinder/internal/models"
)

func TestPickEmpty(t *testing.T) {
	_, err := Pick(nil)
	if err == nil {
		t.Fatal("Pick(nil) should return error")
	}
}

func TestPickSingle(t *testing.T) {
	hikes := []models.Hike{{Name: "Test Hike", Distance: 1.0}}
	h, err := Pick(hikes)
	if err != nil {
		t.Fatalf("Pick() error: %v", err)
	}
	if h.Name != "Test Hike" {
		t.Errorf("Pick() returned %q, want %q", h.Name, "Test Hike")
	}
}

func TestPickMultiple(t *testing.T) {
	hikes := []models.Hike{
		{Name: "Hike A", Distance: 1.0},
		{Name: "Hike B", Distance: 2.0},
		{Name: "Hike C", Distance: 3.0},
	}

	// Run multiple times to verify it returns valid hikes
	for i := 0; i < 20; i++ {
		h, err := Pick(hikes)
		if err != nil {
			t.Fatalf("Pick() error: %v", err)
		}
		found := false
		for _, orig := range hikes {
			if h.Name == orig.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Pick() returned unknown hike %q", h.Name)
		}
	}
}
