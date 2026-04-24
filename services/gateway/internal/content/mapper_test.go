package content

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestEnumMapping(t *testing.T) {
	t.Parallel()

	contentType, err := ContentTypeToProto("timeline_event")
	if err != nil || contentType != pb.ContentType_CONTENT_TYPE_TIMELINE_EVENT {
		t.Fatalf("unexpected content type: %s %v", contentType, err)
	}
	status, err := StatusToProto("published")
	if err != nil || status != pb.ContentStatus_CONTENT_STATUS_PUBLISHED {
		t.Fatalf("unexpected status: %s %v", status, err)
	}
	visibility, err := VisibilityToProto("member")
	if err != nil || visibility != pb.ContentVisibility_CONTENT_VISIBILITY_MEMBER {
		t.Fatalf("unexpected visibility: %s %v", visibility, err)
	}
	aiAccess, err := AIAccessToProto("allowed")
	if err != nil || aiAccess != pb.AIAccess_AI_ACCESS_ALLOWED {
		t.Fatalf("unexpected ai access: %s %v", aiAccess, err)
	}
	sourceType, err := SourceTypeToProtoDefault("agent_assisted")
	if err != nil || sourceType != pb.SourceType_SOURCE_TYPE_AGENT_ASSISTED {
		t.Fatalf("unexpected source type: %s %v", sourceType, err)
	}
	relationType, err := RelationTypeToProto("related_to")
	if err != nil || relationType != pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO {
		t.Fatalf("unexpected relation type: %s %v", relationType, err)
	}
	if EditorTypeToString(pb.EditorType_EDITOR_TYPE_SYSTEM) != "system" {
		t.Fatalf("unexpected editor type mapping")
	}
}

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

func TestRPCContexts(t *testing.T) {
	t.Parallel()

	conf := config.ContentRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"}
	ctx := middleware.WithRequestMeta(context.Background(), ctxmeta.RequestMeta{RequestID: "req-1"})
	ctx = middleware.WithAuthContext(ctx, middleware.AuthContext{UserID: "1", SessionID: "s1", Role: "admin"})
	rpcCtx, err := StudioRPCContext(ctx, conf)
	if err != nil {
		t.Fatalf("unexpected studio context error: %v", err)
	}
	md, _ := metadata.FromOutgoingContext(rpcCtx)
	if got := md.Get(ctxmeta.MetadataKeyUserID); len(got) != 1 || got[0] != "1" {
		t.Fatalf("expected user id metadata, got %v", got)
	}
	if got := md.Get(ctxmeta.MetadataKeyInternalCaller); len(got) != 1 || got[0] != "gateway" {
		t.Fatalf("expected caller metadata, got %v", got)
	}

	publicCtx := PublicRPCContext(ctx, conf)
	publicMD, _ := metadata.FromOutgoingContext(publicCtx)
	if got := publicMD.Get(ctxmeta.MetadataKeyUserID); len(got) != 0 {
		t.Fatalf("did not expect public user metadata, got %v", got)
	}
	if got := publicMD.Get(ctxmeta.MetadataKeyInternalAuthToken); len(got) != 1 || got[0] != "secret" {
		t.Fatalf("expected public internal auth metadata, got %v", got)
	}
}

func TestMapUpstreamErrorPreservesContentDomainCode(t *testing.T) {
	t.Parallel()

	err := errgrpcx.ToStatus(errs.New(errs.CodeContentSlugAlreadyExists, "slug already exists"), "fallback")
	mapped := MapUpstreamError(context.Background(), "content_create", "/route", err)
	if !errors.Is(mapped, errs.E(errs.CodeContentSlugAlreadyExists)) {
		t.Fatalf("expected slug conflict to be preserved, got %v", mapped)
	}
}

func TestMapUpstreamErrorMasksContentInternalCallerUnauthorized(t *testing.T) {
	t.Parallel()

	err := errgrpcx.ToStatus(errs.New(errs.CodeContentInternalCallerUnauthorized, "internal caller authentication failed"), "fallback")
	mapped := MapUpstreamError(context.Background(), "content_create", "/route", err)
	if !errors.Is(mapped, errs.E(errs.CodeGatewayNotReady)) {
		t.Fatalf("expected gateway not ready error, got %v", mapped)
	}
	if errors.Is(mapped, errs.E(errs.CodeContentInternalCallerUnauthorized)) {
		t.Fatalf("expected internal auth code to be masked, got %v", mapped)
	}
}

func TestMapUpstreamErrorMapsContentUnavailable(t *testing.T) {
	t.Parallel()

	mapped := MapUpstreamError(context.Background(), "content_list", "/route", status.Error(codes.Unavailable, "unavailable"))
	if !errors.Is(mapped, errs.E(errs.CodeGatewayContentServiceUnavailable)) {
		t.Fatalf("expected content service unavailable, got %v", mapped)
	}
	if errors.Is(mapped, errs.E(errs.CodeGatewayAuthServiceUnavailable)) {
		t.Fatalf("expected content unavailable not auth unavailable, got %v", mapped)
	}
}

func TestMapUpstreamErrorPreservesTagInUse(t *testing.T) {
	t.Parallel()

	err := errgrpcx.ToStatus(errs.New(errs.CodeContentTagInUse, "tag is in use"), "fallback")
	mapped := MapUpstreamError(context.Background(), "content_tag_delete", "/route", err)
	if !errors.Is(mapped, errs.E(errs.CodeContentTagInUse)) {
		t.Fatalf("expected tag in use error to be preserved, got %v", mapped)
	}
}

func TestMapUpstreamErrorPreservesRelationConflict(t *testing.T) {
	t.Parallel()

	err := errgrpcx.ToStatus(errs.New(errs.CodeContentRelationAlreadyExists, "content relation already exists"), "fallback")
	mapped := MapUpstreamError(context.Background(), "content_relation_create", "/route", err)
	if !errors.Is(mapped, errs.E(errs.CodeContentRelationAlreadyExists)) {
		t.Fatalf("expected relation conflict to be preserved, got %v", mapped)
	}
}

func TestMapUpstreamErrorMapsAlreadyExistsFallbackByRoute(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		action string
		route  string
		code   errs.Code
	}{
		{
			name:   "relation conflict",
			action: "content_relation_create",
			route:  "/api/v3/studio/content/items/:content_id/relations",
			code:   errs.CodeContentRelationAlreadyExists,
		},
		{
			name:   "tag conflict",
			action: "content_tag_create",
			route:  "/api/v3/studio/content/tags",
			code:   errs.CodeContentTagAlreadyExists,
		},
		{
			name:   "slug conflict",
			action: "content_create",
			route:  "/api/v3/studio/content/items",
			code:   errs.CodeContentSlugAlreadyExists,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mapped := MapUpstreamError(context.Background(), tc.action, tc.route, status.Error(codes.AlreadyExists, "already exists"))
			if !errors.Is(mapped, errs.E(tc.code)) {
				t.Fatalf("expected %d, got %v", tc.code, mapped)
			}
		})
	}
}

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
