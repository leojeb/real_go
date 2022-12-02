package main

import (
	"errors"
	"fmt"
)

// 1. 自定义错误结构体, 实现Error()方法
type PathError struct {
	Op   string // "operations", e.g. "open"
	Path string
	Err  error // 可有可无, 包装已有错误时可以加上
}

func (e *PathError) Error() string {
	return e.Op + " " + e.Path + ": " + e.Err.Error()
}

func main() {
	// 2. new 一个
	errNotFound := errors.New("error not found")
	fmt.Println(errNotFound)
}

// 3. 抽象出共用错误接口
type Error interface {
	Timeout() bool   // Is the error a timeout?
	Temporary() bool // Is the error temporary?
}

// 4. 用fmt.Errorf()创建error对象
func fmtCreateError() (int, error) {
	f := -1
	if f < 0 {
		return 0, fmt.Errorf("math: square root of negative number %g", f)
	}
	return 0, nil
}
