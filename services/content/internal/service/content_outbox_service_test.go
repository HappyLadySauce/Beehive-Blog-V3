package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	contentservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/testkit"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

func TestContentWritesOutboxEvents(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	manager := contentservice.NewManager(deps)
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}
	tag, err := manager.CreateTag.Execute(context.Background(), actor, &pb.CreateTagRequest{Name: "Event Tag", Slug: "event-tag"})
	if err != nil {
		t.Fatalf("create tag failed: %v", err)
	}
	created, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:   pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:  "Event Source",
		Slug:   "event-source",
		TagIds: []string{tag.Tag.TagId},
	})
	if err != nil {
		t.Fatalf("create content failed: %v", err)
	}
	contentID := mustParseID(t, created.Content.ContentId)

	if _, err := manager.UpdateContent.Execute(context.Background(), actor, &pb.UpdateContentRequest{
		ContentId:      created.Content.ContentId,
		Type:           pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:          "Event Source v2",
		Slug:           "event-source",
		Status:         pb.ContentStatus_CONTENT_STATUS_PUBLISHED,
		Visibility:     pb.ContentVisibility_CONTENT_VISIBILITY_PUBLIC,
		AiAccess:       pb.AIAccess_AI_ACCESS_ALLOWED,
		CommentEnabled: true,
		ChangeSummary:  "publish",
	}); err != nil {
		t.Fatalf("update content failed: %v", err)
	}
	related := createTestContent(t, manager, actor, "Event Target", "event-target")
	relation, err := manager.CreateContentRelation.Execute(context.Background(), actor, &pb.CreateContentRelationRequest{
		ContentId:    created.Content.ContentId,
		ToContentId:  related.Content.ContentId,
		RelationType: pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO,
	})
	if err != nil {
		t.Fatalf("create relation failed: %v", err)
	}
	if _, err := manager.DeleteContentRelation.Execute(context.Background(), actor, &pb.DeleteContentRelationRequest{ContentId: created.Content.ContentId, RelationId: relation.Relation.RelationId}); err != nil {
		t.Fatalf("delete relation failed: %v", err)
	}
	if _, err := manager.ArchiveContent.Execute(context.Background(), actor, &pb.ArchiveContentRequest{ContentId: created.Content.ContentId}); err != nil {
		t.Fatalf("archive content failed: %v", err)
	}

	events, err := deps.Store.Outbox.ListByResource(context.Background(), contentservice.EventResourceContentItem, contentID)
	if err != nil {
		t.Fatalf("list outbox events failed: %v", err)
	}
	got := eventTypes(events)
	for _, want := range []string{
		contentservice.EventContentCreated,
		contentservice.EventContentUpdated,
		contentservice.EventContentStatusChanged,
		contentservice.EventContentVisibilityChanged,
		contentservice.EventContentAIAccessChanged,
		contentservice.EventContentTagChanged,
		contentservice.EventContentRelationChanged,
		contentservice.EventContentArchived,
	} {
		if got[want] == 0 {
			t.Fatalf("expected event %s, got %+v", want, got)
		}
	}
}

func TestBusinessFailureDoesNotWriteOutboxEvent(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	manager := contentservice.NewManager(deps)
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}
	if _, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:  pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title: "Bad event",
		Slug:  "bad-event",
	}); err != nil {
		t.Fatalf("create content failed: %v", err)
	}
	_, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:  pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title: "Bad event duplicate",
		Slug:  "bad-event",
	})
	if !errors.Is(err, errs.E(errs.CodeContentSlugAlreadyExists)) {
		t.Fatalf("expected duplicate slug error, got %v", err)
	}
	var total int64
	if err := deps.Store.DB().WithContext(context.Background()).Model(&entity.OutboxEvent{}).Where("event_type = ?", contentservice.EventContentCreated).Count(&total).Error; err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected only first create event, got %d", total)
	}
}

func TestOutboxDispatcherPublishesAndMarksDone(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	event := createOutboxEvent(t, deps.Store, 101, contentservice.EventContentCreated)
	publisher := &fakePublisher{}
	dispatcher := contentservice.NewOutboxDispatcher(deps.Store, publisher, func() time.Time {
		return time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	}, contentservice.OutboxDispatcherConfig{BatchSize: 10, MaxAttempts: 3, RetryDelay: time.Second})

	if processed, err := dispatcher.DispatchOnce(context.Background()); err != nil || processed != 1 {
		t.Fatalf("dispatch once failed: processed=%d err=%v", processed, err)
	}
	events, err := deps.Store.Outbox.ListByResource(context.Background(), "content_item", 101)
	if err != nil {
		t.Fatalf("list events failed: %v", err)
	}
	if events[0].ID != event.ID || events[0].Status != repo.OutboxStatusDone || events[0].PublishedAt == nil {
		t.Fatalf("expected event done, got %+v", events[0])
	}
	if len(publisher.messages) != 1 || publisher.messages[0].RoutingKey != contentservice.EventContentCreated {
		t.Fatalf("unexpected published messages: %+v", publisher.messages)
	}
}

