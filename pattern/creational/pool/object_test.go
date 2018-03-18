package pool

import (
	"testing"

	"github.com/mantzas/gocloud/metrics"
	"github.com/stretchr/testify/assert"
)

type testGauge struct {
	size       float64
	key, value string
}

func (tc *testGauge) Set(value float64, tags ...metrics.Tag) {
	tc.size = value
	tc.key = tags[0].Key
	tc.value = tags[0].Value
}

func TestNewObjectPool(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	p, err := NewChannelPool(10, tf.testFactory, tf.testSanitizer, nil, "test")
	assert.NoError(err)
	assert.NotNil(p)
	assert.Equal(0, tf.factoryCalled)
	assert.Equal(0, tf.sanitizerCalled)
}

func TestNewObjectPoolZeroSize(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	p, err := NewChannelPool(0, tf.testFactory, tf.testSanitizer, nil, "test")
	assert.Error(err)
	assert.Nil(p)
}

func TestNewObjectPoolNilFactoryError(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	p, err := NewChannelPool(10, nil, tf.testSanitizer, nil, "test")
	assert.Error(err)
	assert.Nil(p)
}

func TestNewObjectPoolWithNilSanitizer(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	p, err := NewChannelPool(10, tf.testFactory, nil, nil, "test")
	assert.NoError(err)
	assert.NotNil(p)
	assert.Equal(0, tf.factoryCalled)
	assert.Equal(0, tf.sanitizerCalled)
}

func TestObjectPool_Rent(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	p, err := NewChannelPool(10, tf.testFactory, tf.testSanitizer, nil, "test")
	assert.NoError(err)
	assert.NotNil(p)

	o := p.Rent().(testObject)
	assert.Equal("test", o.name, "Expected 'name' but got %s", o.name)
	assert.Equal(1, tf.factoryCalled)
	assert.Equal(0, tf.sanitizerCalled)
}

func TestObjectPool_ReturnNilSanitizer(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	p, err := NewChannelPool(10, tf.testFactory, nil, nil, "test")
	assert.NoError(err)
	assert.NotNil(p)

	o := p.Rent().(testObject)
	assert.Equal("test", o.name, "Expected 'name' but got %s", o.name)

	p.Return(o)

	assert.Equal(1, tf.factoryCalled)
}

func TestObjectPool_Test(t *testing.T) {
	assert := assert.New(t)
	tf := testFunc{}
	tg := testGauge{}
	p, err := NewChannelPool(1, tf.testFactory, tf.testSanitizer, &tg, "test")
	assert.NoError(err)
	assert.NotNil(p)

	o := p.Rent().(testObject)
	assert.Equal("test", o.name, "Expected 'name' but got %s", o.name)

	p.Return(o)
	p.Return(o)

	o = p.Rent().(testObject)
	assert.Equal("test", o.name, "Expected 'name' but got %s", o.name)

	assert.Equal(1, tf.factoryCalled)
	assert.Equal(2, tf.sanitizerCalled)

	assert.Equal(1.0, tg.size)
	assert.Equal("type", tg.key)
	assert.Equal("test", tg.value)
}

var item interface{}

func BenchmarkState_ChannelPool(b *testing.B) {

	tf := testFunc{}
	p, _ := NewChannelPool(2, tf.testFactory, tf.testSanitizer, nil, "test")
	p.Return(p.Rent())

	for i := 0; i < b.N; i++ {
		item = p.Rent()
		p.Return(item)
	}
}

type testObject struct {
	name string
}

type testFunc struct {
	factoryCalled   int
	sanitizerCalled int
}

func (tf *testFunc) testFactory() interface{} {
	tf.factoryCalled++
	return testObject{name: "test"}
}

func (tf *testFunc) testSanitizer(o interface{}) interface{} {
	tf.sanitizerCalled++
	t := o.(testObject)
	t.name = ""
	return o
}
