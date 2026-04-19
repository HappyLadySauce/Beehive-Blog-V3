package entity

import "time"

// FederatedIdentity maps the identity.federated_identities table.
// FederatedIdentity 映射 identity.federated_identities 表。
type FederatedIdentity struct {
	ID                  int64      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID              int64      `gorm:"column:user_id"`
	Provider            string     `gorm:"column:provider"`
	ProviderSubject     string     `gorm:"column:provider_subject"`
	ProviderSubjectType string     `gorm:"column:provider_subject_type"`
	UnionID             *string    `gorm:"column:unionid"`
	OpenID              *string    `gorm:"column:openid"`
	ProviderLogin       *string    `gorm:"column:provider_login"`
	ProviderEmail       *string    `gorm:"column:provider_email"`
	ProviderDisplayName *string    `gorm:"column:provider_display_name"`
	AvatarURL           *string    `gorm:"column:avatar_url"`
	AppIDOrClientID     *string    `gorm:"column:app_id_or_client_id"`
	AccessScope         *string    `gorm:"column:access_scope"`
	RawProfile          []byte     `gorm:"column:raw_profile;type:jsonb"`
	LastLoginAt         *time.Time `gorm:"column:last_login_at"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"`
}

// TableName returns the fully-qualified table name.
// TableName 返回完整限定表名。
func (FederatedIdentity) TableName() string {
	return "identity.federated_identities"
}
