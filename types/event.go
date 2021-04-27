package types

import "github.com/hack-fan/jq"

// Event will publish to redis queue
const (
	EventMessage   = "message"
	EventJobStatus = "job:status"
)

// NewEventMessageQueue is stateless in local, you can new it anywhere.
func NewEventMessageQueue(rdb jq.RedisClient) *jq.Queue {
	return jq.NewQueue("skadi:event:"+EventMessage, rdb)
}

// NewEventJobStatusQueue is stateless in local, you can new it anywhere.
func NewEventJobStatusQueue(rdb jq.RedisClient) *jq.Queue {
	return jq.NewQueue("skadi:event:"+EventJobStatus, rdb)
}