func TestOutboxDispatcherRetriesPublishFailure(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	createOutboxEvent(t, deps.Store, 102, contentservice.EventContentUpdated)
	publisher := &fakePublisher{err: errors.New("broker down")}
	dispatcher := contentservice.NewOutboxDispatcher(deps.Store, publisher, func() time.Time {
		return time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	}, contentservice.OutboxDispatcherConfig{BatchSize: 10, MaxAttempts: 3, RetryDelay: time.Second})

	if processed, err := dispatcher.DispatchOnce(context.Background()); err != nil || processed != 1 {
		t.Fatalf("dispatch once failed: processed=%d err=%v", processed, err)
	}
	events, err := deps.Store.Outbox.ListByResource(context.Background(), "content_item", 102)
	if err != nil {
		t.Fatalf("list events failed: %v", err)
	}
	if events[0].Status != repo.OutboxStatusPending || events[0].Attempts != 1 || events[0].LastError == "" {
		t.Fatalf("expected retryable failed publish, got %+v", events[0])
	}
}

func TestOutboxDispatcherMarksTransientPublishFailureFailedAtMaxAttempts(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	createOutboxEvent(t, deps.Store, 104, contentservice.EventContentUpdated)
	publisher := &fakePublisher{err: errors.New("broker down")}
	dispatcher := contentservice.NewOutboxDispatcher(deps.Store, publisher, func() time.Time {
		return time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	}, contentservice.OutboxDispatcherConfig{BatchSize: 10, MaxAttempts: 1, RetryDelay: time.Second})

	if processed, err := dispatcher.DispatchOnce(context.Background()); err != nil || processed != 1 {
		t.Fatalf("dispatch once failed: processed=%d err=%v", processed, err)
	}
	events, err := deps.Store.Outbox.ListByResource(context.Background(), "content_item", 104)
	if err != nil {
		t.Fatalf("list events failed: %v", err)
	}
	if events[0].Status != repo.OutboxStatusFailed || events[0].Attempts != 1 {
		t.Fatalf("expected transient publish failure to become failed at max attempts, got %+v", events[0])
	}
}

func TestOutboxDispatcherKeepsTransientPublishFailurePendingBeforeDefaultMaxAttempts(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	event := createOutboxEvent(t, deps.Store, 106, contentservice.EventContentUpdated)
	if err := deps.Store.DB().WithContext(context.Background()).
		Model(&entity.OutboxEvent{}).
		Where("id = ?", event.ID).
		Update("attempts", 358).Error; err != nil {
		t.Fatalf("seed outbox attempts failed: %v", err)
	}
	publisher := &fakePublisher{err: errors.New("broker down")}
	dispatcher := contentservice.NewOutboxDispatcher(deps.Store, publisher, func() time.Time {
		return time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	}, contentservice.OutboxDispatcherConfig{BatchSize: 10, RetryDelay: time.Second})

	if processed, err := dispatcher.DispatchOnce(context.Background()); err != nil || processed != 1 {
		t.Fatalf("dispatch once failed: processed=%d err=%v", processed, err)
	}
	events, err := deps.Store.Outbox.ListByResource(context.Background(), "content_item", 106)
	if err != nil {
		t.Fatalf("list events failed: %v", err)
	}
	if events[0].Status != repo.OutboxStatusPending || events[0].Attempts != 359 {
		t.Fatalf("expected transient publish failure to stay pending before default max attempts, got %+v", events[0])
	}
}

func TestOutboxDispatcherMarksPermanentPublishFailureFailedAtMaxAttempts(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	createOutboxEvent(t, deps.Store, 105, "")
	publisher := &fakePublisher{}
	dispatcher := contentservice.NewOutboxDispatcher(deps.Store, publisher, func() time.Time {
		return time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	}, contentservice.OutboxDispatcherConfig{BatchSize: 10, MaxAttempts: 360, RetryDelay: time.Second})

	if processed, err := dispatcher.DispatchOnce(context.Background()); err != nil || processed != 1 {
		t.Fatalf("dispatch once failed: processed=%d err=%v", processed, err)
	}
	events, err := deps.Store.Outbox.ListByResource(context.Background(), "content_item", 105)
	if err != nil {
		t.Fatalf("list events failed: %v", err)
	}
	if events[0].Status != repo.OutboxStatusFailed || events[0].Attempts != 1 {
		t.Fatalf("expected permanent publish failure to become failed, got %+v", events[0])
	}
	if len(publisher.messages) != 0 {
		t.Fatalf("expected no published messages for invalid event, got %+v", publisher.messages)
	}
}

