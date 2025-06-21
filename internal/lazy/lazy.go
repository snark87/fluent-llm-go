package lazy

import (
	"context"
	"sync"
)

type Lazy[T any] struct {
	value func() T
}

func New[T any](value func() T) Lazy[T] {
	return Lazy[T]{
		value: sync.OnceValue(value),
	}
}

func (l *Lazy[T]) Get() T {
	return l.value()
}

type LazyWithError[T any] struct {
	valueFn func(context.Context) (T, error)
	value   *T
}

func NewWithError[T any](value func(context.Context) (T, error)) LazyWithError[T] {
	return LazyWithError[T]{
		valueFn: value,
	}
}

func (l *LazyWithError[T]) Get(ctx context.Context) (T, error) {
	if l.value == nil {
		value, err := l.valueFn(ctx)
		if err != nil {
			return *new(T), err
		}
		l.value = &value
	}
	return *l.value, nil
}
