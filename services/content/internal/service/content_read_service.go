package service

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type GetContentService struct{ deps Dependencies }
type ListStudioContentsService struct{ deps Dependencies }
type ListPublicContentsService struct{ deps Dependencies }
type GetPublicContentBySlugService struct{ deps Dependencies }

func (s *GetContentService) Execute(ctx context.Context, actor Actor, req *pb.GetContentRequest) (*pb.GetContentResponse, error) {
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
	item, err := s.deps.Store.Items.GetByID(ctx, contentID)
	if err != nil {
		return nil, mapRepoErr(err, errs.CodeContentNotFound, "content not found")
	}
	detail, err := itemDetail(ctx, s.deps.Store, item)
	if err != nil {
		return nil, err
	}
	return &pb.GetContentResponse{Content: detail}, nil
}

func (s *ListStudioContentsService) Execute(ctx context.Context, actor Actor, req *pb.ListStudioContentsRequest) (*pb.ListStudioContentsResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	page, pageSize := pageValues(req.GetPage(), req.GetPageSize())
	filter, err := studioFilter(req)
	if err != nil {
		return nil, err
	}
	items, total, err := s.deps.Store.Items.List(ctx, filter, page, pageSize)
	if err != nil {
		return nil, internalErr(err)
	}
	summaries, err := contentSummaries(ctx, s.deps.Store, items)
	if err != nil {
		return nil, err
	}
	return &pb.ListStudioContentsResponse{Items: summaries, Total: total, Page: int32(page), PageSize: int32(pageSize)}, nil
}

func (s *ListPublicContentsService) Execute(ctx context.Context, req *pb.ListPublicContentsRequest) (*pb.ListPublicContentsResponse, error) {
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	page, pageSize := pageValues(req.GetPage(), req.GetPageSize())
	filter := repo.ItemFilter{PublicOnly: true, Keyword: strings.TrimSpace(req.GetKeyword())}
	if req.GetType() != pb.ContentType_CONTENT_TYPE_UNSPECIFIED {
		contentType, err := contentTypeName(req.GetType())
		if err != nil {
			return nil, err
		}
		filter.Type = contentType
	}
	items, total, err := s.deps.Store.Items.List(ctx, filter, page, pageSize)
	if err != nil {
		return nil, internalErr(err)
	}
	summaries, err := contentSummaries(ctx, s.deps.Store, items)
	if err != nil {
		return nil, err
	}
	return &pb.ListPublicContentsResponse{Items: summaries, Total: total, Page: int32(page), PageSize: int32(pageSize)}, nil
}

func (s *GetPublicContentBySlugService) Execute(ctx context.Context, req *pb.GetPublicContentBySlugRequest) (*pb.GetPublicContentBySlugResponse, error) {
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	slug := strings.TrimSpace(req.GetSlug())
	if slug == "" {
		return nil, errInvalidArgument("slug is required")
	}
	item, err := s.deps.Store.Items.GetBySlug(ctx, slug)
	if err != nil {
		return nil, mapRepoErr(err, errs.CodeContentNotFound, "content not found")
	}
	if item.Status != StatusPublished || item.Visibility != VisibilityPublic {
		return nil, errs.New(errs.CodeContentNotFound, "content not found")
	}
	detail, err := itemDetail(ctx, s.deps.Store, item)
	if err != nil {
		return nil, err
	}
	return &pb.GetPublicContentBySlugResponse{Content: detail}, nil
}
