// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ops

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type WsLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WsLogic {
	return &WsLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WsLogic) Ws() (resp *types.WsStubResp, err error) {
	l.logger.Info(
		"ws_stub",
		logs.String("route", "/ws"),
		logs.String("message", "realtime_not_implemented_yet"),
	)
	return &types.WsStubResp{Ok: false}, nil
}
