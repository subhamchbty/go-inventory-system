# Inventory System (Go)

This is a small inventory system I built while learning Go. It's based on a
practice exercise (see `challenge.md`).

## What it does

- Keeps track of items (SKU, name, price, quantity)
- You can add items, restock them and sell them
- Has proper errors so you know *why* something failed (custom `InventoryError` type)
- Saves everything to a JSON file and can load it back
- Reports total value, low stock items etc
- Also supports perishable items (with an expiry date) using struct embedding

## Files

- `item.go` - the `Item` struct, plus `PerishableItem` which embeds `Item`
- `inventory.go` - the `Inventory` struct and all its methods + the error type
- `report.go` - the `Reporter` interface and `printReport`
- `main.go` - the `main` function, just wires everything together

## How to run

```
go run .
```

That builds and runs all the files in the folder. It will create an
`inventory.json` file in the same directory.

If you want to build a binary instead:

```
go build .
```

## Notes to self / things I learned

- Use **pointer receivers** when the method changes the struct (like Restock,
  Sell) otherwise you're just editing a copy and nothing happens.
- The inventory map is `map[string]*Item` (pointers!) so I can change the
  quantity directly.
- `errors.As` is the proper way to check for my custom error type instead of
  a type assertion, it still works even if the error gets wrapped later.
- Embedding `Item` inside `PerishableItem` means `PerishableItem` gets
  `TotalValue()` and `String()` for free without writing any forwarding code.
  That was pretty cool.

## Packages used

Only the standard library: `fmt`, `os`, `errors`, `encoding/json`, `sort`, `time`.
No third party stuff.
