package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"github.com/google/uuid"
)

const (
	EventContentCreated           = "content.created"
	EventContentUpdated           = "content.updated"
	EventContentArchived          = "content.archived"
	EventContentStatusChanged     = "content.status_changed"
	EventContentVisibilityChanged = "content.visibility_changed"
	EventContentAIAccessChanged   = "content.ai_access_changed"
	EventContentTagChanged        = "content.tag_changed"
	EventContentRelationChanged   = "content.relation_changed"

	EventResourceContentItem = "content_item"
)

type outboxEventInput struct {
	EventType    string
	ResourceType string
	ResourceID   int64
	Payload      map[string]any
	OccurredAt   time.Time
}

func writeOutboxEvent(ctx context.Context, store *repo.Store, input outboxEventInput) error {
	if store == nil || store.Outbox == nil {
		return serviceNotInitialized()
	}
	if input.OccurredAt.IsZero() {
		input.OccurredAt = time.Now().UTC()
	}
	if input.ResourceType == "" {
		input.ResourceType = EventResourceContentItem
	}
	payload := clonePayload(input.Payload)
	payload["event_id"] = uuid.NewString()
	payload["event_type"] = input.EventType
	payload["resource_type"] = input.ResourceType
	payload["resource_id"] = strconv.FormatInt(input.ResourceID, 10)
	payload["occurred_at"] = input.OccurredAt.UTC().Format(time.RFC3339Nano)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return errs.Wrap(err, errs.CodeContentInternal, "content internal error")
	}
	event := &entity.OutboxEvent{
		EventID:      payload["event_id"].(string),
		EventType:    input.EventType,
		ResourceType: input.ResourceType,
		ResourceID:   input.ResourceID,
		PayloadJSON:  string(payloadJSON),
		Status:       repo.OutboxStatusPending,
		NextRetryAt:  input.OccurredAt,
	}
	if err := store.Outbox.Create(ctx, event); err != nil {
		return internalErr(err)
	}
	return nil
}

func clonePayload(input map[string]any) map[string]any {
	result := make(map[string]any, len(input)+4)
	for key, value := range input {
		result[key] = value
	}
	return result
}

func baseContentPayload(contentID, actorUserID int64, occurredAt time.Time) map[string]any {
	return map[string]any{
		"content_id":    strconv.FormatInt(contentID, 10),
		"actor_user_id": strconv.FormatInt(actorUserID, 10),
		"occurred_at":   occurredAt.UTC().Format(time.RFC3339Nano),
	}
}

func tagIDsChanged(current []entity.Tag, next []int64) bool {
	if len(current) != len(next) {
		return true
	}
	seen := make(map[int64]struct{}, len(current))
	for _, tag := range current {
		seen[tag.ID] = struct{}{}
	}
	for _, id := range next {
		if _, ok := seen[id]; !ok {
			return true
		}
	}
	return false
}

func stringIDs(values []int64) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, strconv.FormatInt(value, 10))
	}
	return result
}

func tagStringIDs(tags []entity.Tag) []string {
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		result = append(result, strconv.FormatInt(tag.ID, 10))
	}
	return result
}
