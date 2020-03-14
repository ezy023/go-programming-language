// Retrieve comics from the storage device (filesystem in this case)
package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"strconv"
)

func LoadComic(num int) (*Comic, error) {
	log.Printf("Loading comic: %d\n", num)
	file, err := os.Open(path.Join(DataDir, strconv.Itoa(num)))
	if err != nil {
		return nil, err
	}
	var comic Comic
	err = json.NewDecoder(file).Decode(&comic)
	if err != nil {
		return nil, err
	}
	return &comic, nil
}
