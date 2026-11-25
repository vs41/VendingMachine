package vendingmachine

import (
	"errors"
	"fmt"

	inventoryservice "vendingmachine/InvetoryService"
	"vendingmachine/cartservice"
	dispenseservice "vendingmachine/despenseService"
	orderservice "vendingmachine/ordderservice"
	"vendingmachine/paymentservice"
	customer "vendingmachine/userservice"
)

var (
	ErrInvalidOperation     = errors.New("operation not allowed in current state")
	ErrSessionAlreadyActive = errors.New("customer session already active")
	ErrNoActiveSession      = errors.New("no active customer session")
	ErrEmptyCart            = errors.New("cart is empty")
	ErrNoPendingRefund      = errors.New("no pending refund")
	ErrMachineRecovery      = errors.New("machine requires manual recovery")
)

type VendingMachine struct {
	state          VendingStates
	currentSession *customer.CustomerSession

	pendingOrder        *orderservice.Order
	pendingRefundAmount int64
	paymentMethod       string

	cartService      *cartservice.CartService
	inventoryService *inventoryservice.InventoryService
	orderService     *orderservice.OrderService
	paymentService   *paymentservice.PaymentService
	dispenseService  *dispenseservice.DispenseService
}

func NewVendingMachine(
	cartService *cartservice.CartService,
	inventoryService *inventoryservice.InventoryService,
	orderService *orderservice.OrderService,
	paymentService *paymentservice.PaymentService,
	dispenseService *dispenseservice.DispenseService,
) *VendingMachine {
	machine := &VendingMachine{
		cartService:      cartService,
		inventoryService: inventoryService,
		orderService:     orderService,
		paymentService:   paymentService,
		dispenseService:  dispenseService,
		paymentMethod:    "UPI",
	}
	machine.setState(NewIdleState(machine))
	return machine
}

func (v *VendingMachine) SetPaymentMethod(method string) error {
	if method != "UPI" && method != "Card" && method != "Cash" {
		return fmt.Errorf("unsupported payment method: %s", method)
	}
	v.paymentMethod = method
	v.paymentService.SetStrategy(paymentservice.NewPaymentStrategy(method))
	return nil
}

func (v *VendingMachine) setState(state VendingStates) { v.state = state }

// Public API: all customer operations pass through the current state.
func (v *VendingMachine) StartSession(userID string) error { return v.state.StartSession(userID) }
func (v *VendingMachine) AddProduct(slotID string, quantity int) error {
	return v.state.AddProduct(slotID, quantity)
}
func (v *VendingMachine) RemoveProduct(itemID string) error {
	return v.state.RemoveProduct(itemID)
}
func (v *VendingMachine) Checkout() error      { return v.state.Checkout() }
func (v *VendingMachine) CancelSession() error { return v.state.Cancel() }
func (v *VendingMachine) RetryRefund() error   { return v.state.RetryRefund() }

func (v *VendingMachine) startSession(userID string) error {
	if v.currentSession != nil {
		return ErrSessionAlreadyActive
	}
	session, err := customer.NewCustomerSession(userID)
	if err != nil {
		return err
	}
	if err := v.cartService.StartCart(userID); err != nil {
		return err
	}
	v.currentSession = session
	v.setState(NewReadyState(v))
	fmt.Printf("Session started for customer %s\n", userID)
	return nil
}

func (v *VendingMachine) addProduct(slotID string, quantity int) error {
	if v.currentSession == nil || !v.currentSession.IsActive() {
		return ErrNoActiveSession
	}
	return v.cartService.AddItemToCart(slotID, quantity)
}

func (v *VendingMachine) removeProduct(itemID string) error {
	if v.currentSession == nil || !v.currentSession.IsActive() {
		return ErrNoActiveSession
	}
	return v.cartService.RemoveItemFromCart(itemID)
}

