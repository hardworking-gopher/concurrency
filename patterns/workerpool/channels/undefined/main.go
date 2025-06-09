package main

import (
	"fmt"
	"time"
)

// Worker pool pattern using channels with an undefined amount of jobs that need to be performed.
// We use channels when we care about results, or when we depend on values returned from goroutines.

func main() {
	var (
		workerAmount = 5
		jobs         = make(chan string, workerAmount)
		done         = make(chan int, workerAmount)
	)

	for i := 1; i <= workerAmount; i++ {
		go worker(i, jobs, done)
	}

	t := time.After(time.Second * 5)
	for {
		select {
		case jobs <- "job":
			fmt.Println("job send")
			time.Sleep(time.Millisecond * 250)
		case <-t:
			fmt.Println("job sender finished")
			close(jobs)

			for i := 0; i < workerAmount; i++ {
				fmt.Println("worker", <-done, "returned")
			}
			close(done)

			fmt.Println("all workers returned, exiting")

			return
		}
	}

}

func worker(id int, jobs <-chan string, done chan<- int) {
	defer func(done chan<- int) { done <- id }(done)

	for j := range jobs {
		fmt.Println("worker", id, "doing", j)
		time.Sleep(time.Millisecond * 300)
		fmt.Println("worker", id, "done doing", j)
	}

	fmt.Println("worker", id, "done")
}
