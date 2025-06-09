package main

import (
	"fmt"
	"math/rand"
	"time"
)

// We can use select to read from whichever channel is ready

func main() {
	c1 := make(chan struct{})
	c2 := make(chan struct{})

	go func() {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
		c1 <- struct{}{}
	}()

	go func() {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
		c2 <- struct{}{}
	}()

	select {
	case <-c1:
		fmt.Println("read from the first channel")
	case <-c2:
		fmt.Println("read from the second channel")
	}
}
