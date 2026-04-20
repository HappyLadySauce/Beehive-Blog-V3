package logic

import (
	"strconv"
	"time"

	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
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
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	switch {
	case identityservice.IsKind(err, identityservice.ErrorKindInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case identityservice.IsKind(err, identityservice.ErrorKindUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case identityservice.IsKind(err, identityservice.ErrorKindAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case identityservice.IsKind(err, identityservice.ErrorKindFailedPrecondition):
		return status.Error(codes.FailedPrecondition, err.Error())
	case identityservice.IsKind(err, identityservice.ErrorKindNotFound):
		return status.Error(codes.NotFound, err.Error())
	case identityservice.IsKind(err, identityservice.ErrorKindUnimplemented):
		return status.Error(codes.Unimplemented, err.Error())
	default:
		if fallbackMessage == "" {
			fallbackMessage = "internal error"
		}
		return status.Errorf(codes.Internal, "%s: %v", fallbackMessage, err)
	}
}
