package main

import (
	"errors"
	"fmt"
)

/*
错误链,
*/
func database() error {
	return errors.New("data")
}
func cache() error {
	if err := database(); err != nil {
		fmt.Printf("类型%T", err)
		println()
		return fmt.Errorf("cache missing: ", err)
	}
	return nil
}

func handle() error {
	return cache()
}

func main() {
	//file, err := os.Open(".\\main.go")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer file.Close()
	//
	//content, err := io.ReadAll(file)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//println(content)

	var e error = ErrOne // 这里要拿指针因为TestError的Error是绑定的指针
	println(e == ErrOne)
	if t, ok := e.(*TestError); ok {
		println(t.x)
	}

	err := test()
	println(err == nil)

	fmt.Println(handle())

	err1 := errors.New("err1")
	err2 := fmt.Errorf("err2:", err1)
	err3 := fmt.Errorf("err3:", err2)
	fmt.Println(err3)
	//
	println(errors.Unwrap(err3) == err2)
	// 递归检查
	println(errors.Is(err3, err1))

	//------
	x := &TestError{1}
	y := fmt.Errorf("y:", x)
	z := fmt.Errorf("z:", y)
	var x2 *TestError
	if errors.As(z, x2) {
		fmt.Println(x == x2)
	}
}

type TestError struct {
	x int
}

func (e *TestError) Error() string {
	return fmt.Sprintf("测试%d", e.x)
}

var ErrOne = &TestError{1}

func test() error {
	var err *TestError
	// 这里err是个默认struct, 只要值为空, 就是nil
	println(err == nil)
	//println(err.x) 报nil pointer reference, 空指针异常
	err1 := error(err) // 返回值这里把TestError转换成了error, error是接口类型, 而接口只有当类型和值都为nil时才等于nil, 所以这里err1不等于nil
	println(err1 == nil)
	return err1
}
