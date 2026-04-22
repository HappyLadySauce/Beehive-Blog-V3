package repo

import (
	"context"
	"strings"
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

// GetByProviderIdentity fetches a federated identity using the canonical subject plus compatibility identifiers.
// GetByProviderIdentity 使用规范 subject 与兼容标识查询联邦身份。
//
// Provider subject remains the canonical identity key.
// unionid/openid are only compatibility lookups for providers such as WeChat when the stable subject strategy evolves.
// provider_subject 仍然是规范身份键。
// unionid/openid 仅作为微信等 provider 在主体标识策略演进时的兼容查找路径。
func (r *FederatedIdentityRepository) GetByProviderIdentity(ctx context.Context, provider, subject string, openID, unionID *string) (*entity.FederatedIdentity, error) {
	query := r.db.WithContext(ctx).Where("provider = ?", provider)

	conditions := []string{"provider_subject = ?"}
	args := []any{strings.TrimSpace(subject)}
	if unionID != nil && strings.TrimSpace(*unionID) != "" {
		conditions = append(conditions, "unionid = ?")
		args = append(args, strings.TrimSpace(*unionID))
	}
	if openID != nil && strings.TrimSpace(*openID) != "" {
		conditions = append(conditions, "openid = ?")
		args = append(args, strings.TrimSpace(*openID))
	}

	var fed entity.FederatedIdentity
	if err := query.Where("("+strings.Join(conditions, " OR ")+")", args...).First(&fed).Error; err != nil {
		return nil, err
	}

	return &fed, nil
}

// TouchLogin updates profile snapshots and last login time.
// TouchLogin 更新资料快照与最近登录时间。
func (r *FederatedIdentityRepository) TouchLogin(ctx context.Context, id int64, providerSubject, providerSubjectType *string, providerEmail, providerDisplayName, avatarURL, providerLogin, openID, unionID *string, rawProfile []byte, lastLoginAt time.Time) error {
	updates := map[string]any{
		"last_login_at": lastLoginAt,
		"updated_at":    lastLoginAt,
	}
	if providerSubject != nil {
		updates["provider_subject"] = providerSubject
	}
	if providerSubjectType != nil {
		updates["provider_subject_type"] = providerSubjectType
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
	if providerLogin != nil {
		updates["provider_login"] = providerLogin
	}
	if openID != nil {
		updates["openid"] = openID
	}
	if unionID != nil {
		updates["unionid"] = unionID
	}
	if len(rawProfile) > 0 {
		updates["raw_profile"] = rawProfile
	}

	return r.db.WithContext(ctx).
		Model(&entity.FederatedIdentity{}).
		Where("id = ?", id).
		Updates(updates).Error
}
