package file

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
	fileservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/storage"
)

type DataPlaneError struct {
	Status  int
	Message string
}

func (e DataPlaneError) Error() string {
	return e.Message
}

func (e DataPlaneError) HTTPStatus() int {
	return e.Status
}

func (e DataPlaneError) PublicMessage() string {
	return e.Message
}

type AssetStream struct {
	Reader      io.ReadCloser
	ContentType string
	ByteSize    int64
}

type AssetHead struct {
	ContentType string
	ByteSize    int64
}

func dataPlaneError(status int, message string) error {
	return DataPlaneError{Status: status, Message: message}
}

func storageUnavailable() error {
	return dataPlaneError(http.StatusServiceUnavailable, "file storage is unavailable")
}

func normalizeContentType(contentType string) string {
	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

func loadPublicUploadedAsset(ctx context.Context, store interface {
	FindByAssetID(context.Context, string) (*entity.FileAsset, error)
}, assetID string) (*entity.FileAsset, error) {
	asset, err := store.FindByAssetID(ctx, assetID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, dataPlaneError(http.StatusNotFound, "asset not found")
		}
		return nil, dataPlaneError(http.StatusInternalServerError, "get asset failed")
	}
	if asset.Status != fileservice.StatusUploaded ||
		asset.DeletedAt != nil ||
		asset.Visibility != fileservice.VisibilityPublic ||
		strings.TrimSpace(asset.PublicURL) == "" {
		return nil, dataPlaneError(http.StatusNotFound, "asset not found")
	}
	return asset, nil
}

func mapStorageReadError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, storage.ErrStorageDisabled) {
		return storageUnavailable()
	}
	return dataPlaneError(http.StatusNotFound, "asset not found")
}

func validatePendingUpload(now time.Time, asset *entity.FileAsset) error {
	if asset.Status != fileservice.StatusPending || now.UTC().After(asset.ExpiresAt) {
		return dataPlaneError(http.StatusPreconditionFailed, "upload state is invalid")
	}
	return nil
}
