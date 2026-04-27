package service

import (
	"errors"
	"time"

	"github.com/robfig/cron/v3"
)

var ErrInvalidCron = errors.New("invalid cron expression")

// NextOccurrence computes the next time after `from` for the given cron expression.
// Pure function — no side effects, suitable for PBT.
func NextOccurrence(cronExpr string, from time.Time) (time.Time, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return time.Time{}, ErrInvalidCron
	}
	return schedule.Next(from), nil
}
