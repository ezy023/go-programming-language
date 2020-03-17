// Exercise 5.2: Count the number of each element on a webpage.
package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func elemCount(m map[string]int, n *html.Node) {
	if n == nil {
		return
	}
	if n.Type == html.ElementNode {
		m[n.Data]++
	}
	elemCount(m, n.FirstChild)
	elemCount(m, n.NextSibling)
}

func outline(m map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		m[n.Data]++
	}

	i := 1
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		fmt.Printf("call %d\n", i)
		i++
		outline(m, c)
	}
}

func main() {
	s := "<erik><p></p></erik><kim><p></p></kim><hank><p></p></hank>"
	root, err := html.Parse(strings.NewReader(s))
	if err != nil {
		fmt.Printf("Unable to parse HTML from stdin %v\n", err)
		os.Exit(1)
	}
	m := make(map[string]int)
	// elemCount(m, root)
	outline(m, root)
	fmt.Println(m)
}
