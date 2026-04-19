package entity

import "time"

// IdentityAudit maps the identity.identity_audits table.
// IdentityAudit 映射 identity.identity_audits 表。
type IdentityAudit struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     *int64    `gorm:"column:user_id"`
	SessionID  *int64    `gorm:"column:session_id"`
	Provider   *string   `gorm:"column:provider"`
	AuthSource *string   `gorm:"column:auth_source"`
	EventType  string    `gorm:"column:event_type"`
	Result     string    `gorm:"column:result"`
	ClientIP   *string   `gorm:"column:client_ip"`
	UserAgent  *string   `gorm:"column:user_agent"`
	Detail     []byte    `gorm:"column:detail;type:jsonb"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

// TableName returns the fully-qualified table name.
// TableName 返回完整限定表名。
func (IdentityAudit) TableName() string {
	return "identity.identity_audits"
}
