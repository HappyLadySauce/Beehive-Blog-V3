// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package upload

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/logic/upload"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadAvatarPresignHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadAvatarPresignReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := upload.NewUploadAvatarPresignLogic(r.Context(), svcCtx)
		resp, err := l.UploadAvatarPresign(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
