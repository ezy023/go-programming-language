package main

import (
	"fmt"
	"log"
	"os"

	"ch5/links"
)

// buffered channel modeling a counting semaphore, vacant slots means a goroutine can do work
var tokens = make(chan struct{}, 20)

// Exercise 8.6: Add 'depth' specification to crawl so URLs only reachable within depth are crawled
// Exercise 8.7: Write concurrent web crawler that craws a website and stores it locally
func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire token
	list, err := links.Extract(url)
	<-tokens // release token
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)
	var n int // number of pending sends to worklist

	// Start with command line args
	n++
	go func() { worklist <- os.Args[1:] }()

	// Crawl web concurrently
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}
