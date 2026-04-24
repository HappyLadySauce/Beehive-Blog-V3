package entity

import "time"

type Item struct {
	ID                int64      `gorm:"column:id;primaryKey"`
	Type              string     `gorm:"column:type"`
	Title             string     `gorm:"column:title"`
	Slug              string     `gorm:"column:slug"`
	Status            string     `gorm:"column:status"`
	Visibility        string     `gorm:"column:visibility"`
	AIAccess          string     `gorm:"column:ai_access"`
	Summary           *string    `gorm:"column:summary"`
	CoverImageURL     *string    `gorm:"column:cover_image_url"`
	OwnerUserID       int64      `gorm:"column:owner_user_id"`
	AuthorUserID      int64      `gorm:"column:author_user_id"`
	SourceType        string     `gorm:"column:source_type"`
	CurrentRevisionID *int64     `gorm:"column:current_revision_id"`
	CommentEnabled    bool       `gorm:"column:comment_enabled"`
	IsFeatured        bool       `gorm:"column:is_featured"`
	SortOrder         int32      `gorm:"column:sort_order"`
	PublishedAt       *time.Time `gorm:"column:published_at"`
	ArchivedAt        *time.Time `gorm:"column:archived_at"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`
}

func (Item) TableName() string {
	return "content.items"
}
