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

var RESERVED = xtype.Strings{"agent", "delay", "status", "help", "delay", "group", "plan", "poster",
	"link", "unlink", "follow", "unfollow", "sub", "unsub"}

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
	// callback url, if exists, will be called after sent/expired/succeeded/failed
	// left empty will notify in your default IM,
	// set to "disable" will disable any notify or callback
	Callback string `json:"callback" gorm:"type:varchar(255)" validate:"omitempty,lte=255"`
}

// Job will be saved in db, it's a gorm mysql model
type Job struct {
	ID string `json:"id" gorm:"type:varchar(20);primaryKey"`
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

// DelayedJob will be stored in db for future run.
type DelayedJob struct {
	ID string `json:"id" gorm:"type:varchar(20);primaryKey"`
	JobInput
	ActiveAt  time.Time `json:"active_at"`
	CreatedAt time.Time `json:"created_at"`
}

// DelayedJobInput just for api input
type DelayedJobInput struct {
	Message  string `validate:"required,lte=255" json:"message"`
	Duration string `json:"duration,omitempty"`
	Minutes  int    `json:"minutes,omitempty"`
	Hours    int    `json:"hours,omitempty"`
	Days     int    `json:"days,omitempty"`
	// an url, will be called after job status changed,
	// left empty will notify in your default IM,
	// set to "disable" will disable any notify or callback
	Callback string `validate:"omitempty,lte=255" json:"callback,omitempty"`
}

// JobResult is reported by agent
// The length of result is not limited, if it is longer than 1024, will be cut.
type JobResult struct {
	// agent returned, job log or other message
	Result string `json:"result,omitempty"`
}
