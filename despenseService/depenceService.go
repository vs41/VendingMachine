package dispenseservice

import (
	"errors"
	orderservice "vendingmachine/ordderservice"
)

type ItemDispenseResult struct {
	SlotID            string
	RequestedQuantity int
	DispensedQuantity int
	FailedQuantity    int
}

type DispenseResult struct {
	Items          []ItemDispenseResult
	TotalRequested int
	TotalDispensed int
	TotalFailed    int
}

type DispenseService struct {
	controller DispenseController
}

var (
	ErrControllerNotFound    = errors.New("dispense controller not configured")
	ErrInvalidDispenseResult = errors.New("invalid dispense controller result")
)

func NewDispenseService(controller DispenseController) *DispenseService {
	return &DispenseService{
		controller: controller,
	}
}

func (ds *DispenseService) DispenseOrder(
	order *orderservice.Order,
) (*DispenseResult, error) {

	if ds.controller == nil {
		return nil, ErrControllerNotFound
	}

	if order == nil {
		return nil, errors.New("order is nil")
	}

	result := &DispenseResult{}

	for _, item := range order.GetItems() {
		if item == nil || item.GetQuantity() <= 0 || item.GetSlotID() == "" {
			return result, errors.New("invalid order item")
		}

		requested := item.GetQuantity()

		dispensed, err := ds.controller.Dispense(
			item.GetSlotID(),
			requested,
		)

		if dispensed < 0 || dispensed > requested {
			return result, ErrInvalidDispenseResult
		}

		failed := requested - dispensed

		itemResult := ItemDispenseResult{
			SlotID:            item.GetSlotID(),
			RequestedQuantity: requested,
			DispensedQuantity: dispensed,
			FailedQuantity:    failed,
		}

		result.Items = append(result.Items, itemResult)
		result.TotalRequested += requested
		result.TotalDispensed += dispensed
		result.TotalFailed += failed

		if err != nil {
			// Continue with other items so we can collect
			// all confirmed dispensing results.
			continue
		}
	}

	return result, nil
}
