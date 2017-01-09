package cloud

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCircuitOpenError_Error(t *testing.T) {

	require := require.New(t)

	err := NewCircuitOpenError("Test")

	require.Equal("Circuit is open for [Test]", err.Error())
}

func TestNewCircuitOpenError_New(t *testing.T) {

	require := require.New(t)

	err := NewCircuitOpenError("Test")

	require.Equal("Test", err.key)
}
