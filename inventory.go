package main

import (
	"fmt"
	"sort"
)

type Inventory struct {
	Items map[string]*Item // keyed by SKU
}

type InventoryError struct {
	SKU string
	Op  string // e.g. "restock", "sell", "add"
	Msg string
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

func (ierr InventoryError) Error() string {
	return fmt.Sprintf("Error: %s. SKU=%q, Op=%s", ierr.Msg, ierr.SKU, ierr.Op)
}
