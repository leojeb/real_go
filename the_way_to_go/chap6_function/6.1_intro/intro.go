package main

import "fmt"

type binOp func(int, int) int

func main() {
	/*
		go不支持函数重载
	*/
	//函数可以生命使用

	var f1 binOp = func(i int, i2 int) int {
		fmt.Println(i + i2)
		return 1
	}
	f1(1, 2)

	//add := binOp

}
