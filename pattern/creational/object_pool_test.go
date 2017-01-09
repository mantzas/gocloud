package creational

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewObjectPool(t *testing.T) {
	req := require.New(t)
	tf := testFunc{}
	p := NewObjectPool(tf.testObjectFactory, tf.testObjectSanitizer)
	req.NotNil(p)
	req.Equal(0, tf.factoryCalled)
	req.Equal(0, tf.sanitizerCalled)
}

func TestObjectPool_Rent(t *testing.T) {
	req := require.New(t)
	tf := testFunc{}
	p := NewObjectPool(tf.testObjectFactory, tf.testObjectSanitizer)
	req.NotNil(p)

	o := p.Rent().(testObject)
	req.Equal("test", o.name, "Expected 'name' but got %s", o.name)
	req.Equal(1, tf.factoryCalled)
	req.Equal(0, tf.sanitizerCalled)
}

func TestObjectPool_Return(t *testing.T) {
	req := require.New(t)
	tf := testFunc{}
	p := NewObjectPool(tf.testObjectFactory, tf.testObjectSanitizer)
	req.NotNil(p)

	o := p.Rent().(testObject)
	req.Equal("test", o.name, "Expected 'name' but got %s", o.name)

	p.Return(o)

	req.Equal(1, tf.factoryCalled)
	req.Equal(1, tf.sanitizerCalled)
}

type testObject struct {
	name string
}

type testFunc struct {
	factoryCalled   int
	sanitizerCalled int
}

func (tf *testFunc) testObjectFactory() interface{} {
	tf.factoryCalled += 1
	return testObject{name: "test"}
}

func (tf *testFunc) testObjectSanitizer(o interface{}) interface{} {
	tf.sanitizerCalled += 1
	t := o.(testObject)
	t.name = ""
	return o
}
