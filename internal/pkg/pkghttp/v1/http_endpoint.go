package pkghttp

import (
	"context"
	"net/http"
)

type (
	EndpointHandler[T any] func(ctx context.Context, r *Request[T]) (response interface{}, err error)

	EndpointHandlerInputParam interface {
		Exec(ctx context.Context, r RequestInputParam) (response interface{}, err error)
	}

	Endpoint[T any] struct {
		handler EndpointHandler[T]

		endpointOptionParam EndpointOptionParam
	}
)

func (e EndpointHandler[T]) Exec(ctx context.Context, r RequestInputParam) (response interface{}, err error) {
	if req, ok := r.(*Request[T]); ok {
		return e(ctx, req)
	}

	return nil, nil
}

func (e *Endpoint[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := NewRequest[T](r)
	if err != nil {
		e.endpointOptionParam.errorResponseEncoder(ctx, err, w)

		return
	}

	if e.endpointOptionParam.requestDecoderParams != nil {
		for _, rd := range e.endpointOptionParam.requestDecoderParams {
			var err error
			ctx, err = rd.Exec(ctx, req)
			if err != nil {
				e.endpointOptionParam.errorResponseEncoder(ctx, err, w)

				return
			}
		}
	}

	handler := e.handler
	for _, m := range e.endpointOptionParam.middlewares {
		handler = m.Exec(handler).(EndpointHandler[T])
	}

	res, err := handler(ctx, req)
	if err != nil {
		e.endpointOptionParam.errorResponseEncoder(ctx, err, w)

		return
	}

	if err := e.endpointOptionParam.responseEncoder(ctx, w, res); err != nil {
		e.endpointOptionParam.errorResponseEncoder(ctx, err, w)

		return
	}
}

func (e *Endpoint[T]) getOptions() EndpointOptionParam {
	return e.endpointOptionParam
}
