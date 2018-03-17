package circuitbreaker

import (
	"fmt"

	"github.com/mantzas/gocloud/metrics"
	"github.com/pkg/errors"
)

// Status of the circuit breaker
type Status int

const (
	// Closed allow execution
	Closed Status = iota
	// HalfOpen allowing execution to check if resource works again
	HalfOpen
	// Open disallowing execution
	Open
)

// Action function to execute in circuit breaker
type Action func() (interface{}, error)

// CircuitBreaker implementation
type CircuitBreaker struct {
	sr     SettingsRetriever
	states map[string]*State
}

// NewCircuitBreaker constructor
func NewCircuitBreaker(sr SettingsRetriever, m metrics.Counter) *CircuitBreaker {

	states := make(map[string]*State, 0)

	for _, key := range sr.GetKeys() {
		states[key] = NewState(m, key)
	}

	return &CircuitBreaker{sr, states}
}

// Execute the function enclosed
func (cb *CircuitBreaker) Execute(key string, act Action) (interface{}, error) {

	sett, err := cb.sr.Get(key)
	if err != nil {
		return nil, fmt.Errorf("Failed to get setting of %s", key)
	}

	state, ok := cb.states[key]
	if !ok {
		return nil, fmt.Errorf("Failed to get state of %s", key)
	}

	status := state.GetStatus(sett)
	if status == Open {
		return nil, NewCircuitOpenError(key)
	}

	state.IncreaseExecutions()
	defer state.DecreaseExecutions()

	resp, err := act()
	if err != nil {
		state.IncreaseFailure()
		return nil, errors.Wrap(err, "Execution return error")
	}

	if status == HalfOpen {
		state.IncrementRetrySuccessCount()
	}

	return resp, nil
}
