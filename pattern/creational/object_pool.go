package creational

import (
	"sync"
)

type Pool interface {
	Rent() interface{}
	Return(interface{})
}

type ObjectPool struct {
	pool sync.Pool
	s    ObjectSanitizer
}

type ObjectFactory func() interface{}
type ObjectSanitizer func(interface{}) interface{}

func NewObjectPool(f ObjectFactory, s ObjectSanitizer) *ObjectPool {
	return &ObjectPool{sync.Pool{
		New: f,
	}, s}
}

func (op *ObjectPool) Rent() interface{} {
	return op.pool.Get()
}

func (op *ObjectPool) Return(o interface{}) {
	op.pool.Put(op.s(o))
}
