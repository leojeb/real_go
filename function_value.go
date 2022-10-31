package main

import "fmt"

type N int

func (n N) copy() {
	fmt.Printf("%p , %v\n", &n, n)
}

func (n *N) ref() {
	fmt.Printf("%p, %v \n", n, *n)
}
func main() {
	var n N = 1
	// 表达式
	var e func(*N) = (*N).ref // type func(*N)
	e(&n)

	// 值传递
	var v func() = n.ref // type func()
	v()

	var p *N
	p.ref()
	(*N)(nil).ref()
	(*N).ref(nil)
}
