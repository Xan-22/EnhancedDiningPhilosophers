package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type CookManager struct {
	cooks      []*Cook
	orderQueue chan *Order
}

var cookManager = &CookManager{
	orderQueue: make(chan *Order, 1),
}

func init() {
	cookManager.cooks = make([]*Cook, len(COOK_NAMES))
	for i, name := range COOK_NAMES {
		cookManager.cooks[i] = NewCook(name)
	}
}

func (cm *CookManager) List() []*Cook {
	result := make([]*Cook, len(cm.cooks))
	copy(result, cm.cooks)
	return result
}

type Cook struct {
	name          string
	mealsPrepared int
	isOnBreak     bool
	mutex         sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewCook(name Name) *Cook {
	ctx, cancel := context.WithCancel(context.Background())
	return &Cook{name: string(name), ctx: ctx, cancel: cancel}
}

func (c *Cook) Name() string { return c.name }
func (c *Cook) IsOnBreak() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isOnBreak
}

func (c *Cook) Run() {
	defer c.wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			fmt.Printf("Chef %s is waiting for an order.\n", c.name)
			select {
			case order := <-cookManager.orderQueue:
				if order != nil {
					c.cook(order)
					c.mutex.Lock()
					c.mealsPrepared++
					if c.mealsPrepared%4 == 0 {
						c.mutex.Unlock()
						c.takeCoffeeBreak()
						c.mutex.Lock()
					}
					c.mutex.Unlock()
				}
			case <-c.ctx.Done():
				return
			}
		}
	}
}

func (c *Cook) cook(order *Order) {
	fmt.Printf("Chef %s is cooking the %s for Philosopher %s.\n", c.name, order.MealString(),
		order.Philosopher().Name())
	time.Sleep(COOKING_TIME)
	counter.PlaceCompletedMeal(order)
	fmt.Printf("Chef %s has finished cooking the %s for Philosopher %s.\n", c.name,
		order.MealString(), order.Philosopher().Name())
	
	select {
	case waiterManager.cookSemaphore <- struct{}{}:
	default:
	}
}

func (c *Cook) takeCoffeeBreak() {
	c.mutex.Lock()
	c.isOnBreak = true
	c.mutex.Unlock()
	
	fmt.Printf("Chef %s is taking a coffee break.\n", c.name)
	time.Sleep(COFFEE_BREAK_TIME)
	
	c.mutex.Lock()
	c.isOnBreak = false
	c.mutex.Unlock()
	fmt.Printf("Chef %s has returned from a coffee break.\n", c.name)
} 