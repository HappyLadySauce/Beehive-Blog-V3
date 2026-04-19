// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ops

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadyzLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadyzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadyzLogic {
	return &ReadyzLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadyzLogic) Readyz() (resp *types.ReadyzResp, err error) {
	// todo: add your logic here and delete this line

	return
}
