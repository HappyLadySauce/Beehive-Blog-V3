package service

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

// IntrospectService handles access token introspection.
// IntrospectService 处理 access token introspection。
type IntrospectService struct {
	deps Dependencies
}

// NewIntrospectService creates an IntrospectService instance.
// NewIntrospectService 创建 IntrospectService 实例。
func NewIntrospectService(deps Dependencies) *IntrospectService {
	return &IntrospectService{deps: deps}
}

// Execute validates an access token against JWT claims and database state.
// Execute 结合 JWT claims 与数据库状态校验 access token。
func (s *IntrospectService) Execute(ctx context.Context, in IntrospectAccessTokenInput) (*IntrospectionResult, error) {
	if in.AccessToken == "" {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "access_token is required")
	}

	claims, err := auth.ParseAccessToken(s.deps.Config.Security.AccessTokenSecret, in.AccessToken)
	if err != nil {
		return &IntrospectionResult{Active: false}, nil
	}
	now := s.deps.Clock()
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(now) {
		return &IntrospectionResult{Active: false}, nil
	}
	if claims.SessionID <= 0 || claims.UserID <= 0 {
		return &IntrospectionResult{Active: false}, nil
	}

	session, err := s.deps.Store.UserSessions.GetByID(ctx, claims.SessionID)
	if err != nil {
		if repo.IsNotFound(err) {
			return &IntrospectionResult{Active: false}, nil
		}
		return nil, err
	}
	if session.Status != auth.SessionStatusActive || session.ExpiresAt.Before(now) {
		return &IntrospectionResult{Active: false}, nil
	}

	user, err := s.deps.Store.Users.GetByID(ctx, claims.UserID)
	if err != nil {
		if repo.IsNotFound(err) {
			return &IntrospectionResult{Active: false}, nil
		}
		return nil, err
	}
	if user.Status != auth.UserStatusActive {
		return &IntrospectionResult{Active: false}, nil
	}

	return &IntrospectionResult{
		Active:    true,
		User:      user,
		Session:   session,
		ExpiresAt: claims.ExpiresAt.Time.Unix(),
	}, nil
}
