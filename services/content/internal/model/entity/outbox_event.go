package entity

import "time"

type OutboxEvent struct {
	ID           int64      `gorm:"column:id;primaryKey"`
	EventID      string     `gorm:"column:event_id"`
	EventType    string     `gorm:"column:event_type"`
	ResourceType string     `gorm:"column:resource_type"`
	ResourceID   int64      `gorm:"column:resource_id"`
	PayloadJSON  string     `gorm:"column:payload_json"`
	Status       string     `gorm:"column:status"`
	Attempts     int        `gorm:"column:attempts"`
	NextRetryAt  time.Time  `gorm:"column:next_retry_at"`
	LastError    string     `gorm:"column:last_error"`
	PublishedAt  *time.Time `gorm:"column:published_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (OutboxEvent) TableName() string {
	return "content.outbox_events"
}
