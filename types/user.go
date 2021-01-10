package types

import "time"

// User have id and secret, saved in db.
type User struct {
	ID        string    `json:"id" gorm:"type:varchar(20);primary_key"`
	Secret    string    `json:"secret" gorm:"type:varchar(20);index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
