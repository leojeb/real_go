package main

import "fmt"

func main() {
	fmt.Println("start call")
	Call()
	fmt.Println("end call")
}

func badCall() {
	panic("bad end")
}

func Call() {
	defer func() {
		if err := recover(); err != nil {
			println("Panicking: reach bad end")
		}
	}()
	badCall()
	fmt.Println("After bad call ") // cannot reach
}
