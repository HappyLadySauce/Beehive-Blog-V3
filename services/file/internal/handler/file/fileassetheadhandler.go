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

// Read public uploaded asset metadata
func FileAssetHeadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addDataPlaneCORS(w)
		var req types.FileAssetReadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := filelogic.NewFileAssetHeadLogic(r.Context(), svcCtx)
		head, err := l.FileAssetHead(&req)
		if err != nil {
			writeDataPlaneError(w, err)
			return
		}
		writeAssetHeaders(w, head.ContentType, head.ByteSize)
	}
}
