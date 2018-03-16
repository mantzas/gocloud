package cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCircuitOpenError_Error(t *testing.T) {

	assert := assert.New(t)

	err := NewCircuitOpenError("Test")

	assert.Equal("Circuit is open for [Test]", err.Error())
}

func TestNewCircuitOpenError_New(t *testing.T) {

	assert := assert.New(t)

	err := NewCircuitOpenError("Test")

	assert.Equal("Test", err.key)
}
