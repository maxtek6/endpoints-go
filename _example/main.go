package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/maxtek6/endpoints-go"
)

var (
	server    *http.Server
	resources map[string]interface{}
)

func main() {
	setup()

	go func() {
		_ = server.ListenAndServe()
	}()
	defer func() {
		server.Shutdown(context.TODO())
	}()
	runClient()
}

func setup() {
	router := mux.NewRouter()

	resourcesEndpoint := endpoints.New()
	_ = resourcesEndpoint.AddMethod(http.MethodGet, listResources)
	_ = resourcesEndpoint.AddMethod(http.MethodPost, createResource)
	_ = router.Handle("/resources", resourcesEndpoint)

	resourceEndpoint := endpoints.New()
	_ = resourceEndpoint.AddMethod(http.MethodGet, fetchResource)
	_ = resourceEndpoint.AddMethod(http.MethodDelete, deleteResource)
	_ = router.Handle("/resources/{id}", resourcesEndpoint)

	server = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}

type ListResourcesResponse struct {
	Resources []string `json:"resources"`
}

type CreateResourceResponse struct {
	ID string `json:"id"`
}

func listResources(w http.ResponseWriter, r *http.Request) {
	list := []string{}
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	for key := range resources {
		list = append(list, key)
	}
	_ = encoder.Encode(list)
	_, _ = w.Write(buffer.Bytes())
}

func createResource(w http.ResponseWriter, r *http.Request) {
	var resource interface{}
	id := uuid.NewString()
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&resource)
	resources[id] = resource
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	_ = encoder.Encode(CreateResourceResponse{
		ID: id,
	})
	_, _ = w.Write(buffer.Bytes())
}

func fetchResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	resource, ok := resources[id]
	if ok {
		buffer := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buffer)
		_ = encoder.Encode(resource)
		_, _ = w.Write(buffer.Bytes())
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	statusCode := http.StatusNotFound
	_, ok := resources[id]
	if ok {
		statusCode = http.StatusOK
		delete(resources, id)
	}
	w.WriteHeader(statusCode)
}

func runClient() {

}
