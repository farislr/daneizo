package pkghttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Server_Serve(t *testing.T) {
	type exampleBody struct {
		Example string `json:"example"`
	}

	type fields struct {
		responseEncoder       ResponseEncoder
		errorResponseEncoder  ErrorResponseEncoder
		requestDecoders       []RequestDecoderOptionParam
		preRequestMiddlewares []PreRequestMiddlewareParam
	}
	type args struct {
		body any
		e    EndpointHandlerInputParam
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				responseEncoder: func(ctx context.Context, w http.ResponseWriter, response any) error {
					return nil
				},
				errorResponseEncoder: func(ctx context.Context, err error, w http.ResponseWriter) {
				},
			},
			args: args{
				e: EndpointHandler[exampleBody](func(ctx context.Context, r *Request[exampleBody]) (any, error) {
					return nil, nil
				}),
			},
		},
		{
			name: "success with request decoder",
			fields: fields{
				responseEncoder: func(ctx context.Context, w http.ResponseWriter, response any) error {
					return nil
				},
				errorResponseEncoder: func(ctx context.Context, err error, w http.ResponseWriter) {
				},
				requestDecoders: []RequestDecoderOptionParam{
					RequestDecoder[exampleBody](DefaultRequestDecoder[exampleBody]),
				},
			},
			args: args{
				e: EndpointHandler[exampleBody](func(ctx context.Context, r *Request[exampleBody]) (any, error) {
					return nil, nil
				}),
			},
		},
		{
			name: "endpoint error",
			fields: fields{
				responseEncoder: func(ctx context.Context, w http.ResponseWriter, response any) error {
					return nil
				},
				errorResponseEncoder: func(ctx context.Context, err error, w http.ResponseWriter) {
				},
				requestDecoders: []RequestDecoderOptionParam{
					RequestDecoder[exampleBody](DefaultRequestDecoder[exampleBody]),
				},
				preRequestMiddlewares: []PreRequestMiddlewareParam{
					PreRequestMiddleware[exampleBody](func(next EndpointHandler[exampleBody]) EndpointHandler[exampleBody] {
						return func(ctx context.Context, r *Request[exampleBody]) (any, error) {
							fmt.Println("PreRequestMiddleware executed")
							return next(ctx, r)
						}
					}),
				},
			},
			args: args{
				body: exampleBody{
					Example: "example",
				},
				e: EndpointHandler[exampleBody](func(ctx context.Context, req *Request[exampleBody]) (any, error) {

					fmt.Printf("req: %+v\n", req.body)
					return nil, errors.New("test endpoint error")
				}),
			},
		},
		{
			name: "encode response error",
			fields: fields{
				responseEncoder: func(ctx context.Context, w http.ResponseWriter, response any) error {
					return errors.New("error while encode response")
				},
				errorResponseEncoder: func(ctx context.Context, err error, w http.ResponseWriter) {
				},
			},
			args: args{
				e: EndpointHandler[exampleBody](func(ctx context.Context, r *Request[exampleBody]) (any, error) {
					return nil, nil
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			opts := []EndpointOption{}

			for _, dec := range tt.fields.requestDecoders {
				opts = append(opts, WithRequestDecoder(dec.(RequestDecoder[exampleBody])))
			}

			e := NewHandler(tt.args.e.(EndpointHandler[exampleBody]), opts...)

			var req *http.Request
			var err error
			if tt.args.body != nil {
				bodyByte, err := json.Marshal(tt.args.body)
				assert.NoError(t, err)
				bodyBuffer := bytes.NewBuffer(bodyByte)

				req, err = http.NewRequest("POST", "/example", bodyBuffer)
			} else {
				req, err = http.NewRequest("POST", "/example", nil)

			}

			if err != nil {
				assert.Error(t, err)
			}

			rr := httptest.NewRecorder()
			e.ServeHTTP(rr, req)
		})
	}
}
