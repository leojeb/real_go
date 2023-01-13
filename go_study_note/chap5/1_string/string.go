package main

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
	"unsafe"
)

// go 里面string 的 实现
type stringStruct struct {
	str unsafe.Pointer
	len int
}

type stringStructDWARF struct {
	str *byte
	len int
}

// 字符串基础

func main() {
	var s1 = "撒赖咖啡碱撒法"
	var a = []byte(s1)
	var b = []rune(s1)
	println(a[0], b[0], s1[0])

	fmt.Printf("% X ,%d \n", s1, len(s1))
	fmt.Println("ac" > "ab")

	// 按byte遍历
	for i := 0; i < len(s1); i++ {
		fmt.Printf("%d: %X\n", i, s1[i])
	}
	// 按rune遍历
	for i, c := range s1 {
		fmt.Printf("%d: %c\n", i, c)
	}

	//	原始字符串
	println(`println()\\//abc 原始字符串`)
	// +号只能在行末尾, 不能在行首
	println(`asodfi` +
		`asodfiasf`)
	//println(`asdfj`
	//	+ `saodfjas`
	//)
	s3 := `hello, world`
	s2 := s3[:4]

	p1 := (*reflect.StringHeader)(unsafe.Pointer(&s3))
	p2 := (*reflect.StringHeader)(unsafe.Pointer(&s2))
	fmt.Printf("%#v ,\n %#v", p1, p2)
	println()

	// 超过范围按2倍指数增加slice底层数组的空间
	s := `de`
	bs := make([]byte, 0)
	bs = append(bs, "abc"...)
	println(`bs1`, bs)
	bs = append(bs, s...)
	println(`bs`, bs) //
	bs = append(bs, "abcddaosif我的i安师傅鸡丝豆腐扫ID附"...)
	println(`bs`, bs) //
	fmt.Printf("%s \n", bs)

	// 可以覆盖, 已有值, 超过范围不报错, 直接截取.
	buf := make([]byte, 5)
	copy(buf, `abc`)
	copy(buf[2:], `asdfij`)
	println(`buf`, buf)
	fmt.Printf("%s\n", buf)

	// rune, byte, string转换
	var a1 rune = '我'
	var b1 byte = byte(a1)
	var ss string = string(b1)
	var b2 byte = byte(a1)
	println("a1", a1)
	println("b1", b1)
	println("ss", ss)
	println("b2", b2)

	//
	sss := strings.Repeat("a", 1<<10)
	println("sss\n", &sss, "\n", sss)
	// 在从sss转成[]byte类型过程中, 会分配新内存并复制内容
	bss := []byte(sss)
	println("bss \n", &bss, "\n", bss)
	bss[1] = 'B'
	println("bss \n", &bss, "\n", bss)
	fmt.Printf("%s \n", bss)
	// 再将[]byte转换回去也是分配新内存并复制内容
	new_sss := string(bss)
	println("new_sss \n", &new_sss, "\n", new_sss)
	h1 := (*reflect.StringHeader)(unsafe.Pointer(&sss))
	h2 := (*reflect.StringHeader)(unsafe.Pointer(&bss))
	h3 := (*reflect.StringHeader)(unsafe.Pointer(&new_sss))

	fmt.Printf("%#v ,\n %#v, \n %#v", h1, h2, h3)
	println()
	s4 := "刘勇"
	s4c := string(s4[:2] + s4[6:]) //非法拼接
	var s5 []byte = []byte(s4[:])
	fmt.Printf("% X \n", s5)

	println(s4c, len(s4), cap(s5))
	fmt.Printf("%s", s4c)
	fmt.Println(utf8.ValidString(s4))
}
