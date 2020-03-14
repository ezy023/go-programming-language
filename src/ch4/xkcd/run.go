// Go program to run the various pieces of xkcd

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"text/template"
)

var currentComic string = "https://xkcd.com/info.0.json"
var indexStorage string = "comic-index.idx"

func run() {
	_, err := FetchURL(currentComic)
	if err != nil {
		log.Fatalln("Could not fetch current comic")
	}
	// log.Printf("About to request %d comics\n", current.Number)
	// FetchAll(current.Number) // Fetch all comics up to the current comic

	comicIndex, err := getComicIndex()
	if err != nil {
		log.Fatalln("Could not create ComicIndex")
	}
	comicIndex.IndexDir(DataDir)
	comicIndex.PersistIndex()
	log.Printf("ComicIndex %d\n", len(comicIndex.Index))

	term := scanFromStdin()
	comicNums := comicIndex.SearchIndex(term)
	comics := make([]*Comic, 0, len(comicNums))
	for k := range comicNums {
		c, err := LoadComic(k)
		if err != nil {
			log.Printf("Unable to load comic %d. %v", k, err)
			continue
		}
		comics = append(comics, c)
	}
	tempString := `
{{$length := len .}}
{{if gt $length 0}}
{{range . }}
   Comic Number: {{.Number}}
   Title: {{.Title}}
   Transcript: {{.Transcript}}
{{end}}
{{else}}
No Comics match that Query
{{end}}`
	temp := template.Must(template.New("cli").Parse(tempString))
	temp.Execute(os.Stdout, comics)
	// for _, comic := range comics {
	// 	temp.Execute(os.Stdout, comic)
	// }
}

func getComicIndex() (*ComicIndex, error) {
	var index *ComicIndex
	var err error
	if _, err = os.Stat(indexStorage); os.IsNotExist(err) {
		index, err = NewComicIndex(indexStorage)

	} else {
		index, err = LoadComicIndex(indexStorage)

	}
	if err != nil {
		return nil, err
	}
	return index, nil
}

func load() (*ComicIndex, error) {
	comicIndex, err := LoadComicIndex(indexStorage)
	if err != nil {
		log.Fatalf("Could not load index %v", err)
		return nil, err
	}
	return comicIndex, nil
}

func scanFromStdin() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	fmt.Printf("Search> ")
	scanner.Scan()
	return scanner.Text()
}

func main() {
	run()
}
