package circuitbreaker

import (
	"context"
	"github.com/sony/gobreaker"
	"log"
	"time"
)

type CircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

// NewCircuitBreaker creates a new circuit breaker instance
func NewCircuitBreaker(name string) *CircuitBreaker {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     60 * time.Second,

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},

		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit breaker state changed: %s -> %s", from, to)
		},
	})

	return &CircuitBreaker{cb: cb}
}

func (cb *CircuitBreaker) Execute(ctx context.Context, req func() (interface{}, error)) (interface{}, error) {
	return cb.cb.Execute(func() (interface{}, error) {
		return req()
	})
}

func (cb *CircuitBreaker) State() gobreaker.State {
	return cb.cb.State()
}
