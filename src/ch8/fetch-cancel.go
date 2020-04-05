// Implement variant of fetch that requests several URLs concurrently and cancels the requests when the first response arrives
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	urls := os.Args[1:]
	first := make(chan string)
	token := make(chan struct{}, 10)

	client := http.DefaultClient
	ctx, cancelFunc := context.WithCancel(context.Background())
	for _, url := range urls {
		go func(url string) {
			token <- struct{}{}
			defer func() {
				<-token
			}()

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				log.Println(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			io.Copy(os.Stdout, resp.Body)
			first <- url
		}(url)
	}

	select {
	case url := <-first:
		fmt.Printf("Got %s Cancelling the rest\n", url)
		cancelFunc()
	}
}
