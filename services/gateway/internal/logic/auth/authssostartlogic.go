// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthSsoStartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthSsoStartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthSsoStartLogic {
	return &AuthSsoStartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthSsoStartLogic) AuthSsoStart(req *types.AuthSsoStartReq) (resp *types.AuthSsoStartResp, err error) {
	// todo: add your logic here and delete this line

	return
}
