package main

import (
	"bytes"
	"sync"
)

type Info struct {
	mu  sync.Mutex
	Str string
}

func update(info *Info) {
	info.mu.Lock()
	info.Str = "a"
	info.mu.Unlock()
}

type SyncBuffer struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func main() {

}
