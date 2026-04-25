package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/mq"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	contentservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/testkit"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

func TestCreateUpdateArchiveAndPublicRead(t *testing.T) {
	t.Parallel()

	deps := testkit.NewServiceDependencies(t)
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}
	tag, err := contentservice.NewManager(deps).CreateTag.Execute(context.Background(), actor, &pb.CreateTagRequest{
		Name: "Go",
		Slug: "go",
	})
	if err != nil {
		t.Fatalf("create tag failed: %v", err)
	}

	manager := contentservice.NewManager(deps)
	created, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:           pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:          "First post",
		Slug:           "first-post",
		Summary:        "summary",
		BodyMarkdown:   "# First",
		Visibility:     pb.ContentVisibility_CONTENT_VISIBILITY_PUBLIC,
		AiAccess:       pb.AIAccess_AI_ACCESS_ALLOWED,
		SourceType:     pb.SourceType_SOURCE_TYPE_MANUAL,
		CommentEnabled: boolPtr(true),
		TagIds:         []string{tag.Tag.TagId},
	})
	if err != nil {
		t.Fatalf("create content failed: %v", err)
	}
	if created.Content.CurrentRevisionId == "" {
		t.Fatalf("expected current revision id")
	}
	if len(created.Content.Tags) != 1 || created.Content.Tags[0].Slug != "go" {
		t.Fatalf("expected bound tag, got %+v", created.Content.Tags)
	}

	updated, err := manager.UpdateContent.Execute(context.Background(), actor, &pb.UpdateContentRequest{
		ContentId:      created.Content.ContentId,
		Type:           pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:          "First post v2",
		Slug:           "first-post",
		Summary:        "summary v2",
		BodyMarkdown:   "# First v2",
		Status:         pb.ContentStatus_CONTENT_STATUS_PUBLISHED,
		Visibility:     pb.ContentVisibility_CONTENT_VISIBILITY_PUBLIC,
		AiAccess:       pb.AIAccess_AI_ACCESS_ALLOWED,
		CommentEnabled: true,
		TagIds:         []string{tag.Tag.TagId},
		ChangeSummary:  "publish",
	})
	if err != nil {
		t.Fatalf("update content failed: %v", err)
	}
	if updated.Content.PublishedAt == 0 {
		t.Fatalf("expected published_at to be set")
	}

	publicDetail, err := manager.GetPublicContentBySlug.Execute(context.Background(), &pb.GetPublicContentBySlugRequest{Slug: "first-post"})
	if err != nil {
		t.Fatalf("public get failed: %v", err)
	}
	if publicDetail.Content.Title != "First post v2" {
		t.Fatalf("unexpected public title: %s", publicDetail.Content.Title)
	}

	revisions, err := manager.ListContentRevisions.Execute(context.Background(), actor, &pb.ListContentRevisionsRequest{ContentId: created.Content.ContentId})
	if err != nil {
		t.Fatalf("list revisions failed: %v", err)
	}
	if revisions.Total != 2 {
		t.Fatalf("expected two revisions, got %d", revisions.Total)
	}

	if _, err := manager.ArchiveContent.Execute(context.Background(), actor, &pb.ArchiveContentRequest{ContentId: created.Content.ContentId}); err != nil {
		t.Fatalf("archive failed: %v", err)
	}
	if _, err := manager.GetPublicContentBySlug.Execute(context.Background(), &pb.GetPublicContentBySlugRequest{Slug: "first-post"}); !errors.Is(err, errs.E(errs.CodeContentNotFound)) {
		t.Fatalf("expected archived public content to be hidden, got %v", err)
	}
}

func TestCreateContentCanDisableComments(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}
	created, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:           pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:          "Comments off",
		Slug:           "comments-off",
		CommentEnabled: boolPtr(false),
	})
	if err != nil {
		t.Fatalf("create content with comments disabled failed: %v", err)
	}
	if created.Content.CommentEnabled {
		t.Fatalf("expected comments to be disabled")
	}
}

func TestStudioWriteRequiresActor(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	_, err := manager.CreateContent.Execute(context.Background(), contentservice.Actor{}, &pb.CreateContentRequest{
		Type:  pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title: "No actor",
		Slug:  "no-actor",
	})
	if !errors.Is(err, errs.E(errs.CodeContentAccessForbidden)) {
		t.Fatalf("expected access forbidden, got %v", err)
	}
}

