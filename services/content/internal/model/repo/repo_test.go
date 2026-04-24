package repo_test

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/testkit"
)

func TestContentRepositories(t *testing.T) {
	t.Parallel()

	store := testkit.NewServiceDependencies(t).Store
	ctx := context.Background()

	item := &entity.Item{
		Type:         "article",
		Title:        "Repo Post",
		Slug:         "repo-post",
		Status:       "published",
		Visibility:   "public",
		AIAccess:     "allowed",
		OwnerUserID:  1,
		AuthorUserID: 1,
		SourceType:   "manual",
	}
	if err := store.Items.Create(ctx, item); err != nil {
		t.Fatalf("create item failed: %v", err)
	}
	duplicate := *item
	duplicate.ID = 0
	if err := store.Items.Create(ctx, &duplicate); repo.UniqueConstraint(err) != repo.ConstraintItemSlug {
		t.Fatalf("expected slug unique constraint, got %v", err)
	}

	revision := &entity.Revision{
		ContentID:     item.ID,
		RevisionNo:    1,
		TitleSnapshot: item.Title,
		BodyMarkdown:  "body",
		EditorType:    "human",
		SourceType:    "manual",
	}
	if err := store.Revisions.Create(ctx, revision); err != nil {
		t.Fatalf("create revision failed: %v", err)
	}
	item.CurrentRevisionID = &revision.ID
	if err := store.Items.Save(ctx, item); err != nil {
		t.Fatalf("save current revision failed: %v", err)
	}
	current, err := store.Revisions.GetCurrent(ctx, item)
	if err != nil {
		t.Fatalf("get current revision failed: %v", err)
	}
	if current.ID != revision.ID {
		t.Fatalf("unexpected current revision id: %d", current.ID)
	}

	tag := &entity.Tag{Name: "Repo", Slug: "repo"}
	if err := store.Tags.Create(ctx, tag); err != nil {
		t.Fatalf("create tag failed: %v", err)
	}
	if err := store.ContentTags.ReplaceForContent(ctx, item.ID, []int64{tag.ID}); err != nil {
		t.Fatalf("replace tags failed: %v", err)
	}
	total, err := store.ContentTags.CountByTag(ctx, tag.ID)
	if err != nil {
		t.Fatalf("count tag failed: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected one tag binding, got %d", total)
	}
	if err := store.Tags.Delete(ctx, tag.ID); err == nil {
		t.Fatalf("expected bound tag delete to fail")
	}

	relatedItem := &entity.Item{
		Type:         "article",
		Title:        "Related Post",
		Slug:         "related-post",
		Status:       "published",
		Visibility:   "public",
		AIAccess:     "allowed",
		OwnerUserID:  1,
		AuthorUserID: 1,
		SourceType:   "manual",
	}
	if err := store.Items.Create(ctx, relatedItem); err != nil {
		t.Fatalf("create related item failed: %v", err)
	}
	relation := &entity.Relation{
		FromContentID: item.ID,
		ToContentID:   relatedItem.ID,
		RelationType:  "related_to",
		Weight:        10,
	}
	if err := store.Relations.Create(ctx, relation); err != nil {
		t.Fatalf("create relation failed: %v", err)
	}
	duplicateRelation := *relation
	duplicateRelation.ID = 0
	if err := store.Relations.Create(ctx, &duplicateRelation); repo.UniqueConstraint(err) != repo.ConstraintContentRelation {
		t.Fatalf("expected relation unique constraint, got %v", err)
	}
	if err := store.Relations.Create(ctx, &entity.Relation{FromContentID: item.ID, ToContentID: 999, RelationType: "related_to"}); err == nil {
		t.Fatalf("expected dangling relation to fail")
	}
	relations, total, err := store.Relations.List(ctx, repo.RelationFilter{FromContentID: item.ID, RelationType: "related_to"}, 1, 20)
	if err != nil {
		t.Fatalf("list relations failed: %v", err)
	}
	if total != 1 || len(relations) != 1 || relations[0].ID != relation.ID {
		t.Fatalf("unexpected relation list: total=%d relations=%+v", total, relations)
	}
	staleRelation := *relation
	if err := store.Relations.Delete(ctx, relation); err != nil {
		t.Fatalf("delete relation failed: %v", err)
	}
	if err := store.Relations.Delete(ctx, &staleRelation); !repo.IsNotFound(err) {
		t.Fatalf("expected stale relation delete to report not found, got %v", err)
	}
	relation.ID = 0
	if err := store.Relations.Create(ctx, relation); err != nil {
		t.Fatalf("recreate relation failed: %v", err)
	}
	if err := store.DB().WithContext(ctx).Delete(&entity.Item{}, "id = ?", relatedItem.ID).Error; err != nil {
		t.Fatalf("delete related item failed: %v", err)
	}
	relationTotal, err := store.Relations.CountByContent(ctx, item.ID)
	if err != nil {
		t.Fatalf("count relations after cascade failed: %v", err)
	}
	if relationTotal != 0 {
		t.Fatalf("expected cascaded relation delete, got %d", relationTotal)
	}

	privateItem := *item
	privateItem.ID = 0
	privateItem.Slug = "private-post"
	privateItem.Visibility = "private"
	privateItem.CurrentRevisionID = nil
	if err := store.Items.Create(ctx, &privateItem); err != nil {
		t.Fatalf("create private item failed: %v", err)
	}
	items, total, err := store.Items.List(ctx, repo.ItemFilter{PublicOnly: true}, 1, 20)
	if err != nil {
		t.Fatalf("list public items failed: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Slug != "repo-post" {
		t.Fatalf("expected only published public item, total=%d items=%+v", total, items)
	}
}

func TestRepositoryNotFound(t *testing.T) {
	t.Parallel()

	store := testkit.NewServiceDependencies(t).Store
	_, err := store.Items.GetByID(context.Background(), 999)
	if !repo.IsNotFound(err) {
		t.Fatalf("expected record not found, got %v", err)
	}
}
