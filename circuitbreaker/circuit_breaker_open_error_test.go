package circuitbreaker

import (
	"reflect"
	"testing"
)

func TestCircuitOpenError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    CircuitOpenError
		want string
	}{
		{"Error", NewCircuitOpenError("Test"), "Circuit is open for [Test]"},
	}
	for _, tt := range tests {
		if got := tt.e.Error(); got != tt.want {
			t.Errorf("%q. CircuitOpenError.Error() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNewCircuitOpenError(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want CircuitOpenError
	}{
		{"Constructor", args{"Key"}, NewCircuitOpenError("Key")},
	}
	for _, tt := range tests {
		if got := NewCircuitOpenError(tt.args.key); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewCircuitOpenError() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
