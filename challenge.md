# Exercise: Inventory System with Structs in Go

**Difficulty:** Intermediate
**Estimated time:** 90–120 minutes
**Concepts covered:** structs, methods (value vs pointer receivers), slices of structs, maps, interfaces, error handling, custom error types, JSON marshalling/unmarshalling, file I/O, embedding, sorting

---

## Background

Most real-world programs end up modelling some kind of collection of "things with state" — and an inventory system is a clean, realistic way to practice that. You'll build a system that tracks items, supports stock operations with proper error handling, persists itself to disk as JSON, and reports on itself.

You will build this incrementally across six parts. **Do not skip ahead** — each part introduces a concept that the next part depends on, and later parts assume earlier ones are correct.

---

## Setup

Create a new directory and initialise a module:

```
mkdir inventory-system && cd inventory-system
go mod init inventory-system
touch main.go
```

All your work goes in `main.go` for Parts 1–3. Part 4 will ask you to split into multiple files.

---

## Part 1 — Define the core type and value-receiver methods

### Goal
Model a single inventory item as a struct and give it basic behavior using methods.

### Requirements

1. Define a struct:
   ```go
   type Item struct {
       SKU      string
       Name     string
       Price    float64 // price per unit, in dollars
       Quantity int
   }
   ```
2. Write a method with a **value receiver** that computes the total value of an item (`Price * Quantity`):
   ```go
   func (i Item) TotalValue() float64
   ```
3. Write a method with a **value receiver** that returns a human-readable summary string, e.g. `"Widget (SKU: W-001): 12 units @ $4.50"`:
   ```go
   func (i Item) String() string
   ```
4. In `main`, construct two or three `Item` values, print their `String()` output, and print their `TotalValue()`.

### Test input/output
```go
i := Item{SKU: "W-001", Name: "Widget", Price: 4.50, Quantity: 12}
```
- `i.TotalValue()` → `54`
- `i.String()` → `"Widget (SKU: W-001): 12 units @ $4.50"`

### Questions to think about
- You implemented `String()` — what interface does Go's standard library define that this satisfies, and what changes if you pass an `Item` directly to `fmt.Println`?
- If `TotalValue` had a pointer receiver instead, would anything in this part actually break? Why or why not at this stage?

---

## Part 2 — Build an inventory collection with pointer-receiver methods

### Goal
Model the whole inventory as a collection type, and use pointer receivers for any method that mutates state.

### Requirements

1. Define a struct that wraps the collection:
   ```go
   type Inventory struct {
       Items map[string]*Item // keyed by SKU
   }
   ```
   Think carefully about why this is `map[string]*Item` rather than `map[string]Item` before you write the methods below.

2. Write a constructor function:
   ```go
   func NewInventory() *Inventory
   ```
   It should return a ready-to-use `Inventory` with an initialised (non-nil) map.

3. Write a method with a **pointer receiver** to add a brand-new item:
   ```go
   func (inv *Inventory) AddItem(item Item) error
   ```
   It should return an error if an item with that SKU already exists, rather than silently overwriting it.

4. Write a method with a **pointer receiver** to increase stock on an existing SKU:
   ```go
   func (inv *Inventory) Restock(sku string, amount int) error
   ```
   It should return an error if the SKU doesn't exist, or if `amount` is negative or zero.

5. Write a method with a **pointer receiver** to decrease stock (a sale):
   ```go
   func (inv *Inventory) Sell(sku string, amount int) error
   ```
   It should return an error if the SKU doesn't exist, if `amount` is negative or zero, or if there isn't enough stock to fulfil the sale.

6. In `main`, create an `Inventory`, add a few items, call `Restock` and `Sell` on them, and print the resulting state.

### Test input/output
```go
inv := NewInventory()
inv.AddItem(Item{SKU: "W-001", Name: "Widget", Price: 4.50, Quantity: 10})
inv.Restock("W-001", 5)   // Quantity becomes 15
inv.Sell("W-001", 20)     // should return an error, Quantity stays 15
inv.Sell("W-001", 15)     // Quantity becomes 0, no error
```

