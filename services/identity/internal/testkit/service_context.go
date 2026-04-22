package testkit

import (
	"context"
	"testing"
	"time"

	identityprovider "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
)

// NewConfig returns the baseline identity test configuration.
// NewConfig 返回 identity 测试基线配置。
func NewConfig() config.Config {
	return config.Config{
		InternalAuthToken: "test-internal-auth-token",
		AllowedCallers:    []string{"gateway"},
		Postgres: config.PostgresConf{
			Host:   "unused",
			User:   "unused",
			DBName: "unused",
		},
		StateRedis: config.RedisConf{
			Host: "unused",
		},
		Security: config.SecurityConf{
			AccessTokenSecret:      "test-access-token-secret",
			AccessTokenTTLSeconds:  900,
			RefreshTokenTTLSeconds: 3600,
			StateTTLSeconds:        600,
			PasswordHashCost:       4,
		},
		SSO: config.SSOConf{
			GitHub: config.OAuthProviderConf{
				Enabled:      true,
				ClientID:     "github-client-id",
				ClientSecret: "github-client-secret",
				RedirectURL:  "https://example.com/auth/github/callback",
				Scopes:       []string{"read:user", "user:email"},
			},
			QQ: config.OAuthProviderConf{
				Enabled:      true,
				ClientID:     "qq-client-id",
				ClientSecret: "qq-client-secret",
				RedirectURL:  "https://example.com/auth/qq/callback",
				Scopes:       []string{"get_user_info"},
			},
			WeChat: config.OAuthProviderConf{
				Enabled:      true,
				ClientID:     "wechat-appid",
				ClientSecret: "wechat-secret",
				RedirectURL:  "https://example.com/auth/wechat/callback",
				Scopes:       []string{"snsapi_login"},
			},
		},
	}
}

// NewStore returns a clean repository store backed by the shared PostgreSQL container.
// NewStore 返回基于共享 PostgreSQL 容器且已清表的 repository store。
func NewStore(t *testing.T) *repo.Store {
	t.Helper()
	ResetPostgres(t)
	return repo.NewStore(PostgresDB(t))
}

// NewProviderRegistry builds the standard provider registry for tests.
// NewProviderRegistry 构建测试用标准 provider 注册表。
func NewProviderRegistry(conf config.Config) *identityprovider.Registry {
	return identityprovider.NewRegistry(
		identityprovider.NewGitHubClient(conf.SSO.GitHub),
		identityprovider.NewQQClient(conf.SSO.QQ),
		identityprovider.NewWeChatClient(conf.SSO.WeChat),
	)
}

// NewServiceDependencies returns test-ready service dependencies.
// NewServiceDependencies 返回测试可用的 service 依赖。
func NewServiceDependencies(t *testing.T, clock func() time.Time) identityservice.Dependencies {
	t.Helper()

	conf := NewConfig()
	store := NewStore(t)
	ResetRedis(t)
	providers := NewProviderRegistry(conf)
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}

	return identityservice.Dependencies{
		Config:    conf,
		Store:     store,
		Providers: providers,
		Clock:     clock,
	}
}

// NewServiceContext returns a test-ready service context.
// NewServiceContext 返回测试可用的 service context。
func NewServiceContext(t *testing.T) *svc.ServiceContext {
	t.Helper()

	conf := NewConfig()
	store := NewStore(t)
	redisClient := RedisClient(t)
	providers := NewProviderRegistry(conf)
	readinessChecker := func(_ context.Context) error { return nil }
	services := identityservice.NewManager(conf, store, providers, readinessChecker)

	return &svc.ServiceContext{
		Config:    conf,
		DB:        store.DB(),
		Redis:     redisClient,
		Store:     store,
		Providers: providers,
		Services:  services,
	}
}
