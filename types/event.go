package types

import "time"

// Event will be pulled by servers
type Event struct {
	ID          string     `json:"id" gorm:"type:varchar(20);primary_key"`
	UserID      string     `json:"user_id" gorm:"type:varchar(20);index:by_user"`
	ServerID    string     `json:"server_id" gorm:"type:varchar(20);index:by_server"`
	ServerName  string     `json:"server_name" gorm:"type:varchar(20)"`
	Message     string     `json:"message" gorm:"type:varchar(255)"`
	Status      string     `json:"status" gorm:"type:varchar(10)"` // queuing/sent/expired/succeeded/failed
	CreatedAt   time.Time  `json:"created_at" gorm:"index:by_user;index:by_server"`
	SentAt      *time.Time `json:"sent_at"`
	ExpiredAt   *time.Time `json:"expired_at"`
	SucceededAt *time.Time `json:"succeeded_at"`
	FailedAt    *time.Time `json:"failed_at"`
}
