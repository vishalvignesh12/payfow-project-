package statemachine

import (
	"testing"

	"github.com/vishalvignesh12/payflow/internal/models"
)

func TestValidTransitions(t *testing.T) {
	// These should all be valid
	valid := [][2]models.Status{
		{models.StatusCreated, models.StatusPending},
		{models.StatusPending, models.StatusProcessing},
		{models.StatusProcessing, models.StatusSuccess},
		{models.StatusProcessing, models.StatusFailed},
		{models.StatusFailed, models.StatusProcessing},
	}
	for _, pair := range valid {
		if !IsValid(pair[0], pair[1]) {
			t.Errorf("expected valid: %s → %s", pair[0], pair[1])
		}
	}

	// These should all be invalid
	invalid := [][2]models.Status{
		{models.StatusSuccess, models.StatusFailed},
		{models.StatusSuccess, models.StatusPending},
		{models.StatusCreated, models.StatusSuccess},
		{models.StatusPending, models.StatusFailed},
	}
	for _, pair := range invalid {
		if IsValid(pair[0], pair[1]) {
			t.Errorf("expected invalid: %s → %s", pair[0], pair[1])
		}
	}
}
