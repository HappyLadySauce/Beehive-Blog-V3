package content

import (
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

func TestResponseMapping(t *testing.T) {
	t.Parallel()

	detail := ToContentDetail(&pb.ContentDetail{
		ContentId: "1",
		Type:      pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:     "Title",
		Slug:      "title",
		Status:    pb.ContentStatus_CONTENT_STATUS_PUBLISHED,
		Tags:      []*pb.ContentTag{{TagId: "10", Name: "Go", Slug: "go"}},
	})
	if detail.Type != "article" || detail.Status != "published" || len(detail.Tags) != 1 || detail.Tags[0].Slug != "go" {
		t.Fatalf("unexpected detail mapping: %+v", detail)
	}

	revisions := ToRevisionListResp([]*pb.ContentRevisionSummary{{
		RevisionId: "2",
		ContentId:  "1",
		RevisionNo: 1,
		EditorType: pb.EditorType_EDITOR_TYPE_HUMAN,
		SourceType: pb.SourceType_SOURCE_TYPE_MANUAL,
	}}, 1, 1, 20)
	if revisions.Total != 1 || revisions.Items[0].EditorType != "human" || revisions.Items[0].SourceType != "manual" {
		t.Fatalf("unexpected revision mapping: %+v", revisions)
	}

	relations := ToRelationListResp([]*pb.ContentRelationView{{
		RelationId:    "30",
		FromContentId: "1",
		ToContentId:   "2",
		RelationType:  pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO,
	}}, 1, 1, 20)
	if relations.Total != 1 || relations.Items[0].RelationType != "related_to" {
		t.Fatalf("unexpected relation mapping: %+v", relations)
	}
}