func TestOutboxDispatcherReclaimsStaleProcessingEvent(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	event := createOutboxEvent(t, deps.Store, 103, contentservice.EventContentUpdated)
	claimedAt := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	claimed, err := deps.Store.Outbox.ClaimDue(context.Background(), claimedAt, 1, 30*time.Second)
	if err != nil {
		t.Fatalf("claim outbox event failed: %v", err)
	}
	if len(claimed) != 1 || claimed[0].ID != event.ID {
		t.Fatalf("expected one claimed event, got %+v", claimed)
	}

	publisher := &fakePublisher{}
	dispatcher := contentservice.NewOutboxDispatcher(deps.Store, publisher, func() time.Time {
		return claimedAt.Add(31 * time.Second)
	}, contentservice.OutboxDispatcherConfig{BatchSize: 10, MaxAttempts: 3, RetryDelay: time.Second, ProcessingTimeout: 30 * time.Second})

	if processed, err := dispatcher.DispatchOnce(context.Background()); err != nil || processed != 1 {
		t.Fatalf("dispatch stale processing event failed: processed=%d err=%v", processed, err)
	}
	events, err := deps.Store.Outbox.ListByResource(context.Background(), "content_item", 103)
	if err != nil {
		t.Fatalf("list events failed: %v", err)
	}
	if events[0].Status != repo.OutboxStatusDone || events[0].PublishedAt == nil {
		t.Fatalf("expected reclaimed event to be done, got %+v", events[0])
	}
	if len(publisher.messages) != 1 || publisher.messages[0].ID != event.EventID {
		t.Fatalf("unexpected published messages: %+v", publisher.messages)
	}
}

func TestOutboxRepositoryRejectsStaleLeaseMark(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	event := createOutboxEvent(t, deps.Store, 107, contentservice.EventContentUpdated)
	firstClaimAt := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	firstClaim, err := deps.Store.Outbox.ClaimDue(context.Background(), firstClaimAt, 1, 30*time.Second)
	if err != nil {
		t.Fatalf("first claim failed: %v", err)
	}
	if len(firstClaim) != 1 || firstClaim[0].ID != event.ID {
		t.Fatalf("expected first claim to lock event, got %+v", firstClaim)
	}

	secondClaimAt := firstClaimAt.Add(31 * time.Second)
	secondClaim, err := deps.Store.Outbox.ClaimDue(context.Background(), secondClaimAt, 1, 30*time.Second)
	if err != nil {
		t.Fatalf("second claim failed: %v", err)
	}
	if len(secondClaim) != 1 || secondClaim[0].ID != event.ID {
		t.Fatalf("expected second claim to reclaim event, got %+v", secondClaim)
	}
	if err := deps.Store.Outbox.MarkDone(context.Background(), secondClaim[0].ID, secondClaim[0].UpdatedAt, secondClaimAt.Add(time.Second)); err != nil {
		t.Fatalf("second worker mark done failed: %v", err)
	}

	err = deps.Store.Outbox.MarkPublishFailed(
		context.Background(),
		firstClaim[0].ID,
		firstClaim[0].UpdatedAt,
		1,
		3,
		secondClaimAt.Add(time.Minute),
		"old worker failure",
		secondClaimAt.Add(2*time.Second),
	)
	if err == nil {
		t.Fatalf("expected stale lease mark failed to be rejected")
	}

	events, err := deps.Store.Outbox.ListByResource(context.Background(), "content_item", 107)
	if err != nil {
		t.Fatalf("list events failed: %v", err)
	}
	if events[0].Status != repo.OutboxStatusDone || events[0].LastError == "old worker failure" {
		t.Fatalf("expected done event to survive stale worker result, got %+v", events[0])
	}
}

func TestOutboxDispatcherConcurrentClaimDoesNotDuplicate(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	for i := 0; i < 5; i++ {
		createOutboxEvent(t, deps.Store, int64(200+i), contentservice.EventContentUpdated)
	}
	publisher := &fakePublisher{}
	dispatcher := contentservice.NewOutboxDispatcher(deps.Store, publisher, func() time.Time {
		return time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	}, contentservice.OutboxDispatcherConfig{BatchSize: 3, MaxAttempts: 3, RetryDelay: time.Second})

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := dispatcher.DispatchOnce(context.Background()); err != nil {
				t.Errorf("dispatch once failed: %v", err)
			}
		}()
	}
	wg.Wait()
	if len(publisher.messages) != 5 {
		t.Fatalf("expected 5 unique published messages, got %d", len(publisher.messages))
	}
	seen := map[string]struct{}{}
	for _, message := range publisher.messages {
		if _, ok := seen[message.ID]; ok {
			t.Fatalf("duplicate message published: %s", message.ID)
		}
		seen[message.ID] = struct{}{}
	}
}