func (v *VendingMachine) checkout() error {
	if v.currentSession == nil || !v.currentSession.IsActive() {
		return ErrNoActiveSession
	}

	cartItems, err := v.cartService.GetCartItems()
	if err != nil {
		return err
	}
	if len(cartItems) == 0 {
		return ErrEmptyCart
	}

	order, err := v.orderService.CreateOrder(cartItems)
	if err != nil {
		return err
	}
	v.pendingOrder = order

	reservations := make([]inventoryservice.ReservationItem, 0, len(order.GetItems()))
	for _, item := range order.GetItems() {
		reservations = append(reservations, inventoryservice.ReservationItem{
			SlotID: item.GetSlotID(), Quantity: item.GetQuantity(),
		})
	}
	if err := v.inventoryService.Reserve(reservations); err != nil {
		v.orderService.ClearOrder()
		v.pendingOrder = nil
		return err
	}
	if err := v.orderService.UpdateOrderState(orderservice.OrderCreated); err != nil {
		_ = v.inventoryService.Release(reservations)
		v.orderService.ClearOrder()
		v.pendingOrder = nil
		return err
	}

	v.setState(NewPaymentPendingState(v))
	payment, paymentErr := v.paymentService.ProcessPayment(
		order.GetTotalAmount(), v.paymentMethod, order.GetOrderID(),
	)
	if paymentErr != nil {
		releaseErr := v.inventoryService.Release(reservations)
		_ = v.orderService.UpdateOrderState(orderservice.OrderFailed)
		v.paymentService.ClearPayment()
		v.orderService.ClearOrder()
		v.pendingOrder = nil
		if releaseErr != nil {
			return fmt.Errorf("payment failed (%v) and inventory release failed: %w", paymentErr, releaseErr)
		}
		v.setState(NewReadyState(v))
		return paymentErr
	}
	if payment == nil {
		return errors.New("payment service returned an empty payment")
	}

	v.setState(NewDispensingState(v))
	dispenseResult, dispenseErr := v.dispenseService.DispenseOrder(order)
	if dispenseResult == nil {
		// Payment is captured but the machine cannot confirm the hardware result.
		// Keep the order/reservation details rather than falsely declaring success.
		return fmt.Errorf("%w: dispense result unavailable: %v", ErrMachineRecovery, dispenseErr)
	}

	settlements := make([]inventoryservice.DispenseSettlement, 0, len(dispenseResult.Items))
	for _, item := range dispenseResult.Items {
		settlements = append(settlements, inventoryservice.DispenseSettlement{
			SlotID:            item.SlotID,
			ReservedQuantity:  item.RequestedQuantity,
			DispensedQuantity: item.DispensedQuantity,
		})
	}
	if err := v.inventoryService.SettleDispensed(settlements); err != nil {
		return fmt.Errorf("%w: inventory settlement failed: %v", ErrMachineRecovery, err)
	}

	refundAmount := calculateRefundAmount(order, dispenseResult)
	if refundAmount > 0 {
		v.pendingRefundAmount = refundAmount
		if err := v.paymentService.RefundPayment(refundAmount); err != nil {
			v.setState(NewRefundPendingState(v))
			return fmt.Errorf("refund pending: %w", err)
		}
		if err := v.orderService.UpdateOrderState(orderservice.OrderPartiallyCompleted); err != nil {
			return err
		}
	} else {
		if dispenseErr != nil {
			// All units may be recorded as dispensed, but a controller error is still surfaced.
			fmt.Printf("Dispense controller reported: %v\n", dispenseErr)
		}
		if err := v.orderService.UpdateOrderState(orderservice.OrderCompleted); err != nil {
			return err
		}
	}

	v.completeSession()
	return nil
}

func calculateRefundAmount(order *orderservice.Order, result *dispenseservice.DispenseResult) int64 {
	failedBySlot := make(map[string]int)
	for _, item := range result.Items {
		failedBySlot[item.SlotID] += item.FailedQuantity
	}
	var refund int64
	for _, item := range order.GetItems() {
		refund += int64(failedBySlot[item.GetSlotID()]) * item.GetUnitPrice()
	}
	return refund
}

func (v *VendingMachine) cancelSession() error {
	if v.currentSession == nil || !v.currentSession.IsActive() {
		return ErrNoActiveSession
	}
	if err := v.currentSession.Cancel(); err != nil {
		return err
	}
	v.cartService.EndCart()
	v.currentSession = nil
	v.setState(NewIdleState(v))
	return nil
}

func (v *VendingMachine) retryRefund() error {
	if v.pendingOrder == nil || v.pendingRefundAmount <= 0 {
		return ErrNoPendingRefund
	}
	if err := v.paymentService.RefundPayment(v.pendingRefundAmount); err != nil {
		return fmt.Errorf("refund retry failed: %w", err)
	}
	if err := v.orderService.UpdateOrderState(orderservice.OrderPartiallyCompleted); err != nil {
		return err
	}
	v.completeSession()
	return nil
}

func (v *VendingMachine) completeSession() {
	if v.currentSession != nil && v.currentSession.IsActive() {
		_ = v.currentSession.Complete()
	}
	v.cartService.EndCart()
	v.orderService.ClearOrder()
	v.paymentService.ClearPayment()
	v.currentSession = nil
	v.pendingOrder = nil
	v.pendingRefundAmount = 0
	v.setState(NewIdleState(v))
}
