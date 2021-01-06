package types

import "time"

// Agent belongs to user, have id and secret, saved in db.
type Agent struct {
	ID        string     `json:"id" gorm:"type:varchar(20)"`
	UserID    string     `json:"user_id" gorm:"type:varchar(20);index"`
	Secret    string     `json:"secret" gorm:"type:varchar(20);index"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
