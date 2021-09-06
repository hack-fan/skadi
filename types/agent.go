package types

import "time"

// Agent belongs to user, have id and secret, saved in db.
type Agent struct {
	ID     string `json:"id" gorm:"type:varchar(20);primary_key"`
	Secret string `json:"-" gorm:"type:varchar(20);index"`
	UserID string `json:"user_id" gorm:"type:varchar(20);index:idx_uid;uniqueIndex:idx_un1;index:idx_un2"`
	Name   string `json:"name" gorm:"type:varchar(50);uniqueIndex:idx_un1"`
	Alias  string `json:"alias" gorm:"type:varchar(50);index:idx_un2"`
	Remark string `json:"remark"`
	// IP address last connected from
	IP string `json:"ip" gorm:"type:varchar(50)"`
	// Only update after offline, means last activated at that time.
	ActivatedAt *time.Time `json:"activated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AgentBasic used for agent creation ui
type AgentBasic struct {
	Name   string `json:"name"`
	Alias  string `json:"alias"`
	Remark string `json:"remark"`
}
