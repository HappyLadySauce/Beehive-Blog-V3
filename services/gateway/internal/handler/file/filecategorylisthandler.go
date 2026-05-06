package file

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/logic/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileCategoryListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewFileCategoryListLogic(r.Context(), svcCtx)
		resp, err := l.FileCategoryList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
