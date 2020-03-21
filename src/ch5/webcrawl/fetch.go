package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Please supply a url")
	}
	fetch(os.Args[1])
}

// Exercise 5.18
// Fetch downloads the URL and returns the name and length of the local file
func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	n, err = io.Copy(f, resp.Body)
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()
	// Close file, but prefer error from Copy, if any
	// if closeErr := f.Close(); err == nil {
	// 	err = closeErr
	// }
	return local, n, err
}
