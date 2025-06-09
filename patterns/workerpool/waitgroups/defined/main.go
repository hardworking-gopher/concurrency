package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker pool pattern using wait groups with a pre-defined amount of jobs that need to be performed.
// We use wait groups when we don't care about propagating results or errors to the top, or when we need to
// synchronize multiple goroutines before continuing with further processing or some kind of graceful shutdown.

func main() {
	workerAmount := 5

	wg := &sync.WaitGroup{}
	for i := 1; i <= workerAmount; i++ {
		wg.Add(1)

		go worker(wg, i)
	}

	fmt.Println("Waiting for all workers to finish")
	wg.Wait()

	fmt.Println("All workers finished")
}

func worker(wg *sync.WaitGroup, i int) {
	defer wg.Done()

	fmt.Printf("Worker %d is doing its job\n", i)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d finished\n", i)
}
