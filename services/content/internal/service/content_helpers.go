package service

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

func studioFilter(req *pb.ListStudioContentsRequest) (repo.ItemFilter, error) {
	filter := repo.ItemFilter{Keyword: strings.TrimSpace(req.GetKeyword())}
	if req.GetType() != pb.ContentType_CONTENT_TYPE_UNSPECIFIED {
		contentType, err := contentTypeName(req.GetType())
		if err != nil {
			return filter, err
		}
		filter.Type = contentType
	}
	if req.GetStatus() != pb.ContentStatus_CONTENT_STATUS_UNSPECIFIED {
		statusValue, err := statusName(req.GetStatus(), "")
		if err != nil {
			return filter, err
		}
		filter.Status = statusValue
	}
	if req.GetVisibility() != pb.ContentVisibility_CONTENT_VISIBILITY_UNSPECIFIED {
		visibility, err := visibilityName(req.GetVisibility(), "")
		if err != nil {
			return filter, err
		}
		filter.Visibility = visibility
	}
	return filter, nil
}

func contentSummaries(ctx context.Context, store *repo.Store, items []entity.Item) ([]*pb.ContentSummary, error) {
	result := make([]*pb.ContentSummary, 0, len(items))
	for i := range items {
		tags, err := store.ContentTags.ListTags(ctx, items[i].ID)
		if err != nil {
			return nil, internalErr(err)
		}
		result = append(result, toContentSummary(&items[i], tags))
	}
	return result, nil
}
