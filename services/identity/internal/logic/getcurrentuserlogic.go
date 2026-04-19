package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetCurrentUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewGetCurrentUserLogic creates a GetCurrentUserLogic instance.
// NewGetCurrentUserLogic 创建 GetCurrentUserLogic 实例。
func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCurrentUser returns the trusted current user snapshot.
// GetCurrentUser 返回可信的当前用户快照。
func (l *GetCurrentUserLogic) GetCurrentUser(in *pb.GetCurrentUserRequest) (*pb.GetCurrentUserResponse, error) {
	// Parse the trusted user identifier from the upstream transport layer.
	// 从上游传输层解析可信用户标识。
	userID, err := parseID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.Store.Users.GetByID(l.ctx, userID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "lookup current user failed: %v", err)
	}

	return &pb.GetCurrentUserResponse{
		CurrentUser: auth.ToCurrentUser(user),
	}, nil
}
