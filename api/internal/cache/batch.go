package cache

import (
	"context"
	"time"
)

type Batch interface {
	Get(ctx context.Context, key string) *Result[[]byte]
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) *Result[struct{}]
	Delete(ctx context.Context, keys ...string) *Result[int64]
	Increment(ctx context.Context, key string) *Result[int64]
	TTL(ctx context.Context, key string) *Result[time.Duration]

	Exec(ctx context.Context) error
}

func NewResult[T any](value T, err error) *Result[T] {
	return &Result[T]{value: value, err: err}
}

type Result[T any] struct {
	// campos privados; o adapter os preenche no Exec
	value T
	err   error
}

func (r *Result[T]) Resolve(value T, err error) {
	r.value = value
	r.err = err
}

func (r *Result[T]) Value() (T, error) {
	return r.value, r.err
}
