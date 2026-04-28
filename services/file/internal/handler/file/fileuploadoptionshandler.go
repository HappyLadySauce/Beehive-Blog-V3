// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
)

// Preflight local file upload
func FileUploadOptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !addUploadCORS(w, r, localStorageConf(svcCtx)) {
			http.Error(w, "origin is not allowed", http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
