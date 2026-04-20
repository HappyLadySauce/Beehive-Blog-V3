package service

import (
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

// Manager aggregates all identity use-case services.
// Manager 聚合 identity 的全部用例服务。
type Manager struct {
	Register    *RegisterService
	Login       *LoginService
	Refresh     *RefreshService
	Logout      *LogoutService
	CurrentUser *CurrentUserService
	Introspect  *IntrospectService
	SSOStart    *SSOStartService
	SSOFinish   *SSOFinishService
}

// NewManager builds the complete identity service graph.
// NewManager 构建完整的 identity service 依赖图。
func NewManager(c config.Config, store *repo.Store, providers *provider.Registry) *Manager {
	deps := newDependencies(c, store, providers)
	return &Manager{
		Register:    NewRegisterService(deps),
		Login:       NewLoginService(deps),
		Refresh:     NewRefreshService(deps),
		Logout:      NewLogoutService(deps),
		CurrentUser: NewCurrentUserService(deps),
		Introspect:  NewIntrospectService(deps),
		SSOStart:    NewSSOStartService(deps),
		SSOFinish:   NewSSOFinishService(deps),
	}
}
