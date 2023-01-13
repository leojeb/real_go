package main

import (
	"context"
	"time"
)

func main() {
	//request(context.Background())
	//fmt.Scanln()
	request1(context.Background())
}

// case 1
func request(ctx context.Context) {
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second)
	// 必须显式的调用cancelFunc避免及时回收
	defer cancelFunc()

	resp := make(chan int)
	go handle(ctx, resp)

	// do something...
	select {
	case v := <-resp:
		println(v)
	case <-ctx.Done():
		println(ctx.Err().Error())
	}

}

func handle(ctx context.Context, resp chan<- int) {
	println("1/3 handle")
	cache(ctx, resp)
}

func cache(ctx context.Context, resp chan<- int) {
	println("2/3 cache")

	// long-time working...
	time.Sleep(time.Second * 3)

	database(ctx, resp)
}

func database(ctx context.Context, resp chan<- int) {

	// check done!
	select {
	case <-ctx.Done():
		println("3/3 database: timeout!")
		return
	default:

	}

	println("3/3 database")
	resp <- 100
}

// case 2
type contextKey string

func request1(ctx context.Context) {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	// 避免key的冲突， 自定义类型来设置context的key
	ctx = context.WithValue(ctx, contextKey("A"), "a")

	go handle1(ctx)

	// do sth ...
	time.Sleep(time.Second)

}

func handle1(ctx context.Context) {
	println("1/3 handled")
	ctx = context.WithValue(ctx, contextKey("B"), "b")
	cache1(ctx)
}

func cache1(ctx context.Context) {
	println("2/3 cached")
	println("\t", ctx.Value(contextKey("A")).(string))
	println("\t", ctx.Value(contextKey("B")).(string))

	ctx = context.WithValue(ctx, contextKey("A"), "a2")

	database1(ctx)
}

func database1(ctx context.Context) {
	println("3/3 database")
	println("\t", ctx.Value(contextKey("A")).(string))
	println("\t", ctx.Value(contextKey("B")).(string))

}
