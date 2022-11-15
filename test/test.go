package main

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type logger int

const loggerID = "logger_id"

func (l logger) Printf(s string, args ...interface{}) {
	log.Printf("[id=%d] %s", l, fmt.Sprintf(s, args...))
}

func (l logger) Println(s string) {
	log.Printf("[id=%d] %s\n", l, s)
}

func ctxWithLoggerID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, loggerID, id)
}

func getLogger(ctx context.Context) logger {
	return logger(ctx.Value(loggerID).(int))
}

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go hello(ctxWithLoggerID(ctx, i), &wg)
	}
	wg.Wait()
}

func hello(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log := getLogger(ctx)
	log.Printf("hello begin")
	log.Println("hello, processing the hello message")
	log.Printf("hello, end")
}
