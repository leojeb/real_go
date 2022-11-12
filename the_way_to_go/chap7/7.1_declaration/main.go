package main

import "fmt"

func main() {
	a := [...]string{"a", "b", "c", "d"}
	for i := range a {
		fmt.Println("Array item", i, "is", a[i])
	}

	var arr1 = new([5]int) //指针
	var arr2 [5]int        // 值类型
	fmt.Printf("arr1类型%T, arr2类型%T", arr1, arr2)
	// arr2是值, arr1是指针, 将arr1的内容赋值给arr2, 做了一次内存拷贝. 改变arr2, 不改变arr1
	arr2 = *arr1
	arr2[1] = 3
	fmt.Println(*arr1, arr2)

	var arr3 = [5]string{3: "aaa", 4: "asfa"}
	println(arr3)
}
