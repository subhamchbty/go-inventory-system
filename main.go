package main

import (
	"fmt"
)

const inventoryStorePath string = "inventory.json"

type SimpleCounter struct {
	count int
}

func (c *SimpleCounter) Report() string {
	c.count++
	return fmt.Sprintf("Reports generated: %d", c.count)
}

func main() {
	inv := NewInventory()
	inv.AddItem(Item{SKU: "W-001", Name: "Widget", Price: 4.50, Quantity: 3})
	inv.AddItem(Item{SKU: "G-002", Name: "Gadget", Price: 9.00, Quantity: 1})
	inv.AddItem(Item{SKU: "T-003", Name: "Thingamajig", Price: 2.00, Quantity: 50})

	err := inv.Save(inventoryStorePath)
	if err != nil {
		fmt.Println(err)
	}

	inv2 := NewInventory()
	err = inv2.Load(inventoryStorePath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(inv.TotalValue())
	fmt.Println(inv2.TotalValue())

	var r Reporter = inv
	printReport(r)

	var r2 Reporter = &SimpleCounter{}
	printReport(r2)
}
