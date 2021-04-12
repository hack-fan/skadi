package types

import (
	"time"

	"github.com/hack-fan/x/xtype"
)

const (
	JobStatusQueuing   = "queuing"
	JobStatusCanceled  = "canceled"
	JobStatusSent      = "sent"
	JobStatusRunning   = "running"
	JobStatusExpired   = "expired"
	JobStatusSucceeded = "succeeded"
	JobStatusFailed    = "failed"
)

var RESERVED = xtype.Strings{"agent", "status", "help", "delay", "group"}

// JobBasic will be pulled by agent
type JobBasic struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// JobInput is input fields which upstream service passed
type JobInput struct {
	UserID  string `json:"user_id" gorm:"type:varchar(20);index:by_user"`
	AgentID string `json:"agent_id" gorm:"type:varchar(20);index:by_agent"`
	// job payload
	Message string `json:"message" gorm:"type:varchar(255)" validate:"required,lte=255"`
	// source context, any string defined by source
	Source string `json:"source" gorm:"type:varchar(255)" validate:"omitempty,lte=255"`
	// callback url, will be called after expired/succeeded/failed
	Callback string `json:"callback" gorm:"type:varchar(255)" validate:"omitempty,lte=255"`
}

// Job will be saved in db, it's a gorm mysql model
type Job struct {
	ID string `json:"id" gorm:"type:varchar(20);primary_key"`
	JobInput
	// queuing/canceled/sent/expired/succeeded/failed
	Status string `json:"status" gorm:"type:varchar(10)"`
	// agent returned, job log or other message
	Result      string     `json:"result" gorm:"type:varchar(1024)"`
	CreatedAt   time.Time  `json:"created_at" gorm:"index:by_user;index:by_agent"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	CanceledAt  *time.Time `json:"canceled_at,omitempty"`
	ExpiredAt   *time.Time `json:"expired_at,omitempty"`
	SucceededAt *time.Time `json:"succeeded_at,omitempty"`
	FailedAt    *time.Time `json:"failed_at,omitempty"`
}

// JobResult is reported by agent
// The length of result is not limited, if it longer than 1024, will be cut.
type JobResult struct {
	// agent returned, job log or other message
	Result string `json:"result,omitempty"`
}
