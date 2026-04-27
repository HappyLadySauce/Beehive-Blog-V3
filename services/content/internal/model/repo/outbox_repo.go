package repo

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	OutboxStatusPending    = "pending"
	OutboxStatusProcessing = "processing"
	OutboxStatusDone       = "done"
	OutboxStatusFailed     = "failed"
)

type OutboxRepository struct {
	db *gorm.DB
}

func (r *OutboxRepository) Create(ctx context.Context, event *entity.OutboxEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *OutboxRepository) ClaimDue(ctx context.Context, now time.Time, batchSize int, processingTimeout time.Duration) ([]entity.OutboxEvent, error) {
	var events []entity.OutboxEvent
	if processingTimeout <= 0 {
		processingTimeout = time.Minute
	}
	claimedAt := normalizeOutboxLeaseTime(now)
	staleBefore := claimedAt.Add(-processingTimeout)
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where(
				"(status = ? AND next_retry_at <= ?) OR (status = ? AND updated_at <= ?)",
				OutboxStatusPending,
				claimedAt,
				OutboxStatusProcessing,
				staleBefore,
			).
			Order("next_retry_at ASC, id ASC").
			Limit(batchSize).
			Find(&events).Error; err != nil {
			return err
		}
		if len(events) == 0 {
			return nil
		}
		ids := make([]int64, 0, len(events))
		for _, event := range events {
			ids = append(ids, event.ID)
		}
		return tx.Model(&entity.OutboxEvent{}).
			Where("id IN ?", ids).
			Updates(map[string]any{
				"status":     OutboxStatusProcessing,
				"updated_at": claimedAt,
			}).Error
	})
	for i := range events {
		events[i].Status = OutboxStatusProcessing
		events[i].UpdatedAt = claimedAt
	}
	return events, err
}

func (r *OutboxRepository) MarkDone(ctx context.Context, id int64, claimedAt time.Time, now time.Time) error {
	result := r.db.WithContext(ctx).Model(&entity.OutboxEvent{}).
		Where("id = ? AND status = ? AND updated_at = ?", id, OutboxStatusProcessing, normalizeOutboxLeaseTime(claimedAt)).
		Updates(map[string]any{
			"status":       OutboxStatusDone,
			"published_at": now,
			"updated_at":   now,
			"last_error":   "",
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OutboxRepository) MarkPublishFailed(ctx context.Context, id int64, claimedAt time.Time, attempts int, maxAttempts int, nextRetryAt time.Time, lastError string, now time.Time) error {
	status := OutboxStatusPending
	if maxAttempts > 0 && attempts >= maxAttempts {
		status = OutboxStatusFailed
	}
	result := r.db.WithContext(ctx).Model(&entity.OutboxEvent{}).
		Where("id = ? AND status = ? AND updated_at = ?", id, OutboxStatusProcessing, normalizeOutboxLeaseTime(claimedAt)).
		Updates(map[string]any{
			"status":        status,
			"attempts":      attempts,
			"next_retry_at": nextRetryAt,
			"last_error":    lastError,
			"updated_at":    now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OutboxRepository) ListByResource(ctx context.Context, resourceType string, resourceID int64) ([]entity.OutboxEvent, error) {
	var events []entity.OutboxEvent
	err := r.db.WithContext(ctx).
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("id ASC").
		Find(&events).Error
	return events, err
}

func (r *OutboxRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.OutboxEvent{}).Where("status = ?", status).Count(&total).Error
	return total, err
}

func normalizeOutboxLeaseTime(value time.Time) time.Time {
	return value.UTC().Round(time.Microsecond)
}
