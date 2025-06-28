package pkghttp

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EndpointOption(t *testing.T) {

	tests := []struct {
		name     string
		endpoint interface {
			http.Handler
			getOptions() EndpointOptionParam
		}
		option    func(opt *EndpointOptionParam)
		assertion func(opt *EndpointOptionParam) bool
	}{
		{
			name: "default",
			endpoint: NewHandler(func(ctx context.Context, r *Request[struct {
				ID string
			}]) (response interface{}, err error) {
				return nil, nil
			}),
			option: func(opt *EndpointOptionParam) {

			},
			assertion: func(opt *EndpointOptionParam) bool {
				return opt.responseEncoder != nil &&
					opt.errorResponseEncoder != nil
			},
		},
		{
			name: "with endpoint response encoder",
			endpoint: NewHandler(func(ctx context.Context, r *Request[struct {
				ID string
			}]) (response interface{}, err error) {
				return nil, nil
			}),
			option: WithEndpointResponseEncoder(nil),
			assertion: func(endpoint *EndpointOptionParam) bool {
				return endpoint.responseEncoder == nil
			},
		},
		{
			name: "with endpoint error response encoder",
			endpoint: NewHandler(func(ctx context.Context, r *Request[struct {
				ID string
			}]) (response interface{}, err error) {
				return nil, nil
			}),
			option: WithEndpointErrorResponseEncoder(nil),
			assertion: func(opt *EndpointOptionParam) bool {
				return opt.errorResponseEncoder == nil
			},
		},
		{
			name: "with endpoint middlewares",
			endpoint: NewHandler(func(ctx context.Context, r *Request[struct {
				ID string
			}]) (response interface{}, err error) {
				return nil, nil
			}),
			option: WithPreRequestMiddleware(
				func(next EndpointHandler[struct{ ID string }]) EndpointHandler[struct{ ID string }] {
					return func(ctx context.Context, r *Request[struct{ ID string }]) (response interface{}, err error) {
						return next(ctx, r)
					}
				},
			),
			assertion: func(endpoint *EndpointOptionParam) bool {
				return len(endpoint.middlewares) == 1
			},
		},
		{
			name: "with request decoder", endpoint: NewHandler(func(ctx context.Context, r *Request[struct {
				ID string
			}]) (response interface{}, err error) {
				return nil, nil
			}),
			option: WithRequestDecoder[struct{ ID string }](WithPopulateContextFromHeader),
			assertion: func(opt *EndpointOptionParam) bool {
				return opt.requestDecoderParams != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.endpoint.getOptions()
			tt.option(&opts)

			assert.True(t, tt.assertion(&opts))
		})
	}
}
