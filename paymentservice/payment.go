package paymentservice

import "vendingmachine/utiles"

type PaymentState string

const (
	PaymentPending   PaymentState = "PENDING"
	PaymentFailed    PaymentState = "FAILED"
	PaymentCompleted PaymentState = "COMPLETED"
	RefundPending    PaymentState = "REFUND_PENDING"
	RefundCompleted  PaymentState = "REFUND_COMPLETED"
	RefundFailed     PaymentState = "REFUND_FAILED"
)

type Payment struct {
	id            string
	amount        int64
	method        string
	status        PaymentState
	refundAmount  int64
	orderID       string
	transactionID string
}

func NewPayment(amount int64, method, orderID string) *Payment {
	return &Payment{
		id:           utiles.UniqueIdGeneretor(),
		amount:       amount,
		method:       method,
		orderID:      orderID,
		status:       PaymentPending,
		refundAmount: 0,
	}
}

func (p *Payment) GetPaymentID() string {
	return p.id
}

func (p *Payment) GetPaymentMethod() string {
	return p.method
}

func (p *Payment) GetOrderID() string {
	return p.orderID
}

func (p *Payment) GetStatus() PaymentState {
	return p.status
}

func (p *Payment) GetAmount() int64 {
	return p.amount
}

func (p *Payment) GetRefundAmount() int64 {
	return p.refundAmount
}

func (p *Payment) SetTransactionID(id string) {
	p.transactionID = id
}

func (p *Payment) SetStatus(status PaymentState) {
	p.status = status
}

func (p *Payment) AddRefundAmount(amount int64) {
	p.refundAmount += amount
}

func (p *Payment) GetTransactionID() string {
	return p.transactionID
}
