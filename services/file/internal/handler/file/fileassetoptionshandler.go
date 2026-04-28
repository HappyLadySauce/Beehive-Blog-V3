// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
)

// Preflight public asset read
func FileAssetOptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addDataPlaneCORS(w)
		w.WriteHeader(http.StatusNoContent)
	}
}
