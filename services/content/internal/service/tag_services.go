package service

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type CreateTagService struct{ deps Dependencies }
type UpdateTagService struct{ deps Dependencies }
type DeleteTagService struct{ deps Dependencies }
type ListTagsService struct{ deps Dependencies }

func (s *CreateTagService) Execute(ctx context.Context, actor Actor, req *pb.CreateTagRequest) (*pb.CreateTagResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	if err := validateTag(req.GetName(), req.GetSlug()); err != nil {
		return nil, err
	}
	tag := &entity.Tag{
		Name:        strings.TrimSpace(req.GetName()),
		Slug:        strings.TrimSpace(req.GetSlug()),
		Description: stringPtr(req.GetDescription()),
		Color:       stringPtr(req.GetColor()),
	}
	if err := s.deps.Store.Tags.Create(ctx, tag); err != nil {
		return nil, mapRepoErr(err, errs.CodeContentTagNotFound, "tag not found")
	}
	return &pb.CreateTagResponse{Tag: toTag(tag)}, nil
}

func (s *UpdateTagService) Execute(ctx context.Context, actor Actor, req *pb.UpdateTagRequest) (*pb.UpdateTagResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	if err := validateTag(req.GetName(), req.GetSlug()); err != nil {
		return nil, err
	}
	tagID, err := parseID(req.GetTagId(), errs.CodeContentInvalidArgument, "invalid tag id")
	if err != nil {
		return nil, err
	}
	tag, err := s.deps.Store.Tags.GetByID(ctx, tagID)
	if err != nil {
		return nil, mapRepoErr(err, errs.CodeContentTagNotFound, "tag not found")
	}
	tag.Name = strings.TrimSpace(req.GetName())
	tag.Slug = strings.TrimSpace(req.GetSlug())
	tag.Description = stringPtr(req.GetDescription())
	tag.Color = stringPtr(req.GetColor())
	if err := s.deps.Store.Tags.Save(ctx, tag); err != nil {
		return nil, mapRepoErr(err, errs.CodeContentTagNotFound, "tag not found")
	}
	return &pb.UpdateTagResponse{Tag: toTag(tag)}, nil
}

func (s *DeleteTagService) Execute(ctx context.Context, actor Actor, req *pb.DeleteTagRequest) (*pb.DeleteTagResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	tagID, err := parseID(req.GetTagId(), errs.CodeContentInvalidArgument, "invalid tag id")
	if err != nil {
		return nil, err
	}
	err = withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		if _, err := txStore.Tags.LockByID(ctx, tagID); err != nil {
			return mapRepoErr(err, errs.CodeContentTagNotFound, "tag not found")
		}
		total, err := txStore.ContentTags.CountByTag(ctx, tagID)
		if err != nil {
			return internalErr(err)
		}
		if total > 0 {
			return errs.New(errs.CodeContentTagInUse, "tag is in use")
		}
		if err := txStore.Tags.Delete(ctx, tagID); err != nil {
			return mapRepoErr(err, errs.CodeContentTagNotFound, "tag not found")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTagResponse{Ok: true}, nil
}

func (s *ListTagsService) Execute(ctx context.Context, actor Actor, req *pb.ListTagsRequest) (*pb.ListTagsResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	page, pageSize := pageValues(req.GetPage(), req.GetPageSize())
	tags, total, err := s.deps.Store.Tags.List(ctx, strings.TrimSpace(req.GetKeyword()), page, pageSize)
	if err != nil {
		return nil, internalErr(err)
	}
	return &pb.ListTagsResponse{Items: toTags(tags), Total: total, Page: int32(page), PageSize: int32(pageSize)}, nil
}

func validateTag(name, slug string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(slug) == "" {
		return errs.New(errs.CodeContentInvalidArgument, "tag name and slug are required")
	}
	return nil
}
