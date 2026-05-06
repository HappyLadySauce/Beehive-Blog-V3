package entity

import "time"

type FileCategory struct {
	CategoryKey string    `gorm:"column:category_key;primaryKey"`
	DisplayName string    `gorm:"column:display_name"`
	Description string    `gorm:"column:description"`
	Enabled     bool      `gorm:"column:enabled"`
	IsDefault   bool      `gorm:"column:is_default"`
	SortOrder   int32     `gorm:"column:sort_order"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (FileCategory) TableName() string {
	return "file_categories"
}

type FileCategoryExtension struct {
	CategoryKey string    `gorm:"column:category_key;primaryKey"`
	Extension   string    `gorm:"column:extension;primaryKey"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (FileCategoryExtension) TableName() string {
	return "file_category_extensions"
}
