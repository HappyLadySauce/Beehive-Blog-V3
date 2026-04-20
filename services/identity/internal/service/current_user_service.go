package service

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

// CurrentUserService handles trusted current-user lookups.
// CurrentUserService 处理可信当前用户查询。
type CurrentUserService struct {
	deps Dependencies
}

// NewCurrentUserService creates a CurrentUserService instance.
// NewCurrentUserService 创建 CurrentUserService 实例。
func NewCurrentUserService(deps Dependencies) *CurrentUserService {
	return &CurrentUserService{deps: deps}
}

// Execute returns the trusted current user snapshot.
// Execute 返回可信当前用户快照。
func (s *CurrentUserService) Execute(ctx context.Context, in GetCurrentUserInput) (*CurrentUserResult, error) {
	if in.UserID <= 0 {
		return nil, NewError(ErrorKindInvalidArgument, "user_id is invalid", nil)
	}

	user, err := s.deps.Store.Users.GetByID(ctx, in.UserID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, NewError(ErrorKindNotFound, "user not found", nil)
		}
		return nil, err
	}

	return &CurrentUserResult{User: user}, nil
}
