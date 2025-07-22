package main

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	TS := time.Now()

	var tests = []struct {
		name        string
		time        time.Time
		eventData   string
		shouldError bool
	}{
		{
			name:        "valid data",
			time:        TS,
			eventData:   `{"ID":"0124e053-3580-7000-b158-3401dd4f2d37", "Source": "application", "Client": "application", "Type": "transaction"}`,
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

			err := handler(ctx, events)

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
