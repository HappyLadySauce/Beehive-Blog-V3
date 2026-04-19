package repo

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OAuthLoginStateRepository encapsulates OAuth state persistence.
// OAuthLoginStateRepository 封装 OAuth 状态持久化访问。
type OAuthLoginStateRepository struct {
	db *gorm.DB
}

// Create persists a new OAuth login state row.
// Create 持久化新的 OAuth 登录状态记录。
func (r *OAuthLoginStateRepository) Create(ctx context.Context, state *entity.OAuthLoginState) error {
	return r.db.WithContext(ctx).Create(state).Error
}

// GetForUpdateByProviderState loads a state row with FOR UPDATE.
// GetForUpdateByProviderState 使用 FOR UPDATE 加锁查询状态记录。
func (r *OAuthLoginStateRepository) GetForUpdateByProviderState(ctx context.Context, provider, state string) (*entity.OAuthLoginState, error) {
	var row entity.OAuthLoginState
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("provider = ? AND state = ?", provider, state).
		First(&row).Error; err != nil {
		return nil, err
	}

	return &row, nil
}

// Consume marks a state row as consumed.
// Consume 将状态记录标记为已消费。
func (r *OAuthLoginStateRepository) Consume(ctx context.Context, id int64, consumedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.OAuthLoginState{}).
		Where("id = ? AND consumed_at IS NULL", id).
		Update("consumed_at", consumedAt).Error
}
