package persist

import (
	"context"
	"encoding/json"
	"log"

	"github.com/SimonTanner/go-event-processor/lambda/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type EventEntry struct {
	Client     string `dynamo:"Client"`
	CustomerID string `dynamo:"CustomerID"`
	AccountID  string
	EventID    string
	Source     string
	Timestamp  int
	Type       string
	Payload    []byte
}

type IPersistenceLayer interface {
	Persist(context.Context, types.Message) error
}

type PersistenceLayer struct {
	DynamoClient *dynamodb.Client
	TableName    string
}

func NewPersistenceLayer(conf aws.Config, tableName string) PersistenceLayer {
	return PersistenceLayer{
		DynamoClient: dynamodb.NewFromConfig(conf),
		TableName:    tableName,
	}
}

func (p PersistenceLayer) Persist(ctx context.Context, msg types.Message) error {
	eventEntry, err := toDynamoModel(msg)
	if err != nil {
		log.Printf("error converting message to EventEntry: %s", err)
		return err
	}

	item, err := attributevalue.MarshalMap(eventEntry)
	if err != nil {
		log.Printf("error marshalling message: %s", err)
		return err
	}

	log.Printf("attempting to save data %v", item)
	_, err = p.DynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(p.TableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v", err)
		return err
	}

	return nil
}

func toDynamoModel(msg types.Message) (EventEntry, error) {
	payload, err := json.Marshal(msg.Event)
	if err != nil {
		return EventEntry{}, err
	}

	eventEntry := EventEntry{
		Client:     msg.Client,
		CustomerID: msg.Event.CustomerID.String(),
		AccountID:  msg.Event.AccountID.String(),
		EventID:    msg.ID.String(),
		Source:     *aws.String(string(msg.Source)),
		Timestamp:  int(msg.Timestamp.Unix()),
		Type:       *aws.String(string(msg.Type)),
		Payload:    payload,
	}

	return eventEntry, nil
}
