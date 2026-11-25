package paymentservice

import (
	"errors"
	"fmt"

	"vendingmachine/utiles"
)

type PaymentStrategy interface {
	Pay(amount int64) (string, error)
	Refund(transactionID string, amount int64) error
}

// UPI Strategy

type UPI struct{}

func NewUPI() *UPI {
	return &UPI{}
}

func (u *UPI) Pay(amount int64) (string, error) {
	if amount <= 0 {
		return "", errors.New("invalid payment amount")
	}

	transactionID := utiles.UniqueIdGeneretor()
	fmt.Println("UPI payment successful")

	return transactionID, nil
}

func (u *UPI) Refund(transactionID string, amount int64) error {
	if transactionID == "" || amount <= 0 {
		return errors.New("invalid refund details")
	}

	fmt.Println("UPI refund successful")
	return nil
}

// Card Strategy

type Card struct{}

func NewCard() *Card {
	return &Card{}
}

func (c *Card) Pay(amount int64) (string, error) {
	if amount <= 0 {
		return "", errors.New("invalid payment amount")
	}

	transactionID := utiles.UniqueIdGeneretor()
	fmt.Println("Card payment successful")

	return transactionID, nil
}

func (c *Card) Refund(transactionID string, amount int64) error {
	if transactionID == "" || amount <= 0 {
		return errors.New("invalid refund details")
	}

	fmt.Println("Card refund successful")
	return nil
}

// Cash Strategy

type Cash struct{}

func NewCash() *Cash {
	return &Cash{}
}

func (c *Cash) Pay(amount int64) (string, error) {
	if amount <= 0 {
		return "", errors.New("invalid payment amount")
	}

	transactionID := utiles.UniqueIdGeneretor()
	fmt.Println("Cash payment successful")

	return transactionID, nil
}

func (c *Cash) Refund(transactionID string, amount int64) error {
	if transactionID == "" || amount <= 0 {
		return errors.New("invalid refund details")
	}

	fmt.Println("Cash refund successful")
	return nil
}

func NewPaymentStrategy(method string) PaymentStrategy {
	if method == "UPI" {
		return &UPI{}
	} else if method == "Card" {
		return &Card{}
	} else {
		return &Cash{}
	}
}
