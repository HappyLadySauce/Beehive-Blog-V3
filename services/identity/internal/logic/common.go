package logic

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// parseID parses a decimal identifier from trusted upstream input.
// parseID 从可信上游输入中解析十进制标识。
func parseID(fieldName, raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, status.Errorf(codes.InvalidArgument, "%s is invalid", fieldName)
	}

	return id, nil
}

// parseOptionalID parses an optional decimal identifier.
// parseOptionalID 解析可选的十进制标识。
func parseOptionalID(fieldName, raw string) (*int64, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	id, err := parseID(fieldName, raw)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// unixSecondsPtr converts optional unix seconds to a time pointer.
// unixSecondsPtr 将可选 Unix 秒转换为时间指针。
func unixSecondsPtr(fieldName string, value int64) (*time.Time, error) {
	if value == 0 {
		return nil, nil
	}
	if value < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s is invalid", fieldName)
	}

	t := time.Unix(value, 0).UTC()
	return &t, nil
}

// expiresIn returns the remaining seconds until expiration.
// expiresIn 返回距离过期的剩余秒数。
func expiresIn(expiresAt time.Time) int64 {
	seconds := int64(time.Until(expiresAt).Seconds())
	if seconds < 0 {
		return 0
	}

	return seconds
}

// toStatusError converts a service-layer error into a gRPC status error.
// toStatusError 将 service 层错误转换为 gRPC status 错误。
func toStatusError(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}
	return errgrpcx.ToStatus(err, fallbackMessage)
}

// withLogContext enriches a logic context with request-scoped log metadata.
// withLogContext 使用请求级日志元数据增强 logic 上下文。
func withLogContext(ctx context.Context) context.Context {
	return logs.WithRequestID(ctx, ctxmeta.GetRequestIDFromIncomingContext(ctx))
}

// parseActorTarget parses admin actor and target user identifiers.
// parseActorTarget 解析管理员操作者与目标用户标识。
func parseActorTarget(actorRaw, targetRaw string) (int64, int64, error) {
	actorUserID, err := parseID("actor_user_id", actorRaw)
	if err != nil {
		return 0, 0, err
	}
	targetUserID, err := parseID("target_user_id", targetRaw)
	if err != nil {
		return 0, 0, err
	}
	return actorUserID, targetUserID, nil
}

// roleString maps a proto role to the service role string.
// roleString 将 proto 角色映射为 service 层角色字符串。
func roleString(value pb.Role) string {
	switch value {
	case pb.Role_ROLE_ADMIN:
		return auth.UserRoleAdmin
	case pb.Role_ROLE_MEMBER:
		return auth.UserRoleMember
	default:
		return ""
	}
}

// accountStatusString maps a proto account status to the service status string.
// accountStatusString 将 proto 账号状态映射为 service 层状态字符串。
func accountStatusString(value pb.AccountStatus) string {
	switch value {
	case pb.AccountStatus_ACCOUNT_STATUS_PENDING:
		return auth.UserStatusPending
	case pb.AccountStatus_ACCOUNT_STATUS_ACTIVE:
		return auth.UserStatusActive
	case pb.AccountStatus_ACCOUNT_STATUS_DISABLED:
		return auth.UserStatusDisabled
	case pb.AccountStatus_ACCOUNT_STATUS_LOCKED:
		return auth.UserStatusLocked
	case pb.AccountStatus_ACCOUNT_STATUS_DELETED:
		return auth.UserStatusDeleted
	default:
		return ""
	}
}
