package service

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker"
)

// NewCircuitBreakerClient returns a resty client wrapped with a gobreaker circuit breaker.
func NewCircuitBreakerClient(name string) *resty.Client {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 1,
		Interval:    0,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
	})

	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(2 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return err != nil || r.StatusCode() >= 500
		})

	// Wrap Execute with circuit breaker
	client.OnBeforeRequest(func(_ *resty.Client, req *resty.Request) error {
		_, err := cb.Execute(func() (interface{}, error) { return nil, nil })
		return err
	})

	return client
}
