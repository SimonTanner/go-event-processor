package types

import (
	"time"

	"github.com/google/uuid"
)

type Source string

const (
	SourceMonitoring  Source = "monitoring"
	SourceApplication Source = "application"
	SourceAuthorizer  Source = "authorizer"
)

type Message struct {
	ID        uuid.UUID `json:"ID"`
	Source    Source    `json:"Source"`
	Timestamp time.Time `json:"Timestamp"`
	Message   string    `json:"Message"`
	Client    string    `json:"Client"`
}
