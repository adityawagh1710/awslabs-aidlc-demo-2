package service_test

import (
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"github.com/todo-app/scheduler-service/internal/service"
)

// PBT: NextOccurrence must always return a time strictly after `from`.
func TestPBT_NextOccurrenceIsAfterFrom(t *testing.T) {
	properties := gopter.NewProperties(nil)

	validCrons := []string{
		"0 9 * * 1",     // every Monday 9am
		"*/5 * * * *",   // every 5 minutes
		"0 0 1 * *",     // first of month midnight
		"30 8 * * 1-5",  // weekdays 8:30am
	}

	properties.Property("next occurrence is always after from", prop.ForAll(
		func(cronIdx int, offsetHours int) bool {
			cron := validCrons[cronIdx%len(validCrons)]
			from := time.Now().Add(time.Duration(offsetHours) * time.Hour)
			next, err := service.NextOccurrence(cron, from)
			if err != nil {
				return false
			}
			return next.After(from)
		},
		gen.IntRange(0, 3),
		gen.IntRange(0, 8760), // up to 1 year offset
	))

	properties.TestingRun(t)
}

// PBT: NextOccurrence must be deterministic — same inputs yield same output.
func TestPBT_NextOccurrenceDeterministic(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cron := "0 9 * * 1"

	t1, err1 := service.NextOccurrence(cron, from)
	t2, err2 := service.NextOccurrence(cron, from)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, t1, t2)
}

// Invalid cron must return ErrInvalidCron.
func TestNextOccurrence_InvalidCron(t *testing.T) {
	_, err := service.NextOccurrence("not-a-cron", time.Now())
	assert.ErrorIs(t, err, service.ErrInvalidCron)
}
