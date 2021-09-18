package types

// UserSetting stores user's setting. This package does not manage user profile.
// The user can be a person, or a team, in your app.
type UserSetting struct {
	UserID string `json:"user_id" gorm:"type:varchar(20);primaryKey"`
	// The default channels
	InfoChannel    string `json:"info_channel" gorm:"type:varchar(20)"`
	WarningChannel string `json:"warning_channel" gorm:"type:varchar(20)"`
	ErrorChannel   string `json:"error_channel" gorm:"type:varchar(20)"`
}
