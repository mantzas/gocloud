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

// ObjectPool definition
type ObjectPool struct {
	pool sync.Pool
	s    ObjectSanitizer
}

// ObjectFactory function definition
type ObjectFactory func() interface{}

// ObjectSanitizer function definition
type ObjectSanitizer func(interface{}) interface{}

// NewObjectPool constructor
func NewObjectPool(f ObjectFactory, s ObjectSanitizer) (*ObjectPool, error) {

	if f == nil {
		return nil, errors.New("object factory is nil")
	}

	if s == nil {
		s = func(in interface{}) interface{} {
			return in
		}
	}

	return &ObjectPool{sync.Pool{
		New: f,
	}, s}, nil
}

// Rent returns a object from pool
func (op *ObjectPool) Rent() interface{} {
	return op.pool.Get()
}

// Return object to the pool
func (op *ObjectPool) Return(o interface{}) {
	op.pool.Put(op.s(o))
}
