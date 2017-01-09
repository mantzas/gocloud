package cloud

import "time"

// Setting definition
type Setting struct {
	// The key for this setting
	Key string
	// The threshold for the circuit to open
	FailureThreshold int
	// The timeout after which we set the state to half-open and allow a retry
	RetryTimeout time.Duration
	// The threshold of the retry successes which returns the state to open
	RetrySuccessThreshold int
	// The threshold of how many retry executions are allowed when the status is half-open
	MaxRetryExecutionThreshold int
}
