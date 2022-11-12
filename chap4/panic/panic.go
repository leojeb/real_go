package main

import (
	"fmt"
	"log"
)

func main() {
	func(v any) {
		fmt.Printf("%T %v", v, v)
		println(v)
	}("wo1")
	defer func() {
		// 拦截数据, r未必是error, 也可能是nil
		if r := recover(); r != nil {
			log.Fatalln(r)
		}
		recover()
	}()

	func() {
		panic(nil)
	}()

	println("panic后就算recover, panic后续的语句也无法执行")
}
