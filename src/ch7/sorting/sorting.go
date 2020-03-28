package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

type byArtist []*Track

func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byYear []*Track

func (x byYear) Len() int           { return len(x) }
func (x byYear) Less(i, j int) bool { return x[i].Year < x[j].Year }
func (x byYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

type smartSort struct {
	t     []*Track
	stack []string
}

func (x smartSort) Len() int { return len(x.t) }
func (x smartSort) Less(i, j int) bool {
	iTrack := x.t[i]
	jTrack := x.t[j]
	for i := len(x.stack) - 1; i >= 0; i-- {
		field := x.stack[i]
		switch field {
		case "title":
			if eq := iTrack.Title != jTrack.Title; eq {
				return iTrack.Title < jTrack.Title
			}
		case "year":
			if eq := iTrack.Year != jTrack.Year; eq {
				return iTrack.Year < jTrack.Year
			}
		case "length":
			if eq := iTrack.Length != jTrack.Length; eq {
				return iTrack.Length < jTrack.Length
			}
		case "album":
			if eq := iTrack.Album != jTrack.Album; eq {
				return iTrack.Album < jTrack.Album
			}
		case "artist":
			if eq := iTrack.Artist != jTrack.Artist; eq {
				return iTrack.Artist < jTrack.Artist
			}
		default:
			return false
		}
	}
	return false
}

func getInput() string {
	scan := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, "Sort by> ")
	scan.Scan()
	return scan.Text()
}

func (x smartSort) Swap(i, j int) { x.t[i], x.t[j] = x.t[j], x.t[i] }

func main() {
	// printTracks(tracks)

	var smart = smartSort{tracks, make([]string, 0, 5)}

	for s := getInput(); s != "quit"; s = getInput() {
		smart.stack = append(smart.stack, strings.ToLower(s))
		sort.Sort(smart)
		printTracks(tracks)
	}

	// sort.Sort(customSort{tracks, func(x, y *Track) bool {
	// 	if x.Title != y.Title {
	// 		return x.Title < y.Title
	// 	}

	// 	if x.Year != y.Year {
	// 		return x.Year < y.Year
	// 	}

	// 	if x.Length != y.Length {
	// 		return x.Length < y.Length
	// 	}
	// 	return false
	// }})

	fmt.Println()
	printTracks(tracks)
}
