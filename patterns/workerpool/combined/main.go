package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker pool pattern using channels with an undefined amount of jobs that need to be performed.
// We use channels when we care about results, or when we depend on values returned from goroutines.
// We often combine channels and wait group for seamless implementation of data transmission and graceful shutdowns.

func main() {
	var (
		workerAmount   = 5
		jobs           = make(chan string, workerAmount)
		result         = make(chan string, workerAmount)
		wgWorkers      = &sync.WaitGroup{}
		wgResultReader = &sync.WaitGroup{}
	)

	for i := 1; i <= workerAmount; i++ {
		wgWorkers.Add(1)

		go worker(wgWorkers, i, jobs, result)
	}

	wgResultReader.Add(1)
	go func() {
		defer wgResultReader.Done()

		fmt.Println("listening for a results")

		for range result {
			fmt.Println("received result from a goroutine")
		}

		fmt.Println("result receiver exited")
	}()

	t := time.After(time.Second * 3)
	for {
		select {
		case jobs <- "job":
			fmt.Println("job send")
			time.Sleep(time.Millisecond * 250)
		case <-t:
			fmt.Println("job sender finished")
			close(jobs)

			wgWorkers.Wait()
			fmt.Println("all workers returned")

			close(result)
			wgResultReader.Wait()
			fmt.Println("all results were read")

			return
		}
	}

}

func worker(wg *sync.WaitGroup, id int, jobs <-chan string, res chan string) {
	defer wg.Done()

	for j := range jobs {
		fmt.Println("worker", id, "doing", j)
		time.Sleep(time.Second)
		res <- "finished job"
		fmt.Println("worker", id, "done doing", j)
	}

	fmt.Println("worker", id, "returned")
}
