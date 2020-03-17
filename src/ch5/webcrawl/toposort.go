// prereqs maps computer science courses to their prerequisites
package main

import (
	"fmt"
	"sort"
)

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"linear algebra":        {"math"},
	"math":                  {"trigonometry"},
	"trigonometry":          {"calculus"},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func contains(s []string, other string) bool {
	for _, val := range s {
		if val == other {
			return true
		}
	}
	return false
}

func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items, stack []string)

	visitAll = func(items, stack []string) {
		if stack == nil {
			stack = make([]string, 0)
		}
		for _, item := range items {
			if contains(stack, item) {
				fmt.Printf("Cycle Detected in graph: %s %#v\n", item, stack)
				continue
			}
			if !seen[item] {
				seen[item] = true
				stack = append(stack, item)
				visitAll(m[item], stack)
				order = append(order, item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys, nil)
	return order
}