func boolPtr(value bool) *bool {
	return &value
}

func TestStudioServicesRequireAdmin(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	admin := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}
	member := contentservice.Actor{UserID: 2, SessionID: 20, Role: "member"}

	tag, err := manager.CreateTag.Execute(context.Background(), admin, &pb.CreateTagRequest{Name: "Admin", Slug: "admin"})
	if err != nil {
		t.Fatalf("create tag failed: %v", err)
	}
	created, err := manager.CreateContent.Execute(context.Background(), admin, &pb.CreateContentRequest{
		Type:  pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title: "Admin only",
		Slug:  "admin-only",
	})
	if err != nil {
		t.Fatalf("create content failed: %v", err)
	}

	cases := []struct {
		name string
		run  func() error
	}{
		{
			name: "create content",
			run: func() error {
				_, err := manager.CreateContent.Execute(context.Background(), member, &pb.CreateContentRequest{Type: pb.ContentType_CONTENT_TYPE_ARTICLE, Title: "Nope", Slug: "nope"})
				return err
			},
		},
		{
			name: "update content",
			run: func() error {
				_, err := manager.UpdateContent.Execute(context.Background(), member, &pb.UpdateContentRequest{ContentId: created.Content.ContentId, Type: pb.ContentType_CONTENT_TYPE_ARTICLE, Title: "Nope", Slug: "admin-only"})
				return err
			},
		},
		{
			name: "get content",
			run: func() error {
				_, err := manager.GetContent.Execute(context.Background(), member, &pb.GetContentRequest{ContentId: created.Content.ContentId})
				return err
			},
		},
		{
			name: "list content",
			run: func() error {
				_, err := manager.ListStudioContents.Execute(context.Background(), member, &pb.ListStudioContentsRequest{})
				return err
			},
		},
		{
			name: "archive content",
			run: func() error {
				_, err := manager.ArchiveContent.Execute(context.Background(), member, &pb.ArchiveContentRequest{ContentId: created.Content.ContentId})
				return err
			},
		},
		{
			name: "list revisions",
			run: func() error {
				_, err := manager.ListContentRevisions.Execute(context.Background(), member, &pb.ListContentRevisionsRequest{ContentId: created.Content.ContentId})
				return err
			},
		},
		{
			name: "get revision",
			run: func() error {
				_, err := manager.GetContentRevision.Execute(context.Background(), member, &pb.GetContentRevisionRequest{ContentId: created.Content.ContentId, RevisionId: created.Content.CurrentRevisionId})
				return err
			},
		},
		{
			name: "create relation",
			run: func() error {
				_, err := manager.CreateContentRelation.Execute(context.Background(), member, &pb.CreateContentRelationRequest{ContentId: created.Content.ContentId, ToContentId: created.Content.ContentId, RelationType: pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO})
				return err
			},
		},
		{
			name: "list relations",
			run: func() error {
				_, err := manager.ListContentRelations.Execute(context.Background(), member, &pb.ListContentRelationsRequest{ContentId: created.Content.ContentId})
				return err
			},
		},
		{
			name: "delete relation",
			run: func() error {
				_, err := manager.DeleteContentRelation.Execute(context.Background(), member, &pb.DeleteContentRelationRequest{ContentId: created.Content.ContentId, RelationId: "1"})
				return err
			},
		},
		{
			name: "create tag",
			run: func() error {
				_, err := manager.CreateTag.Execute(context.Background(), member, &pb.CreateTagRequest{Name: "Nope", Slug: "nope"})
				return err
			},
		},
		{
			name: "update tag",
			run: func() error {
				_, err := manager.UpdateTag.Execute(context.Background(), member, &pb.UpdateTagRequest{TagId: tag.Tag.TagId, Name: "Nope", Slug: "nope"})
				return err
			},
		},
		{
			name: "delete tag",
			run: func() error {
				_, err := manager.DeleteTag.Execute(context.Background(), member, &pb.DeleteTagRequest{TagId: tag.Tag.TagId})
				return err
			},
		},
		{
			name: "list tags",
			run: func() error {
				_, err := manager.ListTags.Execute(context.Background(), member, &pb.ListTagsRequest{})
				return err
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.run(); !errors.Is(err, errs.E(errs.CodeContentAccessForbidden)) {
				t.Fatalf("expected access forbidden, got %v", err)
			}
		})
	}
}

