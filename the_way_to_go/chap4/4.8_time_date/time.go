package main

import (
	"fmt"
	"time"
)

/*
time使用
*/
func main() {
	var t = time.Now()
	var t1 = time.Time{} // new(time.Time)
	fmt.Println(time.Now(), t.Day(), t.Minute(), t.Month())
	fmt.Println(time.Now(), t1.Day(), t1.Minute(), t1.Month())
	fmt.Printf("%02d.%02d.%4d\n", t.Day(), t.Month(), t.Year())

	fmt.Println(t.Format("02 Jan 2006 15:04"))
	fmt.Println(t.Format(time.RFC1123))
	fmt.Println(t.Format(time.ANSIC))
	fmt.Println(t.Format("20060102"))
	fmt.Println(t.Format("2006/01/02 15:04"))
}
