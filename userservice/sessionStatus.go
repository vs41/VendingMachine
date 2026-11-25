package customer

type SessionStatus string

const (
	SessionActive    SessionStatus = "ACTIVE"
	SessionCompleted SessionStatus = "COMPLETED"
	SessionCancelled SessionStatus = "CANCELLED"
)
