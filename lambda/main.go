package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/SimonTanner/go-event-processor/lambda/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	if len(kinesisEvent.Records) == 0 {
		log.Printf("Kinesis event record lenght is 0")
		return nil
	}

	for _, record := range kinesisEvent.Records {
		log.Printf("processed Kinesis event with EventId: %v", record.EventID)
		recordDataBytes := record.Kinesis.Data
		recordDataText := string(recordDataBytes)
		log.Printf("record data: %v", recordDataText)

		message := types.Message{}
		err := json.Unmarshal(recordDataBytes, &message)
		if err != nil {
			log.Printf("error unmarshalling record: %v", err)
			return err
		}

		message.Timestamp = record.Kinesis.ApproximateArrivalTimestamp.Time
		log.Printf("successfully converted record to struct: %v", message)

		validate = validator.New()

		// Validate the Message
		err = validate.Struct(message)
		if err != nil {
			// Validation failed, handle the error
			errors := err.(validator.ValidationErrors)
			log.Printf("errors validating struct %s", errors)
			return errors
		}
	}

	log.Printf("successfully processed %v records", len(kinesisEvent.Records))
	return nil
}

func main() {
	lambda.Start(handler)
}
