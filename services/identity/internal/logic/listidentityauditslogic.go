package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

type ListIdentityAuditsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewListIdentityAuditsLogic creates a ListIdentityAuditsLogic instance.
// NewListIdentityAuditsLogic 创建 ListIdentityAuditsLogic 实例。
func NewListIdentityAuditsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListIdentityAuditsLogic {
	logCtx := withLogContext(ctx)
	return &ListIdentityAuditsLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// ListIdentityAudits adapts the gRPC request to the audit query service.
// ListIdentityAudits 将 gRPC 请求适配到审计查询 service。
func (l *ListIdentityAuditsLogic) ListIdentityAudits(in *pb.ListIdentityAuditsRequest) (*pb.ListIdentityAuditsResponse, error) {
	actorUserID, err := parseID("actor_user_id", in.GetActorUserId())
	if err != nil {
		return nil, err
	}
	userID, err := parseOptionalID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}
	startedAt, err := unixSecondsPtr("started_at", in.GetStartedAt())
	if err != nil {
		return nil, err
	}
	endedAt, err := unixSecondsPtr("ended_at", in.GetEndedAt())
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Services.Users.ListIdentityAudits(l.ctx, identityservice.ListIdentityAuditsInput{
		ActorUserID: actorUserID,
		EventType:   in.GetEventType(),
		Result:      in.GetResult(),
		UserID:      userID,
		StartedAt:   startedAt,
		EndedAt:     endedAt,
		Page:        int(in.GetPage()),
		PageSize:    int(in.GetPageSize()),
	})
	if err != nil {
		return nil, toStatusError(err, "list identity audits failed")
	}

	items := make([]*pb.IdentityAuditView, 0, len(result.Items))
	for i := range result.Items {
		items = append(items, auth.ToIdentityAuditView(&result.Items[i]))
	}
	return &pb.ListIdentityAuditsResponse{Items: items, Total: result.Total, Page: int32(result.Page), PageSize: int32(result.PageSize)}, nil
}
