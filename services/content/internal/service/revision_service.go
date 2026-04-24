package service

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type ListContentRevisionsService struct{ deps Dependencies }
type GetContentRevisionService struct{ deps Dependencies }

func (s *ListContentRevisionsService) Execute(ctx context.Context, actor Actor, req *pb.ListContentRevisionsRequest) (*pb.ListContentRevisionsResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	contentID, err := parseID(req.GetContentId(), errs.CodeContentInvalidArgument, "invalid content id")
	if err != nil {
		return nil, err
	}
	if _, err := s.deps.Store.Items.GetByID(ctx, contentID); err != nil {
		return nil, mapRepoErr(err, errs.CodeContentNotFound, "content not found")
	}
	page, pageSize := pageValues(req.GetPage(), req.GetPageSize())
	revisions, total, err := s.deps.Store.Revisions.List(ctx, contentID, page, pageSize)
	if err != nil {
		return nil, internalErr(err)
	}
	items := make([]*pb.ContentRevisionSummary, 0, len(revisions))
	for _, revision := range revisions {
		items = append(items, toRevisionSummary(revision))
	}
	return &pb.ListContentRevisionsResponse{Items: items, Total: total, Page: int32(page), PageSize: int32(pageSize)}, nil
}

func (s *GetContentRevisionService) Execute(ctx context.Context, actor Actor, req *pb.GetContentRevisionRequest) (*pb.GetContentRevisionResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	contentID, err := parseID(req.GetContentId(), errs.CodeContentInvalidArgument, "invalid content id")
	if err != nil {
		return nil, err
	}
	revisionID, err := parseID(req.GetRevisionId(), errs.CodeContentInvalidArgument, "invalid revision id")
	if err != nil {
		return nil, err
	}
	revision, err := s.deps.Store.Revisions.GetByID(ctx, contentID, revisionID)
	if err != nil {
		return nil, mapRepoErr(err, errs.CodeContentRevisionNotFound, "content revision not found")
	}
	return &pb.GetContentRevisionResponse{Revision: toRevisionDetail(revision)}, nil
}
