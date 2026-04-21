package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errhttpx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/httpx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const authorizationHeader = "Authorization"

var (
	errAuthorizationRequired = errors.New("authorization is required")
	errInvalidAuthorization  = errors.New("invalid authorization scheme")
	errAccessTokenRequired   = errors.New("access token is required")
)

// AuthMiddleware validates bearer access token through identity introspection.
// AuthMiddleware 通过 identity introspect 校验 Bearer access token。
type AuthMiddleware struct {
	identityClient pb.IdentityClient
	tokenPrefix    string
	internalAuth   ctxmeta.InternalRPCAuth
}

// NewAuthMiddleware creates auth middleware.
// NewAuthMiddleware 创建鉴权中间件。
func NewAuthMiddleware(identityClient pb.IdentityClient, securityConf config.GatewaySecurityConf, rpcConf config.IdentityRPCConf) *AuthMiddleware {
	tokenPrefix := strings.TrimSpace(securityConf.TokenPrefix)
	if tokenPrefix == "" {
		tokenPrefix = "Bearer"
	}
	return &AuthMiddleware{
		identityClient: identityClient,
		tokenPrefix:    tokenPrefix,
		internalAuth: ctxmeta.InternalRPCAuth{
			Token:  rpcConf.InternalAuthToken,
			Caller: rpcConf.InternalCallerName,
		},
	}
}

// Handle authenticates request and injects trusted identity context.
// Handle 完成鉴权并注入可信身份上下文。
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		headerValue := strings.TrimSpace(r.Header.Get(authorizationHeader))
		accessToken, err := parseBearerToken(headerValue, m.tokenPrefix)
		if err != nil {
			errhttpx.WriteError(ctx, w, mapAuthMiddlewareError(err), requestIDFromContext(ctx))
			return
		}

		requestMeta, _ := RequestMetaFrom(ctx)
		rpcCtx := ctxmeta.BuildIdentityOutgoingContext(ctx, requestMeta, m.internalAuth)

		introspectResp, introspectErr := m.identityClient.IntrospectAccessToken(rpcCtx, &pb.IntrospectAccessTokenRequest{
			AccessToken: accessToken,
		})
		if introspectErr != nil {
			upstreamCode := "unknown"
			if st, ok := status.FromError(introspectErr); ok {
				upstreamCode = st.Code().String()
			}
			logs.Ctx(ctx).Error(
				"auth_introspect",
				introspectErr,
				logs.String("route", r.URL.Path),
				logs.String("upstream_code", upstreamCode),
			)
			errhttpx.WriteError(ctx, w, mapAuthMiddlewareError(introspectErr), requestIDFromContext(ctx))
			return
		}
		if introspectResp == nil || !introspectResp.GetActive() {
			errhttpx.WriteError(ctx, w, errs.New(errs.CodeGatewayAccessTokenInactive, "access token is inactive"), requestIDFromContext(ctx))
			return
		}

		authCtx := AuthContext{
			UserID:        introspectResp.GetUserId(),
			SessionID:     introspectResp.GetSessionId(),
			Role:          introspectResp.GetRole().String(),
			AccountStatus: introspectResp.GetAccountStatus().String(),
			AuthSource:    introspectResp.GetAuthSource().String(),
			AccessToken:   accessToken,
		}

		logs.Ctx(ctx).Info(
			"auth_introspect",
			logs.String("route", r.URL.Path),
			logs.UserID(authCtx.UserID),
			logs.SessionID(authCtx.SessionID),
		)

		next(w, r.WithContext(WithAuthContext(ctx, authCtx)))
	}
}

func parseBearerToken(headerValue, tokenPrefix string) (string, error) {
	if headerValue == "" {
		return "", errAuthorizationRequired
	}

	expectedPrefix := tokenPrefix + " "
	if !strings.HasPrefix(strings.ToLower(headerValue), strings.ToLower(expectedPrefix)) {
		return "", errInvalidAuthorization
	}

	token := strings.TrimSpace(headerValue[len(expectedPrefix):])
	if token == "" {
		return "", errAccessTokenRequired
	}
	return token, nil
}

func mapAuthMiddlewareError(err error) error {
	switch {
	case errors.Is(err, errAuthorizationRequired):
		return errs.New(errs.CodeGatewayAuthorizationRequired, "authorization is required")
	case errors.Is(err, errInvalidAuthorization):
		return errs.New(errs.CodeGatewayInvalidAuthorizationScheme, "invalid authorization scheme")
	case errors.Is(err, errAccessTokenRequired):
		return errs.New(errs.CodeGatewayAuthorizationRequired, "access token is required")
	}

	st, ok := status.FromError(err)
	if !ok {
		return errs.New(errs.CodeGatewayAuthServiceUnavailable, "authentication service is unavailable")
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return errs.New(errs.CodeGatewayAccessTokenInvalid, "access token is invalid")
	case codes.PermissionDenied:
		return errs.New(errs.CodeGatewayAccessForbidden, "access is forbidden")
	case codes.Unavailable, codes.DeadlineExceeded:
		return errs.New(errs.CodeGatewayAuthServiceUnavailable, "authentication service is unavailable")
	default:
		return errs.New(errs.CodeGatewayAuthServiceUnavailable, "authentication service is unavailable")
	}
}

func requestIDFromContext(ctx context.Context) string {
	requestMeta, ok := RequestMetaFrom(ctx)
	if !ok {
		return ""
	}

	return requestMeta.RequestID
}
