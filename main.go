package main

import (
	"fmt"
	"log"
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
	err := inv.AddItem(Item{SKU: "W-001", Name: "Widget", Price: 4.50, Quantity: 3})
	if err != nil {
		log.Fatal(err)
	}
	err = inv.AddItem(Item{SKU: "G-002", Name: "Gadget", Price: 9.00, Quantity: 1})
	if err != nil {
		log.Fatal(err)
	}

	err = inv.AddItem(Item{SKU: "T-003", Name: "Thingamajig", Price: 2.00, Quantity: 50})
	if err != nil {
		log.Fatal(err)
	}

	err = inv.AddPerishable(PerishableItem{
		Item:       Item{SKU: "T-005", Name: "Milk", Price: 3.20, Quantity: 8},
		ExpiryDate: "2026-06-17",
	})

	if err != nil {
		log.Fatal(err)
	}

	err = inv.Save(inventoryStorePath)
	if err != nil {
		log.Fatal(err)
	}

	inv2 := NewInventory()
	err = inv2.Load(inventoryStorePath)
	if err != nil {
		log.Fatal(err)
	}

	var r Reporter = inv
	printReport(r)

	var r2 Reporter = &SimpleCounter{}
	printReport(r2)
}
