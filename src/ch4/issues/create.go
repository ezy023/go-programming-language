// Logic for creating new issues on github
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"ch4/github"
)

const (
	baseURL  = "https://api.github.com/repos/ezy023/go-programming-language/issues"
	username = "username"
	token    = "token"
)

type IssueCreate struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
	/* assignee, milestone,  assignees */
}

type IssueCreateResponse struct {
	Id      int
	HTMLURL string `json:"html_url"`
	User    *github.User
}

func main() {
	var i *IssueCreate = createIssue()
	icr := submitIssue(i)
	log.Printf("Issue Created: %v\n", icr)
}

func getEditorPath() (string, error) {
	editor, ok := os.LookupEnv("EDITOR")
	if !ok {
		log.Println("EDITOR environment variable not set, defaulting to vi")
		var err error
		editor, err = exec.LookPath("vi")
		if err != nil {
			log.Fatalln("Could not find default vi editor. Aborting, set 'EDITOR' to desired editor")
		}
	}
	return editor, nil
}

func scanFromStdin(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s ", prompt)
	scanner.Scan()
	return scanner.Text()
}

func editorInput() string {
	editor, err := getEditorPath()
	if err != nil {
		log.Fatalf("Could not find editor to use %v\n", err)
	}
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatalf("Unable to open temp file for editor %v\n", err)
	}

	defer tmp.Close()
	defer os.Remove(tmp.Name())

	cmd := exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	b, err := ioutil.ReadAll(tmp)
	if err != nil {
		log.Fatalf("Could not read contents from file %v %v\n", tmp.Name(), err)
	}
	return string(b)
}

func submitIssue(i *IssueCreate) *IssueCreateResponse {
	body, err := json.Marshal(i)
	if err != nil {
		log.Fatalf("Could not encode request body for issue %v %v\n", i, err)
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	req.SetBasicAuth(username, token)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("POST request to %s failed %v\n", baseURL, err)
	}

	defer resp.Body.Close()

	log.Printf("Response code %d\n", resp.StatusCode)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed reading response Body %v\n", err)
	}
	var icr IssueCreateResponse
	err = json.Unmarshal(b, &icr)
	if err != nil {
		log.Printf("Failed unmarshaling JSON to IssueCreateResponse %v %v\n", string(b), err)
	}

	return &icr
}

func printJson(i *IssueCreate) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Printf("Error marshalling JSON %v\n", err)
	}
	log.Println(string(b))
}

func createIssue() *IssueCreate {
	title := scanFromStdin("Title>")
	body := editorInput()
	labels := strings.Split(scanFromStdin("Labels (as comma separated list)>"), ",")
	return &IssueCreate{Title: title, Body: body, Labels: labels}
}
