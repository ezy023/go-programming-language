// Package github provides a Go API for Github issue tracker
package github

import "time"
import "fmt"

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // In Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func (i *Issue) String() string {
	return fmt.Sprintf("%v %v", i.Number, i.CreatedAt)
}
