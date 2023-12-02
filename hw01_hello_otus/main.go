package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	helloString := "Hello, OTUS!"
	fmt.Println(reverse.String(helloString))
}
