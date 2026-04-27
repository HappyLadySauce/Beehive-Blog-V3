package content

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"google.golang.org/grpc/metadata"
)

func TestContentLogicRequiresTrustedContext(t *testing.T) {
	t.Parallel()

	_, err := NewContentCreateLogic(context.Background(), testServiceContext(&fakeContentClient{})).ContentCreate(&types.ContentCreateReq{Type: "article"})
	if !errors.Is(err, errs.E(errs.CodeGatewayAuthorizationRequired)) {
		t.Fatalf("expected gateway authorization required, got %v", err)
	}
}

func assertActorMetadata(t *testing.T, ctx context.Context) {
	t.Helper()

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}
	if got := md.Get(ctxmeta.MetadataKeyUserID); len(got) != 1 || got[0] != "1" {
		t.Fatalf("expected user id metadata, got %v", got)
	}
	if got := md.Get(ctxmeta.MetadataKeySessionID); len(got) != 1 || got[0] != "s1" {
		t.Fatalf("expected session id metadata, got %v", got)
	}
	if got := md.Get(ctxmeta.MetadataKeyUserRole); len(got) != 1 || got[0] != "admin" {
		t.Fatalf("expected role metadata, got %v", got)
	}
}

func trustedContext() context.Context {
	ctx := middleware.WithRequestMeta(context.Background(), ctxmeta.RequestMeta{RequestID: "req-1"})
	return middleware.WithAuthContext(ctx, middleware.AuthContext{UserID: "1", SessionID: "s1", Role: "admin"})
}

func testServiceContext(client pb.ContentClient) *svc.ServiceContext {
	return &svc.ServiceContext{
		Config:        config.Config{ContentRPC: config.ContentRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"}},
		ContentClient: client,
	}
}

func contentDetail(id string) *pb.ContentDetail {
	return &pb.ContentDetail{
		ContentId: id,
		Type:      pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:     "Title",
		Slug:      "title",
		Status:    pb.ContentStatus_CONTENT_STATUS_DRAFT,
	}
}

func contentRelation(id, fromID, toID string) *pb.ContentRelationView {
	return &pb.ContentRelationView{
		RelationId:    id,
		FromContentId: fromID,
		ToContentId:   toID,
		RelationType:  pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO,
	}
}
