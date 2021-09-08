package types

import (
	"time"
)

const (
	MessageTypeInfo    = "info"
	MessageTypeWarning = "warning"
	MessageTypeText    = "text"
	MessageTypeAuto    = "auto"
)

// Message will send to IM, gen by agent api, or system event.
// Not stored in the database.
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
