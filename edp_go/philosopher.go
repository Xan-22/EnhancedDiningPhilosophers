package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	name                Name
	money               float64
	seat                *Seat
	order               *Order
	foodServedChan      chan bool
	shouldReceiveCoupon bool
	ctx                 context.Context
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
}

var philosophers []*Philosopher
var philosophersMutex sync.RWMutex

func init() {
	philosophers = make([]*Philosopher, len(PHILOSOPHER_NAMES))
	for i, name := range PHILOSOPHER_NAMES {
		ctx, cancel := context.WithCancel(context.Background())
		philosophers[i] = &Philosopher{
			name:           name,
			money:          STARTING_MONEY,
			foodServedChan: make(chan bool, 1),
			ctx:            ctx,
			cancel:         cancel,
		}
	}
}

func ListPhilosophers() []*Philosopher {
	philosophersMutex.RLock()
	defer philosophersMutex.RUnlock()
	result := make([]*Philosopher, len(philosophers))
	copy(result, philosophers)
	return result
}

func GetHeldChopsticks() map[int]bool {
	heldChopsticksMutex.RLock()
	defer heldChopsticksMutex.RUnlock()
	result := make(map[int]bool)
	for k, v := range heldChopsticks {
		result[k] = v
	}
	return result
}

func (p *Philosopher) Name() Name            { return p.name }
func (p *Philosopher) Order() *Order         { return p.order }
func (p *Philosopher) SetOrder(order *Order) { p.order = order }
func (p *Philosopher) ClearOrder()           { p.order = nil }

func (p *Philosopher) think() { time.Sleep(THINKING_TIME) }

func (p *Philosopher) eat() {
	seatNumber := p.seat.Number()
	leftChopstick := seatNumber
	rightChopstick := (seatNumber + 1) % len(philosophers)

	var firstChopstick, secondChopstick int
	if seatNumber == len(philosophers)-1 {
		firstChopstick, secondChopstick = rightChopstick, leftChopstick
	} else {
		firstChopstick, secondChopstick = leftChopstick, rightChopstick
	}

	CHOPSTICKS[firstChopstick].Lock()
	heldChopsticksMutex.Lock()
	heldChopsticks[firstChopstick] = true
	heldChopsticksMutex.Unlock()

	CHOPSTICKS[secondChopstick].Lock()
	heldChopsticksMutex.Lock()
	heldChopsticks[secondChopstick] = true
	heldChopsticksMutex.Unlock()

	time.Sleep(EATING_TIME)

	heldChopsticksMutex.Lock()
	delete(heldChopsticks, secondChopstick)
	heldChopsticksMutex.Unlock()
	CHOPSTICKS[secondChopstick].Unlock()

	heldChopsticksMutex.Lock()
	delete(heldChopsticks, firstChopstick)
	heldChopsticksMutex.Unlock()
	CHOPSTICKS[firstChopstick].Unlock()
}

func (p *Philosopher) pay() {
	mealCost := p.order.Cost()
	p.money -= mealCost
	if p.money < 0 {
		fmt.Printf("Philosopher %s cannot afford the meal ($%.2f) and is leaving for good. Balance: $%.2f\n",
			p.name, mealCost, p.money)
		p.VacateSeat()
		p.money = 0
	} else {
		fmt.Printf("Philosopher %s has paid $%.2f and left the restaurant.\n", p.name, mealCost)
		p.VacateSeat()
	}
}

func (p *Philosopher) waitForWaiter() bool {
	fmt.Printf("Philosopher %s is waiting for a waiter.\n", p.name)
	waiterManager.AddPhilosopherToQueue(p)

	startTime := time.Now()
	for time.Since(startTime) <= TIMEOUT {
		if p.order != nil {
			fmt.Printf("Philosopher %s got an order from waiter.\n", p.name)
			return true
		}
		time.Sleep(WAITING_TIME)
	}

	fmt.Printf("Philosopher %s gave up waiting for a waiter.\n", p.name)
	return false
}

func (p *Philosopher) NotifyFoodServed() {
	select {
	case p.foodServedChan <- true:
	default:
	}
}

func (p *Philosopher) SetShouldReceiveCoupon(shouldReceive bool) {
	p.shouldReceiveCoupon = shouldReceive
}

func (p *Philosopher) giveCoupon(amount float64) {
	p.money += amount
	fmt.Printf("Philosopher %s received a $%.2f coupon. New balance: $%.2f\n", p.name, amount, p.money)
}

func (p *Philosopher) Run() {
	defer p.wg.Done()
	defer p.cancel()
	for p.money > 0 {
		p.attemptToDine()
		p.think()
	}
	p.VacateSeat()
	fmt.Printf("Philosopher %s has left the restaurant for good.\n", p.name)
}

func (p *Philosopher) attemptToDine() {
	p.seat = seatManager.AvailableSeat()
	if p.seat != nil && p.seat.AttemptToOccupy() {
		fmt.Printf("Philosopher %s is being seated in chair %d.\n", p.name, p.seat.Number())
		p.think()
		fmt.Printf("Philosopher %s is about to call for a waiter.\n", p.name)
		hasWaiter := p.waitForWaiter()
		if !hasWaiter {
			fmt.Printf("Philosopher %s has left the restaurant without being served.\n", p.name)
			p.seat.Vacate()
			return
		}
		fmt.Printf("Philosopher %s got waiter, waiting for food.\n", p.name)
		p.shouldReceiveCoupon = false
		select {
		case <-p.foodServedChan:
			if p.shouldReceiveCoupon {
				p.giveCoupon(COUPON_VALUE)
				return
			}
		case <-p.ctx.Done():
			return
		}
		fmt.Printf("Philosopher %s got food, about to eat.\n", p.name)
		p.eat()
		if p.order != nil {
			p.pay()
			p.ClearOrder()
		}
	} else {
		fmt.Printf("Philosopher %s could not get a seat.\n", p.name)
	}
}

func (p *Philosopher) VacateSeat() {
	if p.seat != nil {
		p.seat.Vacate()
		p.seat = nil
	}
}