func TestContentRelationsLifecycleAndValidation(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}
	first := createTestContent(t, manager, actor, "Relation From", "relation-from")
	second := createTestContent(t, manager, actor, "Relation To", "relation-to")

	created, err := manager.CreateContentRelation.Execute(context.Background(), actor, &pb.CreateContentRelationRequest{
		ContentId:    first.Content.ContentId,
		ToContentId:  second.Content.ContentId,
		RelationType: pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO,
		Weight:       10,
		SortOrder:    2,
		MetadataJson: ` {"reason":"same topic"} `,
	})
	if err != nil {
		t.Fatalf("create relation failed: %v", err)
	}
	if created.Relation.RelationType != pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO || created.Relation.Weight != 10 {
		t.Fatalf("unexpected relation: %+v", created.Relation)
	}
	var metadata map[string]string
	if err := json.Unmarshal([]byte(created.Relation.MetadataJson), &metadata); err != nil {
		t.Fatalf("expected valid metadata_json, got %q: %v", created.Relation.MetadataJson, err)
	}
	if metadata["reason"] != "same topic" {
		t.Fatalf("unexpected metadata_json: %q", created.Relation.MetadataJson)
	}

	listed, err := manager.ListContentRelations.Execute(context.Background(), actor, &pb.ListContentRelationsRequest{
		ContentId:    first.Content.ContentId,
		RelationType: pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO,
	})
	if err != nil {
		t.Fatalf("list relations failed: %v", err)
	}
	if listed.Total != 1 || len(listed.Items) != 1 || listed.Items[0].RelationId != created.Relation.RelationId {
		t.Fatalf("unexpected relation list: %+v", listed)
	}

	_, err = manager.CreateContentRelation.Execute(context.Background(), actor, &pb.CreateContentRelationRequest{
		ContentId:    first.Content.ContentId,
		ToContentId:  second.Content.ContentId,
		RelationType: pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO,
	})
	if !errors.Is(err, errs.E(errs.CodeContentRelationAlreadyExists)) {
		t.Fatalf("expected duplicate relation error, got %v", err)
	}

	_, err = manager.CreateContentRelation.Execute(context.Background(), actor, &pb.CreateContentRelationRequest{
		ContentId:    first.Content.ContentId,
		ToContentId:  first.Content.ContentId,
		RelationType: pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO,
	})
	if !errors.Is(err, errs.E(errs.CodeContentInvalidArgument)) {
		t.Fatalf("expected self relation invalid argument, got %v", err)
	}

	_, err = manager.CreateContentRelation.Execute(context.Background(), actor, &pb.CreateContentRelationRequest{
		ContentId:    first.Content.ContentId,
		ToContentId:  second.Content.ContentId,
		RelationType: pb.ContentRelationType_CONTENT_RELATION_TYPE_REFERENCES,
		MetadataJson: `{bad`,
	})
	if !errors.Is(err, errs.E(errs.CodeContentInvalidArgument)) {
		t.Fatalf("expected invalid metadata_json argument, got %v", err)
	}

	_, err = manager.CreateContentRelation.Execute(context.Background(), actor, &pb.CreateContentRelationRequest{
		ContentId:    first.Content.ContentId,
		ToContentId:  "999999",
		RelationType: pb.ContentRelationType_CONTENT_RELATION_TYPE_REFERENCES,
	})
	if !errors.Is(err, errs.E(errs.CodeContentNotFound)) {
		t.Fatalf("expected missing related content error, got %v", err)
	}

	_, err = manager.DeleteContentRelation.Execute(context.Background(), actor, &pb.DeleteContentRelationRequest{
		ContentId:  second.Content.ContentId,
		RelationId: created.Relation.RelationId,
	})
	if !errors.Is(err, errs.E(errs.CodeContentRelationNotFound)) {
		t.Fatalf("expected wrong owner relation not found, got %v", err)
	}

	deleted, err := manager.DeleteContentRelation.Execute(context.Background(), actor, &pb.DeleteContentRelationRequest{
		ContentId:  first.Content.ContentId,
		RelationId: created.Relation.RelationId,
	})
	if err != nil {
		t.Fatalf("delete relation failed: %v", err)
	}
	if !deleted.Ok {
		t.Fatalf("expected delete relation ok")
	}
	listed, err = manager.ListContentRelations.Execute(context.Background(), actor, &pb.ListContentRelationsRequest{ContentId: first.Content.ContentId})
	if err != nil {
		t.Fatalf("list relations after delete failed: %v", err)
	}
	if listed.Total != 0 {
		t.Fatalf("expected no relations after delete, got %d", listed.Total)
	}
}

