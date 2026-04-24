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
			return alreadyExistsFallback(action, route)
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

var alreadyExistsCodesByAction = map[string]errs.Code{
	"content_create":          errs.CodeContentSlugAlreadyExists,
	"content_update":          errs.CodeContentSlugAlreadyExists,
	"content_relation_create": errs.CodeContentRelationAlreadyExists,
	"content_tag_create":      errs.CodeContentTagAlreadyExists,
	"content_tag_update":      errs.CodeContentTagAlreadyExists,
}

var alreadyExistsCodesByRoute = map[string]errs.Code{
	"/api/v3/studio/content/items":                       errs.CodeContentSlugAlreadyExists,
	"/api/v3/studio/content/items/:content_id":           errs.CodeContentSlugAlreadyExists,
	"/api/v3/studio/content/items/:content_id/relations": errs.CodeContentRelationAlreadyExists,
	"/api/v3/studio/content/tags":                        errs.CodeContentTagAlreadyExists,
	"/api/v3/studio/content/tags/:tag_id":                errs.CodeContentTagAlreadyExists,
}

func alreadyExistsFallback(action string, route string) error {
	if code, ok := alreadyExistsCodesByAction[action]; ok {
		return alreadyExistsError(code)
	}
	if code, ok := alreadyExistsCodesByRoute[route]; ok {
		return alreadyExistsError(code)
	}
	return alreadyExistsError(errs.CodeContentSlugAlreadyExists)
}

func alreadyExistsError(code errs.Code) error {
	switch code {
	case errs.CodeContentRelationAlreadyExists:
		return errs.New(code, "content relation already exists")
	case errs.CodeContentTagAlreadyExists:
		return errs.New(code, "content tag already exists")
	default:
		return errs.New(errs.CodeContentSlugAlreadyExists, "content slug already exists")
	}
}
