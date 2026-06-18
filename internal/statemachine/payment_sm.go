package statemachine

import (
	"fmt"

	"github.com/vishalvignesh12/payflow/internal/models"
)

var validTransition = map[models.Status][]models.Status{
	models.StatusCreated:    {models.StatusPending},
	models.StatusPending:    {models.StatusProcessing},
	models.StatusProcessing: {models.StatusSuccess, models.StatusFailed},
	models.StatusFailed:     {models.StatusProcessing},
}

func IsValid(from, to models.Status) bool {
	allowed, ok := validTransition[from]

	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func ValidTransition(from, to models.Status) error {
	if !IsValid(from, to) {
		return fmt.Errorf("invalid transition: %s → %s", from, to)
	}
	return nil
}
