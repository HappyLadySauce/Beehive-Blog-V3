package file

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
)

type publicStatusError interface {
	error
	HTTPStatus() int
	PublicMessage() string
}

func writeDataPlaneError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	var dataErr publicStatusError
	if errors.As(err, &dataErr) {
		http.Error(w, dataErr.PublicMessage(), dataErr.HTTPStatus())
		return
	}
	http.Error(w, "file service failed", http.StatusInternalServerError)
}

func addPublicReadCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func addUploadCORS(w http.ResponseWriter, r *http.Request, c config.LocalStorageConf) bool {
	w.Header().Add("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Upload-Token")

	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}
	for _, allowed := range c.AllowedOrigins {
		if origin == strings.TrimSpace(allowed) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			return true
		}
	}
	return false
}

func localStorageConf(svcCtx *svc.ServiceContext) config.LocalStorageConf {
	if svcCtx == nil {
		return config.LocalStorageConf{}
	}
	return svcCtx.Config.Storage.Local
}

func writeAssetHeaders(w http.ResponseWriter, contentType string, byteSize int64) {
	w.Header().Set("Content-Type", strings.TrimSpace(contentType))
	w.Header().Set("Content-Length", strconv.FormatInt(byteSize, 10))
}

// logStreamCopyError keeps client aborts low-noise while surfacing source read failures.
// logStreamCopyError 低噪处理客户端断连，同时暴露源文件读取失败。
func logStreamCopyError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	if isClientStreamAbort(err) {
		logs.Ctx(ctx).Warn("file_asset_stream_client_aborted")
		return
	}
	logs.Ctx(ctx).Error("file_asset_stream_failed", err)
}

func isClientStreamAbort(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, net.ErrClosed) ||
		isNetworkStreamAbort(err)
}

func isNetworkStreamAbort(err error) bool {
	var netErr *net.OpError
	if !errors.As(err, &netErr) {
		return false
	}
	var syscallErr *os.SyscallError
	if !errors.As(netErr, &syscallErr) {
		return false
	}
	return isStreamAbortSyscall(syscallErr.Err)
}
