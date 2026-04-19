package repo

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RefreshTokenRepository encapsulates refresh token persistence.
// RefreshTokenRepository 封装 refresh token 持久化访问。
type RefreshTokenRepository struct {
	db *gorm.DB
}

// Create persists a refresh token row.
// Create 持久化 refresh token 记录。
func (r *RefreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetActiveForUpdateByHash fetches an active refresh token row with FOR UPDATE.
// GetActiveForUpdateByHash 使用 FOR UPDATE 加锁查询活跃 refresh token。
func (r *RefreshTokenRepository) GetActiveForUpdateByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

// Revoke marks a refresh token as revoked.
// Revoke 将 refresh token 标记为已吊销。
func (r *RefreshTokenRepository) Revoke(ctx context.Context, id int64, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", revokedAt).Error
}

// RevokeActiveBySessionID revokes all active refresh tokens for a session.
// RevokeActiveBySessionID 吊销指定会话的全部活跃 refresh token。
func (r *RefreshTokenRepository) RevokeActiveBySessionID(ctx context.Context, sessionID int64, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("session_id = ? AND revoked_at IS NULL", sessionID).
		Update("revoked_at", revokedAt).Error
}
