package service

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type CreateContentService struct{ deps Dependencies }
type UpdateContentService struct{ deps Dependencies }
type ArchiveContentService struct{ deps Dependencies }

func (s *CreateContentService) Execute(ctx context.Context, actor Actor, req *pb.CreateContentRequest) (*pb.CreateContentResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errInvalidArgument("request is required")
	}
	if err := validateTitleSlug(req.Title, req.Slug); err != nil {
		return nil, err
	}
	contentType, err := contentTypeName(req.Type)
	if err != nil {
		return nil, err
	}
	visibility, err := visibilityName(req.Visibility, VisibilityPrivate)
	if err != nil {
		return nil, err
	}
	aiAccess, err := aiAccessName(req.AiAccess, AIAccessDenied)
	if err != nil {
		return nil, err
	}
	sourceType, err := sourceName(req.SourceType, "manual")
	if err != nil {
		return nil, err
	}
	bodyJSON, err := bodyJSONPtr(req.BodyJson)
	if err != nil {
		return nil, err
	}

	var detail *pb.ContentDetail
	err = withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		tagIDs, err := loadTags(ctx, txStore, req.TagIds)
		if err != nil {
			return err
		}
		item := &entity.Item{
			Type:           contentType,
			Title:          strings.TrimSpace(req.Title),
			Slug:           strings.TrimSpace(req.Slug),
			Status:         StatusDraft,
			Visibility:     visibility,
			AIAccess:       aiAccess,
			Summary:        stringPtr(req.Summary),
			CoverImageURL:  stringPtr(req.CoverImageUrl),
			OwnerUserID:    actor.UserID,
			AuthorUserID:   actor.UserID,
			SourceType:     sourceType,
			CommentEnabled: req.CommentEnabled,
			IsFeatured:     req.IsFeatured,
			SortOrder:      req.SortOrder,
		}
		if err := txStore.Items.Create(ctx, item); err != nil {
			return mapRepoErr(err, errs.CodeContentNotFound, "content not found")
		}
		editorUserID := actor.UserID
		revision := &entity.Revision{
			ContentID:       item.ID,
			RevisionNo:      1,
			TitleSnapshot:   item.Title,
			SummarySnapshot: item.Summary,
			BodyMarkdown:    req.BodyMarkdown,
			BodyJSON:        bodyJSON,
			EditorType:      EditorHuman,
			EditorUserID:    &editorUserID,
			ChangeSummary:   stringPtr(req.ChangeSummary),
			SourceType:      sourceType,
		}
		if err := txStore.Revisions.Create(ctx, revision); err != nil {
			return mapRepoErr(err, errs.CodeContentRevisionNotFound, "content revision not found")
		}
		item.CurrentRevisionID = &revision.ID
		if err := txStore.Items.Save(ctx, item); err != nil {
			return mapRepoErr(err, errs.CodeContentNotFound, "content not found")
		}
		if err := txStore.ContentTags.ReplaceForContent(ctx, item.ID, tagIDs); err != nil {
			return mapRepoErr(err, errs.CodeContentTagNotFound, "tag not found")
		}
		detail, err = itemDetail(ctx, txStore, item)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateContentResponse{Content: detail}, nil
}

