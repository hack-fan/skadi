package types

import "time"

// ChannelType is provider+way, some provider have more than one way to send message
type ChannelType string

// All channel types
const (
	ChannelTypeWechat   ChannelType = "mp"
	ChannelTypeTelegram ChannelType = "tg"
	ChannelTypeFeishu   ChannelType = "feishu"
)

// Channel is the way of sending message to users.
type Channel struct {
	ID     string `json:"id" gorm:"type:varchar(20);primaryKey"`
	UserID string `json:"user_id" gorm:"type:varchar(20);index:idx_uid"`
	// Name can't contain space.
	Name      string      `json:"name" gorm:"type:varchar(50)"`
	Remark    string      `json:"remark"`
	Type      ChannelType `json:"type" gorm:"type:varchar(10)"`
	Address   string      `json:"address" gorm:"type:varchar(50)"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
