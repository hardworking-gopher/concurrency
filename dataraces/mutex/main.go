package main

import (
	"fmt"
	"sync"
)

// A mutex is a synchronisation primitive that grants access to a shared resource,
// ensuring that only one goroutine can operate on it at a time. It is used for
// more complex data structures, as opposed to atomic package which provides
// synchronisation for primitive data types

type Container struct {
	mu sync.Mutex
	m  map[string]int
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
