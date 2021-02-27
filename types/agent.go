package types

import "time"

// Agent belongs to user, have id and secret, saved in db.
type Agent struct {
	ID     string `json:"id" gorm:"type:varchar(20);primary_key"`
	UserID string `json:"user_id" gorm:"type:varchar(20);index"`
	Secret string `json:"secret" gorm:"type:varchar(20);index"`
	// IP address last connected from
	IP string `json:"ip" gorm:"type:varchar(50)"`
	// Only update after a agent offline
	ActivatedAt *time.Time `json:"activated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
