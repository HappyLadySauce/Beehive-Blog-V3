package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
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
		CommentEnabled: true,
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
