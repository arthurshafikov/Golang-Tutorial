package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	s := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(s)
}
