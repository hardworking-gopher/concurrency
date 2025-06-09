package main

import (
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	Name      string
	LeftFork  int
	RightFork int
	Satiation int
}

var (
	hunger  = 3
	eatTime = time.Second * 1

	finishOrder = make([]*Philosopher, 0)
	finishMutex = new(sync.Mutex)
)

func main() {
	fmt.Println("Gathering the philosophers")

	philosophers := []*Philosopher{
		{"Plato", 0, 1, 0},
		{"Democritus", 1, 2, 0},
		{"Confucius", 2, 3, 0},
		{"Aristotle", 3, 0, 0},
	}

	fmt.Println("Getting the philosophers to eat")

	fmt.Println("Order of philosophers:")
	for i, p := range philosophers {
		fmt.Printf("\t%d: %s\n", i+1, p.Name)
	}

	dinner(philosophers)

	fmt.Println("Dinner has been finished")
	fmt.Println("Order of finished philosophers:")

	for _, p := range finishOrder {
		fmt.Println(p.Name)
	}
}

func dinner(philosophers []*Philosopher) {
	wg := new(sync.WaitGroup)
	wg.Add(len(philosophers))

	forks := map[int]*sync.Mutex{}
	for i := range philosophers {
		// express each fork in the form of mutex
		forks[i] = new(sync.Mutex)
	}

	fmt.Println("Making philosophers to eat")

	for _, p := range philosophers {
		go dine(p, forks, wg)
	}

	wg.Wait()
}

func dine(philosopher *Philosopher, forks map[int]*sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < hunger; i++ {

		if philosopher.LeftFork > philosopher.RightFork {
			forks[philosopher.RightFork].Lock()
			fmt.Println(philosopher.Name, "took the right fork")

			forks[philosopher.LeftFork].Lock()
			fmt.Println(philosopher.Name, "took the left fork")
		} else {
			forks[philosopher.LeftFork].Lock()
			fmt.Println(philosopher.Name, "took the left fork")

			forks[philosopher.RightFork].Lock()
			fmt.Println(philosopher.Name, "took the right fork")
		}

		fmt.Println(philosopher.Name, "has both forks")
		fmt.Println(philosopher.Name, "is eating")
		time.Sleep(eatTime)

		forks[philosopher.LeftFork].Unlock()
		forks[philosopher.RightFork].Unlock()

		fmt.Println(philosopher.Name, "put down the forks")
	}

	finishMutex.Lock()
	finishOrder = append(finishOrder, philosopher)
	finishMutex.Unlock()

	fmt.Println("Finished philosopher was added to the list:", philosopher.Name)

	fmt.Println(philosopher.Name, "is satiated")
}
