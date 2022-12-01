package main

import (
	"fmt"
	"regexp"
)

func main() {
	searchIn := "John: 2578.34 William: 4567.23 Steve: 5632.18"
	ok, _ := regexp.Match("[0-9]+.[0-9]", []byte(searchIn))
	ok, _ = regexp.MatchString("[0-9]+.[0-9]", searchIn)
	fmt.Println(ok)
}
