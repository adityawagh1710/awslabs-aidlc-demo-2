package service_test

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/todo-app/todo-service/internal/model"
	"github.com/todo-app/todo-service/internal/service"
)

// PBT: the only valid transitions are pending→in_progress and in_progress→done.
// Any other combination must be rejected.
func TestPBT_StatusTransitions(t *testing.T) {
	allStatuses := []model.TodoStatus{model.StatusPending, model.StatusInProgress, model.StatusDone}

	properties := gopter.NewProperties(nil)

	properties.Property("only valid transitions are allowed", prop.ForAll(
		func(fromIdx, toIdx int) bool {
			from := allStatuses[fromIdx%3]
			to := allStatuses[toIdx%3]

			// Same status is always allowed (no-op).
			if from == to {
				return service.ValidateTransition(from, to) == nil
			}

			allowed, ok := model.ValidTransitions[from]
			shouldBeValid := false
			if ok {
				for _, s := range allowed {
					if s == to {
						shouldBeValid = true
						break
					}
				}
			}

			err := service.ValidateTransition(from, to)
			isValid := err == nil

			return shouldBeValid == isValid
		},
		gen.IntRange(0, 2),
		gen.IntRange(0, 2),
	))

	properties.TestingRun(t)
}
