package main

import (
	"fmt"
)

// We can use select to either send a message to a channel or continue if the sending operation will be blocked

func main() {
	c := make(chan int)

	go func() {
		c <- 1 // fills empty channel, blocking call
	}()

	select {
	case c <- 2: // channel is filled, cannot send message right away
	default:
		fmt.Println("channel is full") // triggered
	}
}
