package main

import (
	"fmt"
)

// We can use idiomatic way of handling access to data in
// concurrent environment by using goroutine that owns the data

type WriteOp struct {
	Key  string
	Val  string
	Done chan bool
}

type ReadOp struct {
	Key string
	Val string
	Res chan int
}

func main() {

	var (
		// main is the owner of the data
		state = map[string]int{"a": 0}

		writeOps = make(chan WriteOp)
		readOps  = make(chan ReadOp)
	)

	inc := func(key string, n int) {

		for range n {
			wOp := WriteOp{
				Key:  "a",
				Done: make(chan bool),
			}
			writeOps <- wOp
			<-wOp.Done // signal that data was written if needed
		}
	}

	read := func(key string, n int) {

		for range n {
			rOp := ReadOp{
				Key: "a",
				Res: make(chan int),
			}
			readOps <- rOp
			fmt.Println("read data", <-rOp.Res)
		}
	}

	// other goroutines send data to the owner rather than try to write them on their own
	go inc("a", 1000)
	go read("a", 100)
	go inc("a", 1000)

	// select waits for either a read or a write operation to be performed, one at a time
	for {
		select {
		case w := <-writeOps:
			state[w.Key]++
			w.Done <- true

			if state[w.Key] == 2000 {
				close(writeOps)
				close(readOps)

				fmt.Println(state)

				return
			}
		case r := <-readOps:
			r.Res <- state[r.Key]
		}
	}

}
