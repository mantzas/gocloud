package circuitbreaker

import (
	"reflect"
	"testing"
	"time"
)

func TestState_Reset(t *testing.T) {
	state := NewState()
	state.IncreaseFailure()
	state.IncrementRetrySuccessCount()
	tests := []struct {
		name string
		s    *State
	}{
		{"Reset", state},
	}
	for _, tt := range tests {
		tt.s.Reset()
		if tt.s.currentFailureCount != 0 {
			t.Errorf("%q. State.Reset() not reset", tt.name)
		}
	}
}

func TestState_IncreaseFailure(t *testing.T) {
	tests := []struct {
		name string
		s    *State
	}{
		{"IncreaseFailure", NewState()},
	}
	for _, tt := range tests {
		tt.s.IncreaseFailure()
		if tt.s.currentFailureCount != 1 {
			t.Errorf("%q. State.IncreaseFailure() not incremented", tt.name)
		}
	}
}

func TestState_IncrementRetrySuccessCount(t *testing.T) {
	tests := []struct {
		name string
		s    *State
	}{
		{"IncrementRetrySuccessCount", NewState()},
	}
	for _, tt := range tests {
		tt.s.IncrementRetrySuccessCount()
		if tt.s.retrySuccessCount != 1 {
			t.Errorf("%q. State.IncrementRetrySuccessCount() not incremented", tt.name)
		}
	}
}

func TestState_IncreaseExecutions(t *testing.T) {
	tests := []struct {
		name string
		s    *State
	}{
		{"IncreaseExecutions", NewState()},
	}
	for _, tt := range tests {
		tt.s.IncreaseExecutions()
		if tt.s.currentExecutions != 1 {
			t.Errorf("%q. State.IncreaseExecutions() not incremented", tt.name)
		}
	}
}

func TestState_DecreaseExecutions(t *testing.T) {
	state := NewState()
	state.IncreaseExecutions()
	tests := []struct {
		name string
		s    *State
	}{
		{"DecreaseExecutions", state},
	}
	for _, tt := range tests {
		tt.s.DecreaseExecutions()
		if tt.s.currentExecutions != 0 {
			t.Errorf("%q. State.DecreaseExecutions() not decremented", tt.name)
		}
	}
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

	type args struct {
		sett *Setting
	}
	tests := []struct {
		name string
		s    *State
		args args
		want Status
	}{
		{"Closed", NewState(), args{sett: &setting}, Closed},
		{"Open", stateClosed, args{sett: &setting}, Open},
		{"HalfOpen", stateHalf, args{sett: &setting}, HalfOpen},
		{"Open Max Retry", stateOpenMaxRetry, args{sett: &setting}, Open},
		{"Closes after retry success", stateClosedRetrySuccess, args{sett: &setting}, Closed},
	}
	for _, tt := range tests {
		if got := tt.s.GetStatus(tt.args.sett); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. State.GetStatus() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNewState(t *testing.T) {
	tests := []struct {
		name string
		want *State
	}{
		{"NewState", NewState()},
	}
	for _, tt := range tests {
		if got := NewState(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewState() = %v, want %v", tt.name, got, tt.want)
		}

	}
}

func BenchmarkState_GetStatus(b *testing.B) {

	setting := Setting{"Name", 1, time.Second, 1, 1}
	state := NewState()

	for i := 0; i < b.N; i++ {
		state.GetStatus(&setting)
	}
}
