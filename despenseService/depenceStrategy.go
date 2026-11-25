package dispenseservice

import (
	"errors"
	"fmt"
)

var (
	ErrDispenseFailed = errors.New("product dispensing failed")
)

type DispenseController interface {
	Dispense(slotID string, quantity int) (int, error)
}

// Mock implementation for LLD interview

type MockDispenseController struct{}

func NewMockDispenseController() *MockDispenseController {
	return &MockDispenseController{}
}

func (dc *MockDispenseController) Dispense(
	slotID string,
	quantity int,
) (int, error) {

	if slotID == "" || quantity <= 0 {
		return 0, errors.New("invalid dispensing request")
	}

	// Simulate successful dispensing.
	fmt.Printf("Dispensing %d products from slot %s\n", quantity, slotID)

	return quantity, nil
}
