// List any packages that transitively depend on the supplied args
package main

import (
	_ "bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// map of pkg to a list of packages that import it
var index map[string][]string = make(map[string][]string)

type Info struct {
	ImportPath string
	Deps       []string
	Imports    []string
}

func (l Info) String() string {
	return fmt.Sprintf("ImportPath: %s\nDeps: %v\nImports: %v", l.ImportPath, l.Deps, l.Imports)
}

func listPackages(info chan<- Info) {
	cmdArgs := append([]string{"list", "-json", "..."})
	cmd := exec.Command("go", cmdArgs...)
	cmdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	decoder := json.NewDecoder(cmdout)
	for decoder.More() {
		var i Info
		if err := decoder.Decode(&i); err != nil {
			log.Fatalln(err)
		}
		info <- i
	}
	close(info)

	if err := cmd.Wait(); err != nil {
		log.Fatalln(err)
	}
}

func buildIndex(pkginfo <-chan Info, done chan<- struct{}) {
	for info := range pkginfo {
		for _, d := range info.Deps {
			index[d] = append(index[d], info.ImportPath)
		}

	}
	done <- struct{}{}
}

func main() {
	pkgs := make(chan Info)
	done := make(chan struct{})
	go listPackages(pkgs)
	go buildIndex(pkgs, done)
	<-done
	for _, a := range os.Args[1:] {
		fmt.Printf("%s: %s\n", a, index[a])
	}
}
