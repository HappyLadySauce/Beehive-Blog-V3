package logic

import (
	"context"
	"strconv"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
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
