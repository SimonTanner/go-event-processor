package types

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Source string

const (
	SourceMonitoring  Source = "monitoring"
	SourceApplication Source = "application"
	SourceAuthorizer  Source = "authorizer"
)

type EventType string

const (
	TypeTransaction    EventType = "transaction"
	TypeFraudDetection EventType = "fraud_detection"
	TypeCheckAccount   EventType = "check_account"
)

type Message struct {
	ID        uuid.UUID `json:"ID" validate:"required"`
	Source    Source    `json:"Source" validate:"required"`
	Timestamp time.Time `json:"Timestamp" validate:"required"`
	// Event     Event     `json:"Event" validate:"required"`
	Client string    `json:"Client" validate:"required"`
	Type   EventType `json:"Type" validate:"required"`
}

type SharedData struct {
	CustomerID uuid.UUID `json:"CustomerID" validate:"required"`
	AccountID  uuid.UUID `json:"AccountID" validate:"required"`
	Time       time.Time `json:"Time" validate:"required"`
}

type Event struct {
	SharedData
	TransactionID uuid.UUID `json:"TransactionID" validate:"required"`
	FraudScore    int       `json:"FraudScore" validate:"required"`
}

func MessageStructLevelValidation(sl validator.StructLevel) {

	message := sl.Current().Interface().(Message)

	if message.Type != TypeTransaction && message.Type != TypeFraudDetection && message.Type != TypeCheckAccount {
		sl.ReportError(message.Type, "Type", "Type", "", "")
	}
}
