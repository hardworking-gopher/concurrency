package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Container struct {
	mu sync.Mutex
	m  map[string]int
	m2 map[string]atomic.Int32
}

func main() {
	c := Container{
		m: map[string]int{"a": 0, "b": 0},
	}

	wg := sync.WaitGroup{}

	inc := func(k string, v int) {
		defer wg.Done()

		for range v {
			c.mu.Lock()
			c.m[k]++
			c.mu.Unlock()
		}
	}

	wg.Add(3)
	go inc("a", 1000)
	go inc("b", 1000)
	go inc("b", 1000)

	wg.Wait()

	fmt.Println(c.m)
}
