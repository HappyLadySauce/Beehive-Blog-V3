package service

import (
	"strconv"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/google/uuid"
)

func objectKey(categoryKey string, ownerUserID int64, extension string) string {
	prefix := strings.ToLower(strings.TrimSpace(categoryKey))
	if prefix == "" {
		prefix = "default"
	}
	return prefix + "/" + strconv.FormatInt(ownerUserID, 10) + "/" + uuid.NewString() + normalizeExtension(extension)
}

func storageBucket(conf config.StorageConf) string {
	if strings.EqualFold(strings.TrimSpace(conf.Driver), "s3") {
		return strings.TrimSpace(conf.S3.Bucket)
	}
	return strings.TrimSpace(conf.Local.Bucket)
}

func publicURLForVisibility(conf config.StorageConf, visibility string, assetID string, objectKey string) string {
	if visibility != VisibilityPublic {
		return ""
	}
	baseURL := strings.TrimRight(strings.TrimSpace(conf.PublicBaseURL), "/")
	if baseURL == "" {
		return ""
	}
	if strings.EqualFold(strings.TrimSpace(conf.Driver), "s3") {
		objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
		if objectKey == "" {
			return ""
		}
		return baseURL + "/" + objectKey
	}
	return baseURL + "/" + strings.TrimSpace(assetID)
}
