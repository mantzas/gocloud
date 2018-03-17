package circuitbreaker

import (
	"sync"
	"time"

	"github.com/mantzas/gocloud/metrics"
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
	counter              metrics.Counter
	keyTag               metrics.Tag
	failureTag           metrics.Tag
	totalTag             metrics.Tag
	m                    *sync.Mutex
}

// NewState creates a new state. If no metric counter is provided
// the default null counter is used.
func NewState(c metrics.Counter, key string) *State {
	k := metrics.NewTag("key", key)
	f := metrics.NewTag("status", "failure")
	t := metrics.NewTag("status", "executions")

	if c == nil {
		c = &metrics.NullCounter{}
	}

	return &State{0, 0, 0, utcFuture, c, k, f, t, &sync.Mutex{}}
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
	s.counter.Increase(1, s.keyTag, s.failureTag)
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
	s.counter.Increase(1, s.keyTag, s.totalTag)
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
