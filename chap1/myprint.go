package chap1

import (
	"fmt"
)

type A interface{}

func Myprint(var_name string, e A) {
	print(var_name)
	fmt.Printf("的值是: %v, 类型是: %T", e, e)
	println()
}
