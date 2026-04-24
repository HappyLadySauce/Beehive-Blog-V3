package entity

import "time"

type Relation struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	FromContentID int64     `gorm:"column:from_content_id"`
	ToContentID   int64     `gorm:"column:to_content_id"`
	RelationType  string    `gorm:"column:relation_type"`
	Weight        int32     `gorm:"column:weight"`
	SortOrder     int32     `gorm:"column:sort_order"`
	MetadataJSON  *string   `gorm:"column:metadata_json"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (Relation) TableName() string {
	return "content.relations"
}
