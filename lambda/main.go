package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/SimonTanner/go-event-processor/lambda/persist"
	"github.com/SimonTanner/go-event-processor/lambda/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Handler struct {
	PersistenceLayer persist.IPersistenceLayer
}

func NewHandler(persistenceLayer persist.IPersistenceLayer) Handler {
	return Handler{
		PersistenceLayer: persistenceLayer,
	}
}

func (h *Handler) handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	if len(kinesisEvent.Records) == 0 {
		log.Printf("Kinesis event record length is 0, exiting")
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
		validate.RegisterStructValidation(types.MessageStructLevelValidation, types.Message{})

		// Validate the Message
		err = validate.Struct(message)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			log.Printf("errors validating struct %s", errors)
			return errors
		}

		err = h.PersistenceLayer.Persist(ctx, message)
		if err != nil {
			log.Printf("errors saving message to DB %s", err)
			return err
		}
	}

	log.Printf("successfully processed %v records", len(kinesisEvent.Records))
	return nil
}

func main() {
	region := os.Getenv("AWS_REGION")
	tableName := os.Getenv("DYNAMO_TABLE_NAME")
	if tableName == "" {
		tableName = "GoEvents"
	}
	log.Printf("creating dynamodb client for table: %s, in region: %s", tableName, region)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Printf("errors loading aws config: %s", err)
		return
	}

	pl := persist.NewPersistenceLayer(cfg, tableName)

	h := NewHandler(pl)
	lambda.Start(h.handler)
}
