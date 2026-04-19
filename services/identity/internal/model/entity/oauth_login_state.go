package entity

import "time"

// OAuthLoginState maps the identity.oauth_login_states table.
// OAuthLoginState 映射 identity.oauth_login_states 表。
type OAuthLoginState struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Provider        string     `gorm:"column:provider"`
	State           string     `gorm:"column:state"`
	RedirectURI     string     `gorm:"column:redirect_uri"`
	ClientType      *string    `gorm:"column:client_type"`
	DeviceID        *string    `gorm:"column:device_id"`
	CodeVerifier    *string    `gorm:"column:code_verifier"`
	RequestedScopes *string    `gorm:"column:requested_scopes"`
	ExpiresAt       time.Time  `gorm:"column:expires_at"`
	ConsumedAt      *time.Time `gorm:"column:consumed_at"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
}

// TableName returns the fully-qualified table name.
// TableName 返回完整限定表名。
func (OAuthLoginState) TableName() string {
	return "identity.oauth_login_states"
}