### Questions to think about
- Why does using `map[string]*Item` instead of `map[string]Item` matter for how `Restock` and `Sell` are implemented? What would go wrong (or get more awkward) with the non-pointer version?
- `AddItem` takes an `Item` by value, not a pointer — why is that the right choice here even though the other two methods deal in pointers internally?
- What is the Go rule about mixing value and pointer receivers on the same type, and does your `Item`/`Inventory` design risk violating it?

---

## Part 3 — Custom error types and reporting

### Goal
Replace generic `errors.New` calls with a custom error type so callers can distinguish *why* an operation failed, and add read-only reporting methods.

### Requirements

1. Define a custom error type:
   ```go
   type InventoryError struct {
       SKU string
       Op  string // e.g. "restock", "sell", "add"
       Msg string
   }
   ```
   Implement the `error` interface for it (i.e. give it an `Error() string` method) so it can be returned anywhere `error` is expected.

2. Update `AddItem`, `Restock`, and `Sell` from Part 2 to return `*InventoryError` (wrapped as `error`) instead of plain `errors.New` errors, with a distinct, descriptive `Msg` for each failure case.

3. In `main` (or a small test you write), use `errors.As` to detect an `*InventoryError` returned from one of these calls and print its `SKU` and `Op` fields specifically — not just the generic error string.

4. Write a read-only method (value or pointer receiver — decide which makes sense and be ready to justify it):
   ```go
   func (inv *Inventory) TotalValue() float64
   ```
   Sum of `TotalValue()` across all items.

5. Write another reporting method:
   ```go
   func (inv *Inventory) LowStock(threshold int) []*Item
   ```
   Returns a slice of items whose `Quantity` is at or below `threshold`, sorted by `Quantity` ascending.

### Test input/output
```go
inv := NewInventory()
inv.AddItem(Item{SKU: "W-001", Name: "Widget", Price: 4.50, Quantity: 3})
inv.AddItem(Item{SKU: "G-002", Name: "Gadget", Price: 9.00, Quantity: 1})
inv.AddItem(Item{SKU: "T-003", Name: "Thingamajig", Price: 2.00, Quantity: 50})

err := inv.Sell("Z-999", 1) // unknown SKU
// err should be a *InventoryError with SKU="Z-999", Op="sell"

inv.LowStock(5) // should return G-002 (qty 1) then W-001 (qty 3), in that order
```

### Questions to think about
- Why is `errors.As` the right tool here instead of a type assertion like `err.(*InventoryError)` directly? What does `errors.As` give you that a bare assertion doesn't, especially if errors get wrapped later with `fmt.Errorf("...: %w", err)`?
- Your `InventoryError` has an `Error() string` method with a value receiver almost certainly — what would change about how it's used if you returned `InventoryError` (not `*InventoryError`) from your methods instead?
- In `LowStock`, you're returning `[]*Item` (pointers into the map's underlying storage). What's a potential danger of handing these pointers back to the caller, and how would changing to `[]Item` change that tradeoff?

---

## Part 4 — Split into files and add JSON persistence

### Goal
Reorganise the code into multiple files, and make the inventory's state survive across program runs by saving/loading it as JSON.

### Requirements

1. Split your single `main.go` into at least these files (same package, `package main`):
   - `item.go` — the `Item` struct and its methods
   - `inventory.go` — the `Inventory` struct, its methods, and `InventoryError`
   - `main.go` — just the `main` function and wiring

2. Add JSON struct tags to `Item` so that field names serialize in `snake_case` rather than Go's default `PascalCase`, e.g.:
   ```go
   SKU string `json:"sku"`
   ```
   Decide on tags for all four fields.

3. Write a method to serialize the inventory:
   ```go
   func (inv *Inventory) Save(path string) error
   ```
   Use `encoding/json` and `os.WriteFile`. Think about whether you want `json.Marshal` or `json.MarshalIndent` for a file a human might open and read.

