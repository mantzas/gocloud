package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCircuitBreaker(t *testing.T) {

	assert := assert.New(t)

	c := NewCircuitBreaker(NewLocalSettingsProvider(), &testMetric{})

	assert.NotNil(c)
}

func TestExecute_MissingKey(t *testing.T) {

	assert := assert.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	_, err := NewCircuitBreaker(pr, &testMetric{}).Execute("test1", testSuccessAction)

	assert.NotNil(err)

}

func TestExecute_MissingState(t *testing.T) {

	assert := assert.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	cb := NewCircuitBreaker(pr, &testMetric{})

	delete(cb.states, "test")

	_, err := cb.Execute("test", testSuccessAction)

	assert.NotNil(err)
}

func TestExecute_Closed(t *testing.T) {

	assert := assert.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	res, err := NewCircuitBreaker(pr, &testMetric{}).Execute("test", testSuccessAction)

	assert.Nil(err)
	assert.Equal("test", res)
}

func TestExecute_Open(t *testing.T) {

	assert := assert.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	cb := NewCircuitBreaker(pr, &testMetric{})
	cb.states["test"].currentFailureCount = 1

	_, err := cb.Execute("test", testSuccessAction)

	assert.NotNil(err)
}

func TestExecute_Failed(t *testing.T) {

	assert := assert.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	_, err := NewCircuitBreaker(pr, &testMetric{}).Execute("test", testFailureAction)

	assert.NotNil(err)
}

func TestExecute_SuccessAfterFailed(t *testing.T) {

	assert := assert.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	cb := NewCircuitBreaker(pr, &testMetric{})
	_, err := cb.Execute("test", testFailureAction)
	time.Sleep(2 * time.Second)
	_, err = cb.Execute("test", testSuccessAction)

	assert.Nil(err)
}

func BenchmarkCircuitBreaker_Execute(b *testing.B) {

	c := NewCircuitBreaker(NewLocalSettingsProvider(), &testMetric{})

	for i := 0; i < b.N; i++ {
		c.Execute("Test", testSuccessAction)
	}
}

func testSuccessAction() (interface{}, error) {
	return "test", nil
}

func testFailureAction() (interface{}, error) {
	return "", errors.New("Test error")
}
