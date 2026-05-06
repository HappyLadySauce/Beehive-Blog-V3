package service

import "time"

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"

	StatusPending  = "pending"
	StatusUploaded = "uploaded"
	StatusDeleted  = "deleted"
)

type CreateUploadInput struct {
	ActorUserID string
	CategoryKey string
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
	CategoryKey string
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
	CategoryKey string
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

type FileCategoryView struct {
	CategoryKey       string
	DisplayName       string
	Description       string
	Enabled           bool
	IsDefault         bool
	SortOrder         int32
	AllowedExtensions []string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ListFileCategoriesInput struct {
	IncludeDisabled bool
}

type CreateFileCategoryInput struct {
	CategoryKey       string
	DisplayName       string
	Description       string
	Enabled           bool
	IsDefault         bool
	SortOrder         int32
	AllowedExtensions []string
}

type UpdateFileCategoryInput struct {
	CategoryKey string
	DisplayName string
	Description string
	Enabled     bool
	SortOrder   int32
}

type UpdateFileCategoryExtensionsInput struct {
	CategoryKey       string
	AllowedExtensions []string
}
