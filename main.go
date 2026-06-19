package main

import (
	"errors"
	"fmt"
	"sort"
)

type Item struct {
	SKU      string
	Name     string
	Price    float64 // price per unit, in dollars
	Quantity int
}

func (i Item) TotalValue() float64 {
	return float64(i.Quantity) * i.Price
}

func (i Item) String() string {
	return fmt.Sprintf("Widget (SKU: %s): %d units @ $%g", i.SKU, i.Quantity, i.Price)
}

type Inventory struct {
	Items map[string]*Item // keyed by SKU
}

func NewInventory() *Inventory {
	return &Inventory{Items: map[string]*Item{}}
}

func (inv *Inventory) AddItem(item Item) error {
	if inv.Items[item.SKU] != nil {
		return &InventoryError{
			SKU: item.SKU,
			Op:  "add",
			Msg: "item already exists",
		}
	}

	inv.Items[item.SKU] = &item

	return nil
}

func (inv *Inventory) Restock(sku string, amount int) error {
	if inv.Items[sku] == nil {
		return &InventoryError{
			SKU: sku,
			Op:  "restock",
			Msg: "item does not exist",
		}
	}

	if amount <= 0 {
		return &InventoryError{
			SKU: sku,
			Op:  "restock",
			Msg: "quantity cannot equal or less than 0",
		}
	}

	inv.Items[sku].Quantity = inv.Items[sku].Quantity + amount

	return nil
}

func (inv *Inventory) Sell(sku string, amount int) error {
	if inv.Items[sku] == nil {
		return &InventoryError{
			SKU: sku,
			Op:  "sell",
			Msg: "item does not exist",
		}
	}

	if amount <= 0 {
		return &InventoryError{
			SKU: sku,
			Op:  "sell",
			Msg: "quantity cannot equal or less than 0",
		}
	}

	if inv.Items[sku].Quantity < amount {
		return &InventoryError{
			SKU: inv.Items[sku].SKU,
			Op:  "sell",
			Msg: "not enough stock",
		}
	}

	inv.Items[sku].Quantity = inv.Items[sku].Quantity - amount

	return nil
}

func (inv *Inventory) TotalValue() float64 {
	var totalValue float64

	for _, value := range inv.Items {
		totalValue += value.TotalValue()
	}

	return totalValue
}

func (inv *Inventory) LowStock(threshold int) []*Item {
	itemPtr := []*Item{}

	for _, value := range inv.Items {
		if value.Quantity < threshold {
			itemPtr = append(itemPtr, value)
		}
	}

	sort.Slice(itemPtr, func(i, j int) bool {
		return itemPtr[i].Quantity < itemPtr[j].Quantity
	})

	return itemPtr
}

type InventoryError struct {
	SKU string
	Op  string // e.g. "restock", "sell", "add"
	Msg string
}

func (ierr InventoryError) Error() string {
	return fmt.Sprintf("Error: %s. SKU=%q, Op=%s", ierr.Msg, ierr.SKU, ierr.Op)
}

func main() {
	inv := NewInventory()
	inv.AddItem(Item{SKU: "W-001", Name: "Widget", Price: 4.50, Quantity: 3})
	inv.AddItem(Item{SKU: "G-002", Name: "Gadget", Price: 9.00, Quantity: 1})
	inv.AddItem(Item{SKU: "T-003", Name: "Thingamajig", Price: 2.00, Quantity: 50})

	err := inv.Sell("Z-999", 1) // unknown SKU
	// err should be a *InventoryError with SKU="Z-999", Op="sell"
	if invErr, ok := errors.AsType[*InventoryError](err); ok {
		fmt.Println(invErr.Error())
	}

	fmt.Println(inv.TotalValue())

	lowStockItems := inv.LowStock(5) // should return G-002 (qty 1) then W-001 (qty 3), in that order

	fmt.Println(lowStockItems)
}
