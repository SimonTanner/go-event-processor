package types

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func TestValidation(t *testing.T) {
	UUID, uuidErr := uuid.NewV7()
	if uuidErr != nil {
		t.Fatalf("unable to generate UUID %s", uuidErr)
	}

	var tests = []struct {
		name        string
		message     Message
		shouldError bool
	}{
		{
			name: "valid transaction message",
			message: Message{
				ID:        UUID,
				Source:    SourceApplication,
				Timestamp: time.Now(),
				Event: Event{
					SharedData: SharedData{
						CustomerID: UUID,
						AccountID:  UUID,
						Time:       time.Now(),
					},
					TransactionID: UUID,
				},
				Client: "Some Company",
				Type:   TypeTransaction,
			},
			shouldError: false,
		},
		{
			name: "should error if the event doesn't contain customer details",
			message: Message{
				ID:        UUID,
				Source:    SourceApplication,
				Timestamp: time.Now(),
				Event: Event{
					SharedData:    SharedData{},
					TransactionID: UUID,
				},
				Client: "Some Company",
				Type:   TypeFraudDetection,
			},
			shouldError: true,
		},
		{
			name: "should error if the message type and event don't match",
			message: Message{
				ID:        UUID,
				Source:    SourceApplication,
				Timestamp: time.Now(),
				Event: Event{
					SharedData: SharedData{
						CustomerID: UUID,
						AccountID:  UUID,
						Time:       time.Now(),
					},
					TransactionID: UUID,
				},
				Client: "Some Company",
				Type:   TypeFraudDetection,
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := validator.New()
			validate.RegisterStructValidation(MessageStructLevelValidation, Message{})
			err := validate.Struct(tt.message)
			if !tt.shouldError {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
			} else {
				if err == nil {
					t.Error("expected an error to be raised")
				}
			}
		})

	}
}
