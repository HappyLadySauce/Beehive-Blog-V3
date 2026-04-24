package entity

import "time"

type Revision struct {
	ID                  int64     `gorm:"column:id;primaryKey"`
	ContentID           int64     `gorm:"column:content_id"`
	RevisionNo          int32     `gorm:"column:revision_no"`
	TitleSnapshot       string    `gorm:"column:title_snapshot"`
	SummarySnapshot     *string   `gorm:"column:summary_snapshot"`
	BodyMarkdown        string    `gorm:"column:body_markdown"`
	BodyJSON            *string   `gorm:"column:body_json"`
	EditorType          string    `gorm:"column:editor_type"`
	EditorUserID        *int64    `gorm:"column:editor_user_id"`
	EditorAgentClientID *int64    `gorm:"column:editor_agent_client_id"`
	ChangeSummary       *string   `gorm:"column:change_summary"`
	SourceType          string    `gorm:"column:source_type"`
	CreatedAt           time.Time `gorm:"column:created_at"`
}

func (Revision) TableName() string {
	return "content.revisions"
}
