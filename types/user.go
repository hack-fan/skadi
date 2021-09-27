package types

// UserSetting stores user's setting. This package does not manage user profile.
// The user can be a person, or a team, in your app.
type UserSetting struct {
	UserID string `json:"user_id" gorm:"type:varchar(20);primaryKey"`
	// The default channels
	DefaultInfoChannel    string `json:"default_info_channel" gorm:"type:varchar(20)"`
	DefaultWarningChannel string `json:"default_warning_channel" gorm:"type:varchar(20)"`
	DefaultErrorChannel   string `json:"default_error_channel" gorm:"type:varchar(20)"`
}
