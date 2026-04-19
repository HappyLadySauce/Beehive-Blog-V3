package entity

import "time"

// CredentialLocal maps the identity.credential_locals table.
// CredentialLocal 映射 identity.credential_locals 表。
type CredentialLocal struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID            int64     `gorm:"column:user_id"`
	PasswordHash      string    `gorm:"column:password_hash"`
	PasswordUpdatedAt time.Time `gorm:"column:password_updated_at"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

// TableName returns the fully-qualified table name.
// TableName 返回完整限定表名。
func (CredentialLocal) TableName() string {
	return "identity.credential_locals"
}
