package pkghttp

type (
	PreRequestMiddlewareParam interface {
		Exec(next EndpointHandlerInputParam) EndpointHandlerInputParam
	}

	PreRequestMiddleware[T any] func(next EndpointHandler[T]) EndpointHandler[T]
)

func (m PreRequestMiddleware[T]) Exec(next EndpointHandlerInputParam) EndpointHandlerInputParam {
	if handler, ok := next.(EndpointHandler[T]); ok {
		return m(handler)
	}

	return next
}
