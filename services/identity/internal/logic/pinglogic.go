package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

// PingLogic adapts readiness probing to the service layer.
// PingLogic 将就绪探测适配到 service 层。
type PingLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPingLogic creates a PingLogic instance.
// NewPingLogic 创建 PingLogic 实例。
func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	logCtx := withLogContext(ctx)
	return &PingLogic{
		logger: logs.Ctx(logCtx),
		ctx:    logCtx,
		svcCtx: svcCtx,
	}
}

// Ping executes the identity readiness probe.
// Ping 执行 identity 就绪探测。
func (l *PingLogic) Ping(_ *pb.PingRequest) (*pb.PingResponse, error) {
	result, err := l.svcCtx.Services.Ping.Execute(l.ctx)
	if err != nil {
		return nil, toStatusError(err, "identity readiness probe failed")
	}

	return &pb.PingResponse{
		Ok:      result.OK,
		Service: result.Service,
		Version: result.Version,
	}, nil
}
