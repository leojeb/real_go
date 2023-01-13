package main

import (
	. "go_project0/go_study_note/chap1"
	"math"
	. "unsafe"
)

var I int

func main() {
	println("hello world")

	/*
		基本类型 int8/16/32/64, float8/16/32/64, rune, byte
	*/
	println(math.MaxInt8, math.MinInt8)
	var a = 11_222
	println(a)

	var m map[int]string = nil
	println(m)
	var c []int = nil
	println(m, c)

	var d rune = '我'
	println(d)
	var e byte = 'e'
	println("e\n", e)

	/*
		引用类型 slice, map , channel
	*/
	// new只按照类型大小分配内存, 返回一个指针
	var s1 = *new([]int)
	Myprint("s1", s1)
	println(Sizeof(s1))
	// make分配大小并初始化
	var s2 = make([]int, 10, 10)
	Myprint("s2", s2)
	println()

	// slice, map, channel
	var s3 []int          // 三个字段组成的结构体 24
	var m1 map[int]string // 指针 8
	var c1 chan int       // 指针 8

	Myprint("s3", s3)
	println(Sizeof(s3))
	Myprint("m1", m1)
	println(Sizeof(m1))
	Myprint("c1", c1)
	println(Sizeof(c1))

	/*
		自定义类型

	*/
	type B int   //定义新类型
	type C = int // 别名

	var g int = 1
	var f = B(g)
	//var f1 B = g  Cannot use 'g' (type int) as the type B
	var r = C(g)
	Myprint("f", f)
	Myprint("r", r)
}
