package types

type ChannelType string

// Channel is the way can send message to users.
type Channel struct {
	ID     string      `json:"id" gorm:"type:varchar(20);primaryKey"`
	UserID string      `json:"user_id" gorm:"type:varchar(20);index:idx_uid"`
	Name   string      `json:"name" gorm:"type:varchar(50)"`
	Type   ChannelType `json:"type" gorm:"type:varchar(10)"`
}
