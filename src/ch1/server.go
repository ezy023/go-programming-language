// HTTP server - very basic
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Method %s\n", req.Method)
	fmt.Fprintf(w, "Host: %s\n", req.Host)
	fmt.Fprintf(w, "Requested URL: %s\n", req.URL.Path)

	for key, val := range req.Header {
		fmt.Fprintf(w, "Header[%q]: %q\n", key, val)
	}

	if err := req.ParseForm(); err != nil {
		log.Printf("ParseForm Error: %v\n", err)
	}

	for key, val := range req.Form {
		fmt.Fprintf(w, "Form[%s]: %q\n", key, val)
	}
}

func main() {
	http.HandleFunc("/handle", handler)
	http.HandleFunc("/gif", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "Error parsing form %v", err)
		}
		cycles, _ := strconv.ParseFloat(req.Form["cycles"][0], 64)
		lissajous(w, cycles)
	})
	http.ListenAndServe("localhost:8000", nil)

}
