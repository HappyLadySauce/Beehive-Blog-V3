package content

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"google.golang.org/grpc"
)

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
