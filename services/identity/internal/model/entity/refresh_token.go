package entity

import "time"

// RefreshToken maps the identity.refresh_tokens table.
// RefreshToken 映射 identity.refresh_tokens 表。
type RefreshToken struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement"`
	SessionID          int64      `gorm:"column:session_id"`
	TokenHash          string     `gorm:"column:token_hash"`
	IssuedAt           time.Time  `gorm:"column:issued_at"`
	ExpiresAt          time.Time  `gorm:"column:expires_at"`
	RotatedFromTokenID *int64     `gorm:"column:rotated_from_token_id"`
	RevokedAt          *time.Time `gorm:"column:revoked_at"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
}

// TableName returns the fully-qualified table name.
// TableName 返回完整限定表名。
func (RefreshToken) TableName() string {
	return "identity.refresh_tokens"
}
