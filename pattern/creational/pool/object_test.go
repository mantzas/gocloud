package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewObjectPool(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	p := NewObjectPool(tf.testObjectFactory, tf.testObjectSanitizer)
	assert.NotNil(p)
	assert.Equal(0, tf.factoryCalled)
	assert.Equal(0, tf.sanitizerCalled)
}

func TestObjectPool_Rent(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	p := NewObjectPool(tf.testObjectFactory, tf.testObjectSanitizer)
	assert.NotNil(p)

	o := p.Rent().(testObject)
	assert.Equal("test", o.name, "Expected 'name' but got %s", o.name)
	assert.Equal(1, tf.factoryCalled)
	assert.Equal(0, tf.sanitizerCalled)
}

func TestObjectPool_Return(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	p := NewObjectPool(tf.testObjectFactory, tf.testObjectSanitizer)
	assert.NotNil(p)

	o := p.Rent().(testObject)
	assert.Equal("test", o.name, "Expected 'name' but got %s", o.name)

	p.Return(o)

	assert.Equal(1, tf.factoryCalled)
	assert.Equal(1, tf.sanitizerCalled)
}

type testObject struct {
	name string
}

type testFunc struct {
	factoryCalled   int
	sanitizerCalled int
}

func (tf *testFunc) testObjectFactory() interface{} {
	tf.factoryCalled++
	return testObject{name: "test"}
}

func (tf *testFunc) testObjectSanitizer(o interface{}) interface{} {
	tf.sanitizerCalled++
	t := o.(testObject)
	t.name = ""
	return o
}
