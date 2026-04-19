package entity

import "time"

// User maps the identity.users table.
// User 映射 identity.users 表。
type User struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Username    string     `gorm:"column:username"`
	Email       *string    `gorm:"column:email"`
	Nickname    *string    `gorm:"column:nickname"`
	AvatarURL   *string    `gorm:"column:avatar_url"`
	Role        string     `gorm:"column:role"`
	Status      string     `gorm:"column:status"`
	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

// TableName returns the fully-qualified table name.
// TableName 返回完整限定表名。
func (User) TableName() string {
	return "identity.users"
}
