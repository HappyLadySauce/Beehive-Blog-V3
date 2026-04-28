// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"io"
	"net/http"

	filelogic "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/logic/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Read a public uploaded asset
func FileAssetGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addDataPlaneCORS(w)
		var req types.FileAssetReadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := filelogic.NewFileAssetGetLogic(r.Context(), svcCtx)
		stream, err := l.FileAssetGet(&req)
		if err != nil {
			writeDataPlaneError(w, err)
			return
		}
		defer stream.Reader.Close()
		writeAssetHeaders(w, stream.ContentType, stream.ByteSize)
		_, _ = io.Copy(w, stream.Reader)
	}
}
