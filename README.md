# Endpoints

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://pkg.go.dev/github.com/maxtek6/endpoints-go)
[![codecov](https://codecov.io/gh/maxtek6/endpoints-go/branch/master/graph/badge.svg)](https://codecov.io/gh/maxtek6/endpoints-go)

Go module for serving multiple methods on a single HTTP path.

## Usage

The `Endpoint` struct is an implementation of the `http.Handler`
interface. Once an endpoint has been created, it can be served
using `http.Handle()`, or any other function that accepts the
`http.Handler` interface:

```go
// create handler function
handleGet := func(w http.ResponseWriter, r *http.Request) {
    // handle GET method
}

// create new endpoint
endpoint := endpoints.New()

// map support for GET method
_ = endpoint.AddMethod(http.MethodGet, handleGet)

// mount endpoint to /path
http.Handle("/path", endpoint)
```