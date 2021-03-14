package types

import "time"

type Event struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type EventCenter interface {
	// Pub publish a event to a queue or pool
	Pub(e *Event) error
	// Get a event, nil if no event found
	Get() (*Event, error)
}
