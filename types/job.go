package types

import "time"

// JobBasic will be pulled by agent
type JobBasic struct {
	ID      string `json:"id" gorm:"type:varchar(20);primary_key"`
	Message string `json:"message" gorm:"type:varchar(255)"`
}

// Job will be saved
type Job struct {
	JobBasic
	UserID  string `json:"user_id" gorm:"type:varchar(20);index:by_user"`
	AgentID string `json:"agent_id" gorm:"type:varchar(20);index:by_agent"`
	// queuing/sent/expired/succeeded/failed
	Status string `json:"status" gorm:"type:varchar(10)"`
	// agent returned, job log or other message
	Result string `json:"result" gorm:"type:varchar(1024)"`
	// source context, any string defined by source
	Source string `json:"source" gorm:"type:varchar(255)"`
	// callback url, will be called after expired/succeeded/failed
	Callback    string     `json:"callback" gorm:"type:varchar(255)"`
	CreatedAt   time.Time  `json:"created_at" gorm:"index:by_user;index:by_agent"`
	SentAt      *time.Time `json:"sent_at"`
	ExpiredAt   *time.Time `json:"expired_at"`
	SucceededAt *time.Time `json:"succeeded_at"`
	FailedAt    *time.Time `json:"failed_at"`
	CallbackAt  *time.Time `json:"callback_at"`
}
