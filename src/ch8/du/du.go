package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Exercise 8.9 Write a version of du that computes and periodically dispalys separate totals for each of the root directories

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}

	}
}

var sema = make(chan struct{}, 20)

// dirents returns the entries of directory dir
func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token, if you can send to this channel that means there is availability to do work (the channel buffer isn't full), otherwise you'll block until another goroutine receives from the channel and opens up a slot in the buffer
	case <-done:
		return nil
	}
	defer func() { <-sema }() // release token
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries

}

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse file tree
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Cancel traversal when input is detected
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(done)
	}()

	// Print the results periodically
	var tick <-chan time.Time
	// flag.Bool() returns a pointer to a bool so we have to dereference here
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64

loop:
	for {
		select {
		case <-done:
			// Drain fileSizes to allow existing goroutines to finish
			for range fileSizes {
				// Do nothing
			}
			return
		case size, ok := <-fileSizes:
			if !ok {
				break loop
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes/1e9))
}
