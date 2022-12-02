package main

import (
	"fmt"
	"strconv"
	"strings"
)

/*
*
这是所有自定义包实现者应该遵守的最佳实践：

1）在包内部，总是应该从 panic 中 recover：不允许显式的超出包范围的 panic()

2）向包的调用者返回错误值（而不是 panic）。
*/
func main() {
	numbers, err := Parse("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(numbers)
}

type ParseError struct {
	Index int
	Word  string
	Err   error
}

func (e *ParseError) String() string {
	return fmt.Sprintf("pkg parse: error parsing %q as int", e.Word)
}

func Parse(input string) (numbers []int, err error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("", r)

			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
		}
	}()

	fields := strings.Fields(input)
	numbers = fields2numbers(fields)

	return
}

// panic 接收任意参数, 并且会在recover时将其返回, recover截取做操作

func fields2numbers(fields []string) (numbers []int) {

	if len(fields) == 0 {
		panic(" no word to parse")
	}

	for idx, field := range fields {
		num, err := strconv.Atoi(field)
		if err != nil {
			panic(&ParseError{idx, field, err})
		}
		numbers = append(numbers, num)
	}
	return

}
