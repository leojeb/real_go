package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

var where = func() {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("%s:%d", file, line)
}

func main() {
	// 计算函数执行时间
	start := time.Now()
	var f = Adder()
	fmt.Println(f(1))
	fmt.Println(f(20))
	fmt.Println(f(300))
	/*
		不管外部函数是否退出, 闭包函数都能够继续操作外部函数中的内部变量.
	*/

	// 使用闭包调试, 可以通过runtime包, 获取代码执行的文件路径和行数等信息

	where()
	println("哔哔哔")
	where()

	println("哔哔哔")
	f1()
	end := time.Now()
	deltaT := end.Sub(start)
	println(deltaT)
}

var where1 = log.Print

func f1() {
	println("f1函数里面")
	where1("aaa")
}

func Add1() func(b int) int {
	return func(b int) int {
		return b + 2
	}
}

func Add2(a int) func(b int) int {
	return func(b int) int {
		return a + b
	}
}

func Adder() func(int) int {
	where()
	var x int
	return func(i int) int {
		where()
		x += i
		return x
	}
}
