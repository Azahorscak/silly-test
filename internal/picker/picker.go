package picker

import (
	"fmt"
	"math/rand"

	"github.com/user/hikefinder/internal/models"
)

// Pick selects a random hike from the given slice.
func Pick(hikes []models.Hike) (models.Hike, error) {
	if len(hikes) == 0 {
		return models.Hike{}, fmt.Errorf("no hikes available — try a larger radius or different location")
	}
	return hikes[rand.Intn(len(hikes))], nil
}
