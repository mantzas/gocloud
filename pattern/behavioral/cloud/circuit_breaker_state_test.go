package cloud

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestState_Reset(t *testing.T) {

	assert := assert.New(t)

	state := NewState()
	state.IncreaseFailure()
	state.IncrementRetrySuccessCount()

	state.Reset()

	assert.Equal(0, state.currentFailureCount)
	assert.Equal(0, state.retrySuccessCount)
}

func TestState_IncreaseFailure(t *testing.T) {

	assert := assert.New(t)

	state := NewState()
	state.IncreaseFailure()

	assert.Equal(1, state.currentFailureCount)
}

func TestState_IncrementRetrySuccessCount(t *testing.T) {

	assert := assert.New(t)

	state := NewState()
	state.IncrementRetrySuccessCount()

	assert.Equal(1, state.retrySuccessCount)
}

func TestState_IncreaseExecutions(t *testing.T) {

	assert := assert.New(t)

	state := NewState()
	state.IncreaseExecutions()

	assert.Equal(1, state.currentExecutions)
}

func TestState_DecreaseExecutions(t *testing.T) {
	assert := assert.New(t)

	state := NewState()
	state.DecreaseExecutions()

	assert.Equal(-1, state.currentExecutions)
}

func TestState_GetStatus(t *testing.T) {

	setting := Setting{"Name", 1, time.Second, 1, 1}
	stateClosed := NewState()
	stateClosed.IncreaseFailure()
	stateHalf := NewState()
	stateHalf.IncreaseFailure()
	stateHalf.lastFailureTimestamp = stateHalf.lastFailureTimestamp.Add(-2 * time.Second)

	stateOpenMaxRetry := NewState()
	stateOpenMaxRetry.IncreaseFailure()
	stateOpenMaxRetry.lastFailureTimestamp = stateHalf.lastFailureTimestamp.Add(-2 * time.Second)
	stateOpenMaxRetry.IncreaseExecutions()
	stateOpenMaxRetry.IncreaseExecutions()

	stateClosedRetrySuccess := NewState()
	stateClosedRetrySuccess.IncreaseFailure()
	stateClosedRetrySuccess.lastFailureTimestamp = stateHalf.lastFailureTimestamp.Add(-2 * time.Second)
	stateClosedRetrySuccess.IncrementRetrySuccessCount()

	tests := []struct {
		name string
		s    *State
		sett *Setting
		want Status
	}{
		{"Closed", NewState(), &setting, Closed},
		{"Open", stateClosed, &setting, Open},
		{"HalfOpen", stateHalf, &setting, HalfOpen},
		{"Open Max Retry", stateOpenMaxRetry, &setting, Open},
		{"Closes after retry success", stateClosedRetrySuccess, &setting, Closed},
	}
	for _, tt := range tests {
		if got := tt.s.GetStatus(tt.sett); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. State.GetStatus() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNewState(t *testing.T) {

	assert := assert.New(t)

	state := NewState()

	assert.Equal(0, state.currentExecutions)
	assert.Equal(0, state.currentFailureCount)
	assert.Equal(0, state.retrySuccessCount)
	assert.Equal(time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC), state.lastFailureTimestamp)
}

func BenchmarkState_GetStatus(b *testing.B) {

	setting := Setting{"Name", 1, time.Second, 1, 1}
	state := NewState()

	for i := 0; i < b.N; i++ {
		state.GetStatus(&setting)
	}
}
