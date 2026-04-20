package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityprovider "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	"gorm.io/gorm"
)

var githubUsernameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// Dependencies defines the infrastructure and helper set used by services.
// Dependencies 定义 service 层使用的基础设施与 helper 集合。
type Dependencies struct {
	Config         config.Config
	Store          *repo.Store
	Providers      *identityprovider.Registry
	Clock          func() time.Time
	CheckReadiness func(ctx context.Context) error
}

// auditInput describes a single identity audit write.
// auditInput 描述单次 identity 审计写入。
type auditInput struct {
	UserID     *int64
	SessionID  *int64
	Provider   *string
	AuthSource *string
	EventType  string
	Result     string
	ClientIP   *string
	UserAgent  *string
	Detail     []byte
}

// newDependencies creates normalized service dependencies.
// newDependencies 创建规范化后的 service 依赖集合。
func newDependencies(c config.Config, store *repo.Store, providers *identityprovider.Registry, readinessChecker func(ctx context.Context) error) Dependencies {
	return Dependencies{
		Config:    c,
		Store:     store,
		Providers: providers,
		Clock: func() time.Time {
			return time.Now().UTC()
		},
		CheckReadiness: readinessChecker,
	}
}

// buildAuthResult assembles a service auth result.
// buildAuthResult 组装 service 认证结果。
func buildAuthResult(user *entity.User, session *entity.UserSession, accessToken, refreshToken string, accessExpiresAt time.Time) *AuthResult {
	return &AuthResult{
		User:            user,
		Session:         session,
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresAt: accessExpiresAt,
	}
}

// issueAccessToken issues a short-lived access token for the given session.
// issueAccessToken 为给定会话签发短期 access token。
func issueAccessToken(secret string, ttlSeconds int64, user *entity.User, session *entity.UserSession, now time.Time) (string, time.Time, error) {
	return auth.IssueAccessToken(
		secret,
		time.Duration(ttlSeconds)*time.Second,
		user.ID,
		user.Role,
		user.Status,
		session.ID,
		session.AuthSource,
		now,
	)
}

// validateActiveUserStatus ensures the account status can authenticate.
// validateActiveUserStatus 确保账号状态允许认证。
func validateActiveUserStatus(accountStatus string) error {
	switch accountStatus {
	case auth.UserStatusActive:
		return nil
	case auth.UserStatusPending:
		return errs.New(errs.CodeIdentityAccountPending, "account pending")
	case auth.UserStatusDisabled:
		return errs.New(errs.CodeIdentityAccountDisabled, "account disabled")
	case auth.UserStatusLocked:
		return errs.New(errs.CodeIdentityAccountLocked, "account locked")
	default:
		return errs.New(errs.CodeIdentityInvalidArgument, "account status is invalid")
	}
}

// writeAudit writes an identity audit record on a best-effort basis.
// writeAudit 以尽力而为方式写入 identity 审计记录。
func writeAudit(ctx context.Context, store *repo.Store, input auditInput) {
	if store == nil {
		return
	}

	record := &entity.IdentityAudit{
		UserID:     input.UserID,
		SessionID:  input.SessionID,
		Provider:   input.Provider,
		AuthSource: input.AuthSource,
		EventType:  input.EventType,
		Result:     input.Result,
		ClientIP:   input.ClientIP,
		UserAgent:  input.UserAgent,
		Detail:     input.Detail,
	}
	if err := store.IdentityAudits.Create(ctx, record); err != nil {
		logs.Ctx(ctx).Error(
			"identity_audit_write_failed",
			err,
			logs.String("event", input.EventType),
			logs.String("result", input.Result),
		)
	}
}

// stringPtr returns a pointer for non-empty strings.
// stringPtr 为非空字符串返回指针。
func stringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

// buildUniqueGitHubUsername builds a collision-free username candidate for GitHub sign-ins.
// buildUniqueGitHubUsername 为 GitHub 登录构建无冲突用户名候选值。
func buildUniqueGitHubUsername(ctx context.Context, store *repo.Store, preferred string) (string, error) {
	candidate := strings.TrimSpace(preferred)
	candidate = githubUsernameSanitizer.ReplaceAllString(candidate, "_")
	candidate = strings.Trim(candidate, "_")
	if len(candidate) < 3 {
		candidate = "github_user"
	}
	if len(candidate) > 32 {
		candidate = candidate[:32]
	}

	base := candidate
	for i := 0; i < 20; i++ {
		if _, err := auth.NormalizeUsername(candidate); err == nil {
			if _, err := store.Users.GetByUsername(ctx, candidate); err == nil {
				// occupied, continue below
			} else if repo.IsNotFound(err) {
				return candidate, nil
			} else {
				return "", err
			}
		}

		suffix := fmt.Sprintf("_%02d", i+1)
		maxBaseLen := 32 - len(suffix)
		truncatedBase := base
		if len(truncatedBase) > maxBaseLen {
			truncatedBase = truncatedBase[:maxBaseLen]
		}
		candidate = truncatedBase + suffix
	}

	return "", fmt.Errorf("unable to allocate unique username")
}

// withTransaction executes a service transaction using the shared store.
// withTransaction 使用共享 store 执行 service 事务。
func withTransaction(ctx context.Context, store *repo.Store, fn func(txStore *repo.Store) error) error {
	if store == nil || store.DB() == nil {
		return fmt.Errorf("service store is not initialized")
	}

	return store.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(store.WithTx(tx))
	})
}
