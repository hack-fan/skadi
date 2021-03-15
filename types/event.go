package types

import (
	"context"
	"time"
)

const (
	EventTypeInfo    = "info"
	EventTypeWarning = "warning"
)

type Event struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type EventInput struct {
	Message string `json:"message"`
}

type EventHandler func(e *Event) error

type EventCenter interface {
	// Pub publish a event to a queue or pool
	Pub(e *Event) error
	// Get a event, nil if no event found
	Get() (*Event, error)
	// StartWorker to check and get new event periodically in background
	StartWorker(ctx context.Context, handler EventHandler)
}
