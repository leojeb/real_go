package main

import "fmt"

func main() {
	var s1 = "撒赖咖啡碱撒法"
	var a = []byte(s1)
	var b = []rune(s1)
	println(a[0], b[0], s1[0])
	fmt.Printf("% X ,%d \n", s1, len(s1))
	fmt.Println("ac" > "ab")

	// 按byte遍历
	for i := 0; i < len(s1); i++ {
		fmt.Printf("%d: %X\n", i, s1[i])

	}
	// 按rune遍历
	for i, c := range s1 {
		fmt.Printf("%d: %c\n", i, c)
	}
}
