package logic

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IntrospectAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewIntrospectAccessTokenLogic creates an IntrospectAccessTokenLogic instance.
// NewIntrospectAccessTokenLogic 创建 IntrospectAccessTokenLogic 实例。
func NewIntrospectAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IntrospectAccessTokenLogic {
	return &IntrospectAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// IntrospectAccessToken validates an access token against JWT claims and database state.
// IntrospectAccessToken 结合 JWT claims 和数据库状态校验 access token。
func (l *IntrospectAccessTokenLogic) IntrospectAccessToken(in *pb.IntrospectAccessTokenRequest) (*pb.IntrospectAccessTokenResponse, error) {
	// Validate the raw access token input.
	// 校验原始 access token 输入。
	accessToken := strings.TrimSpace(in.GetAccessToken())
	if accessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token is required")
	}

	claims, err := auth.ParseAccessToken(l.svcCtx.Config.Security.AccessTokenSecret, accessToken)
	if err != nil {
		l.Infof("access token introspection rejected invalid token: err=%v", err)
		return statusInactiveResponse(), nil
	}
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now().UTC()) {
		return statusInactiveResponse(), nil
	}
	if claims.SessionID <= 0 || claims.UserID <= 0 {
		return statusInactiveResponse(), nil
	}

	session, err := l.svcCtx.Store.UserSessions.GetByID(l.ctx, claims.SessionID)
	if err != nil {
		if repo.IsNotFound(err) {
			return statusInactiveResponse(), nil
		}
		return nil, status.Errorf(codes.Internal, "lookup session failed: %v", err)
	}
	if session.Status != auth.SessionStatusActive || session.ExpiresAt.Before(time.Now().UTC()) {
		return statusInactiveResponse(), nil
	}

	user, err := l.svcCtx.Store.Users.GetByID(l.ctx, claims.UserID)
	if err != nil {
		if repo.IsNotFound(err) {
			return statusInactiveResponse(), nil
		}
		return nil, status.Errorf(codes.Internal, "lookup user failed: %v", err)
	}
	if user.Status != auth.UserStatusActive {
		return statusInactiveResponse(), nil
	}

	return &pb.IntrospectAccessTokenResponse{
		Active:        true,
		UserId:        strconv.FormatInt(user.ID, 10),
		Role:          auth.ToProtoRole(user.Role),
		AccountStatus: auth.ToProtoAccountStatus(user.Status),
		SessionId:     strconv.FormatInt(session.ID, 10),
		AuthSource:    auth.ToProtoAuthSource(session.AuthSource),
		ExpiresAt:     claims.ExpiresAt.Time.Unix(),
	}, nil
}
