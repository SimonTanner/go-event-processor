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
	ID        uuid.UUID `json:"ID" validate:"required"`
	Source    Source    `json:"Source" validate:"required"`
	Timestamp time.Time `json:"Timestamp" validate:"required"`
	// Event     Event     `json:"Event" validate:"required"`
	Client string `json:"Client" validate:"required"`
}
