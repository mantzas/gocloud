package circuitbreaker

import (
	"fmt"
)

// CircuitOpenError defines a error when circuit is open
type CircuitOpenError struct {
	key string
}

func (e CircuitOpenError) Error() string {
	return fmt.Sprintf("Circuit is open for [%s]", e.key)
}

// NewCircuitOpenError constructor
func NewCircuitOpenError(key string) CircuitOpenError {
	return CircuitOpenError{key}
}
