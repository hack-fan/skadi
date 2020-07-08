package types

import "time"

// ServerJob will be pulled by servers
type ServerJob struct {
	ID         string `json:"id" gorm:"type:varchar(20);primary_key"`
	ServerName string `json:"server_name" gorm:"type:varchar(20)"`
	Message    string `json:"message" gorm:"type:varchar(255)"`
}

// Job will be saved
type Job struct {
	ServerJob
	UserID      string     `json:"user_id" gorm:"type:varchar(20);index:by_user"`
	Source      SourceType `json:"source" gorm:"type:varchar(10)"`
	ServerID    string     `json:"server_id" gorm:"type:varchar(20);index:by_server"`
	Status      string     `json:"status" gorm:"type:varchar(10)"` // queuing/sent/expired/succeeded/failed
	Result      string     `json:"result"`
	CreatedAt   time.Time  `json:"created_at" gorm:"index:by_user;index:by_server"`
	SentAt      *time.Time `json:"sent_at"`
	ExpiredAt   *time.Time `json:"expired_at"`
	SucceededAt *time.Time `json:"succeeded_at"`
	FailedAt    *time.Time `json:"failed_at"`
}
