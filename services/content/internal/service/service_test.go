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
	if created.Content.BodyJson != `{"type":"doc"}` {
		t.Fatalf("expected trimmed body_json, got %q", created.Content.BodyJson)
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
