package main

import (
	"bytes"
	"reflect"
)

type Xer interface {
	B()
	C()
}

/*
1. 类型有个方法集, 这决定了它是否实现某个接口, 根据接收参数receiver的不同可以分为(T和*T)两种,
*/
type T int

func (T) a()  {} // 必须为大写(导出成员)否则反射无法获取到
func (T) B()  {}
func (*T) C() {}
func (*T) D() {}

func show(i any) {
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumMethod(); i++ {
		println(t.Method(i).Name)
	}
}

// 2. 别名扩展
type G = T    // 别名 G和T就是一个东西
func (*G) E() {} // 为T类型扩展一个新方法, 在T类中也可看到
// 当别名和本来的类型不位于同一个包里面时, 可以定义别名, 但是不能扩展方法
type buf = bytes.Buffer

//func (*buf) C() {}  Cannot define new methods on the non-local type 'bytes.Buffer'

func main() {
	var n T = 1
	var p *T = &n
	show(n)
	println("---------")
	show(p)

	// 3. 直接方法调用的话, 不涉及方法集, 编译器可以自动转换接收参数receiver
	n.B()
	n.C()
	// 就是自己的实例对象可以随便调自己类的任何方法

	// 4. 将T类型变量 通过 类型转换 成接口变量时, 需要进行方法集检查, 否则报错

	// Cannot use 'n' (type T) as the type Xer Type does not implement 'Xer' as the 'C' method has a pointer receiver
	// var x Xer = n

	var x1 Xer = p
	x1.B()
	x1.C()

	println("T{ E }")
	show(struct{ T }{})

	println("T{ *E }")
	show(struct{ *T }{})

	println("*T{ E }")
	show(&struct{ T }{})

	println("*T{ *E }")
	show(&struct{ *T }{})

	var t T
	var g G
	t.a()
	g.a()
	g.E()

}
