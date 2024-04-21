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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddMethod(t *testing.T) {
	endpoint := New()
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {

	}

	err := endpoint.AddMethod("BADMETHOD", dummyHandler)
	assert.Error(t, err)

	err = endpoint.AddMethod(http.MethodConnect, dummyHandler)
	assert.NoError(t, err)
}

func TestRemoveMethod(t *testing.T) {
	endpoint := New()
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {

	}

	_ = endpoint.AddMethod(http.MethodGet, dummyHandler)

	err := endpoint.RemoveMethod(http.MethodPost)
	assert.Error(t, err)

	err = endpoint.RemoveMethod(http.MethodGet)
	assert.NoError(t, err)
}

type TestResponseWriter struct {
}

func NewTestResponseWriter() *TestResponseWriter {
	return nil
}

func (w *TestResponseWriter) Header() http.Header {
	return nil
}

func (w *TestResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w *TestResponseWriter) WriteHeader(statusCode int) {}

func TestServerHTTP(t *testing.T) {
	postHandled := false
	unsupportedHandled := false

	postRequest, _ := http.NewRequest(http.MethodPost, "http://example.com/test", nil)
	getRequest, _ := http.NewRequest(http.MethodGet, "http://example.com/test", nil)

	responseWriter := NewTestResponseWriter()

	endpoint := New()
	_ = endpoint.AddMethod(http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		postHandled = true
	})

	endpoint.ServeHTTP(responseWriter, postRequest)
	assert.True(t, postHandled)

	endpoint.ServeHTTP(responseWriter, getRequest)
	assert.False(t, unsupportedHandled)

	endpoint.HandleUnsupportedMethod(func(w http.ResponseWriter, r *http.Request) {
		unsupportedHandled = true
	})
	endpoint.ServeHTTP(responseWriter, getRequest)
	assert.True(t, unsupportedHandled)
}
