package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("The restaurant is now open for business.")

	for _, cook := range cookManager.List() {
		cook.wg.Add(1)
		go cook.Run()
	}

	for _, waiter := range waiterManager.List() {
		waiter.wg.Add(1)
		go waiter.Run()
	}

	for _, philosopher := range ListPhilosophers() {
		philosopher.wg.Add(1)
		go philosopher.Run()
	}

	monitorRestaurant()
}

func monitorRestaurant() {
	done := false
	for !done {
		time.Sleep(1000 * time.Millisecond)

		if allPhilosophersLeft() {
			fmt.Println("The restaurant has closed down.")
			done = true
		} else {
			logRestaurantStatus()
		}
	}
}

func allPhilosophersLeft() bool {
	for _, philosopher := range ListPhilosophers() {
		select {
		case <-philosopher.ctx.Done():
		default:
			return false
		}
	}
	return true
}

func logRestaurantStatus() {
	activePhilosophers := countActivePhilosophers()
	if activePhilosophers == 0 {
		return
	}

	occupiedSeats := countOccupiedSeats()
	heldChopsticks := GetHeldChopsticks()

	seatInfo := buildSeatInfo(occupiedSeats)
	chopstickInfo := buildChopstickInfo(heldChopsticks)

	fmt.Printf("\nRestaurant status: %d active philosophers, %s, %s, %d orders on counter\n\n",
		activePhilosophers, seatInfo, chopstickInfo, counter.OrderCount())
}

func countActivePhilosophers() int {
	activePhilosophers := 0
	for _, philosopher := range ListPhilosophers() {
		select {
		case <-philosopher.ctx.Done():
		default:
			activePhilosophers++
		}
	}
	return activePhilosophers
}

func countOccupiedSeats() int {
	occupiedSeats := 0
	for _, seat := range seatManager.Seats() {
		if seat.IsOccupied() {
			occupiedSeats++
		}
	}
	return occupiedSeats
}

func buildSeatInfo(occupiedSeats int) string {
	var occupiedSeatNumbers []int
	for _, seat := range seatManager.Seats() {
		if seat.IsOccupied() {
			occupiedSeatNumbers = append(occupiedSeatNumbers, seat.Number())
		}
	}

	var seatInfo strings.Builder
	seatInfo.WriteString(fmt.Sprintf("%d seats taken (", occupiedSeats))
	for i, seatNum := range occupiedSeatNumbers {
		if i > 0 {
			seatInfo.WriteString(", ")
		}
		seatInfo.WriteString(strconv.Itoa(seatNum))
	}
	seatInfo.WriteString(")")
	return seatInfo.String()
}

func buildChopstickInfo(heldChopsticks map[int]bool) string {
	var chopstickNumbers []int
	for chopstickNum := range heldChopsticks {
		chopstickNumbers = append(chopstickNumbers, chopstickNum)
	}
	sort.Ints(chopstickNumbers)

	var chopstickInfo strings.Builder
	chopstickInfo.WriteString(fmt.Sprintf("%d chopsticks taken (", len(heldChopsticks)))
	for i, chopstickNum := range chopstickNumbers {
		if i > 0 {
			chopstickInfo.WriteString(", ")
		}
		chopstickInfo.WriteString(strconv.Itoa(chopstickNum))
	}
	chopstickInfo.WriteString(")")
	return chopstickInfo.String()
}
