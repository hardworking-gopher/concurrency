package main

import (
	"fmt"
	"time"
)

// We can time out either the entire for loop or single iteration by using time.After

func main() {
	c1 := make(chan struct{})

	go func() {
		time.Sleep(time.Second * 5)
		c1 <- struct{}{}
	}()

	t := time.After(time.Second * 6)

	for {
		select {
		case <-c1:
			fmt.Println("read for c1")
		case <-time.After(time.Second):
			fmt.Println("iteration timeout")
		case <-t:
			fmt.Println("overall timeout")
			return
		}
	}
}
