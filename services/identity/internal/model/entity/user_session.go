package entity

import "time"

// UserSession maps the identity.user_sessions table.
// UserSession 映射 identity.user_sessions 表。
type UserSession struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     int64      `gorm:"column:user_id"`
	AuthSource string     `gorm:"column:auth_source"`
	ClientType *string    `gorm:"column:client_type"`
	DeviceID   *string    `gorm:"column:device_id"`
	DeviceName *string    `gorm:"column:device_name"`
	IPAddress  *string    `gorm:"column:ip_address"`
	UserAgent  *string    `gorm:"column:user_agent"`
	Status     string     `gorm:"column:status"`
	LastSeenAt *time.Time `gorm:"column:last_seen_at"`
	ExpiresAt  time.Time  `gorm:"column:expires_at"`
	RevokedAt  *time.Time `gorm:"column:revoked_at"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
}

// TableName returns the fully-qualified table name.
// TableName 返回完整限定表名。
func (UserSession) TableName() string {
	return "identity.user_sessions"
}
