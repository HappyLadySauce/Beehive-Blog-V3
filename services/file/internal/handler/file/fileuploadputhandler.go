// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"net/http"

	filelogic "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/logic/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Upload a pending local file object
func FileUploadPutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !addUploadCORS(w, r, localStorageConf(svcCtx)) {
			http.Error(w, "origin is not allowed", http.StatusForbidden)
			return
		}
		var req types.FileUploadPutReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := filelogic.NewFileUploadPutLogic(r.Context(), svcCtx)
		if err := l.FileUploadPut(&req, r.Body, r.Header.Get("Content-Type"), r.ContentLength); err != nil {
			writeDataPlaneError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
