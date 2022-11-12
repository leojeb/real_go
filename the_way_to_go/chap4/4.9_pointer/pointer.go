package main

import "fmt"

func main() {
	// 不能获取字面量或常量的地址, 如
	const i = 5
	//a := &i
	//a1 := &10
	/*
		go 语言不允许指针运算, 所以pointer++是不允许的, 会导致内存泄漏
		对空指针的反向引用是不合法的, 如:
	*/
	//var p *int = nil
	//*p = 0 //试图给一个空指针指向的内容赋值, 试图反向引用来赋值是不允许的

	for pos, c := range "wode好吗" {
		fmt.Print(pos, c)
		fmt.Printf("\t%c", c)
		println()
	}
	for pos, rune1 := range "wode号码" {
		fmt.Print(pos, rune1)
		fmt.Printf("\t%c", rune1)
		println()
	}
}
