package main

import (
	"fmt"
)

// Pipeline pattern allows us to chain multiple stages.

func main() {
	nums := []int{2, 4, 8, 9, 10}

	for s := range square(sender(nums)) {
		fmt.Println(s)
	}
}

func sender(nums []int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()

	return out
}

func square(nums <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for n := range nums {
			out <- n * n
		}
		close(out)
	}()

	return out
}
