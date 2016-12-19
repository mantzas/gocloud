package circuitbreaker

import (
	"sync"
	"time"
)

var utcFuture time.Time

func init() {
	utcFuture = time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC)
}

// State definition
type State struct {
	currentFailureCount  int
	retrySuccessCount    int
	currentExecutions    int
	lastFailureTimestamp time.Time
	m                    *sync.Mutex
}

// Reset the state
func (s *State) Reset() {
	s.m.Lock()
	defer s.m.Unlock()

	s.innerReset()
}

func (s *State) innerReset() {
	s.currentFailureCount = 0
	s.retrySuccessCount = 0
	s.lastFailureTimestamp = utcFuture
}

// IncreaseFailure increases the failure count
func (s *State) IncreaseFailure() {
	s.m.Lock()
	defer s.m.Unlock()

	s.currentFailureCount++
	s.lastFailureTimestamp = time.Now().UTC()
}

// IncrementRetrySuccessCount increments the retry success count
func (s *State) IncrementRetrySuccessCount() {
	s.m.Lock()
	defer s.m.Unlock()

	s.retrySuccessCount++
}

// IncreaseExecutions increases the current execution count
func (s *State) IncreaseExecutions() {
	s.m.Lock()
	defer s.m.Unlock()

	s.currentExecutions++
}

// DecreaseExecutions decreases the current execution count
func (s *State) DecreaseExecutions() {
	s.m.Lock()
	defer s.m.Unlock()

	s.currentExecutions--
}

// GetStatus returns the status of the circuit
func (s *State) GetStatus(sett *Setting) Status {
	s.m.Lock()
	defer s.m.Unlock()

	if sett.FailureThreshold > s.currentFailureCount {
		return Closed
	}

	retry := s.lastFailureTimestamp.Add(sett.RetryTimeout)
	now := time.Now().UTC()

	if retry.Before(now) || retry.Equal(now) {

		if s.retrySuccessCount >= sett.RetrySuccessThreshold {
			s.innerReset()
			return Closed
		}

		if s.currentExecutions > sett.MaxRetryExecutionThreshold {
			return Open
		}

		return HalfOpen
	}

	return Open
}

// NewState creates a new state
func NewState() *State {
	return &State{0, 0, 0, utcFuture, &sync.Mutex{}}
}
