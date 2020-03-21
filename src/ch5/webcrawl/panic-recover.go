package main

import "fmt"

// Exercise 5.19
// Function that returns a value with no return statemnt, using panic and recover.
func panrec(x int) (y int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover() = %v\n", r)
			y = 42
		}
	}()
	y = x * 2
	panic(21)
}

func main() {
	fmt.Println(panrec(4))
}
