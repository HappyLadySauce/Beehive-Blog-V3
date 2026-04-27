package repo

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"gorm.io/gorm"
)

// CredentialLocalRepository encapsulates local credential access.
// CredentialLocalRepository 封装本地凭证访问。
type CredentialLocalRepository struct {
	db *gorm.DB
}

// Create persists a local credential row.
// Create 持久化本地凭证记录。
func (r *CredentialLocalRepository) Create(ctx context.Context, credential *entity.CredentialLocal) error {
	return r.db.WithContext(ctx).Create(credential).Error
}

// GetByUserID fetches local credentials for a user.
// GetByUserID 按用户 ID 查询本地凭证。
func (r *CredentialLocalRepository) GetByUserID(ctx context.Context, userID int64) (*entity.CredentialLocal, error) {
	var credential entity.CredentialLocal
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&credential).Error; err != nil {
		return nil, err
	}

	return &credential, nil
}

// UpdatePasswordHash replaces the local password hash for a user.
// UpdatePasswordHash 替换用户本地密码哈希。
func (r *CredentialLocalRepository) UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string, at time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&entity.CredentialLocal{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"password_hash":       passwordHash,
			"password_updated_at": at,
			"updated_at":          at,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
