package content

import (
	"context"
	"errors"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapUpstreamError(ctx context.Context, action string, route string, err error) error {
	if err == nil {
		return nil
	}

	if st, ok := status.FromError(err); ok {
		logs.Ctx(ctx).Error(
			action,
			err,
			logs.String("route", route),
			logs.String("upstream_code", st.Code().String()),
		)
		if parsed := errgrpcx.FromStatus(err, 0, ""); parsed != nil {
			if domainErr := errs.Parse(parsed); domainErr != nil {
				if errors.Is(domainErr, errs.E(errs.CodeContentInternalCallerUnauthorized)) {
					return errs.New(errs.CodeGatewayNotReady, "content service is not ready")
				}
				return domainErr
			}
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return errs.New(errs.CodeContentInvalidArgument, "invalid argument")
		case codes.Unauthenticated:
			return errs.New(errs.CodeGatewayAccessTokenInvalid, "access token is invalid")
		case codes.PermissionDenied:
			return errs.New(errs.CodeContentAccessForbidden, "content access forbidden")
		case codes.NotFound:
			return errs.New(errs.CodeContentNotFound, "content not found")
		case codes.AlreadyExists:
			return errs.New(errs.CodeContentSlugAlreadyExists, "resource already exists")
		case codes.FailedPrecondition:
			return errs.New(errs.CodeContentInvalidTransition, "request precondition failed")
		case codes.Unavailable:
			return errs.New(errs.CodeGatewayContentServiceUnavailable, "content service is unavailable")
		case codes.DeadlineExceeded:
			return errs.New(errs.CodeGatewayUpstreamTimeout, "upstream service timed out")
		default:
			return errs.New(errs.CodeContentInternal, "content internal error")
		}
	}

	logs.Ctx(ctx).Error(action, err, logs.String("route", route))
	return errs.New(errs.CodeContentInternal, "content internal error")
}
