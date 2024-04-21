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
	endpoint.AddMethod(http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
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
