package repo

import (
	"context"

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
