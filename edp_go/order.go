package main

import (
	"math/rand"
	"strings"
)

func randomMeal() []Food {
	// Safety checks to prevent index out of range
	if len(entrees) == 0 || len(soups) == 0 || len(desserts) == 0 {
		// Return a default meal if arrays are not initialized
		return []Food{
			{ENTREE, "Default Entree", 10.00},
			{SOUP, "No Soup", 0.00},
			{DESSERT, "No Dessert", 0.00},
		}
	}

	return []Food{
		entrees[rand.Intn(len(entrees))],
		soups[rand.Intn(len(soups))],
		desserts[rand.Intn(len(desserts))],
	}
}

type Order struct {
	philosopher *Philosopher
	meal        []Food
	cost        float64
}

func NewOrder(philosopher *Philosopher) *Order {
	meal := randomMeal()
	cost := 0.0
	for _, food := range meal {
		cost += food.price
	}
	return &Order{philosopher: philosopher, meal: meal, cost: cost}
}

func (o *Order) Philosopher() *Philosopher { return o.philosopher }
func (o *Order) Cost() float64             { return o.cost }
func (o *Order) Meal() []Food              { return o.meal }

func (o *Order) MealString() string {
	var mealItems []string
	for _, food := range o.meal {
		if (food.foodType == SOUP && food.name == "No Soup") ||
			(food.foodType == DESSERT && food.name == "No Dessert") {
			continue
		}
		mealItems = append(mealItems, food.String())
	}
	return strings.Join(mealItems, " and ")
}
