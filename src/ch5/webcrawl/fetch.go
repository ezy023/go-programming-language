package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("HTTP request failed %v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to read all content from response %v\n", err)
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, string(b))
}
