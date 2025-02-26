package main

import (
	"github.com/brianvoe/gofakeit/v6"
	"sync"
	"testing"
)

func Test_dine(t *testing.T) {

	philosophers := make([]*Philosopher, 4)
	forks := make(map[int]*sync.Mutex)

	for i := range philosophers {
		forks[i] = &sync.Mutex{}

		leftFork, rightFork := i, i+1

		if i == len(philosophers)-1 {
			rightFork = 0
		}

		philosophers[i] = &Philosopher{
			Name:      gofakeit.Name(),
			LeftFork:  leftFork,
			RightFork: rightFork,
		}
	}

	wg := new(sync.WaitGroup)

	t.Run("success", func(t *testing.T) {
		for _, philosopher := range philosophers {
			wg.Add(1)

			dine(philosopher, forks, wg)
		}
	})

	wg.Wait()
}
