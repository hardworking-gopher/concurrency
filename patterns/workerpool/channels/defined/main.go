package main

import (
	"fmt"
	"time"
)

// Worker pool pattern using channels with a pre-defined amount of jobs that need to be performed.
// We use channels when we care about results, or when we depend on values returned from goroutines.

func main() {

	var (
		workerAmount = 5
		jobs         = make(chan int, workerAmount)
		results      = make(chan int, workerAmount)
	)

	for i := 1; i <= workerAmount; i++ {
		go worker(jobs, results, i)
	}

	for i := 1; i <= workerAmount; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 1; i <= workerAmount; i++ {
		fmt.Printf("Main received result: %d\n", <-results)
	}

	fmt.Println("All workers finished")
}

func worker(jobs <-chan int, results chan<- int, i int) {
	for job := range jobs {
		fmt.Printf("Worker %d received job %d\n", i, job)
		time.Sleep(time.Second)
		results <- i
		fmt.Printf("Worker %d finished job %d\n", i, job)
	}
}
