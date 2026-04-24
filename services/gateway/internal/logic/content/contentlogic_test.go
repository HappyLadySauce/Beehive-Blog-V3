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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fakeContentClient struct {
	createFn     func(context.Context, *pb.CreateContentRequest, ...grpc.CallOption) (*pb.CreateContentResponse, error)
	updateFn     func(context.Context, *pb.UpdateContentRequest, ...grpc.CallOption) (*pb.UpdateContentResponse, error)
	getFn        func(context.Context, *pb.GetContentRequest, ...grpc.CallOption) (*pb.GetContentResponse, error)
	listFn       func(context.Context, *pb.ListStudioContentsRequest, ...grpc.CallOption) (*pb.ListStudioContentsResponse, error)
	archiveFn    func(context.Context, *pb.ArchiveContentRequest, ...grpc.CallOption) (*pb.ArchiveContentResponse, error)
	listRevsFn   func(context.Context, *pb.ListContentRevisionsRequest, ...grpc.CallOption) (*pb.ListContentRevisionsResponse, error)
	getRevFn     func(context.Context, *pb.GetContentRevisionRequest, ...grpc.CallOption) (*pb.GetContentRevisionResponse, error)
	createRelFn  func(context.Context, *pb.CreateContentRelationRequest, ...grpc.CallOption) (*pb.CreateContentRelationResponse, error)
	deleteRelFn  func(context.Context, *pb.DeleteContentRelationRequest, ...grpc.CallOption) (*pb.DeleteContentRelationResponse, error)
	listRelsFn   func(context.Context, *pb.ListContentRelationsRequest, ...grpc.CallOption) (*pb.ListContentRelationsResponse, error)
	createTagFn  func(context.Context, *pb.CreateTagRequest, ...grpc.CallOption) (*pb.CreateTagResponse, error)
	updateTagFn  func(context.Context, *pb.UpdateTagRequest, ...grpc.CallOption) (*pb.UpdateTagResponse, error)
	deleteTagFn  func(context.Context, *pb.DeleteTagRequest, ...grpc.CallOption) (*pb.DeleteTagResponse, error)
	listTagsFn   func(context.Context, *pb.ListTagsRequest, ...grpc.CallOption) (*pb.ListTagsResponse, error)
	listPublicFn func(context.Context, *pb.ListPublicContentsRequest, ...grpc.CallOption) (*pb.ListPublicContentsResponse, error)
	getPublicFn  func(context.Context, *pb.GetPublicContentBySlugRequest, ...grpc.CallOption) (*pb.GetPublicContentBySlugResponse, error)
	pingFn       func(context.Context, *pb.PingRequest, ...grpc.CallOption) (*pb.PingResponse, error)
}

func (f *fakeContentClient) CreateContent(ctx context.Context, in *pb.CreateContentRequest, opts ...grpc.CallOption) (*pb.CreateContentResponse, error) {
	return f.createFn(ctx, in, opts...)
}
func (f *fakeContentClient) UpdateContent(ctx context.Context, in *pb.UpdateContentRequest, opts ...grpc.CallOption) (*pb.UpdateContentResponse, error) {
	return f.updateFn(ctx, in, opts...)
}
func (f *fakeContentClient) GetContent(ctx context.Context, in *pb.GetContentRequest, opts ...grpc.CallOption) (*pb.GetContentResponse, error) {
	return f.getFn(ctx, in, opts...)
}
func (f *fakeContentClient) ListStudioContents(ctx context.Context, in *pb.ListStudioContentsRequest, opts ...grpc.CallOption) (*pb.ListStudioContentsResponse, error) {
	return f.listFn(ctx, in, opts...)
}
func (f *fakeContentClient) ArchiveContent(ctx context.Context, in *pb.ArchiveContentRequest, opts ...grpc.CallOption) (*pb.ArchiveContentResponse, error) {
	return f.archiveFn(ctx, in, opts...)
}
func (f *fakeContentClient) ListContentRevisions(ctx context.Context, in *pb.ListContentRevisionsRequest, opts ...grpc.CallOption) (*pb.ListContentRevisionsResponse, error) {
	return f.listRevsFn(ctx, in, opts...)
}
func (f *fakeContentClient) GetContentRevision(ctx context.Context, in *pb.GetContentRevisionRequest, opts ...grpc.CallOption) (*pb.GetContentRevisionResponse, error) {
	return f.getRevFn(ctx, in, opts...)
}
func (f *fakeContentClient) CreateContentRelation(ctx context.Context, in *pb.CreateContentRelationRequest, opts ...grpc.CallOption) (*pb.CreateContentRelationResponse, error) {
	return f.createRelFn(ctx, in, opts...)
}
func (f *fakeContentClient) DeleteContentRelation(ctx context.Context, in *pb.DeleteContentRelationRequest, opts ...grpc.CallOption) (*pb.DeleteContentRelationResponse, error) {
	return f.deleteRelFn(ctx, in, opts...)
}
func (f *fakeContentClient) ListContentRelations(ctx context.Context, in *pb.ListContentRelationsRequest, opts ...grpc.CallOption) (*pb.ListContentRelationsResponse, error) {
	return f.listRelsFn(ctx, in, opts...)
}
func (f *fakeContentClient) CreateTag(ctx context.Context, in *pb.CreateTagRequest, opts ...grpc.CallOption) (*pb.CreateTagResponse, error) {
	return f.createTagFn(ctx, in, opts...)
}
func (f *fakeContentClient) UpdateTag(ctx context.Context, in *pb.UpdateTagRequest, opts ...grpc.CallOption) (*pb.UpdateTagResponse, error) {
	return f.updateTagFn(ctx, in, opts...)
}
func (f *fakeContentClient) DeleteTag(ctx context.Context, in *pb.DeleteTagRequest, opts ...grpc.CallOption) (*pb.DeleteTagResponse, error) {
	return f.deleteTagFn(ctx, in, opts...)
}
func (f *fakeContentClient) ListTags(ctx context.Context, in *pb.ListTagsRequest, opts ...grpc.CallOption) (*pb.ListTagsResponse, error) {
	return f.listTagsFn(ctx, in, opts...)
}
func (f *fakeContentClient) ListPublicContents(ctx context.Context, in *pb.ListPublicContentsRequest, opts ...grpc.CallOption) (*pb.ListPublicContentsResponse, error) {
	return f.listPublicFn(ctx, in, opts...)
}
func (f *fakeContentClient) GetPublicContentBySlug(ctx context.Context, in *pb.GetPublicContentBySlugRequest, opts ...grpc.CallOption) (*pb.GetPublicContentBySlugResponse, error) {
	return f.getPublicFn(ctx, in, opts...)
}
func (f *fakeContentClient) Ping(ctx context.Context, in *pb.PingRequest, opts ...grpc.CallOption) (*pb.PingResponse, error) {
	return f.pingFn(ctx, in, opts...)
}

