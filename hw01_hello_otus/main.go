package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const str = "Hello, OTUS!"

func main() {
	fmt.Println(stringutil.Reverse(str))
}
