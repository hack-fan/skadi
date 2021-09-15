package types

import (
	"time"
)

const (
	MessageTypeText = "text"
	// MessageTypeInfo Warning had been moved to MessageLevel
	MessageTypeInfo    = "info"
	MessageTypeWarning = "warning"
	// MessageTypeAuto have been deprecated
	MessageTypeAuto = "auto"
)

const (
	MessageLevelInfo    = "info"
	MessageLevelWarning = "warning"
	MessageLevelError   = "error"
)

// Message will send to IM, gen by agent api, or system event.
// Not stored in the database.
type Message struct {
	ID      string `json:"id"`
	AgentID string `json:"agent_id"`
	UserID  string `json:"user_id"`
	MessageInput
	CreatedAt time.Time `json:"created_at"`
}

// MessageInput from api
// TODO: check it
type MessageInput struct {
	Type    string `json:"type"`
	Level   string `json:"level"`
	Message string `json:"message"`
}
