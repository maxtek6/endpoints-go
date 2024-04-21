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
