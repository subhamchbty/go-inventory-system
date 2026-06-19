package main

import "fmt"

type Item struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"` // price per unit, in dollars
	Quantity int     `json:"quantity"`
}

func (i Item) TotalValue() float64 {
	return float64(i.Quantity) * i.Price
}

func (i Item) String() string {
	return fmt.Sprintf("Widget (SKU: %s): %d units @ $%g", i.SKU, i.Quantity, i.Price)
}
