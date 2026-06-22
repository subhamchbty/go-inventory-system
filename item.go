package main

import (
	"fmt"
	"time"
)

type Item struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"` // price per unit, in dollars
	Quantity int     `json:"quantity"`
}

type PerishableItem struct {
	Item
	ExpiryDate string `json:"expiry_date"` // format: "2006-01-02"
}

func (i Item) TotalValue() float64 {
	return float64(i.Quantity) * i.Price
}

func (i Item) String() string {
	return fmt.Sprintf("%s (SKU: %s): %d units @ $%g", i.Name, i.SKU, i.Quantity, i.Price)
}

func (p PerishableItem) IsExpired(today string) bool {
	t1, _ := time.Parse("2006-01-02", p.ExpiryDate)
	t2, _ := time.Parse("2006-01-02", today)
	if t1.Before(t2) || t1.Equal(t2) {
		return true
	}

	// fmt.Println(time.Parse("2006-01-02", today))

	return false
}
