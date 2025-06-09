// This is a simple demonstration of how to solve the Sleeping Barber dilemma, a classic computer science problem
// which illustrates the complexities that arise when there are multiple operating system processes. Here, we have
// a finite number of barbers, a finite number of seats in a waiting room, a fixed length of time the barbershop is
// open, and clients arriving at (roughly) regular intervals. When a barber has nothing to do, he or she checks the
// waiting room for new clients, and if one or more is there, a haircut takes place. Otherwise, the barber goes to
// sleep until a new client arrives. So the rules are as follows:
//
//   - if there are no customers, the barber falls asleep in the chair
//   - a customer must wake the barber if he is asleep
//   - if a customer arrives while the barber is working, the customer leaves if all chairs are occupied and
//     sits in an empty chair if it's available
//   - when the barber finishes a haircut, he inspects the waiting room to see if there are any waiting customers
//     and falls asleep if there are none
//   - shop can stop accepting new clients at closing time, but the barbers cannot leave until the waiting room is
//     empty
//   - after the shop is closed and there are no clients left in the waiting area, the barber
//     goes home
//
// The Sleeping Barber was originally proposed in 1965 by computer science pioneer Edsger Dijkstra.
//
// The point of this problem, and its solution, was to make it clear that in a lot of cases, the use of
// semaphores (mutexes) is not needed.
package main

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

// variables

const (
	waitingRoomLength  = 3
	timeBarbershopOpen = 5 * time.Second
	clientArrivalRate  = 500 * time.Millisecond
)

var colors = []color.Attribute{
	color.FgHiBlue,
	color.FgHiMagenta,
	color.FgHiCyan,
	color.FgHiWhite,
}

func main() {
	var (
		r          = rand.New(rand.NewSource(time.Now().UnixNano()))
		barbershop = &BarberShop{
			IsOpen:          true,
			WaitingRoomChan: make(chan string, waitingRoomLength),
			BarberDoneChan:  make(chan string),
		}
	)

	barbershop.AddBarber(gofakeit.FirstName(), time.Second*time.Duration(r.Intn(5)+1))
	barbershop.AddBarber(gofakeit.FirstName(), time.Second*time.Duration(r.Intn(5)+1))

	color.HiYellow("Barbershop is opened!")
	color.HiYellow("Barbers available: %d", barbershop.Barbers)

	var (
		closingChan = make(chan bool)
		closed      = make(chan bool)
	)

	go func() {
		<-time.After(timeBarbershopOpen)
		color.HiRed("Closing barbershop")

		closingChan <- true
		barbershop.CloseBarbershop()
		closed <- true
	}()

	go func() {
		iColor := 0

		for {
			iColor++

			// cycle though some color for each client
			if iColor == len(colors) {
				iColor = 0
			}

			client := gofakeit.FirstName()

			randColor := color.New(colors[iColor])
			randColor.Printf("New client has arrived - %s\n", client)
			randColor.Printf("Waiting room status: %d/%d\n", len(barbershop.WaitingRoomChan), waitingRoomLength)

			select {
			case <-time.After(clientArrivalRate):
				barbershop.AddClient(client, randColor)
			case <-closingChan:
				return
			}

		}
	}()

	<-closed

	color.HiRed("Barbershop is closed!")
}
