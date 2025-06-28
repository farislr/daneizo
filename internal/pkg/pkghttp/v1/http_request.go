package pkghttp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type (
	RequestInputParam interface{}

	RequestReadWriter interface {
		Decode(v any) error
		Encode(v any) error
		Header() http.Header
		URL() *url.URL
	}

	Request[T any] struct {
		httpReq *http.Request
		body    T
	}
)

func NewRequest[T any](r *http.Request) (*Request[T], error) {
	ri := &Request[T]{
		httpReq: r,
	}

	if ri.httpReq.Body != nil && ri.httpReq.ContentLength > 0 {
		if err := json.NewDecoder(ri.httpReq.Body).Decode(&ri.body); err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
	}

	return ri, nil
}

func (r *Request[T]) Body() T {
	return r.body
}

func (r *Request[T]) Header() http.Header {
	return r.httpReq.Header
}

func (r *Request[T]) URL() *url.URL {
	return r.httpReq.URL
}

func (r *Request[T]) Raw() *http.Request {
	return r.httpReq
}
