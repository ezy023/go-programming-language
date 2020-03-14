// This file holds the code to build and read a basic index of words to comics
package main

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type ComicIndex struct {
	Index map[string]map[int]struct{}
	file  *os.File
}

func NewComicIndex(filepath string) (*ComicIndex, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Could not open or create file %s %v", filepath, err)
	}
	return &ComicIndex{make(map[string]map[int]struct{}), file}, nil
}

func LoadComicIndex(filepath string) (*ComicIndex, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return NewComicIndex(filepath)
	}
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Could not load index file %v", err)
	}
	var cin ComicIndex
	err = gob.NewDecoder(file).Decode(&cin.Index)
	if err != nil {
		log.Fatalf("Unable to decode index from file %s %v", filepath, err)
	}
	return &cin, nil
}

func (c *ComicIndex) Add(word string, comicNum int) {
	entry, present := c.Index[word]
	if !present {
		entry = make(map[int]struct{})
	}
	entry[comicNum] = struct{}{}
	c.Index[word] = entry
}

func (c *ComicIndex) IndexDir(dataDir string) error {
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		log.Fatalln("Cannot open data directory %v", err)
		return err
	}
	for _, finfo := range files {
		if finfo.IsDir() {
			continue
		}
		file, err := os.Open(path.Join(dataDir, finfo.Name()))
		if err != nil {
			log.Printf("Index failed to open file %s %v", finfo.Name(), err)
			file.Close()
			continue
		}
		comic, err := readFromFile(file)
		if err != nil {
			log.Printf("Index could not decode comic from stored json %s %v", file.Name(), err)
			file.Close()
			continue
		}
		file.Close()
		addToIndex(c, comic)
	}
	return nil
}

func (c *ComicIndex) PersistIndex() error {
	err := gob.NewEncoder(c.file).Encode(c.Index)
	if err != nil {
		log.Printf("Unable to encode comic index %v", err)
		return err
	}
	return nil
}

func (c *ComicIndex) SearchIndex(term string) map[int]struct{} {
	return c.Index[term]
}

func readFromFile(f *os.File) (*Comic, error) {
	var comic Comic
	err := json.NewDecoder(f).Decode(&comic)
	if err != nil {
		return nil, err
	}
	return &comic, nil
}

func addToIndex(cidx *ComicIndex, comic *Comic) {
	for _, word := range strings.Split(comic.Transcript, " ") {
		cidx.Add(stripChars(word), comic.Number)
	}
	for _, word := range strings.Split(comic.Title, " ") {
		cidx.Add(stripChars(word), comic.Number)
	}
}

func stripChars(word string) string {
	return strings.Trim(word, "\"{}[]:;.\n\t, ")
}
