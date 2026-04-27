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

// ListFilter describes user list filtering and pagination.
// ListFilter 描述用户列表过滤与分页参数。
type ListFilter struct {
	Keyword  string
	Role     string
	Status   string
	Page     int
	PageSize int
}

// List returns users matching the filter and the total count.
// List 返回符合过滤条件的用户列表与总数。
func (r *UserRepository) List(ctx context.Context, filter ListFilter) ([]entity.User, int64, error) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	query := r.db.WithContext(ctx).Model(&entity.User{})

	if keyword := strings.ToLower(strings.TrimSpace(filter.Keyword)); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("LOWER(username) LIKE ? OR LOWER(COALESCE(email, '')) LIKE ? OR LOWER(COALESCE(nickname, '')) LIKE ?", like, like, like)
	}
	if role := strings.TrimSpace(filter.Role); role != "" {
		query = query.Where("role = ?", role)
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []entity.User
	if err := query.
		Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateProfile updates mutable self-service profile fields.
// UpdateProfile 更新用户可自助修改的资料字段。
func (r *UserRepository) UpdateProfile(ctx context.Context, id int64, nickname, avatarURL *string, at time.Time) (*entity.User, error) {
	updates := map[string]any{
		"nickname":   nickname,
		"avatar_url": avatarURL,
		"updated_at": at,
	}
	if err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

// UpdateRole updates a user's role.
// UpdateRole 更新用户角色。
func (r *UserRepository) UpdateRole(ctx context.Context, id int64, role string, at time.Time) (*entity.User, error) {
	return r.updateStringField(ctx, id, "role", role, at)
}

// UpdateStatus updates a user's account status.
// UpdateStatus 更新用户账号状态。
func (r *UserRepository) UpdateStatus(ctx context.Context, id int64, status string, at time.Time) (*entity.User, error) {
	return r.updateStringField(ctx, id, "status", status, at)
}

func (r *UserRepository) updateStringField(ctx context.Context, id int64, column, value string, at time.Time) (*entity.User, error) {
	if err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			column:       value,
			"updated_at": at,
		}).Error; err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
