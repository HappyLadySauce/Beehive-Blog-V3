package logic

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var githubUsernameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// auditInput describes a single identity audit write.
// auditInput 描述一条 identity 审计写入请求。
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

// parseID parses a decimal ID from trusted string input.
// parseID 从可信字符串输入中解析十进制 ID。
func parseID(fieldName, raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, status.Errorf(codes.InvalidArgument, "%s is invalid", fieldName)
	}

	return id, nil
}

// validateActiveUserStatus ensures the account status can authenticate.
// validateActiveUserStatus 确保账号状态允许认证。
func validateActiveUserStatus(accountStatus string) error {
	switch accountStatus {
	case auth.UserStatusActive:
		return nil
	case auth.UserStatusPending:
		return status.Error(codes.FailedPrecondition, "account_pending")
	case auth.UserStatusDisabled:
		return status.Error(codes.FailedPrecondition, "account_disabled")
	case auth.UserStatusLocked:
		return status.Error(codes.FailedPrecondition, "account_locked")
	default:
		return status.Error(codes.FailedPrecondition, "account_status_invalid")
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
		logx.WithContext(ctx).Errorf("identity audit write failed: event=%s result=%s err=%v", input.EventType, input.Result, err)
	}
}

// stringPtr returns a pointer for non-empty strings.
// stringPtr 为非空字符串返回指针。
func stringPtr(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

// statusInactiveResponse returns a consistent inactive token introspection response.
// statusInactiveResponse 返回统一的 inactive introspection 响应。
func statusInactiveResponse() *pb.IntrospectAccessTokenResponse {
	return &pb.IntrospectAccessTokenResponse{Active: false}
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
		if len(base) > maxBaseLen {
			base = base[:maxBaseLen]
		}
		candidate = base + suffix
	}

	return "", fmt.Errorf("unable to allocate unique username")
}
