// Example showing an annonmous function and variable closure
package main

import "fmt"

func squares() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}

func main() {
	var f func() int = squares()
	for i := 0; i < 10; i++ {
		fmt.Printf("Call %d: %d\n", i, f())
	}
}
