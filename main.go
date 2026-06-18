package main

import (
	"errors"
	"fmt"
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

	err := inv.AddItem(Item{SKU: "W-001", Name: "Widget", Price: 4.50, Quantity: 10})
	if invErr, ok := errors.AsType[*InventoryError](err); ok {
		fmt.Println(invErr.Error())
	}

	err = inv.Restock("W-001", 5) // Quantity becomes 15
	if invErr, ok := errors.AsType[*InventoryError](err); ok {
		fmt.Println(invErr.Error())
	}

	err = inv.Sell("W-001", 20) // should return an error, Quantity stays 15
	if invErr, ok := errors.AsType[*InventoryError](err); ok {
		fmt.Println(invErr.Error())
	}

	err = inv.Sell("W-001", 15) // Quantity becomes 0, no error
	if invErr, ok := errors.AsType[*InventoryError](err); ok {
		fmt.Println(invErr.Error())
	}
}
