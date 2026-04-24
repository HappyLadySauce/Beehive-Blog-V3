package service_test

import (
	"context"
	"encoding/json"
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
