package main

import (
	"errors"
	"fmt"

	inventoryservice "vendingmachine/InvetoryService"
	"vendingmachine/cartservice"
	dispenseservice "vendingmachine/despenseService"
	"vendingmachine/ordderservice"
	"vendingmachine/paymentservice"
	"vendingmachine/vendingmachine"
)

// DemoPaymentStrategy lets main demonstrate payment/refund failures.
type DemoPaymentStrategy struct {
	failPayment     bool
	failFirstRefund bool
	refundAttempts  int
}

func (s *DemoPaymentStrategy) Pay(amount int64) (string, error) {
	if s.failPayment {
		return "", errors.New("simulated payment declined")
	}
	return "demo-transaction-001", nil
}

func (s *DemoPaymentStrategy) Refund(transactionID string, amount int64) error {
	s.refundAttempts++
	if s.failFirstRefund && s.refundAttempts == 1 {
		return errors.New("simulated refund provider timeout")
	}
	fmt.Printf("Refund of %d paise successful (transaction %s)\n", amount, transactionID)
	return nil
}

// DemoDispenseController can simulate a partial dispense for a selected slot.
type DemoDispenseController struct {
	partialSlot string
}

func (d *DemoDispenseController) Dispense(slotID string, quantity int) (int, error) {
	if slotID == d.partialSlot {
		dispensed := quantity - 1
		if dispensed < 0 {
			dispensed = 0
		}
		fmt.Printf("Partially dispensed %d of %d products from %s\n", dispensed, quantity, slotID)
		return dispensed, errors.New("simulated dispensing jam")
	}
	fmt.Printf("Dispensed %d products from %s\n", quantity, slotID)
	return quantity, nil
}

func buildMachine(paymentStrategy paymentservice.PaymentStrategy, controller dispenseservice.DispenseController) (*vendingmachine.VendingMachine, *inventoryservice.InventoryService) {
	// Seed products and physical stock in two inventory slots.
	water := inventoryservice.NewProduct("P-001", "Water Bottle", "WATER-500", 2000) // ₹20.00
	chips := inventoryservice.NewProduct("P-002", "Potato Chips", "CHIPS-50", 3000)  // ₹30.00

	slot1 := inventoryservice.NewSlot("SLOT-1", "FLOOR-1", water, 10)
	slot2 := inventoryservice.NewSlot("SLOT-2", "FLOOR-1", chips, 8)
	floor := inventoryservice.NewFloor("FLOOR-1", []*inventoryservice.Slot{slot1, slot2})
	inventory := inventoryservice.NewInventory([]*inventoryservice.Floor{floor})
	inventorySvc := inventoryservice.NewInventoryService(inventory)

	cartSvc := cartservice.NewCartService(nil, inventorySvc)
	orderSvc := orderservice.NewOrderService()
	paymentSvc := paymentservice.NewPaymentService(paymentStrategy)
	dispenseSvc := dispenseservice.NewDispenseService(controller)

	machine := vendingmachine.NewVendingMachine(
		cartSvc, inventorySvc, orderSvc, paymentSvc, dispenseSvc,
	)
	return machine, inventorySvc
}

func printStock(inv *inventoryservice.InventoryService) {
	for _, id := range []string{"SLOT-1", "SLOT-2"} {
		slot := inv.FindSlot(id)
		fmt.Printf("%s: physical=%d reserved=%d available=%d\n",
			id, slot.GetProductQuantity(), slot.GetReservedQuantity(), slot.GetAvailableQuantity())
	}
}

func main() {
	fmt.Println("\n========== SCENARIO 1: SUCCESSFUL PURCHASE ==========")
	strategy1 := &DemoPaymentStrategy{}
	machine1, inventory1 := buildMachine(strategy1, dispenseservice.NewMockDispenseController())
	must(machine1.StartSession("customer-101"))
	must(machine1.AddProduct("SLOT-1", 2))
	must(machine1.AddProduct("SLOT-2", 1))
	must(machine1.Checkout())
	fmt.Println("Checkout completed.")
	printStock(inventory1)

	fmt.Println("\n========== SCENARIO 2: PAYMENT FAILURE ==========")
	strategy2 := &DemoPaymentStrategy{failPayment: true}
	machine2, inventory2 := buildMachine(strategy2, dispenseservice.NewMockDispenseController())
	must(machine2.StartSession("customer-102"))
	must(machine2.AddProduct("SLOT-1", 2))
	if err := machine2.Checkout(); err != nil {
		fmt.Println("Expected checkout error:", err)
	}
	fmt.Println("Reservation released; customer remains in READY state and can retry.")
	printStock(inventory2)
	fmt.Println("Switching to UPI strategy and retrying the same cart...")
	must(machine2.SetPaymentMethod("UPI"))
	must(machine2.Checkout())
	printStock(inventory2)

	fmt.Println("\n========== SCENARIO 3: PARTIAL DISPENSE + REFUND RETRY ==========")
	strategy3 := &DemoPaymentStrategy{failFirstRefund: true}
	controller3 := &DemoDispenseController{partialSlot: "SLOT-1"}
	machine3, inventory3 := buildMachine(strategy3, controller3)
	must(machine3.StartSession("customer-103"))
	must(machine3.AddProduct("SLOT-1", 2))
	must(machine3.AddProduct("SLOT-2", 1))
	if err := machine3.Checkout(); err != nil {
		fmt.Println("Expected refund-pending result:", err)
	}
	fmt.Println("Stock after partial dispense; failed unit is released from reservation:")
	printStock(inventory3)
	fmt.Println("Retrying refund...")
	must(machine3.RetryRefund())
	fmt.Println("Refund resolved; machine returned to IDLE.")

	fmt.Println("\n========== SCENARIO 4: CUSTOMER CANCELS SESSION ==========")
	strategy4 := &DemoPaymentStrategy{}
	machine4, inventory4 := buildMachine(strategy4, dispenseservice.NewMockDispenseController())
	must(machine4.StartSession("customer-104"))
	must(machine4.AddProduct("SLOT-2", 2))
	must(machine4.CancelSession())
	fmt.Println("Session cancelled; cart cleared; no stock was reserved or deducted.")
	printStock(inventory4)

	fmt.Println("\n========== SCENARIO 5: INVALID OPERATION IN IDLE ==========")
	if err := machine4.AddProduct("SLOT-1", 1); err != nil {
		fmt.Println("Expected state validation error:", err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
