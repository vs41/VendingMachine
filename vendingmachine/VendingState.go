package vendingmachine

// State interface
type VendingStates interface {
	StartSession(userID string) error
	AddProduct(slotID string, quantity int) error
	RemoveProduct(itemID string) error
	Checkout() error
	Cancel() error
	RetryRefund() error
}

// Base state
type baseState struct {
	machine *VendingMachine
}

func (s *baseState) StartSession(userID string) error {
	return ErrInvalidOperation
}

func (s *baseState) AddProduct(slotID string, quantity int) error {
	return ErrInvalidOperation
}

func (s *baseState) RemoveProduct(itemID string) error {
	return ErrInvalidOperation
}

func (s *baseState) Checkout() error {
	return ErrInvalidOperation
}

func (s *baseState) Cancel() error {
	return ErrInvalidOperation
}

func (s *baseState) RetryRefund() error {
	return ErrInvalidOperation
}

// Idle State
type IdleState struct {
	baseState
}

func NewIdleState(machine *VendingMachine) *IdleState {
	return &IdleState{
		baseState: baseState{machine: machine},
	}
}

func (s *IdleState) StartSession(userID string) error {
	return s.machine.startSession(userID)
}

// Ready State
type ReadyState struct {
	baseState
}

func NewReadyState(machine *VendingMachine) *ReadyState {
	return &ReadyState{
		baseState: baseState{machine: machine},
	}
}

func (s *ReadyState) AddProduct(slotID string, quantity int) error {
	return s.machine.addProduct(slotID, quantity)
}

func (s *ReadyState) RemoveProduct(itemID string) error {
	return s.machine.removeProduct(itemID)
}

func (s *ReadyState) Checkout() error {
	return s.machine.checkout()
}

func (s *ReadyState) Cancel() error {
	return s.machine.cancelSession()
}

// Payment Pending State
type PaymentPendingState struct {
	baseState
}

func NewPaymentPendingState(machine *VendingMachine) *PaymentPendingState {
	return &PaymentPendingState{
		baseState: baseState{machine: machine},
	}
}

// Dispensing State
type DispensingState struct {
	baseState
}

func NewDispensingState(machine *VendingMachine) *DispensingState {
	return &DispensingState{
		baseState: baseState{machine: machine},
	}
}

// Refund Pending State
type RefundPendingState struct {
	baseState
}

func NewRefundPendingState(machine *VendingMachine) *RefundPendingState {
	return &RefundPendingState{
		baseState: baseState{machine: machine},
	}
}

func (s *RefundPendingState) RetryRefund() error {
	return s.machine.retryRefund()
}
