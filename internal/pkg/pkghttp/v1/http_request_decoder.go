package pkghttp

import (
	"context"
)

// RequestDecoder function defines the contract for decoding request.
type (
	RequestDecoder[T any] func(ctx context.Context, r *Request[T]) (context.Context, error)

	RequestDecoderOptionParam interface {
		Exec(ctx context.Context, r RequestInputParam) (context.Context, error)
	}
)

func (rd RequestDecoder[T]) Exec(ctx context.Context, r RequestInputParam) (context.Context, error) {
	if req, ok := r.(*Request[T]); ok {
		return rd(ctx, req)
	}

	return ctx, nil
}

func DefaultRequestDecoder[T any](ctx context.Context, _ *Request[T]) (context.Context, error) {
	return ctx, nil
}
