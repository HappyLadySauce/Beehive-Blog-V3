// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package identity

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/logic/identity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func IdentityUserPasswordResetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminResetUserPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := identity.NewIdentityUserPasswordResetLogic(r.Context(), svcCtx)
		resp, err := l.IdentityUserPasswordReset(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
