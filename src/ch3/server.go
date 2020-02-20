// Simple HTTP server to serve the surface image
package main

import (
	"net/http"
)

func surface(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	render(w)
}

func main() {
	http.HandleFunc("/", surface)
	http.ListenAndServe("localhost:8000", nil)
}
