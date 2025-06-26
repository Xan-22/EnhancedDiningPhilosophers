package main

import "sync"

type Counter struct {
	orders          chan *Order
	completedMeals  chan *Order
	orderCount      int
	orderCountMutex sync.RWMutex
}

var counter = &Counter{
	orders:         make(chan *Order, 100),
	completedMeals: make(chan *Order, 100),
}

func (c *Counter) PlaceOrder(order *Order) {
	c.orders <- order
	c.orderCountMutex.Lock()
	c.orderCount++
	c.orderCountMutex.Unlock()
}

func (c *Counter) TakeOrder() *Order {
	select {
	case order := <-c.orders:
		c.orderCountMutex.Lock()
		c.orderCount--
		c.orderCountMutex.Unlock()
		return order
	default:
		return nil
	}
}

func (c *Counter) HasOrders() bool { return len(c.orders) > 0 }

func (c *Counter) OrderCount() int {
	c.orderCountMutex.RLock()
	defer c.orderCountMutex.RUnlock()
	return c.orderCount
}

func (c *Counter) PlaceCompletedMeal(order *Order) { c.completedMeals <- order }

func (c *Counter) HasCompletedMeals() bool { return len(c.completedMeals) > 0 }

func (c *Counter) PollCompletedMeal() *Order {
	select {
	case order := <-c.completedMeals:
		return order
	default:
		return nil
	}
} 