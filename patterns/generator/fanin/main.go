package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Fan-in (or multiplexing) utilizes the generator pattern to react to whichever input
// is ready with a result, rather than waiting for all results before repeating the loop.
// This pattern is suitable in cases when we have multiple sources of data, and we cannot
// predict the order at which this data is arrived.

func main() {
	var (
		wg  = &sync.WaitGroup{}
		out = funIn(
			wg,
			generator("data source 1"),
			generator("data source 2"),
		)
	)

	for i := 0; i < 20; i++ {
		fmt.Println("received", <-out)
	}

	wg.Wait()

	fmt.Println("done")
}

func funIn(wg *sync.WaitGroup, inputs ...<-chan string) <-chan string {
	output := make(chan string)

	for _, input := range inputs {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for {
				v, ok := <-input
				if !ok {
					return
				}
				output <- v
			}
		}()
	}

	return output
}

func generator(msg string) <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))

			fmt.Println("sending", i, "to", msg)
			c <- fmt.Sprintf("%d from %s", i, msg)
		}
		close(c)
	}()

	return c
}
