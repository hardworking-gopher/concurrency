package main

import (
	"fmt"
	"time"
)

// Generator pattern is when our function returns a channel.
// It is useful when we want to decompose our function or hide some complexity.
// It also helps to decouple producer from consumer.
func main() {
	ann := generator("Ann")
	joe := generator("Joe")

	for i := 0; i < 10; i++ {
		fmt.Println(<-ann)
		fmt.Println(<-joe)
	}

	fmt.Println("exited")
}

func generator(msg string) <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; i < 10; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Millisecond * 250)
		}
		close(c)
	}()

	return c
}