4. Write a method to deserialize:
   ```go
   func (inv *Inventory) Load(path string) error
   ```
   Use `os.ReadFile` and `json.Unmarshal`. Think carefully about the shape you're unmarshalling into — `Inventory.Items` is a `map[string]*Item`, but JSON objects don't have a built-in notion of "this value is a pointer." Work out what intermediate type, if any, you need.

5. In `main`: build an inventory, add several items, `Save` it to `inventory.json`, then create a **second**, empty `Inventory` and `Load` the same file into it. Print both inventories' `TotalValue()` to confirm they match.

### Test input/output
After saving an inventory with one item (`SKU: "W-001", Name: "Widget", Price: 4.5, Quantity: 12`) and the JSON tags above, `inventory.json` should contain an object whose item entry looks like:
```json
{"sku": "W-001", "name": "Widget", "price": 4.5, "quantity": 12}
```
Loading that file back should produce an `Inventory` where `Items["W-001"].Quantity == 12`.

### Questions to think about
- `json.Unmarshal` needs a pointer to the destination — why? What actually happens if you accidentally pass a non-pointer `Inventory` value to it?
- If two Go files in the same directory both declare `package main`, what determines which one Go treats as the entry point, and does it matter which file `func main()` lives in?
- Why might `json.MarshalIndent` be preferable to `json.Marshal` for this specific file, given who might open it and why?

---

## Part 5 — An interface for reporting and a second implementation

### Goal
Extract a behavior into an interface so that `main` can work with "anything that reports inventory status" without caring about the concrete type underneath.

### Requirements

1. Define an interface:
   ```go
   type Reporter interface {
       Report() string
   }
   ```
2. Give `*Inventory` a `Report() string` method that produces a multi-line summary: total number of distinct SKUs, total units across all items, and total value (reuse `TotalValue` from Part 3). Decide whether this should be a pointer or value receiver and be consistent with your earlier choices.

3. Define a second, unrelated type that also satisfies `Reporter` — for example a `SimpleCounter` that just tracks a running count of "reports generated" and reports that count as a string. This type should have **nothing to do with `Item` or `Inventory`**; the point is to prove the interface is decoupled from your inventory code.

4. Write a function:
   ```go
   func printReport(r Reporter)
   ```
   It should accept *any* `Reporter` and print its `Report()` output. Call it once with your `*Inventory` and once with your `SimpleCounter`.

### Test input/output
Given an inventory with 3 SKUs and total quantity 65 across them, `Report()` should produce output along these lines (exact wording is up to you, but it must include these three figures):
```
Distinct SKUs: 3
Total units: 65
Total value: $XXX.XX
```

### Questions to think about
- `printReport` takes a `Reporter`, not an `*Inventory`. What's the concrete benefit of that choice here — what could you do with `printReport` that you couldn't do if its parameter type were `*Inventory`?
- If `*Inventory` satisfies `Reporter` but `Inventory` (the non-pointer type) does not, why? What's the rule governing which one "has" the method?
- Could `SimpleCounter` satisfy `Reporter` with a method defined in a completely different file, or even conceptually a different package that imports this interface? What does that tell you about how Go decides interface satisfaction, compared to languages with explicit `implements` keywords?

---

## Part 6 — Embedding and a specialized item type

### Goal
Use struct embedding to create a more specific kind of inventory item without duplicating fields, and make sure it still plays well with your existing `Inventory` machinery.

### Requirements

1. Define a new struct that embeds `Item`:
   ```go
   type PerishableItem struct {
       Item
       ExpiryDate string // format: "2006-01-02"
   }
   ```
2. Give `PerishableItem` its own method:
   ```go
   func (p PerishableItem) IsExpired(today string) bool
   ```
   Compare `ExpiryDate` to `today` (both as `"2006-01-02"` strings — you can compare lexically, or use `time.Parse` if you want practice with the `time` package).

