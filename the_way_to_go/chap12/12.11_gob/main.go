package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

func main() {
	EncDec()
	EncodeToFile()
}

func EncDec() {
	// 指代network, 通常enc和dec都是绑定network的, 而且在不同程序中执行
	var network bytes.Buffer
	// var network *bytes.Buffer
	// 不能写成这样, 因为这样会初始化一个指针指向bytes.Buffer对象, 但是并不会初始化Buffer对象, 导致这个指针是一个nil pointer
	// 而var network bytes.Buffer  会实实在在的初始化一个Buffer对象, 再用&network指向它, 此时&network不是空指针
	enc := gob.NewEncoder(&network)
	dec := gob.NewDecoder(&network)

	err := enc.Encode(P{3, 4, 5, "Pythagoras"})
	if err != nil {
		log.Fatal("encode error ", err)
	}
	// 用一个对象来接收解码出来的内容
	var q Q
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error ", err)
	}

	fmt.Printf("%q: {%d,%d}\n", q.Name, q.X, q.Y)

	var c *Q
	var e *Q
	var d Q

	var a int32 = 1
	var b int32 = 1
	c = &Q{&a, &b, "a"}

	fmt.Printf("c=%v\n", c)
	fmt.Printf("*c=%v\n", *c)
	fmt.Printf("e=%v\n", e)
	fmt.Printf("d=%v", d)
	//fmt.Printf("*c=%v",*c) 无法取值, 空指针
	//fmt.Printf("*d=%v",*d) 无法取值, d不是指针
}

type Address struct {
	Type    string
	City    string
	Country string
}

type VCard struct {
	FirstName string
	LastName  string
	Addresses []*Address
	Remark    string
}

func EncodeToFile() {
	pa := &Address{"private", "Aartselaar", "Belgium"}
	wa := &Address{"work", "Boom", "Belgium"}
	vc := VCard{"Jan", "Kersschot", []*Address{pa, wa}, "none"}

	file, _ := os.OpenFile("encode.txt", os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	enc := gob.NewEncoder(file)
	err := enc.Encode(vc)
	if err != nil {
		log.Fatal("encode error ", err)
	}
}
