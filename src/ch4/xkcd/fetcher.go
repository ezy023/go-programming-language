// Tool to request xkcd URLs (once) and store the contents to disk
// package fetcher
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Comic struct {
	Image         string `json:"img"`
	Transcript    string
	Number        int    `json:"num"`
	AlternateText string `json:"alt"`
	Title         string
}

func (c *Comic) String() string {
	return fmt.Sprintf("Title: %s\nTranscript: %s\nImage: %s\n", c.Title, c.Transcript, c.Image)
}

// This will return an struct of the unmarshaled JSON
func fetchURL(url string) (*Comic, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch url %s. %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Instead of reading into a string we'll decode using the streaming decoder
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not read response from url %s. %v", url, err)
		return nil, err
	}
	var comic Comic
	err = json.Unmarshal(b, &comic)
	if resp.StatusCode > 299 {
		log.Printf("Request responded with non-2xx response %d.", resp.StatusCode)
	}
	if err != nil {
		log.Printf("Could not unmarshal json from response %s %v", string(b), err)
		return nil, err
	}

	return &comic, nil
}

func iterURL(num int) {
	for i := 1; i < num; i++ {
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)
		comic, err := fetchURL(url)
		if err != nil {
			log.Printf("Could not fetch URL %s. %v\n", url, err)
			continue
		}
		log.Printf("%d: %s\n", i, comic.Title)
	}
}

func main() {
	// comic, err := fetchURL(url)
	// if err != nil {
	// 	log.Printf("Could not fetch comic from url %s. %v", url, err)
	// }
	iterURL(100)
}
