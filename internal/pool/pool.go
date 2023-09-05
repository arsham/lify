// Package pool provides a generic pool of objects.
package pool

import (
	"sync"
)

// Pool is a generic pool of objects. Before the object is returned to the
// pool, the cleanup function is called.
type Pool[T any] struct {
	pool *sync.Pool
}

// NewPool creates a new pool of objects. The fn function should reset the
// necessary fields of the object to its zero value.
func NewPool[T any](fn func() T) Pool[T] {
	return Pool[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				return fn()
			},
		},
	}
}

// Put returns the object to the pool.
func (p *Pool[T]) Put(v T) {
	p.pool.Put(v)
}

// Get returns an object from the pool.
func (p *Pool[T]) Get() T {
	// nolint:forcetypeassert // we know the type.
	return p.pool.Get().(T)
}
