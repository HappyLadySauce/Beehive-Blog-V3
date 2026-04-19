package repo

import (
	"gorm.io/gorm"
)

// Store aggregates all identity repositories.
// Store 聚合 identity 服务的全部 repository。
type Store struct {
	db *gorm.DB

	Users               *UserRepository
	CredentialLocals    *CredentialLocalRepository
	OAuthLoginStates    *OAuthLoginStateRepository
	FederatedIdentities *FederatedIdentityRepository
	UserSessions        *UserSessionRepository
	RefreshTokens       *RefreshTokenRepository
	IdentityAudits      *IdentityAuditRepository
}

// NewStore creates a store backed by the provided database handle.
// NewStore 使用给定数据库句柄创建 Store。
func NewStore(db *gorm.DB) *Store {
	return &Store{
		db:                  db,
		Users:               &UserRepository{db: db},
		CredentialLocals:    &CredentialLocalRepository{db: db},
		OAuthLoginStates:    &OAuthLoginStateRepository{db: db},
		FederatedIdentities: &FederatedIdentityRepository{db: db},
		UserSessions:        &UserSessionRepository{db: db},
		RefreshTokens:       &RefreshTokenRepository{db: db},
		IdentityAudits:      &IdentityAuditRepository{db: db},
	}
}

// DB exposes the underlying gorm DB for transaction boundaries only.
// DB 暴露底层 gorm DB，仅用于事务边界控制。
func (s *Store) DB() *gorm.DB {
	return s.db
}

// WithTx returns a store bound to the provided transaction.
// WithTx 返回绑定到指定事务的 Store。
func (s *Store) WithTx(tx *gorm.DB) *Store {
	return NewStore(tx)
}
