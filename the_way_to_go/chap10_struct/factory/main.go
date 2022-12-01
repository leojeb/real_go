package main

import "fmt"

// 由于go中只有struct, 不像那些面向对象语言一样, 所以需要自己实现简单的构造函数
type File struct {
	Fd   int
	Name string
}

func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	return &File{fd, name}
}

// 如何强制用户使用工厂模式创建新的对象, 将struct的名字改为小写(变为私有)
type privateType struct {
	a int
}

func NewPrivateType(a int) *privateType {
	return &privateType{a}
}

func (pt *privateType) String() string {
	/**
	演示一下栈溢出, 循环调用, 因为Sprintf底层调用的是类型的String方法,
	main函数里调用Sprintf -> 调用String() -> String调用Sprintf, 形成闭环
	*/
	return fmt.Sprintf("此变量值为%v", pt)
}

func main() {
	file1 := NewFile(1, "file")
	print(file1.Name)

	private1 := NewPrivateType(1)
	fmt.Sprintf("%v", private1)
	// make 不能用于底层struct不是map/slice/channel的结构体.
	type A struct {
		a int
	}
	//make(A) cannot make A
	type B map[string]string
	b := make(B)
	fmt.Println(b)
}
