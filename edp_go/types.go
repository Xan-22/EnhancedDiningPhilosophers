package main

import (
	"sync"
	"time"
)

const (
	THINKING_TIME         = 1000 * time.Millisecond
	EATING_TIME           = 2500 * time.Millisecond
	WAITING_TIME          = 100 * time.Millisecond
	TIMEOUT               = 2500 * time.Millisecond
	STARTING_MONEY        = 200.0
	COUPON_VALUE          = 5.0
	COOKING_TIME          = 6000 * time.Millisecond
	COFFEE_BREAK_TIME     = 3000 * time.Millisecond
	CHECK_ORDERS_INTERVAL = 500 * time.Millisecond
)

type Name string

var (
	PHILOSOPHER_NAMES = []Name{"Susan Haack", "Zhaozhou", "David Hume", "Omar Khayyám", "Kaṇāda"}
	COOK_NAMES        = []Name{"Eren", "Mikasa", "Armin"}
	WAITER_NAMES      = []Name{"Miria", "Isaac"}
)

type FoodType int

const (
	ENTREE FoodType = iota
	SOUP
	DESSERT
)

type Food struct {
	foodType FoodType
	name     string
	price    float64
}

func (f Food) String() string { return f.name }

var (
	entrees = []Food{
		{ENTREE, "Paella", 13.25}, {ENTREE, "Wu Hsiang Chi", 10.00},
		{ENTREE, "Bogrács Gulyás", 11.25}, {ENTREE, "Spanakopita", 6.50},
		{ENTREE, "Moui Nagden", 12.95}, {ENTREE, "Sambal Goreng Udang", 14.95},
	}
	soups    = []Food{{SOUP, "No Soup", 0.00}, {SOUP, "Albóndigas", 3.00}}
	desserts = []Food{{DESSERT, "No Dessert", 0.00}, {DESSERT, "Berog", 3.50}}
)

// Global chopsticks
var CHOPSTICKS []*sync.Mutex
var heldChopsticks = make(map[int]bool)
var heldChopsticksMutex sync.RWMutex

func init() {
	// Initialize chopsticks
	CHOPSTICKS = make([]*sync.Mutex, len(PHILOSOPHER_NAMES))
	for i := range CHOPSTICKS {
		CHOPSTICKS[i] = &sync.Mutex{}
	}
}
