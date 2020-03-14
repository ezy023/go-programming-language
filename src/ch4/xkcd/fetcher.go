// Tool to request xkcd URLs (once) and store the contents to disk
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

var DataDir string = "./comics"

type Comic struct {
	Image         string `json:"img"`
	Transcript    string `json:"transcript"`
	Number        int    `json:"num"`
	AlternateText string `json:"alt"`
	Title         string `json:"title"`
}

func init() {
	err := createDataDir(DataDir)
	if err != nil {
		log.Fatalf("Could not find or create data directory: %v", err)
	}
}

func (c *Comic) String() string {
	return fmt.Sprintf("Title: %s\nTranscript: %s\nImage: %s\n", c.Title, c.Transcript, c.Image)
}

func FetchURL(url string) (*Comic, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch url %s. %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		log.Printf("Request responded with non-2xx response %d.", resp.StatusCode)
	}

	var comic Comic
	err = json.NewDecoder(resp.Body).Decode(&comic)

	if err != nil {
		log.Printf("Could not unmarshal json from response %v", err)
		return nil, err
	}

	return &comic, nil
}

func FetchAll(num int) {
	for i := 1; i < num; i++ {
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)
		comic, err := FetchURL(url)
		if err != nil {
			log.Printf("Could not fetch URL %s. %v\n", url, err)
			continue
		}
		log.Printf("Storing: %s\n", strconv.Itoa(comic.Number))
		err = writeComicToFile(comic)
		if err != nil {
			log.Printf("Error writing comic %d to file %v", comic.Number, err)
			continue
		}
	}
}

func writeComicToFile(c *Comic) error {
	filepath := path.Join(DataDir, strconv.Itoa(c.Number))
	log.Printf("fpath is %s num %s\n", filepath, strconv.Itoa(c.Number))
	f, err := getDataFile(filepath)
	if err != nil {
		return err
	}
	log.Printf("Storing comic at path %s", f.Name())
	return json.NewEncoder(f).Encode(c)
}

func getDataFile(filepath string) (*os.File, error) {
	info, _ := os.Stat(filepath)
	if info != nil {
		return os.Open(filepath)
	}
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func createDataDir(dir string) error {
	info, _ := os.Stat(dir)
	if info != nil && info.IsDir() {
		log.Printf("Data directory %#v is already present.", info.Name())
		return nil
	}

	err := os.Mkdir(dir, 0766)
	if err != nil {
		log.Printf("Could not create directory %s. %v\n", dir, err)
		return err
	}
	return nil
}
