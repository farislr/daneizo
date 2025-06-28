package pkghttp

type (
	EndpointOptionParam struct {
		requestDecoderParams []RequestDecoderOptionParam
		middlewares          []PreRequestMiddlewareParam
		responseEncoder      ResponseEncoder
		errorResponseEncoder ErrorResponseEncoder
	}

	EndpointOption func(*EndpointOptionParam)
)

func WithRequestDecoder[T any](decoder RequestDecoder[T]) EndpointOption {
	return func(eop *EndpointOptionParam) {
		eop.requestDecoderParams = append(eop.requestDecoderParams, decoder)
	}
}

func WithEndpointResponseEncoder(encoder ResponseEncoder) EndpointOption {
	return func(eop *EndpointOptionParam) {
		eop.responseEncoder = encoder
	}
}

func WithEndpointErrorResponseEncoder(encoder ErrorResponseEncoder) EndpointOption {
	return func(eop *EndpointOptionParam) {
		eop.errorResponseEncoder = encoder
	}
}

func WithPreRequestMiddleware[T any](m PreRequestMiddleware[T]) EndpointOption {
	return func(eop *EndpointOptionParam) {
		eop.middlewares = append(eop.middlewares, m)
	}
}
