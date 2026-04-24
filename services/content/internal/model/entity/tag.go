package entity

import "time"

type Tag struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name"`
	Slug        string    `gorm:"column:slug"`
	Description *string   `gorm:"column:description"`
	Color       *string   `gorm:"column:color"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Tag) TableName() string {
	return "content.tags"
}

type ContentTag struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	ContentID int64     `gorm:"column:content_id"`
	TagID     int64     `gorm:"column:tag_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (ContentTag) TableName() string {
	return "content.content_tags"
}
