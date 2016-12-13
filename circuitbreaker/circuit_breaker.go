package circuitbreaker

import "time"

var utcFuture time.Time

func init() {
	utcFuture = time.Now().UTC().AddDate(1, 0, 0)
}

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
