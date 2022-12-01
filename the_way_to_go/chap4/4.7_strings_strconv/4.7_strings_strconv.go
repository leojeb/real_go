package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	example := "This is leon 刘勇"
	fmt.Printf("%t \n", strings.HasPrefix(example, "Th"))
	fmt.Printf("%t ", strings.HasPrefix(example, "Th"))
	log.Printf("%t ", strings.HasPrefix(example, "Th"))
	log.Printf("%t ", strings.HasPrefix(example, "Th"))

	fmt.Println("-------------")

	fmt.Println("Contains: ", strings.Contains(example, "abc"))
	fmt.Println("Index: ", strings.Index(example, "i"))
	fmt.Println("LastIndex: ", strings.LastIndex(example, "刘勇"))
	fmt.Println("IndexRune: ", strings.IndexRune(example, '刘'))
	fmt.Println("Join: ", strings.Join([]string{example, example}, "//"))
	fmt.Println("ReplaceAll: ", strings.ReplaceAll(example, "is", "was"))
	fmt.Println("Count: ", strings.Count(example, "is"))
	fmt.Println("Repeat: ", strings.Repeat(example, 2))
	fmt.Println("ToUpper: ", strings.ToUpper(example))
	fmt.Println("Trim: ", strings.Trim(example, "This"))
	fmt.Println("TrimSpace: ", strings.TrimSpace(example))
	fmt.Println("Fields: ", strings.Fields(example), len(strings.Fields(example)))
	fmt.Println("Split: ", strings.Split(example, "/"), len(strings.Split(example, "/")))

	//	read from string
	buf := make([]byte, 10)
	strReader := strings.NewReader(example)
	_, err := strReader.Read(buf)
	if err != nil {
		log.Println(err.Error())
	}

}
