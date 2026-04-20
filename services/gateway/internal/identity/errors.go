package identity

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MapUpstreamError maps upstream gRPC errors into unified domain errors.
// MapUpstreamError 将上游 gRPC 错误映射为统一领域错误。
func MapUpstreamError(ctx context.Context, action string, route string, err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if ok {
		logs.Ctx(ctx).Error(
			action,
			err,
			logs.String("route", route),
			logs.String("upstream_code", st.Code().String()),
		)
		if parsed := errgrpcx.FromStatus(err, 0, ""); parsed != nil {
			if domainErr := errs.Parse(parsed); domainErr != nil {
				return domainErr
			}
		}

		switch st.Code() {
		case codes.InvalidArgument:
			return errs.New(errs.CodeGatewayBadRequest, "bad request")
		case codes.Unauthenticated:
			return errs.New(errs.CodeGatewayAccessTokenInvalid, "access token is invalid")
		case codes.PermissionDenied:
			return errs.New(errs.CodeGatewayAccessForbidden, "access is forbidden")
		case codes.FailedPrecondition:
			return errs.New(errs.CodeGatewayBadRequest, "request precondition failed")
		case codes.Unavailable:
			return errs.New(errs.CodeGatewayAuthServiceUnavailable, "authentication service is unavailable")
		case codes.DeadlineExceeded:
			return errs.New(errs.CodeGatewayUpstreamTimeout, "upstream service timed out")
		default:
			return errs.New(errs.CodeGatewayInternal, "internal server error")
		}
	}

	logs.Ctx(ctx).Error(
		action,
		err,
		logs.String("route", route),
	)
	return errs.New(errs.CodeGatewayInternal, "internal server error")
}
