package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"ch5/links"
)

const pageStoreRoot = "./page-store/"

// breadFirst calls f for each item in the worklist
// Any items retruend by f are added to the worklist
// f is called at most once for each item
func breadthFirst(f func(item, origHost string) []string, origHost string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item, origHost)...)
			}
		}
	}
}

func crawl(link, origHost string) []string {
	linkURL, err := url.Parse(link)
	if err != nil {
		log.Print(err)
	}
	if origHost == linkURL.Host {
		storePage(linkURL)
	} else {
		fmt.Println(link)
	}

	list, err := links.Extract(link)
	if err != nil {
		log.Print(err)
	}
	return list
}

func genFileName(base string) string {
	if strings.HasSuffix(base, ".html") {
		return base
	}
	return strings.Join([]string{base, ".html"}, "")
}

func storePage(u *url.URL) {
	urlPath := strings.TrimRight(u.Path, "/")
	dirLoc := path.Join(pageStoreRoot, path.Dir(urlPath))
	filename := genFileName(path.Base(urlPath))
	err := os.MkdirAll(dirLoc, 0755)
	if err != nil && !os.IsExist(err) {
		log.Printf("Failed to create directory %s %v", u.String(), err)
	}

	fullpath := path.Join(dirLoc, filename)
	file, err := os.OpenFile(fullpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("Failed to open file %s %v", u.String(), err)
	}
	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("Failed to make request %s %v", u.String(), err)
	}

	defer resp.Body.Close()

	fmt.Printf("Storing page %s at %s\n", u.String(), fullpath)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Failed to store %s %v", u.String(), err)
	}

}

func main() {
	// Crawl the web breadth-first
	// starting from the command line arguments
	initLink, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	origHost := initLink.Host
	breadthFirst(crawl, origHost, os.Args[1:])
}
