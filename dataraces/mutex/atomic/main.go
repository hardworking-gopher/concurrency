package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// The sync/atomic package in Go provides low-level atomic operations on
// primitive data types, offering a lighter-weight alternative to mutexes
// for ensuring thread-safety when updating or reading single values

func main() {
	var (
		at = atomic.Int32{}
		wg = &sync.WaitGroup{}
	)

	for range 50 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range 100 {
				at.Add(1)
			}
		}()
	}

	wg.Wait()

	fmt.Println(at.Load())
}