func TestContentLogicSmoke(t *testing.T) {
	t.Parallel()

	client := &fakeContentClient{
		createFn: func(ctx context.Context, in *pb.CreateContentRequest, _ ...grpc.CallOption) (*pb.CreateContentResponse, error) {
			assertActorMetadata(t, ctx)
			if in.GetType() != pb.ContentType_CONTENT_TYPE_ARTICLE || in.GetTitle() != "Title" {
				t.Fatalf("unexpected create request: %+v", in)
			}
			return &pb.CreateContentResponse{Content: contentDetail("1")}, nil
		},
		updateFn: func(ctx context.Context, in *pb.UpdateContentRequest, _ ...grpc.CallOption) (*pb.UpdateContentResponse, error) {
			assertActorMetadata(t, ctx)
			if in.GetStatus() != pb.ContentStatus_CONTENT_STATUS_PUBLISHED {
				t.Fatalf("unexpected update status: %s", in.GetStatus())
			}
			return &pb.UpdateContentResponse{Content: contentDetail(in.GetContentId())}, nil
		},
		getFn: func(ctx context.Context, in *pb.GetContentRequest, _ ...grpc.CallOption) (*pb.GetContentResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.GetContentResponse{Content: contentDetail(in.GetContentId())}, nil
		},
		listFn: func(ctx context.Context, in *pb.ListStudioContentsRequest, _ ...grpc.CallOption) (*pb.ListStudioContentsResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.ListStudioContentsResponse{Items: []*pb.ContentSummary{{ContentId: "1", Type: pb.ContentType_CONTENT_TYPE_ARTICLE}}, Total: 1, Page: 1, PageSize: 20}, nil
		},
		archiveFn: func(ctx context.Context, in *pb.ArchiveContentRequest, _ ...grpc.CallOption) (*pb.ArchiveContentResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.ArchiveContentResponse{Ok: in.GetContentId() == "1"}, nil
		},
		listRevsFn: func(ctx context.Context, in *pb.ListContentRevisionsRequest, _ ...grpc.CallOption) (*pb.ListContentRevisionsResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.ListContentRevisionsResponse{Items: []*pb.ContentRevisionSummary{{RevisionId: "2", ContentId: in.GetContentId(), RevisionNo: 1}}, Total: 1, Page: 1, PageSize: 20}, nil
		},
		getRevFn: func(ctx context.Context, in *pb.GetContentRevisionRequest, _ ...grpc.CallOption) (*pb.GetContentRevisionResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.GetContentRevisionResponse{Revision: &pb.ContentRevisionDetail{RevisionId: in.GetRevisionId(), ContentId: in.GetContentId(), RevisionNo: 1}}, nil
		},
		createRelFn: func(ctx context.Context, in *pb.CreateContentRelationRequest, _ ...grpc.CallOption) (*pb.CreateContentRelationResponse, error) {
			assertActorMetadata(t, ctx)
			if in.GetRelationType() != pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO {
				t.Fatalf("unexpected relation type: %s", in.GetRelationType())
			}
			return &pb.CreateContentRelationResponse{Relation: contentRelation("30", in.GetContentId(), in.GetToContentId())}, nil
		},
		deleteRelFn: func(ctx context.Context, in *pb.DeleteContentRelationRequest, _ ...grpc.CallOption) (*pb.DeleteContentRelationResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.DeleteContentRelationResponse{Ok: in.GetRelationId() == "30"}, nil
		},
		listRelsFn: func(ctx context.Context, in *pb.ListContentRelationsRequest, _ ...grpc.CallOption) (*pb.ListContentRelationsResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.ListContentRelationsResponse{Items: []*pb.ContentRelationView{contentRelation("30", in.GetContentId(), "2")}, Total: 1, Page: 1, PageSize: 20}, nil
		},
		createTagFn: func(ctx context.Context, in *pb.CreateTagRequest, _ ...grpc.CallOption) (*pb.CreateTagResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.CreateTagResponse{Tag: &pb.ContentTag{TagId: "10", Name: in.GetName(), Slug: in.GetSlug()}}, nil
		},
		updateTagFn: func(ctx context.Context, in *pb.UpdateTagRequest, _ ...grpc.CallOption) (*pb.UpdateTagResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.UpdateTagResponse{Tag: &pb.ContentTag{TagId: in.GetTagId(), Name: in.GetName(), Slug: in.GetSlug()}}, nil
		},
		deleteTagFn: func(ctx context.Context, in *pb.DeleteTagRequest, _ ...grpc.CallOption) (*pb.DeleteTagResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.DeleteTagResponse{Ok: in.GetTagId() == "10"}, nil
		},
		listTagsFn: func(ctx context.Context, in *pb.ListTagsRequest, _ ...grpc.CallOption) (*pb.ListTagsResponse, error) {
			assertActorMetadata(t, ctx)
			return &pb.ListTagsResponse{Items: []*pb.ContentTag{{TagId: "10", Name: "Go", Slug: "go"}}, Total: 1, Page: 1, PageSize: 20}, nil
		},
	}
	svcCtx := testServiceContext(client)
	ctx := trustedContext()

	if resp, err := NewContentCreateLogic(ctx, svcCtx).ContentCreate(&types.ContentCreateReq{Type: "article", Title: "Title", Slug: "title"}); err != nil || resp.Content.ContentId != "1" {
		t.Fatalf("create failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentUpdateLogic(ctx, svcCtx).ContentUpdate(&types.ContentUpdateReq{ContentId: "1", Type: "article", Title: "Title", Slug: "title", Status: "published", Visibility: "public", AiAccess: "allowed"}); err != nil || resp.Content.ContentId != "1" {
		t.Fatalf("update failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentGetLogic(ctx, svcCtx).ContentGet(&types.ContentIdReq{ContentId: "1"}); err != nil || resp.Content.ContentId != "1" {
		t.Fatalf("get failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentListLogic(ctx, svcCtx).ContentList(&types.ContentListReq{Page: 1, PageSize: 20}); err != nil || resp.Total != 1 {
		t.Fatalf("list failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentArchiveLogic(ctx, svcCtx).ContentArchive(&types.ContentIdReq{ContentId: "1"}); err != nil || !resp.Ok {
		t.Fatalf("archive failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentRevisionListLogic(ctx, svcCtx).ContentRevisionList(&types.ContentRevisionListReq{ContentId: "1"}); err != nil || resp.Total != 1 {
		t.Fatalf("revision list failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentRevisionGetLogic(ctx, svcCtx).ContentRevisionGet(&types.ContentRevisionIdReq{ContentId: "1", RevisionId: "2"}); err != nil || resp.Revision.RevisionId != "2" {
		t.Fatalf("revision get failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentRelationCreateLogic(ctx, svcCtx).ContentRelationCreate(&types.ContentRelationCreateReq{ContentId: "1", ToContentId: "2", RelationType: "related_to"}); err != nil || resp.Relation.RelationId != "30" {
		t.Fatalf("relation create failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentRelationListLogic(ctx, svcCtx).ContentRelationList(&types.ContentRelationListReq{ContentId: "1", Page: 1, PageSize: 20}); err != nil || resp.Total != 1 {
		t.Fatalf("relation list failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentRelationDeleteLogic(ctx, svcCtx).ContentRelationDelete(&types.ContentRelationIdReq{ContentId: "1", RelationId: "30"}); err != nil || !resp.Ok {
		t.Fatalf("relation delete failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentTagCreateLogic(ctx, svcCtx).ContentTagCreate(&types.ContentTagCreateReq{Name: "Go", Slug: "go"}); err != nil || resp.Tag.TagId != "10" {
		t.Fatalf("tag create failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentTagUpdateLogic(ctx, svcCtx).ContentTagUpdate(&types.ContentTagUpdateReq{TagId: "10", Name: "Go", Slug: "go"}); err != nil || resp.Tag.TagId != "10" {
		t.Fatalf("tag update failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentTagListLogic(ctx, svcCtx).ContentTagList(&types.ContentTagListReq{Page: 1, PageSize: 20}); err != nil || resp.Total != 1 {
		t.Fatalf("tag list failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewContentTagDeleteLogic(ctx, svcCtx).ContentTagDelete(&types.ContentTagIdReq{TagId: "10"}); err != nil || !resp.Ok {
		t.Fatalf("tag delete failed: resp=%+v err=%v", resp, err)
	}
}

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
