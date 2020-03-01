// Issues prints a table of github issues
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"ch4/github"
)

// Group Issues into three buckets, < month, month < issues < year, > year
func groupIssues(issues []*github.Issue) map[string][]*github.Issue {
	grouped := make(map[string][]*github.Issue, 3)
	for _, issue := range issues {
		created := issue.CreatedAt
		switch {
		case created.After(time.Now().AddDate(0, -1, 0)):
			grouped["< 1 month"] = append(grouped["< 1 month"], issue)
		case created.After(time.Now().AddDate(-1, 0, 0)):
			grouped["< 1 year"] = append(grouped["< 1 year"], issue)
		default:
			grouped["> 1 year"] = append(grouped["> 1 year"], issue)
		}
	}
	return grouped
}

func list() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	grouped := groupIssues(result.Items)
	for key, val := range grouped {
		fmt.Println(key)
		for _, item := range val {
			fmt.Printf("\t#%-5d %9.9s %30.30s %30.20v\n",
				item.Number, item.User.Login, item.Title, item.CreatedAt)
		}
	}
}
