package repo

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"gorm.io/gorm"
)

// UserRepository encapsulates data access for users.
// UserRepository 封装 users 表的数据访问。
type UserRepository struct {
	db *gorm.DB
}

// Create persists a new user row.
// Create 持久化新用户记录。
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID fetches a user by primary key.
// GetByID 按主键查询用户。
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByUsername fetches a user by normalized username.
// GetByUsername 按规范化用户名查询用户。
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail fetches a user by case-insensitive email.
// GetByEmail 按邮箱大小写不敏感查询用户。
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByLoginIdentifier resolves a user by username or email.
// FindByLoginIdentifier 按用户名或邮箱解析用户。
func (r *UserRepository) FindByLoginIdentifier(ctx context.Context, identifier string) (*entity.User, error) {
	if strings.Contains(identifier, "@") {
		return r.GetByEmail(ctx, identifier)
	}

	return r.GetByUsername(ctx, identifier)
}

// TouchLogin updates the user's last login timestamp.
// TouchLogin 更新用户最近登录时间。
func (r *UserRepository) TouchLogin(ctx context.Context, id int64, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"last_login_at": at,
			"updated_at":    at,
		}).Error
}
