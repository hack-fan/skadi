package types

import "time"

type Event struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
