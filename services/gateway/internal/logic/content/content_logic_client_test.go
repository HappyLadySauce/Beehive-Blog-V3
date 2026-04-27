package content

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"google.golang.org/grpc"
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
