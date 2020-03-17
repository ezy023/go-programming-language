// Exercise 5.3: print all text elements on page, not <script> or <style> elements
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func printText(n *html.Node) {
	if n == nil {
		return
	}
	if n.Type == html.ElementNode && n.Data == "script" || n.Data == "style" {
		return
	}
	if n.Type == html.TextNode {
		fmt.Println(n.Data)
	}
	printText(n.FirstChild)
	printText(n.NextSibling)
}

func main() {
	root, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Printf("Unable to parse HTML %v\n", err)
		os.Exit(1)
	}
	printText(root)
}
