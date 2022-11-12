package main

import "fmt"

/*
Shaper 通过go的接口的特性, 可以很方便的实现多态, 增加代码的扩展性
*/
type Shaper interface {
	Area() float64 // 实现此接口的实例的函数签名必须得一模一样
}

// 接口嵌套接口
type Nesting_interface interface {
	Shaper
	what() float64
}

type Square struct {
	side float64
}

func (s *Square) Area() float64 {
	return s.side * s.side
}

type Rectangle struct {
	width  float64
	length float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.length
}

func main() {
	// 方法接收参数是 *Square, 所以*Square类实现了接口, 而Square没有
	var s = &Square{4}
	var a1 = Shaper(s)
	var r = Rectangle{1, 2}
	var a2 = Shaper(r)
	a1.Area()
	a2.Area()

	// 实际上可以直接赋值, 无需显示转换
	a0 := Shaper(s)
	var a3 Shaper = r
	a0.Area()
	a3.Area()
	shapers := []Shaper{s, r}
	for _, shaper := range shapers {
		fmt.Println(shaper.Area())
	}

	// 接口断言
	if v, ok := a1.(*Square); ok {
		//如果实现链包含Square, 做一些操作
		// Process(v)
		fmt.Println(v)
		fmt.Printf("%T", v)
	}

	if _, ok := a2.(*Square); ok {
		println("a2是*Square类型")
	}
	if _, ok := a2.(Rectangle); ok {
		println("a2是Rectangle类型")
	}

	TypeSwitch(a1, a2)

}

// 类型判断type-switch
func TypeSwitch(items ...any) {
	for i, item := range items {
		switch item.(type) {
		case Rectangle:
			fmt.Printf("第%d个的类型是Rectangle", i)
			println()
		case *Square:
			fmt.Printf("第%d个的类型是*Square", i)
			println()
		}

	}
}
