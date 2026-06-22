package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

type Inventory struct {
	Items       map[string]*Item           `json:"items"` // keyed by SKU
	Perishables map[string]*PerishableItem `json:"perishables"`
}

type InventoryError struct {
	SKU string
	Op  string // e.g. "restock", "sell", "add"
	Msg string
}

func NewInventory() *Inventory {
	return &Inventory{Items: map[string]*Item{}, Perishables: map[string]*PerishableItem{}}
}

func (inv *Inventory) AddItem(item Item) error {
	if inv.Perishables[item.SKU] != nil || inv.Items[item.SKU] != nil {
		return &InventoryError{
			SKU: item.SKU,
			Op:  "add",
			Msg: "item already exists",
		}
	}

	inv.Items[item.SKU] = &item

	return nil
}

func (inv *Inventory) AddPerishable(item PerishableItem) error {
	if inv.Perishables[item.SKU] != nil || inv.Items[item.SKU] != nil {
		return &InventoryError{
			SKU: item.SKU,
			Op:  "add",
			Msg: "item already exists",
		}
	}

	inv.Perishables[item.SKU] = &item

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

	for _, p := range inv.Perishables {
		totalValue += p.TotalValue()
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

func (inv *Inventory) Save(path string) error {
	jsonData, err := json.MarshalIndent(inv, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 1000)
	if err != nil {
		return err
	}

	return nil
}

func (inv *Inventory) Load(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, inv)
	if err != nil {
		return err
	}

	return nil
}

func (inv *Inventory) Report() string {
	var totalUnits int
	var totalExpired int
	for _, item := range inv.Items {
		totalUnits += item.Quantity
	}

	today := time.Now()
	for _, p := range inv.Perishables {
		if p.IsExpired(today.Format("2006-01-02")) {
			totalExpired++
		}
		totalUnits += p.Quantity
	}

	return fmt.Sprintf("Distinct SKUs: %d\nTotal units: %d\nTotal value: $%.2f\nPerishable Items: %d\nExpired Items: %d",
		len(inv.Items)+len(inv.Perishables),
		totalUnits,
		inv.TotalValue(),
		len(inv.Perishables),
		totalExpired,
	)
}

func (ierr InventoryError) Error() string {
	return fmt.Sprintf("Error: %s. SKU=%q, Op=%s", ierr.Msg, ierr.SKU, ierr.Op)
}
