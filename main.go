package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const METADATA_URL = "http://metadata.google.internal/computeMetadata/v1/"
const METADATA_HEADERS = "{'Metadata-Flavor': 'Google'}"

func main() {

	// TODO: determine if gorilla/mux is needed or can use http only
	r := mux.NewRouter()
	r.HandleFunc("/", EnvironmentHandler)
	r.HandleFunc("/healthz", HealthHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// EnvironmentHandler returns info about the environment
func EnvironmentHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "{}")
}

// HealthHandler is a health endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
