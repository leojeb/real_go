package main

import (
	"fmt"
	"os"
	"regexp"
)

/*
slice是一个指针, 指向底层数组, cap是指slice的第一个元素到数组末尾的长度
*/
func main() {
	/*
		new和make的区别
		new() 函数分配内存，make() 函数初始化
		new()创建一个[]int类型变量(切片), 它包含一个指针用来指向底层数组, 此处是nil;还包含len和cap的值.
		new()会返回一个指针这个指针存储着切片的地址(此地址不为nil), 切片里则存储着底层数组地址, 是nil

		make()会直接返回切片变量而不是指向切片的指针
		切片本身也是指针, 所以new([]int)实际上是套娃.

		new只是为当前类型变量分配内存, 不会初始化内部结构, 比如slice里面还有一个底层数组需要初始化, 而new只会初始化len, cap属性, 不会
		初始化底层数组.
		make则会初始化变量内部结构, 譬如初始化slice底层数组.
	*/
	var p1 = new([]int)
	fmt.Printf("%T, %T", p1, *p1)
	println()
	println(*p1, *p1 == nil, p1, p1 == nil, len(*p1), cap(*p1))

	var p2 = make([]int, 0)
	fmt.Printf("%T", p2, p2, len(p2), cap(p2))
	println()
	println(p2)

	//type T struct {
	//	a int
	//}
	//t := new(T)
	//println(t, t.a)
	//var s1 = make([]int, 10)
	s1 := []int{1, 2, 3, 4, 5, 6, 7}
	println(s1[0:2:3])
	fmt.Println(s1[0:1:3])
	//var buffer bytes.Buffer
	//buffer.WriteString()
	//buffer.Len()
	//for {
	//	if s, ok := getNextString(); ok {
	//		buffer.WriteString(s)
	//
	//	} else {
	//		break
	//	}
	//}
	//fmt.Print(buffer.String(), "\n")

}

var digitRegex = regexp.MustCompile("[0-9]+")

func FindAllDigits(filename string) []byte {
	b, _ := os.ReadFile(filename)
	bytes_slice := digitRegex.FindAll(b, -1)
	c := make([]byte, 0)
	for _, bytes := range bytes_slice {
		c = append(c, bytes...)
	}
	return c
}
