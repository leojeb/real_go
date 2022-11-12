package main

import "time"

/*
	内存缓存, 提高效率
	fibonacci和fibonacci1的时间占比几乎是1:7. 因为前者使用了内存缓存， 而后者计算了很多重复结果
*/
const LIM = 41

var fibs [LIM]uint64

func main() {
	start := time.Now()
	res := fibonacci(30)
	println(res)
	end := time.Now()
	delta := end.Sub(start)
	println("1耗时:", delta.Nanoseconds())

	start = time.Now()
	res = fibonacci1(30)
	println(res)
	end = time.Now()
	delta = end.Sub(start)
	println("2耗时:", delta.Nanoseconds())
}

func fibonacci(i int) uint64 {
	if fibs[i] != 0 {
		return fibs[i]
	} else if i <= 1 {
		fibs[i] = 1
		return fibs[i]
	} else {
		fibs[i] = fibonacci(i-1) + fibonacci(i-2)
		return fibs[i]
	}
}

func fibonacci1(i int) uint64 {
	if i <= 1 {
		return 1
	} else {
		return fibonacci1(i-1) + fibonacci1(i-2)
	}
}
