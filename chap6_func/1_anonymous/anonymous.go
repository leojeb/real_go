package main

import "sync"

// 1.
type data struct {
	// 匿名嵌入, 相当于继承了sync.Mutex类
	sync.Mutex
	buf [1024]byte
}

// 2.
type E struct {
	e int
}
type F struct {
	E // 继承E里所有的方法和变量
}
type G struct {
	E
	e string
}

func (E) toString() string { return "E" }
func (f F) toString() string {
	println(f.e)
	return "F" + f.E.toString() + string(f.e)
}
func (G) toString(s string) string {
	// 函数同名但是签名(输入输出数量和类型)
	return "G"
}

func main() {

	// 1. 直接调用父类Mutex里的Lock方法
	d := data{}
	d.Lock()
	defer d.Unlock()

	// 2. 子类继承父类, 方法重写(方法完全相同)调子类方法
	var f F
	println(f.toString())
	println(f.E.toString())

	//3.  方法重载(同名但不同签名), 也是优先调子类方法
	var g G
	//println(g.toString()) 如果不填参数会报错
	println(g.toString("111"))
	println(g.E.toString())

	// 4. 父类子类变量重名时, 父类只能访问自己的方法, 子类都可以访问
	println(g.e, g.E.e)
	var e1 E
	println(e1.e)
}
