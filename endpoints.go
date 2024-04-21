// Copyright (c) 2024 Maxtek Consulting
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package endpoints

import (
	"fmt"
	"net/http"
)

// The Endpoint type represents a single HTTP endpoint. All requests
// to this URL will be served through this object.
type Endpoint struct {
	handlers            map[string]http.HandlerFunc
	onUnsupportedMethod http.HandlerFunc
}

// New creates a new Endpoint.
//
// The newly constructed Endpoint has no supported method and cannot
// be used to serve requests until it is mounted using http.Handle or
// a similar function.
func New() *Endpoint {
	return &Endpoint{
		handlers: map[string]http.HandlerFunc{},
	}
}

// The AddMethod function adds support for a specific HTTP method.
//
// Each HTTP method is mapped to an http.HandlerFunc object. If the
// method has been previously added with another http.HandlerFunc,
// the old value will be overwritten. An error will be returned if
// method does not match a valid HTTP method.
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

// Remove function removes support for a specific HTTP method.
//
// If the method has been previously mapped to an http.HandlerFunc,
// the value will be removed from the map and the Endpoint will no
// longer support that method. An error will be returned if there is
// no http.HandlerFunc associated wit the method string.
func (e *Endpoint) RemoveMethod(method string) error {
	_, ok := e.handlers[method]
	if !ok {
		return fmt.Errorf("no handler associated with HTTP method \"%s\"", method)
	}
	delete(e.handlers, method)
	return nil
}

// HandleUnsupportedMethod sets an http.HandlerFunc for unsupported
// methods.
//
// By default, an http.Request that is passed to an endpoint will
// return http.StatusMethodNotAllowed and no response body if the
// method has not been added to the endpoint. Using this function
// will redirect all requests with unsupported request methods to
// a single http.HandlerFunc.
func (e *Endpoint) HandleUnsupportedMethod(f http.HandlerFunc) {
	e.onUnsupportedMethod = f
}

// ServeHTTP implements the http.Handler interface.
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