func TestContentRelationsRejectGuest(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	_, err := manager.ListContentRelations.Execute(context.Background(), contentservice.Actor{}, &pb.ListContentRelationsRequest{ContentId: "1"})
	if !errors.Is(err, errs.E(errs.CodeContentAccessForbidden)) {
		t.Fatalf("expected access forbidden, got %v", err)
	}
}

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
	}, contentservice.OutboxDispatcherConfig{BatchSize: 10, MaxAttempts: 1, RetryDelay: time.Second})

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

func TestCreateContentDefaultsToPrivateAndAIDenied(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "ROLE_ADMIN"}
	created, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:  pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title: "Default private",
		Slug:  "default-private",
	})
	if err != nil {
		t.Fatalf("create content failed: %v", err)
	}
	if created.Content.Visibility != pb.ContentVisibility_CONTENT_VISIBILITY_PRIVATE {
		t.Fatalf("expected private visibility, got %s", created.Content.Visibility)
	}
	if created.Content.AiAccess != pb.AIAccess_AI_ACCESS_DENIED {
		t.Fatalf("expected denied ai access, got %s", created.Content.AiAccess)
	}
	if !created.Content.CommentEnabled {
		t.Fatalf("expected comments to be enabled by default")
	}
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

func TestBodyJSONValidation(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}

	_, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:     pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:    "Bad JSON",
		Slug:     "bad-json",
		BodyJson: `{"type":`,
	})
	if !errors.Is(err, errs.E(errs.CodeContentInvalidArgument)) {
		t.Fatalf("expected invalid argument for create body_json, got %v", err)
	}

	created, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:     pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:    "Good JSON",
		Slug:     "good-json",
		BodyJson: ` {"type":"doc"} `,
	})
	if err != nil {
		t.Fatalf("create valid json content failed: %v", err)
	}
	var body map[string]string
	if err := json.Unmarshal([]byte(created.Content.BodyJson), &body); err != nil {
		t.Fatalf("expected valid body_json response, got %q: %v", created.Content.BodyJson, err)
	}
	if body["type"] != "doc" {
		t.Fatalf("expected body_json type doc, got %q from %q", body["type"], created.Content.BodyJson)
	}

	_, err = manager.UpdateContent.Execute(context.Background(), actor, &pb.UpdateContentRequest{
		ContentId: created.Content.ContentId,
		Type:      pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:     "Good JSON",
		Slug:      "good-json",
		BodyJson:  `{bad`,
	})
	if !errors.Is(err, errs.E(errs.CodeContentInvalidArgument)) {
		t.Fatalf("expected invalid argument for update body_json, got %v", err)
	}
}

func TestBlankBodyJSONIsStoredAsEmptyResponse(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}
	created, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:     pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:    "Blank JSON",
		Slug:     "blank-json",
		BodyJson: "   ",
	})
	if err != nil {
		t.Fatalf("create blank body_json content failed: %v", err)
	}
	if created.Content.BodyJson != "" {
		t.Fatalf("expected empty body_json response, got %q", created.Content.BodyJson)
	}
}

func TestTagDeleteFailsWhenBound(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}
	tag, err := manager.CreateTag.Execute(context.Background(), actor, &pb.CreateTagRequest{Name: "Bound", Slug: "bound"})
	if err != nil {
		t.Fatalf("create tag failed: %v", err)
	}
	if _, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:   pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:  "Tagged",
		Slug:   "tagged",
		TagIds: []string{tag.Tag.TagId},
	}); err != nil {
		t.Fatalf("create content failed: %v", err)
	}
	_, err = manager.DeleteTag.Execute(context.Background(), actor, &pb.DeleteTagRequest{TagId: tag.Tag.TagId})
	if !errors.Is(err, errs.E(errs.CodeContentTagInUse)) {
		t.Fatalf("expected tag in use, got %v", err)
	}
}
