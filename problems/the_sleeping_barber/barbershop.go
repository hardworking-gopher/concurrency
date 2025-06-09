package main

import (
	"github.com/fatih/color"
	"time"
)

type BarberShop struct {
	Barbers         int
	IsOpen          bool
	WaitingRoomChan chan string
	BarberDoneChan  chan string
}

func (b *BarberShop) AddBarber(name string, haircutDuration time.Duration) {
	b.Barbers++

	go func() {
		for {
			client, ok := <-b.WaitingRoomChan
			if !ok {
				color.HiRed("%s sees that waiting room is empty and closed", name)
				color.Green("Sending %s to home", name)

				b.BarberDoneChan <- name
				return
			}

			color.Green("\t%s started working on %s", name, client)
			time.Sleep(haircutDuration)
			color.Green("\t%s has finished working on %s", name, client)
		}
	}()
}

func (b *BarberShop) AddClient(name string, c *color.Color) {
	if b.IsOpen {
		select {
		case b.WaitingRoomChan <- name:
			c.Printf("%s has claimed a seat\n", name)
		default:
			color.Red("Waiting room is full, %s has left", name)
		}
	} else {
		color.Red("Barbershop is already closed so %s left", name)
	}
}

func (b *BarberShop) CloseBarbershop() {
	b.IsOpen = false

	close(b.WaitingRoomChan)

	for i := 1; i <= b.Barbers; i++ {
		name := <-b.BarberDoneChan
		color.Green("%s went home", name)
	}

	color.Red("All barbers went home")

	close(b.BarberDoneChan)
}
