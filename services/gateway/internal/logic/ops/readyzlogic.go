// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ops

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ReadyzLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadyzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadyzLogic {
	return &ReadyzLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadyzLogic) Readyz() (resp *types.ReadyzResp, err error) {
	if l.svcCtx == nil || l.svcCtx.IdentityProbe == nil || l.svcCtx.ContentProbe == nil {
		l.logger.Error(
			"readyz_check",
			errs.New(errs.CodeGatewayNotReady, "service is not ready"),
			logs.String("dependency", "gateway"),
			logs.String("reason", "probe_not_initialized"),
		)
		return &types.ReadyzResp{Status: "not_ready"}, errs.New(errs.CodeGatewayNotReady, "service is not ready")
	}

	probeCtx, cancel := context.WithTimeout(l.ctx, 2*time.Second)
	defer cancel()

	if err := l.svcCtx.IdentityProbe.Check(probeCtx); err != nil {
		l.logger.Error(
			"readyz_check",
			err,
			logs.String("dependency", "identity"),
		)
		return &types.ReadyzResp{Status: "not_ready"}, errs.Wrap(err, errs.CodeGatewayNotReady, "service is not ready")
	}
	if err := l.svcCtx.ContentProbe.Check(probeCtx); err != nil {
		l.logger.Error(
			"readyz_check",
			err,
			logs.String("dependency", "content"),
		)
		return &types.ReadyzResp{Status: "not_ready"}, errs.Wrap(err, errs.CodeGatewayNotReady, "service is not ready")
	}

	return &types.ReadyzResp{Status: "ready"}, nil
}
