package publiccontent

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fakePublicContentClient struct {
	listPublicFn func(context.Context, *pb.ListPublicContentsRequest, ...grpc.CallOption) (*pb.ListPublicContentsResponse, error)
	getPublicFn  func(context.Context, *pb.GetPublicContentBySlugRequest, ...grpc.CallOption) (*pb.GetPublicContentBySlugResponse, error)
}

func (f *fakePublicContentClient) CreateContent(context.Context, *pb.CreateContentRequest, ...grpc.CallOption) (*pb.CreateContentResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) UpdateContent(context.Context, *pb.UpdateContentRequest, ...grpc.CallOption) (*pb.UpdateContentResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) GetContent(context.Context, *pb.GetContentRequest, ...grpc.CallOption) (*pb.GetContentResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) ListStudioContents(context.Context, *pb.ListStudioContentsRequest, ...grpc.CallOption) (*pb.ListStudioContentsResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) ArchiveContent(context.Context, *pb.ArchiveContentRequest, ...grpc.CallOption) (*pb.ArchiveContentResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) ListContentRevisions(context.Context, *pb.ListContentRevisionsRequest, ...grpc.CallOption) (*pb.ListContentRevisionsResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) GetContentRevision(context.Context, *pb.GetContentRevisionRequest, ...grpc.CallOption) (*pb.GetContentRevisionResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) CreateContentRelation(context.Context, *pb.CreateContentRelationRequest, ...grpc.CallOption) (*pb.CreateContentRelationResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) DeleteContentRelation(context.Context, *pb.DeleteContentRelationRequest, ...grpc.CallOption) (*pb.DeleteContentRelationResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) ListContentRelations(context.Context, *pb.ListContentRelationsRequest, ...grpc.CallOption) (*pb.ListContentRelationsResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) CreateTag(context.Context, *pb.CreateTagRequest, ...grpc.CallOption) (*pb.CreateTagResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) UpdateTag(context.Context, *pb.UpdateTagRequest, ...grpc.CallOption) (*pb.UpdateTagResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) DeleteTag(context.Context, *pb.DeleteTagRequest, ...grpc.CallOption) (*pb.DeleteTagResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) ListTags(context.Context, *pb.ListTagsRequest, ...grpc.CallOption) (*pb.ListTagsResponse, error) {
	return nil, nil
}
func (f *fakePublicContentClient) ListPublicContents(ctx context.Context, in *pb.ListPublicContentsRequest, opts ...grpc.CallOption) (*pb.ListPublicContentsResponse, error) {
	return f.listPublicFn(ctx, in, opts...)
}
func (f *fakePublicContentClient) GetPublicContentBySlug(ctx context.Context, in *pb.GetPublicContentBySlugRequest, opts ...grpc.CallOption) (*pb.GetPublicContentBySlugResponse, error) {
	return f.getPublicFn(ctx, in, opts...)
}
func (f *fakePublicContentClient) Ping(context.Context, *pb.PingRequest, ...grpc.CallOption) (*pb.PingResponse, error) {
	return nil, nil
}

func TestPublicContentLogicDoesNotRequireActor(t *testing.T) {
	t.Parallel()

	client := &fakePublicContentClient{
		listPublicFn: func(ctx context.Context, in *pb.ListPublicContentsRequest, _ ...grpc.CallOption) (*pb.ListPublicContentsResponse, error) {
			assertNoActorMetadata(t, ctx)
			return &pb.ListPublicContentsResponse{Items: []*pb.ContentSummary{{ContentId: "1", Type: pb.ContentType_CONTENT_TYPE_ARTICLE}}, Total: 1, Page: 1, PageSize: 20}, nil
		},
		getPublicFn: func(ctx context.Context, in *pb.GetPublicContentBySlugRequest, _ ...grpc.CallOption) (*pb.GetPublicContentBySlugResponse, error) {
			assertNoActorMetadata(t, ctx)
			return &pb.GetPublicContentBySlugResponse{Content: &pb.ContentDetail{ContentId: "1", Slug: in.GetSlug(), Type: pb.ContentType_CONTENT_TYPE_ARTICLE}}, nil
		},
	}
	svcCtx := &svc.ServiceContext{
		Config:        config.Config{ContentRPC: config.ContentRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"}},
		ContentClient: client,
	}
	ctx := middleware.WithRequestMeta(context.Background(), ctxmeta.RequestMeta{RequestID: "req-1"})

	listResp, err := NewPublicContentListLogic(ctx, svcCtx).PublicContentList(&types.PublicContentListReq{Page: 1, PageSize: 20})
	if err != nil || listResp.Total != 1 {
		t.Fatalf("public list failed: resp=%+v err=%v", listResp, err)
	}
	getResp, err := NewPublicContentGetLogic(ctx, svcCtx).PublicContentGet(&types.PublicContentSlugReq{Slug: "post"})
	if err != nil || getResp.Content.Slug != "post" {
		t.Fatalf("public get failed: resp=%+v err=%v", getResp, err)
	}
}

func assertNoActorMetadata(t *testing.T, ctx context.Context) {
	t.Helper()

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}
	if got := md.Get(ctxmeta.MetadataKeyUserID); len(got) != 0 {
		t.Fatalf("did not expect user id metadata, got %v", got)
	}
	if got := md.Get(ctxmeta.MetadataKeyInternalAuthToken); len(got) != 1 || got[0] != "secret" {
		t.Fatalf("expected internal auth token metadata, got %v", got)
	}
}
