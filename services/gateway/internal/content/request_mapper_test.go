package content

import (
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

func TestBuildCreateRequestRejectsInvalidEnum(t *testing.T) {
	t.Parallel()

	_, err := BuildCreateRequest(&types.ContentCreateReq{Type: "bad"})
	if !errors.Is(err, errs.E(errs.CodeContentInvalidType)) {
		t.Fatalf("expected invalid type error, got %v", err)
	}
}

func TestBuildRelationRequests(t *testing.T) {
	t.Parallel()

	createReq, err := BuildCreateRelationRequest(&types.ContentRelationCreateReq{
		ContentId:    "1",
		ToContentId:  "2",
		RelationType: "related_to",
		Weight:       10,
		MetadataJson: `{"reason":"same topic"}`,
	})
	if err != nil {
		t.Fatalf("build create relation failed: %v", err)
	}
	if createReq.GetRelationType() != pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO || createReq.GetWeight() != 10 {
		t.Fatalf("unexpected create relation request: %+v", createReq)
	}

	listReq, err := BuildListRelationsRequest(&types.ContentRelationListReq{ContentId: "1", RelationType: "references", Page: 2, PageSize: 10})
	if err != nil {
		t.Fatalf("build list relation failed: %v", err)
	}
	if listReq.GetRelationType() != pb.ContentRelationType_CONTENT_RELATION_TYPE_REFERENCES || listReq.GetPage() != 2 {
		t.Fatalf("unexpected list relation request: %+v", listReq)
	}

	_, err = BuildCreateRelationRequest(&types.ContentRelationCreateReq{RelationType: "bad"})
	if !errors.Is(err, errs.E(errs.CodeContentInvalidArgument)) {
		t.Fatalf("expected invalid relation type, got %v", err)
	}
}

func TestBuildCreateRequestPreservesOptionalCommentEnabled(t *testing.T) {
	t.Parallel()

	unspecified, err := BuildCreateRequest(&types.ContentCreateReq{
		Type:  "article",
		Title: "Title",
		Slug:  "title",
	})
	if err != nil {
		t.Fatalf("build unspecified comment request failed: %v", err)
	}
	if unspecified.CommentEnabled != nil {
		t.Fatalf("expected nil comment_enabled when omitted, got %v", *unspecified.CommentEnabled)
	}

	disabled := false
	explicit, err := BuildCreateRequest(&types.ContentCreateReq{
		Type:           "article",
		Title:          "Title",
		Slug:           "title",
		CommentEnabled: &disabled,
	})
	if err != nil {
		t.Fatalf("build explicit comment request failed: %v", err)
	}
	if explicit.CommentEnabled == nil || *explicit.CommentEnabled {
		t.Fatalf("expected explicit false comment_enabled, got %v", explicit.CommentEnabled)
	}
}
