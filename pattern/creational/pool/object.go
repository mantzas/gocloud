package pool

import (
	"errors"
	"sync"
)

// Pool interface
type Pool interface {
	Rent() interface{}
	Return(interface{})
}

// ChannelPool definition
type ChannelPool struct {
	size int
	pool chan interface{}
	s    Sanitizer
	f    Factory
	m    *sync.Mutex
}

// Factory function definition
type Factory func() interface{}

// Sanitizer function definition
type Sanitizer func(interface{}) interface{}

// NewChannelPool constructor
func NewChannelPool(size int, f Factory, s Sanitizer) (*ChannelPool, error) {

	if size <= 0 {
		return nil, errors.New("size must be positive")
	}

	if f == nil {
		return nil, errors.New("object factory is nil")
	}

	if s == nil {
		s = func(in interface{}) interface{} {
			return in
		}
	}

	pool := make(chan interface{}, size)

	return &ChannelPool{size, pool, s, f, &sync.Mutex{}}, nil
}

// Rent returns a object from pool
func (op *ChannelPool) Rent() interface{} {
	op.m.Lock()
	defer op.m.Unlock()

	if len(op.pool) == 0 {
		return op.f()
	}

	return <-op.pool
}

// Return object to the pool
func (op *ChannelPool) Return(o interface{}) {
	op.m.Lock()
	defer op.m.Unlock()

	if len(op.pool) == op.size {
		<-op.pool
	}
	op.pool <- op.s(o)
}
