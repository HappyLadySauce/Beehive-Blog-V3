package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

type ListUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewListUsersLogic creates a ListUsersLogic instance.
// NewListUsersLogic 创建 ListUsersLogic 实例。
func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	logCtx := withLogContext(ctx)
	return &ListUsersLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// ListUsers adapts the gRPC request to the user management service.
// ListUsers 将 gRPC 请求适配到用户管理 service。
func (l *ListUsersLogic) ListUsers(in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	actorUserID, err := parseID("actor_user_id", in.GetActorUserId())
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.Services.Users.ListUsers(l.ctx, identityservice.ListUsersInput{
		ActorUserID:    actorUserID,
		Keyword:        in.GetKeyword(),
		Role:           roleString(in.GetRole()),
		Status:         accountStatusString(in.GetStatus()),
		IncludeDeleted: in.GetIncludeDeleted(),
		Page:           int(in.GetPage()),
		PageSize:       int(in.GetPageSize()),
	})
	if err != nil {
		return nil, toStatusError(err, "list users failed")
	}

	items := make([]*pb.AdminUserView, 0, len(result.Items))
	for i := range result.Items {
		items = append(items, auth.ToAdminUserView(&result.Items[i]))
	}
	return &pb.ListUsersResponse{Items: items, Total: result.Total, Page: int32(result.Page), PageSize: int32(result.PageSize)}, nil
}
