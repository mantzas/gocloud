package circuitbreaker

import (
	"fmt"
)

// OpenError defines a error when circuit is open
type OpenError struct {
	key string
}

// Error returns the error text
func (e OpenError) Error() string {
	return fmt.Sprintf("Circuit is open for [%s]", e.key)
}

// NewCircuitOpenError constructor
func NewCircuitOpenError(key string) OpenError {
	return OpenError{key}
}
