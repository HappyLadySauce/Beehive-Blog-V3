package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type PingLogic struct {
	baseLogic
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *PingLogic) Ping(_ *pb.PingRequest) (*pb.PingResponse, error) {
	if err := l.svcCtx.Services.Ping(l.ctx); err != nil {
		return nil, toStatus(err, "file service is not ready")
	}
	return &pb.PingResponse{Ok: true, Service: "file"}, nil
}
