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
