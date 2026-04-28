package entity

import "time"

type FileAsset struct {
	AssetID     string     `gorm:"column:asset_id;primaryKey"`
	UploadID    string     `gorm:"column:upload_id"`
	OwnerUserID int64      `gorm:"column:owner_user_id"`
	Scope       string     `gorm:"column:scope"`
	Visibility  string     `gorm:"column:visibility"`
	Status      string     `gorm:"column:status"`
	Bucket      string     `gorm:"column:bucket"`
	ObjectKey   string     `gorm:"column:object_key"`
	PublicURL   string     `gorm:"column:public_url"`
	FileName    string     `gorm:"column:file_name"`
	ContentType string     `gorm:"column:content_type"`
	ByteSize    int64      `gorm:"column:byte_size"`
	ExpiresAt   time.Time  `gorm:"column:expires_at"`
	UploadedAt  *time.Time `gorm:"column:uploaded_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (FileAsset) TableName() string {
	return "file_assets"
}