func (s *UpdateContentService) Execute(ctx context.Context, actor Actor, req *pb.UpdateContentRequest) (*pb.UpdateContentResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errInvalidArgument("request is required")
	}
	contentID, err := parseID(req.ContentId, errs.CodeContentInvalidArgument, "invalid content id")
	if err != nil {
		return nil, err
	}
	if err := validateTitleSlug(req.Title, req.Slug); err != nil {
		return nil, err
	}
	contentType, err := contentTypeName(req.Type)
	if err != nil {
		return nil, err
	}
	bodyJSON, err := bodyJSONPtr(req.BodyJson)
	if err != nil {
		return nil, err
	}

	var detail *pb.ContentDetail
	err = withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		item, err := txStore.Items.LockByID(ctx, contentID)
		if err != nil {
			return mapRepoErr(err, errs.CodeContentNotFound, "content not found")
		}
		if item.Status == StatusArchived {
			return errs.New(errs.CodeContentInvalidTransition, "archived content cannot be updated")
		}
		statusValue, err := statusName(req.Status, item.Status)
		if err != nil {
			return err
		}
		if statusValue == StatusArchived {
			return errs.New(errs.CodeContentInvalidTransition, "use ArchiveContent to archive content")
		}
		visibility, err := visibilityName(req.Visibility, item.Visibility)
		if err != nil {
			return err
		}
		aiAccess, err := aiAccessName(req.AiAccess, item.AIAccess)
		if err != nil {
			return err
		}
		tagIDs, err := loadTags(ctx, txStore, req.TagIds)
		if err != nil {
			return err
		}
		currentRevision, err := txStore.Revisions.GetCurrent(ctx, item)
		if err != nil {
			return mapRepoErr(err, errs.CodeContentRevisionNotFound, "content revision not found")
		}

		now := s.deps.Clock()
		if item.PublishedAt == nil && statusValue == StatusPublished {
			item.PublishedAt = &now
		}
		item.Type = contentType
		item.Title = strings.TrimSpace(req.Title)
		item.Slug = strings.TrimSpace(req.Slug)
		item.Status = statusValue
		item.Visibility = visibility
		item.AIAccess = aiAccess
		item.Summary = stringPtr(req.Summary)
		item.CoverImageURL = stringPtr(req.CoverImageUrl)
		item.CommentEnabled = req.CommentEnabled
		item.IsFeatured = req.IsFeatured
		item.SortOrder = req.SortOrder

		changedBody := currentRevision.BodyMarkdown != req.BodyMarkdown || deref(currentRevision.BodyJSON) != deref(bodyJSON)
		changedSnapshot := currentRevision.TitleSnapshot != item.Title || deref(currentRevision.SummarySnapshot) != deref(item.Summary)
		changedSummary := strings.TrimSpace(req.ChangeSummary) != ""
		if changedBody || changedSnapshot || changedSummary {
			nextNo, err := txStore.Revisions.NextRevisionNo(ctx, item.ID)
			if err != nil {
				return internalErr(err)
			}
			editorUserID := actor.UserID
			revision := &entity.Revision{
				ContentID:       item.ID,
				RevisionNo:      nextNo,
				TitleSnapshot:   item.Title,
				SummarySnapshot: item.Summary,
				BodyMarkdown:    req.BodyMarkdown,
				BodyJSON:        bodyJSON,
				EditorType:      EditorHuman,
				EditorUserID:    &editorUserID,
				ChangeSummary:   stringPtr(req.ChangeSummary),
				SourceType:      item.SourceType,
			}
			if err := txStore.Revisions.Create(ctx, revision); err != nil {
				return mapRepoErr(err, errs.CodeContentRevisionNotFound, "content revision not found")
			}
			item.CurrentRevisionID = &revision.ID
		}
		if err := txStore.Items.Save(ctx, item); err != nil {
			return mapRepoErr(err, errs.CodeContentNotFound, "content not found")
		}
		if err := txStore.ContentTags.ReplaceForContent(ctx, item.ID, tagIDs); err != nil {
			return mapRepoErr(err, errs.CodeContentTagNotFound, "tag not found")
		}
		detail, err = itemDetail(ctx, txStore, item)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateContentResponse{Content: detail}, nil
}

func (s *ArchiveContentService) Execute(ctx context.Context, actor Actor, req *pb.ArchiveContentRequest) (*pb.ArchiveContentResponse, error) {
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
	err = withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		item, err := txStore.Items.LockByID(ctx, contentID)
		if err != nil {
			return mapRepoErr(err, errs.CodeContentNotFound, "content not found")
		}
		now := s.deps.Clock()
		item.Status = StatusArchived
		item.ArchivedAt = &now
		if err := txStore.Items.Save(ctx, item); err != nil {
			return mapRepoErr(err, errs.CodeContentNotFound, "content not found")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.ArchiveContentResponse{Ok: true}, nil
}
