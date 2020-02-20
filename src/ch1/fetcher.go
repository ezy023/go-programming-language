// HTTP request spammer
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func fetch(url string, ch chan<- string) {
	_, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s %v\n", url, err)
	}
	ch <- "done"
}

func main() {
	var num int
	if len(os.Args) <= 1 {
		num = 100
	} else {
		var err error
		num, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Error converting num %v\n", err)
			os.Exit(1)
		}
	}

	url := "http://localhost:8000/handle"
	var ch chan string = make(chan string)
	for i := 0; i < num; i++ {
		go fetch(url, ch)
	}

	for i := 0; i < num; i++ {
		fmt.Printf("%s\n", <-ch)
	}

	resp, err := http.Get("http://localhost:8000/counter")
	if err != nil {
		fmt.Printf("Error checking counter %v\n", err)
	}
	fmt.Println("Counter Response: ", resp.Status)
	io.Copy(os.Stdout, resp.Body)

}
