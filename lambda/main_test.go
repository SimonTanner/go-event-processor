package main

import (
	"context"
	"testing"
	"time"

	"github.com/SimonTanner/go-event-processor/lambda/types"
	"github.com/aws/aws-lambda-go/events"
)

type MockPersistenceLayer struct {
	PersistCalls []map[string]interface{}
}

func (p MockPersistenceLayer) Persist(ctx context.Context, msg types.Message) error {
	p.PersistCalls = append(p.PersistCalls, map[string]interface{}{
		"ctx":     ctx,
		"message": msg,
	})
	return nil
}

func TestHandler(t *testing.T) {
	TS := time.Now()

	var tests = []struct {
		name        string
		time        time.Time
		eventData   string
		shouldError bool
	}{
		{
			name: "valid data",
			time: TS,
			eventData: `{
			"ID":"0124e053-3580-7000-b158-3401dd4f2d37",
			"Source": "application",
			"Client": "application",
			"Type": "transaction",
			"Event": {"TransactionID":"0124e053-3580-7000-a762-0502e4a1022e","FraudScore":20,"CustomerID":"0124e053-3580-7000-a762-0502e4a1022e","AccountID":"0124e053-3580-7000-a762-0502e4a1022e","Time":"2009-11-10T23:00:00Z"}}`,
			shouldError: false,
		},
		{
			name:        "invalid ID",
			time:        TS,
			eventData:   `{"ID":"0124e053-3580-7000-b158-3401dd4f2d3", "Source": "application", "Client": "application", "Type": "transaction"}`,
			shouldError: true,
		},
		{
			name:        "invalid Type",
			time:        TS,
			eventData:   `{"ID":"0124e053-3580-7000-b158-3401dd4f2d37", "Source": "application", "Client": "application", "Type": "blah"}`,
			shouldError: true,
		},
		{
			name:        "invalid Source",
			time:        TS,
			eventData:   `{"ID":"0124e053-3580-7000-b158-3401dd4f2d37", "Source": "foo", "Client": "application", "Type": "fraud_detection"}`,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			events := events.KinesisEvent{
				Records: []events.KinesisEventRecord{
					{
						Kinesis: events.KinesisRecord{
							ApproximateArrivalTimestamp: events.SecondsEpochTime{
								Time: TS,
							},
							Data: []byte(tt.eventData),
						},
					},
				},
			}

			pl := MockPersistenceLayer{}

			h := NewHandler(pl)

			err := h.handler(ctx, events)

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
