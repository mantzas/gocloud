package circuitbreaker

import (
	"testing"
	"time"

	"github.com/mantzas/gocloud/metrics"
	"github.com/stretchr/testify/assert"
)

type testMetric struct {
	failures   int
	executions int
	key        string
}

func (tm *testMetric) Increase(value int, tags ...metrics.Tag) {
	for _, tag := range tags {
		switch tag.Key {
		case "key":
			tm.key = tag.Value
		case "status":
			if tag.Value == "failure" {
				tm.failures += value
			} else {
				tm.executions += value
			}
		}
	}
}

func TestState_NewWithNilCounter(t *testing.T) {
	assert := assert.New(t)
	state := NewState(nil, "key")
	assert.NotNil(state)
}

func TestState_Reset(t *testing.T) {
	assert := assert.New(t)
	metric := testMetric{}
	state := NewState(&metric, "key")
	state.IncreaseFailure()
	state.IncrementRetrySuccessCount()

	state.Reset()

	assert.Equal(0, state.currentFailureCount)
	assert.Equal(0, state.retrySuccessCount)
	assert.Equal(1, metric.failures)
	assert.Equal(0, metric.executions)
	assert.Equal("key", metric.key)
}

func TestState_IncreaseFailure(t *testing.T) {
	assert := assert.New(t)
	metric := testMetric{}
	state := NewState(&metric, "key")
	state.IncreaseFailure()
	assert.Equal(1, state.currentFailureCount)
	assert.Equal(1, metric.failures)
	assert.Equal(0, metric.executions)
	assert.Equal("key", metric.key)
}

func TestState_IncrementRetrySuccessCount(t *testing.T) {
	assert := assert.New(t)
	metric := testMetric{}
	state := NewState(&metric, "key")
	state.IncrementRetrySuccessCount()
	assert.Equal(1, state.retrySuccessCount)
	assert.Equal(0, metric.failures)
	assert.Equal(0, metric.executions)
	assert.Equal("", metric.key)
}

func TestState_IncreaseExecutions(t *testing.T) {
	assert := assert.New(t)
	metric := testMetric{}
	state := NewState(&metric, "key")
	state.IncreaseExecutions()
	assert.Equal(1, state.currentExecutions)
	assert.Equal(0, metric.failures)
	assert.Equal(1, metric.executions)
	assert.Equal("key", metric.key)
}

func TestState_DecreaseExecutions(t *testing.T) {
	assert := assert.New(t)
	metric := testMetric{}
	state := NewState(&metric, "key")
	state.DecreaseExecutions()
	assert.Equal(-1, state.currentExecutions)
	assert.Equal(0, metric.failures)
	assert.Equal(0, metric.executions)
	assert.Equal("", metric.key)
}

func TestState_GetStatus(t *testing.T) {
	assert := assert.New(t)
	setting := Setting{"Name", 1, time.Second, 1, 1}
	stateClosed := NewState(&testMetric{}, "key")
	stateClosed.IncreaseFailure()
	stateHalf := NewState(&testMetric{}, "key")
	stateHalf.IncreaseFailure()
	stateHalf.lastFailureTimestamp = stateHalf.lastFailureTimestamp.Add(-2 * time.Second)

	stateOpenMaxRetry := NewState(&testMetric{}, "key")
	stateOpenMaxRetry.IncreaseFailure()
	stateOpenMaxRetry.lastFailureTimestamp = stateHalf.lastFailureTimestamp.Add(-2 * time.Second)
	stateOpenMaxRetry.IncreaseExecutions()
	stateOpenMaxRetry.IncreaseExecutions()

	stateClosedRetrySuccess := NewState(&testMetric{}, "key")
	stateClosedRetrySuccess.IncreaseFailure()
	stateClosedRetrySuccess.lastFailureTimestamp = stateHalf.lastFailureTimestamp.Add(-2 * time.Second)
	stateClosedRetrySuccess.IncrementRetrySuccessCount()

	tests := []struct {
		name       string
		s          *State
		sett       *Setting
		want       Status
		wantMetric testMetric
	}{
		{"Closed", NewState(&testMetric{}, "key"), &setting, Closed, testMetric{0, 0, ""}},
		{"Open", stateClosed, &setting, Open, testMetric{1, 0, "key"}},
		{"HalfOpen", stateHalf, &setting, HalfOpen, testMetric{1, 0, "key"}},
		{"Open Max Retry", stateOpenMaxRetry, &setting, Open, testMetric{1, 2, "key"}},
		{"Closes after retry success", stateClosedRetrySuccess, &setting, Closed, testMetric{1, 0, "key"}},
	}
	for _, tt := range tests {

		assert.Equal(tt.want, tt.s.GetStatus(tt.sett), tt.name)
		assert.Equal(&tt.wantMetric, tt.s.counter, tt.name)
	}
}

func TestNewState(t *testing.T) {

	assert := assert.New(t)

	state := NewState(&testMetric{}, "key")

	assert.Equal(0, state.currentExecutions)
	assert.Equal(0, state.currentFailureCount)
	assert.Equal(0, state.retrySuccessCount)
	assert.Equal(time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC), state.lastFailureTimestamp)
}

func BenchmarkState_GetStatus(b *testing.B) {

	setting := Setting{"Name", 1, time.Second, 1, 1}
	state := NewState(&testMetric{}, "key")

	for i := 0; i < b.N; i++ {
		state.GetStatus(&setting)
	}
}
