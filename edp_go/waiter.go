package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type WaiterManager struct {
	waiters          []*Waiter
	philosopherQueue chan *Philosopher
	cookSemaphore    chan struct{}
}

var waiterManager = &WaiterManager{
	philosopherQueue: make(chan *Philosopher, 1000),
	cookSemaphore:    make(chan struct{}, len(COOK_NAMES)),
}

func init() {
	for i := 0; i < len(COOK_NAMES); i++ {
		waiterManager.cookSemaphore <- struct{}{}
	}
	waiterManager.waiters = make([]*Waiter, len(WAITER_NAMES))
	for i, name := range WAITER_NAMES {
		waiterManager.waiters[i] = NewWaiter(name)
	}
}

func (wm *WaiterManager) List() []*Waiter {
	result := make([]*Waiter, len(wm.waiters))
	copy(result, wm.waiters)
	return result
}

func (wm *WaiterManager) AvailableWaiter() *Waiter {
	for _, waiter := range wm.waiters {
		if !waiter.IsProcessingOrder() {
			return waiter
		}
	}
	return nil
}

func (wm *WaiterManager) AddPhilosopherToQueue(philosopher *Philosopher) {
	select {
	case wm.philosopherQueue <- philosopher:
	default:
	}
}

type Waiter struct {
	name              Name
	order             *Order
	patron            *Philosopher
	isProcessingOrder bool
	mutex             sync.RWMutex
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
}

func NewWaiter(name Name) *Waiter {
	ctx, cancel := context.WithCancel(context.Background())
	return &Waiter{name: name, ctx: ctx, cancel: cancel}
}

func (w *Waiter) Name() string { return string(w.name) }
func (w *Waiter) IsProcessingOrder() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.isProcessingOrder
}

func (w *Waiter) TakeOrder(philosopher *Philosopher, order *Order) {
	w.mutex.Lock()
	w.isProcessingOrder = true
	w.order = order
	w.patron = philosopher
	w.mutex.Unlock()
}

func (w *Waiter) Run() {
	defer w.wg.Done()
	fmt.Printf("Waiter %s is ready to take orders.\n", w.name)
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			w.mutex.RLock()
			processing := w.isProcessingOrder
			order := w.order
			patron := w.patron
			w.mutex.RUnlock()

			if processing && order != nil && patron != nil {
				w.processOrder()
			} else {
				if counter.HasCompletedMeals() {
					completedOrder := counter.PollCompletedMeal()
					if completedOrder != nil {
						w.deliverOrder(completedOrder)
					}
				} else {
					w.checkForPhilosophersNeedingService()
				}
			}
			time.Sleep(CHECK_ORDERS_INTERVAL)
		}
	}
}

func (w *Waiter) processOrder() {
	w.mutex.RLock()
	currentOrder := w.order
	currentPatron := w.patron
	w.mutex.RUnlock()

	if currentOrder == nil || currentPatron == nil {
		return
	}

	select {
	case <-waiterManager.cookSemaphore:
		select {
		case cookManager.orderQueue <- currentOrder:
			fmt.Printf("Waiter %s placed order for %s.\n", w.name, currentPatron.Name())
		default:
		}
	default:
		fmt.Printf("Waiter %s cannot place order for %s - all chefs busy. Giving $5.00 coupon.\n",
			w.name, currentPatron.Name())
		currentPatron.SetShouldReceiveCoupon(true)
		currentPatron.NotifyFoodServed()
		currentPatron.ClearOrder()
		currentPatron.VacateSeat()
		fmt.Printf("Philosopher %s has left the restaurant without being served.\n", currentPatron.Name())
	}

	w.mutex.Lock()
	w.order = nil
	w.patron = nil
	w.isProcessingOrder = false
	w.mutex.Unlock()
}

func (w *Waiter) deliverOrder(completedOrder *Order) {
	targetPhilosopher := completedOrder.Philosopher()
	targetPhilosopher.NotifyFoodServed()
	fmt.Printf("Waiter %s is serving philosopher %s %s.\n", w.name,
		targetPhilosopher.Name(), completedOrder.MealString())

	w.mutex.Lock()
	w.order = nil
	w.patron = nil
	w.isProcessingOrder = false
	w.mutex.Unlock()
}

func (w *Waiter) checkForPhilosophersNeedingService() {
	w.mutex.RLock()
	processing := w.isProcessingOrder
	w.mutex.RUnlock()

	if processing {
		return
	}

	select {
	case callingPhilosopher := <-waiterManager.philosopherQueue:
		newOrder := NewOrder(callingPhilosopher)
		callingPhilosopher.SetOrder(newOrder)
		w.TakeOrder(callingPhilosopher, newOrder)

		mealDesc := newOrder.MealString()
		fmt.Printf("Waiter %s has taken order for %s from philosopher %s.\n",
			w.name, mealDesc, callingPhilosopher.Name())
	default:
	}
}
