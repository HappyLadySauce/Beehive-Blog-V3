package repo

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"gorm.io/gorm"
)

// FederatedIdentityRepository encapsulates federated identity access.
// FederatedIdentityRepository 封装联邦身份访问。
type FederatedIdentityRepository struct {
	db *gorm.DB
}

// Create persists a federated identity row.
// Create 持久化联邦身份记录。
func (r *FederatedIdentityRepository) Create(ctx context.Context, fed *entity.FederatedIdentity) error {
	return r.db.WithContext(ctx).Create(fed).Error
}

// GetByProviderSubject fetches a federated identity by stable provider subject.
// GetByProviderSubject 按稳定 provider subject 查询联邦身份。
func (r *FederatedIdentityRepository) GetByProviderSubject(ctx context.Context, provider, subject string) (*entity.FederatedIdentity, error) {
	var fed entity.FederatedIdentity
	if err := r.db.WithContext(ctx).
		Where("provider = ? AND provider_subject = ?", provider, subject).
		First(&fed).Error; err != nil {
		return nil, err
	}

	return &fed, nil
}

// TouchLogin updates profile snapshots and last login time.
// TouchLogin 更新资料快照与最近登录时间。
func (r *FederatedIdentityRepository) TouchLogin(ctx context.Context, id int64, providerEmail, providerDisplayName, avatarURL *string, rawProfile []byte, lastLoginAt time.Time) error {
	updates := map[string]any{
		"last_login_at": lastLoginAt,
		"updated_at":    lastLoginAt,
	}
	if providerEmail != nil {
		updates["provider_email"] = providerEmail
	}
	if providerDisplayName != nil {
		updates["provider_display_name"] = providerDisplayName
	}
	if avatarURL != nil {
		updates["avatar_url"] = avatarURL
	}
	if len(rawProfile) > 0 {
		updates["raw_profile"] = rawProfile
	}

	return r.db.WithContext(ctx).
		Model(&entity.FederatedIdentity{}).
		Where("id = ?", id).
		Updates(updates).Error
}
