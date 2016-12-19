package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewCircuitBreaker(t *testing.T) {

	require := require.New(t)

	c := NewCircuitBreaker(NewLocalSettingsProvider())

	require.NotNil(c)
}

func TestExecute_MissingKey(t *testing.T) {

	require := require.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	_, err := NewCircuitBreaker(pr).Execute("test1", testSuccessAction)

	require.NotNil(err)

}

func TestExecute_MissingState(t *testing.T) {

	require := require.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	cb := NewCircuitBreaker(pr)

	delete(cb.states, "test")

	_, err := cb.Execute("test", testSuccessAction)

	require.NotNil(err)
}

func TestExecute_Closed(t *testing.T) {

	require := require.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	res, err := NewCircuitBreaker(pr).Execute("test", testSuccessAction)

	require.Nil(err)
	require.Equal("test", res)
}

func TestExecute_Open(t *testing.T) {

	require := require.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	cb := NewCircuitBreaker(pr)
	cb.states["test"].currentFailureCount = 1

	_, err := cb.Execute("test", testSuccessAction)

	require.NotNil(err)
}

func TestExecute_Failed(t *testing.T) {

	require := require.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	_, err := NewCircuitBreaker(pr).Execute("test", testFailureAction)

	require.NotNil(err)
}

func TestExecute_SuccessAfterFailed(t *testing.T) {

	require := require.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "test", FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1})

	cb := NewCircuitBreaker(pr)
	_, err := cb.Execute("test", testFailureAction)
	time.Sleep(2 * time.Second)
	_, err = cb.Execute("test", testSuccessAction)

	require.Nil(err)
}

func testSuccessAction() (interface{}, error) {
	return "test", nil
}

func testFailureAction() (interface{}, error) {
	return "", errors.New("Test error")
}
