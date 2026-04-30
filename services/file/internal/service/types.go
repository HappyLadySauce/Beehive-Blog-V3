package service

import "time"

const (
	ScopeAvatar       = "avatar"
	ScopeContentCover = "content_cover"
	ScopeContentImage = "content_image"
	ScopeAttachment   = "attachment"

	VisibilityPublic  = "public"
	VisibilityPrivate = "private"

	StatusPending  = "pending"
	StatusUploaded = "uploaded"
	StatusDeleted  = "deleted"
)

type CreateUploadInput struct {
	ActorUserID string
	Scope       string
	FileName    string
	ContentType string
	ByteSize    int64
	Visibility  string
}

type CreateUploadResult struct {
	Asset     *AssetView
	UploadURL string
	Headers   map[string]string
	ExpiresAt time.Time
	MaxBytes  int64
}

type ListAssetsInput struct {
	ActorUserID string
	Scope       string
	Status      string
	Visibility  string
	OwnerUserID string
	Keyword     string
	Page        int
	PageSize    int
}

type AssetListResult struct {
	Items    []*AssetView
	Total    int64
	Page     int
	PageSize int
}

type AssetView struct {
	AssetID     string
	UploadID    string
	OwnerUserID string
	Scope       string
	Visibility  string
	Status      string
	Bucket      string
	ObjectKey   string
	PublicURL   string
	FileName    string
	ContentType string
	ByteSize    int64
	CreatedAt   time.Time
	ExpiresAt   time.Time
	UploadedAt  *time.Time
	DeletedAt   *time.Time
}
