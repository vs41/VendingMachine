package customer

import (
	"errors"

	"vendingmachine/utiles"
)

var (
	ErrSessionAlreadyCompleted = errors.New("session already completed")
	ErrSessionAlreadyCancelled = errors.New("session already cancelled")
	ErrInvalidUserID           = errors.New("invalid user ID")
)

type CustomerSession struct {
	sessionID string
	userID    string
	status    SessionStatus
}

func NewCustomerSession(userID string) (*CustomerSession, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	return &CustomerSession{
		sessionID: utiles.UniqueIdGeneretor(),
		userID:    userID,
		status:    SessionActive,
	}, nil
}

// Getters

func (cs *CustomerSession) GetSessionID() string {
	return cs.sessionID
}

func (cs *CustomerSession) GetUserID() string {
	return cs.userID
}

func (cs *CustomerSession) GetStatus() SessionStatus {
	return cs.status
}

func (cs *CustomerSession) IsActive() bool {
	return cs.status == SessionActive
}

// Complete the customer session.
func (cs *CustomerSession) Complete() error {
	if cs.status == SessionCompleted {
		return ErrSessionAlreadyCompleted
	}

	if cs.status == SessionCancelled {
		return ErrSessionAlreadyCancelled
	}

	cs.status = SessionCompleted
	return nil
}

// Cancel the customer session.
func (cs *CustomerSession) Cancel() error {
	if cs.status == SessionCompleted {
		return ErrSessionAlreadyCompleted
	}

	if cs.status == SessionCancelled {
		return ErrSessionAlreadyCancelled
	}

	cs.status = SessionCancelled
	return nil
}
