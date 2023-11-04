package event

import "time"

type TransactionCreated struct {
	Name    string
	Payload interface{}
}

func NewTransactionCreated() *TransactionCreated {
	return &TransactionCreated{
		Name: "TransactionCreated",
	}
}

func NewBalanceUpdated() *TransactionCreated {
	return &TransactionCreated{
		Name: "BalancedUpdated",
	}
}

func (e *TransactionCreated) GetName() string {
	return e.Name
}

func (e *TransactionCreated) GetPayload() interface{} {
	return e.Payload
}

// SetPayload implements events.EventInterface.
func (e *TransactionCreated) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *TransactionCreated) GetDateTime() time.Time {
	return time.Now()
}