3. Confirm — by actually calling it in `main`, not just by reasoning about it — that `PerishableItem` has access to `Item`'s methods (`TotalValue`, `String`) without you writing any forwarding code.

4. Your existing `Inventory.Items` is `map[string]*Item`, which can't directly hold `*PerishableItem` values. Add a **second** map to `Inventory` for perishables:
   ```go
   Perishables map[string]*PerishableItem
   ```
   along with an `AddPerishable` method mirroring `AddItem`'s validation logic (no duplicate SKUs — and decide whether a SKU should be allowed to exist in *both* maps, or whether your `AddItem`/`AddPerishable` should check across both and reject collisions).

5. Update `Report()` from Part 5 to also include perishable count and how many are currently expired, given a `today` string passed in.

### Test input/output
```go
p := PerishableItem{
    Item:       Item{SKU: "M-010", Name: "Milk", Price: 3.20, Quantity: 8},
    ExpiryDate: "2026-06-01",
}
p.IsExpired("2026-06-17") // true
p.TotalValue()            // 25.6 — inherited from Item, no override written
```

### Questions to think about
- `PerishableItem` embeds `Item` by value, not `*Item`. What would change about method promotion and mutation semantics if you'd embedded `*Item` instead?
- If you gave `PerishableItem` its *own* `String()` method, what would calling `fmt.Println(p)` print — the embedded `Item.String()` or `PerishableItem`'s own? Why does Go resolve it that way?
- You now have two near-identical "add with duplicate-SKU check" methods (`AddItem`, `AddPerishable`). What would it take to unify them, and is an interface or generics the more natural tool for that — and why might you choose to leave them separate anyway for now?

---

## Stretch Goals

If you finish early, try one or more of these:

**A. Transaction log**
Add a method `Sell` and `Restock` should append a record (timestamp, SKU, op, amount) to a slice of transaction structs on `Inventory`. Add a `History() []Transaction` method and a way to filter history by SKU.

**B. Generic `Filter` helper**
Write a generic function `Filter[T any](items []T, predicate func(T) bool) []T` and use it to reimplement `LowStock` from Part 3 in terms of it.

**C. Concurrency-safe wrapper (later, not now)**
You were asked to keep this exercise sequential — but as a forward-looking stretch, sketch (in comments only, no code) where you'd need a `sync.Mutex` if `Restock`/`Sell` were called from multiple goroutines, without actually implementing it.

**D. CLI with the `flag` package**
Wire up `-load=inventory.json`, `-save=inventory.json`, and subcommands-via-`flag.Args()` like `restock W-001 5` so the program is usable from a shell rather than only from hardcoded `main` logic.

**E. Validation on `Price`**
Add a custom `UnmarshalJSON` method on `Item` that rejects negative prices or quantities during `Load`, returning a descriptive error instead of silently accepting bad data.

---

## Checklist before you consider it done

- [x] `Item` methods use value receivers; `Inventory` mutating methods use pointer receivers, and you can justify each choice
- [ ] `AddItem`/`AddPerishable` reject duplicate SKUs without panicking or overwriting
- [ ] `Restock`/`Sell` reject invalid amounts and insufficient stock, returning `*InventoryError` via the `error` interface
- [ ] `errors.As` is used at least once to recover structured error fields, not just `err.Error()` string-matching
- [ ] JSON round-trip (`Save` then `Load` into a fresh `Inventory`) preserves all field values exactly
- [ ] `Reporter` interface is satisfied by two unrelated types and exercised through `printReport`
- [ ] `PerishableItem` correctly promotes `Item`'s methods with zero forwarding code
- [ ] `go vet ./...` reports no issues
- [ ] `gofmt -l .` reports no unformatted files

---

## Packages you are allowed to use

`fmt`, `os`, `errors`, `encoding/json`, `sort`, `time` (Part 6 only)

No third-party libraries. The standard library is more than enough.
