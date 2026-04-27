package service_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/mq"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	contentservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

func boolPtr(value bool) *bool {
	return &value
}

func createTestContent(t *testing.T, manager *contentservice.Manager, actor contentservice.Actor, title, slug string) *pb.CreateContentResponse {
	t.Helper()

	created, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:  pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title: title,
		Slug:  slug,
	})
	if err != nil {
		t.Fatalf("create test content %s failed: %v", slug, err)
	}
	return created
}

func mustParseID(t *testing.T, value string) int64 {
	t.Helper()

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		t.Fatalf("parse id failed: %v", err)
	}
	return id
}

func eventTypes(events []entity.OutboxEvent) map[string]int {
	result := map[string]int{}
	for _, event := range events {
		result[event.EventType]++
	}
	return result
}

func createOutboxEvent(t *testing.T, store *repo.Store, resourceID int64, eventType string) entity.OutboxEvent {
	t.Helper()

	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	event := entity.OutboxEvent{
		EventID:      fmt.Sprintf("evt-%d-%s", resourceID, eventType),
		EventType:    eventType,
		ResourceType: "content_item",
		ResourceID:   resourceID,
		PayloadJSON:  `{"ok":true}`,
		Status:       repo.OutboxStatusPending,
		NextRetryAt:  now,
	}
	if err := store.Outbox.Create(context.Background(), &event); err != nil {
		t.Fatalf("create outbox event failed: %v", err)
	}
	return event
}

type fakePublisher struct {
	mu       sync.Mutex
	messages []mq.Message
	err      error
}

func (p *fakePublisher) Publish(_ context.Context, message mq.Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.err != nil {
		return p.err
	}
	p.messages = append(p.messages, message)
	return nil
}

func (p *fakePublisher) Health(context.Context) error {
	return p.err
}

func (p *fakePublisher) Close() error {
	return nil
}
