package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CheckPanic(err error) {
	if err != nil {
		//_, file, line, _ := runtime.Caller(1)
		//log.Printf("%s:%d\n", file, line)
		//fmt.Println(err.Error())
		panic(err)
	}

}

func bufioRead() {
	fmt.Println(os.Getwd())
	inputFile, _ := os.Open("1.clip")
	//CheckPanic(inputError)

	defer func() { inputFile.Close() }()

	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readerError := inputReader.ReadString('\n') // inputReader.ReadBytes('\n')
		fmt.Printf("	input is ", inputString)
		if readerError == io.EOF {
			return
		} else {
			CheckPanic(readerError)
		}

	}
}

func readFileLine() {
	inputFile, inputError := os.Open("input.dat")
	CheckPanic(inputError)

	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	for {
		readString, inputError := inputReader.ReadString('\n') // inputReader.ReadBytes('\n')
		fmt.Printf("input was %s: ", readString)

		if inputError == io.EOF {
			return
		}
	}

}

func readWholeFile() {
	inputFile := "a.txt"
	outFile := "o.txt"
	fileBytes, err := os.ReadFile(inputFile)
	CheckPanic(err)
	fmt.Printf("content is %s\n", string(fileBytes))

	err = os.WriteFile(outFile, fileBytes, 0644)
	CheckPanic(err)

}

func bufferReader() {
	inputFile, inputErr := os.Open("input.txt")
	defer inputFile.Close()
	CheckPanic(inputErr)

	inputReader := bufio.NewReader(inputFile)

	bufSize := 1024
	var buf = make([]byte, bufSize)
	for {
		n, readErr := inputReader.Read(buf)
		if n == 0 {
			break
		}
		CheckPanic(readErr)
	}
}

func readByColumn() {
	file, err := os.Open("input.txt")
	CheckPanic(err)
	defer file.Close()

	var col1, col2, col3 []string
	for {
		var c1, c2, c3 string
		_, err := fmt.Fscanln(file, &c1, &c2, &c3)
		CheckPanic(err)
		col1 = append(col1, c1)
		col2 = append(col2, c2)
		col3 = append(col3, c3)

	}
	fmt.Println(col1)
	fmt.Println(col2)
	fmt.Println(col3)

}

func getFileName() {
	path := "/a/b/c"
	filename := filepath.Base(path)
	fmt.Println(filename)
}

func readCompressedFile() {
	/**
	compress 包提供了读取压缩文件的功能，支持的压缩文件格式为：bzip2、flate、gzip、lzw 和 zlib。
	下面的程序展示了如何读取一个 gzip 文件.
	*/

	fName := "input.gz"
	var r *bufio.Reader
	fi, err := os.Open(fName)
	CheckPanic(err)
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		r = bufio.NewReader(fi)
	} else {
		r = bufio.NewReader(fz)
	}

	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Done reading file")
			os.Exit(0)
		} else {
			CheckPanic(err)
		}
		fmt.Println(line)

	}

}

func writeFile() {
	file, err := os.OpenFile("o.txt", os.O_WRONLY|os.O_CREATE, 0777)
	CheckPanic(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	outputString := "hello world"

	for i := 0; i < 10; i++ {
		writer.WriteString(outputString)
	}
	writer.Flush()

}

func otherWriteFile() {
	file, err := os.OpenFile("test", os.O_WRONLY|os.O_CREATE, 0666)
	CheckPanic(err)

	defer file.Close()
	file.WriteString("abcdefg hello ")

}
func main() {
	readCompressedFile()
}
