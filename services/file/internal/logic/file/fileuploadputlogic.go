// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/storage"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/types"
)

type FileUploadPutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Upload a pending local file object
func NewFileUploadPutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadPutLogic {
	return &FileUploadPutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadPutLogic) FileUploadPut(req *types.FileUploadPutReq, body io.Reader, contentType string, contentLength int64) error {
	if l == nil || l.svcCtx == nil || l.svcCtx.Store == nil || l.svcCtx.LocalStorage == nil {
		return storageUnavailable()
	}
	uploadID := strings.TrimSpace(req.UploadId)
	asset, err := l.svcCtx.Store.Assets.FindByUploadID(l.ctx, uploadID)
	if err != nil {
		if repo.IsNotFound(err) {
			return dataPlaneError(http.StatusNotFound, "upload not found")
		}
		return dataPlaneError(http.StatusInternalServerError, "get upload failed")
	}
	if err := validatePendingUpload(time.Now(), asset); err != nil {
		return err
	}
	if !l.svcCtx.LocalStorage.VerifyUploadToken(uploadID, asset.ObjectKey, req.Token) {
		return dataPlaneError(http.StatusForbidden, "upload token is invalid")
	}
	if normalizeContentType(contentType) != asset.ContentType {
		return dataPlaneError(http.StatusPreconditionFailed, "content_type does not match upload")
	}
	if contentLength > asset.ByteSize {
		return dataPlaneError(http.StatusRequestEntityTooLarge, "file is too large")
	}
	info, err := l.svcCtx.LocalStorage.PutPending(l.ctx, asset.ObjectKey, body, asset.ByteSize)
	if err != nil {
		if errors.Is(err, storage.ErrStorageObjectTooLarge) {
			return dataPlaneError(http.StatusRequestEntityTooLarge, "file is too large")
		}
		if errors.Is(err, storage.ErrStorageDisabled) {
			return storageUnavailable()
		}
		return dataPlaneError(http.StatusServiceUnavailable, "write upload failed")
	}
	if info.ByteSize != asset.ByteSize {
		return dataPlaneError(http.StatusPreconditionFailed, "file byte_size does not match upload")
	}
	return nil
}
