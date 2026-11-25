# Vending Machine LLD (Go)

## Run

```bash
go run .
```

Run package checks:

```bash
go test ./...
```

## Included scenarios in `main.go`

1. Successful purchase of multiple products.
2. Payment failure, inventory reservation release, then payment retry.
3. Partial dispensing, failed refund, and successful refund retry.
4. Customer cancels a session before checkout.
5. Invalid operation while the machine is idle (State Pattern validation).

## Inventory setup

`main.go` seeds a floor with two physical slots:

- `SLOT-1`: Water Bottle, 10 units, ₹20.00 each.
- `SLOT-2`: Potato Chips, 8 units, ₹30.00 each.

All prices are stored as paise (`int64`) to avoid floating-point currency calculations.

## Main responsibilities

- `vendingmachine`: Context, state transitions, and transaction orchestration.
- `VendingState.go`: State interface, base state, and concrete states.
- `InvetoryService`: Products, floors, slots, stock, reservations, and settlement.
- `cartservice`: Customer cart and availability checks.
- `ordderservice`: Order snapshot and order lifecycle.
- `paymentservice`: Payment strategy and refund tracking.
- `despenseService`: Dispensing controller and per-slot dispensing results.
- `userservice`: Customer session lifecycle.

The payment and dispensing implementations in `main.go` are test doubles for interview/demo scenarios; they are not integrations with real payment gateways or vending-machine hardware.
