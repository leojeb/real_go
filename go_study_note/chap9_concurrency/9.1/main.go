package main

import "time"

var c int

func inc() int {
	c++
	return c
}

func main() {
	//	1. 参数立即计算并复制
	var a = 100
	go func(a int, b int) {
		println(a, b)
		time.Sleep(time.Second)
	}(a, inc())

	a += 100
	println("main:", a, inc())

	time.Sleep(time.Second * 3)
}
