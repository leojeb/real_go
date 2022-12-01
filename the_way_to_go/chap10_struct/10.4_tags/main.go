package main

import (
	"fmt"
	. "go_project0/the_way_to_go/chap10_struct/factory"
	"reflect"
)

var a = File{3, "测试导包"}

type Tags struct {
	a bool "a的标签"
	b int  "b的标签"
}

type Extends struct {
	Tags "继承Tags"
	b    int
}

func refTag(tt Tags, ix int) {
	ttType := reflect.TypeOf(tt)
	ixField := ttType.Field(ix)
	fmt.Printf("%v \n", ixField.Tag)
}

type A int
type B int

func (i *A) name() {
	fmt.Println(1)
}

func (i *B) name() {
	fmt.Println(1)
}

type C struct {
	A
	B
}

var c = C{}

func main() {
	t1 := Tags{true, 1}
	for i := 0; i < 2; i++ {
		refTag(t1, i)
	}

	e1 := Extends{Tags{false, 2}, 1}
	fmt.Println(e1.b, e1.Tags.b, e1.b)

	c.A.name()
	c.B.name()
	//c.name() ambiguous
}
