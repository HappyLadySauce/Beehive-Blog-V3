package file

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/logic/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileCategoryExtensionsUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileCategoryExtensionsUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewFileCategoryExtensionsUpdateLogic(r.Context(), svcCtx)
		resp, err := l.FileCategoryExtensionsUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
