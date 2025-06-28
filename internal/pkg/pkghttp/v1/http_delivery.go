package pkghttp

func NewHandler[T any](handler EndpointHandler[T], options ...EndpointOption) *Endpoint[T] {
	endpointOptionParam := EndpointOptionParam{
		responseEncoder:      DefaultResponseEncoder,
		errorResponseEncoder: DefaultErrorEncoder,
	}

	for _, option := range options {
		option(&endpointOptionParam)
	}

	return &Endpoint[T]{
		handler: handler,

		endpointOptionParam: endpointOptionParam,
	}
}
