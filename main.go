package main

import (
	"errors"
	"fmt"
)

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
