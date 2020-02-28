// Insert commas into string numbers
package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func commasRec(num string) string {
	if len(num) <= 3 {
		return num
	}

	return commasRec(num[:len(num)-3]) + "," + num[len(num)-3:len(num)]
}

func commasIter(num string) string {
	if len(num) <= 3 {
		return num
	}

	var buf bytes.Buffer
	mod := len(num) % 3
	if mod > 0 {
		fmt.Fprintf(&buf, "%s,", num[:mod])
	}
	var i, j int
	for i, j = mod, mod+3; j < len(num); i, j = j, j+3 {
		fmt.Fprintf(&buf, "%s,", num[i:j])
	}
	fmt.Fprintf(&buf, "%s", num[i:j])
	return buf.String()
}

func commas(num string) string {
	var sign bool
	if num[0] == '-' || num[0] == '+' {
		sign = true
	}

	d := strings.Index(".", num)
	if d >= 0 {
		if sign {
			return string(num[0]) + commasRec(num[1:d]) + num[d:]
		} else {
			return commasRec(num[0:d]) + num[d:]
		}
	}
	if sign {
		return string(num[0]) + commasRec(num[1:])
	} else {
		return commasRec(num)
	}

}

func main() {
	var i string
	if len(os.Args) < 2 {
		i = "123456789"
	} else {
		tmp := os.Args[1]
		_, err := strconv.ParseFloat(tmp, 64)
		if err != nil {
			fmt.Println("Please enter a valid number: ", err)
			os.Exit(1)
		}
		i = tmp
	}
	fmt.Println(commas(i))

}
