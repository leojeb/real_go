package main

import (
	"context"
	"runtime"
	"sync"
)

func task1() {
	go func() {
		defer println("g done")
		//time.Sleep(time.Second)
	}()

	defer println("main done")
	// 主线程任务结束, 不会等待并发任务
	//time.Sleep(time.Second)
}

// 1. 等待
func task2() {
	chan1 := make(chan struct{})
	go func() {
		defer close(chan1)
		println("done")
	}()
	<-chan1 // 如果不关闭通道chan1, 会一直阻塞住
}

func task3() {
	//
	wg := sync.WaitGroup{}
	// 添加等待任务一定要在创建并发任务和等待语句之前
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			println(i, " done")
		}(i)
	}
	wg.Wait()

	var loop_num = 20
	for i := 10; i < loop_num+10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			println(i, " done")
		}(i)
	}
	wg.Wait()
}

func task4() {
	// 上下文实现等待, 和通道实现基本一致
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		println("ctx finished")
	}()
	<-ctx.Done()
}

func task5() {
	// 锁实现同步, go更倾向于用通信(channel)实现
	var lock sync.Mutex
	lock.Lock()

	go func() {
		defer lock.Unlock()
		println("lock finished")
	}()

	lock.Lock()
	lock.Unlock()
	println("exit")
}

// 2. 终止
func task6() {
	chan1 := make(chan struct{})
	go func() {
		defer close(chan1)
		defer println("done.")

		a()
		b()
		c()
	}()
	<-chan1
}

func a() { println("a") }
func b() { println("b"); runtime.Goexit(); println("b 2") }
func c() { println("c") }

func main() {
	//task1()
	task2()
	//task3()
	//task4()
	//task5()
	//task6()
	//var q = make(chan struct{})
	//go func() {
	//	defer close(q)
	//	defer println("done")
	//	time.Sleep(time.Second)
	//	println("这行会执行")
	//}()
	//<-q
	//// 这里是 通信方式的等待, 所以上面的goroutine会执行, 而下面的不会
	////runtime.Goexit() // 等其他执行完后, 崩溃进程
	//go func() {
	//	println("122")
	//}()
	//defer println('1')
	//os.Exit(-1) // 什么都不等, 直接退出

}
