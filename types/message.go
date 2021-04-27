package types

import (
	"context"
	"time"
)

const (
	MessageTypeInfo    = "info"
	MessageTypeWarning = "warning"
	MessageTypeText    = "text"
)

// Message will send to IM, gen by agent api, or system event.
type Message struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageInput struct {
	Message string `json:"message"`
}

// EventHandler handle the event, and please handle the error yourself.
type EventHandler func(e *Message)

type EventCenter interface {
	// Pub publish a event to a queue or pool
	Pub(e *Message) error
	// Get a event, nil if no event found
	Get() (*Message, error)
	// StartWorker to check and get new event periodically in background
	StartWorker(ctx context.Context, handler EventHandler)
}
