package paymentservice

import "errors"

var (
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentNotSuccessful  = errors.New("payment is not successful")
	ErrInvalidRefundAmount   = errors.New("invalid refund amount")
	ErrRefundAlreadyComplete = errors.New("payment is already fully refunded")
)

type PaymentService struct {
	payment  *Payment
	strategy PaymentStrategy
}

func NewPaymentService(strategy PaymentStrategy) *PaymentService {
	return &PaymentService{
		strategy: strategy,
	}
}

func (ps *PaymentService) ProcessPayment(
	amount int64,
	method string,
	orderID string,
) (*Payment, error) {

	if ps.payment != nil {
		return nil, errors.New("payment already exists")
	}

	if amount <= 0 {
		return nil, errors.New("invalid payment amount")
	}

	payment := NewPayment(amount, method, orderID)

	transactionID, err := ps.strategy.Pay(amount)
	if err != nil {
		payment.SetStatus(PaymentFailed)
		ps.payment = payment
		return payment, err
	}

	payment.SetTransactionID(transactionID)
	payment.SetStatus(PaymentCompleted)

	ps.payment = payment

	return payment, nil
}

func (ps *PaymentService) RefundPayment(amount int64) error {
	if ps.payment == nil {
		return ErrPaymentNotFound
	}

	if ps.payment.GetStatus() != PaymentCompleted &&
		ps.payment.GetStatus() != RefundPending &&
		ps.payment.GetStatus() != RefundFailed {
		return ErrPaymentNotSuccessful
	}

	remainingAmount := ps.payment.GetAmount() - ps.payment.GetRefundAmount()

	if amount <= 0 || amount > remainingAmount {
		return ErrInvalidRefundAmount
	}

	if remainingAmount == 0 {
		return ErrRefundAlreadyComplete
	}

	ps.payment.SetStatus(RefundPending)

	err := ps.strategy.Refund(
		ps.payment.GetTransactionID(),
		amount,
	)

	if err != nil {
		ps.payment.SetStatus(RefundFailed)
		return err
	}

	ps.payment.AddRefundAmount(amount)

	if ps.payment.GetRefundAmount() == ps.payment.GetAmount() {
		ps.payment.SetStatus(RefundCompleted)
	} else {
		ps.payment.SetStatus(PaymentCompleted)
	}

	return nil
}

func (ps *PaymentService) GetPayment() (*Payment, error) {
	if ps.payment == nil {
		return nil, ErrPaymentNotFound
	}

	return ps.payment, nil
}

func (ps *PaymentService) SetStrategy(strategy PaymentStrategy) {
	ps.strategy = strategy
}

func (ps *PaymentService) ClearPayment() {
	ps.payment = nil
}
