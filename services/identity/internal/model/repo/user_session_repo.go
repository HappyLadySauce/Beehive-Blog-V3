package repo

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserSessionRepository encapsulates session persistence.
// UserSessionRepository 封装会话持久化访问。
type UserSessionRepository struct {
	db *gorm.DB
}

// Create persists a session row.
// Create 持久化会话记录。
func (r *UserSessionRepository) Create(ctx context.Context, session *entity.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID fetches a session by primary key.
// GetByID 按主键查询会话。
func (r *UserSessionRepository) GetByID(ctx context.Context, id int64) (*entity.UserSession, error) {
	var session entity.UserSession
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

// GetForUpdateByID loads a session row with FOR UPDATE.
// GetForUpdateByID 使用 FOR UPDATE 加锁查询会话。
func (r *UserSessionRepository) GetForUpdateByID(ctx context.Context, id int64) (*entity.UserSession, error) {
	var session entity.UserSession
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

// TouchActive updates the session activity timestamps and expiry.
// TouchActive 更新会话活跃时间与过期时间。
func (r *UserSessionRepository) TouchActive(ctx context.Context, id int64, lastSeenAt, expiresAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserSession{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"last_seen_at": lastSeenAt,
			"expires_at":   expiresAt,
			"updated_at":   lastSeenAt,
		}).Error
}

// Revoke marks a session as revoked.
// Revoke 将会话标记为已吊销。
func (r *UserSessionRepository) Revoke(ctx context.Context, id int64, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserSession{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     "revoked",
			"revoked_at": revokedAt,
			"updated_at": revokedAt,
		}).Error
}

// MarkExpired marks a session as expired.
// MarkExpired 将会话标记为已过期。
func (r *UserSessionRepository) MarkExpired(ctx context.Context, id int64, expiredAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserSession{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     "expired",
			"updated_at": expiredAt,
		}).Error
}
