package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func CheckError(err error) {
	fmt.Println(err.Error())
	panic(err)
}

func main() {
	fmt.Println(os.Getwd())
	inputFile, _ := os.Open("clip1.json")
	//CheckError(inputError)

	defer func() { inputFile.Close() }()

	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readerError := inputReader.ReadString('\n')
		fmt.Printf("	input is ", inputString)
		if readerError == io.EOF {
			return
		}
	}
}
