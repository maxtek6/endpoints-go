package endpoints

import (
	"fmt"
	"net/http"
)

type Endpoint struct {
	handlers            map[string]http.HandlerFunc
	onUnsupportedMethod http.HandlerFunc
}

func New() *Endpoint {
	return &Endpoint{
		handlers: map[string]http.HandlerFunc{},
	}
}

func (e *Endpoint) AddMethod(method string, f http.HandlerFunc) error {
	switch method {
	case http.MethodConnect:
		fallthrough
	case http.MethodDelete:
		fallthrough
	case http.MethodGet:
		fallthrough
	case http.MethodHead:
		fallthrough
	case http.MethodOptions:
		fallthrough
	case http.MethodPatch:
		fallthrough
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodTrace:
		e.handlers[method] = f
	default:
		return fmt.Errorf("invalid HTTP method \"%s\"", method)
	}
	return nil
}

func (e *Endpoint) RemoveMethod(method string) error {
	_, ok := e.handlers[method]
	if !ok {
		return fmt.Errorf("no handler associated with HTTP method \"%s\"", method)
	}
	delete(e.handlers, method)
	return nil
}

func (e *Endpoint) HandleUnsupportedMethod(f http.HandlerFunc) {
	e.onUnsupportedMethod = f
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := e.handlers[r.Method]
	if ok {
		handler(w, r)
	} else {
		if e.onUnsupportedMethod != nil {
			e.onUnsupportedMethod(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
