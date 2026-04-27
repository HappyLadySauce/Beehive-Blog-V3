package content

import (
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

func ToContentDetailResp(resp *pb.ContentDetail) *types.ContentDetailResp {
	return &types.ContentDetailResp{Content: ToContentDetail(resp)}
}

func ToContentListResp(items []*pb.ContentSummary, total int64, page, pageSize int32) *types.ContentListResp {
	result := make([]types.ContentSummaryView, 0, len(items))
	for _, item := range items {
		result = append(result, ToContentSummary(item))
	}
	return &types.ContentListResp{
		Items:    result,
		Total:    total,
		Page:     int(page),
		PageSize: int(pageSize),
	}
}

func ToRevisionListResp(items []*pb.ContentRevisionSummary, total int64, page, pageSize int32) *types.ContentRevisionListResp {
	result := make([]types.ContentRevisionSummaryView, 0, len(items))
	for _, item := range items {
		result = append(result, types.ContentRevisionSummaryView{
			RevisionId:    item.GetRevisionId(),
			ContentId:     item.GetContentId(),
			RevisionNo:    int(item.GetRevisionNo()),
			EditorType:    EditorTypeToString(item.GetEditorType()),
			ChangeSummary: item.GetChangeSummary(),
			SourceType:    SourceTypeToString(item.GetSourceType()),
			CreatedAt:     item.GetCreatedAt(),
		})
	}
	return &types.ContentRevisionListResp{Items: result, Total: total, Page: int(page), PageSize: int(pageSize)}
}

func ToRevisionDetailResp(revision *pb.ContentRevisionDetail) *types.ContentRevisionDetailResp {
	if revision == nil {
		return &types.ContentRevisionDetailResp{}
	}
	return &types.ContentRevisionDetailResp{Revision: types.ContentRevisionDetailView{
		RevisionId:          revision.GetRevisionId(),
		ContentId:           revision.GetContentId(),
		RevisionNo:          int(revision.GetRevisionNo()),
		TitleSnapshot:       revision.GetTitleSnapshot(),
		SummarySnapshot:     revision.GetSummarySnapshot(),
		BodyMarkdown:        revision.GetBodyMarkdown(),
		BodyJson:            revision.GetBodyJson(),
		EditorType:          EditorTypeToString(revision.GetEditorType()),
		EditorUserId:        revision.GetEditorUserId(),
		EditorAgentClientId: revision.GetEditorAgentClientId(),
		ChangeSummary:       revision.GetChangeSummary(),
		SourceType:          SourceTypeToString(revision.GetSourceType()),
		CreatedAt:           revision.GetCreatedAt(),
	}}
}

func ToTagResp(tag *pb.ContentTag) *types.ContentTagResp {
	return &types.ContentTagResp{Tag: ToTag(tag)}
}

func ToTagListResp(items []*pb.ContentTag, total int64, page, pageSize int32) *types.ContentTagListResp {
	result := make([]types.ContentTagView, 0, len(items))
	for _, item := range items {
		result = append(result, ToTag(item))
	}
	return &types.ContentTagListResp{Items: result, Total: total, Page: int(page), PageSize: int(pageSize)}
}

func ToRelationResp(relation *pb.ContentRelationView) *types.ContentRelationResp {
	return &types.ContentRelationResp{Relation: ToRelation(relation)}
}

func ToRelationListResp(items []*pb.ContentRelationView, total int64, page, pageSize int32) *types.ContentRelationListResp {
	result := make([]types.ContentRelationView, 0, len(items))
	for _, item := range items {
		result = append(result, ToRelation(item))
	}
	return &types.ContentRelationListResp{Items: result, Total: total, Page: int(page), PageSize: int(pageSize)}
}

func ToContentDetail(item *pb.ContentDetail) types.ContentDetailView {
	if item == nil {
		return types.ContentDetailView{}
	}
	return types.ContentDetailView{
		ContentId:         item.GetContentId(),
		Type:              ContentTypeToString(item.GetType()),
		Title:             item.GetTitle(),
		Slug:              item.GetSlug(),
		Summary:           item.GetSummary(),
		BodyMarkdown:      item.GetBodyMarkdown(),
		BodyJson:          item.GetBodyJson(),
		CoverImageUrl:     item.GetCoverImageUrl(),
		Status:            StatusToString(item.GetStatus()),
		Visibility:        VisibilityToString(item.GetVisibility()),
		AiAccess:          AIAccessToString(item.GetAiAccess()),
		OwnerUserId:       item.GetOwnerUserId(),
		AuthorUserId:      item.GetAuthorUserId(),
		SourceType:        SourceTypeToString(item.GetSourceType()),
		CurrentRevisionId: item.GetCurrentRevisionId(),
		CommentEnabled:    item.GetCommentEnabled(),
		IsFeatured:        item.GetIsFeatured(),
		SortOrder:         int(item.GetSortOrder()),
		PublishedAt:       item.GetPublishedAt(),
		ArchivedAt:        item.GetArchivedAt(),
		CreatedAt:         item.GetCreatedAt(),
		UpdatedAt:         item.GetUpdatedAt(),
		Tags:              ToTags(item.GetTags()),
	}
}

func ToContentSummary(item *pb.ContentSummary) types.ContentSummaryView {
	if item == nil {
		return types.ContentSummaryView{}
	}
	return types.ContentSummaryView{
		ContentId:     item.GetContentId(),
		Type:          ContentTypeToString(item.GetType()),
		Title:         item.GetTitle(),
		Slug:          item.GetSlug(),
		Summary:       item.GetSummary(),
		CoverImageUrl: item.GetCoverImageUrl(),
		Status:        StatusToString(item.GetStatus()),
		Visibility:    VisibilityToString(item.GetVisibility()),
		AiAccess:      AIAccessToString(item.GetAiAccess()),
		PublishedAt:   item.GetPublishedAt(),
		ArchivedAt:    item.GetArchivedAt(),
		CreatedAt:     item.GetCreatedAt(),
		UpdatedAt:     item.GetUpdatedAt(),
		Tags:          ToTags(item.GetTags()),
	}
}

func ToTags(tags []*pb.ContentTag) []types.ContentTagView {
	result := make([]types.ContentTagView, 0, len(tags))
	for _, tag := range tags {
		result = append(result, ToTag(tag))
	}
	return result
}

func ToTag(tag *pb.ContentTag) types.ContentTagView {
	if tag == nil {
		return types.ContentTagView{}
	}
	return types.ContentTagView{
		TagId:       tag.GetTagId(),
		Name:        tag.GetName(),
		Slug:        tag.GetSlug(),
		Description: tag.GetDescription(),
		Color:       tag.GetColor(),
		CreatedAt:   tag.GetCreatedAt(),
		UpdatedAt:   tag.GetUpdatedAt(),
	}
}

func ToRelation(relation *pb.ContentRelationView) types.ContentRelationView {
	if relation == nil {
		return types.ContentRelationView{}
	}
	return types.ContentRelationView{
		RelationId:    relation.GetRelationId(),
		FromContentId: relation.GetFromContentId(),
		ToContentId:   relation.GetToContentId(),
		RelationType:  RelationTypeToString(relation.GetRelationType()),
		Weight:        int(relation.GetWeight()),
		SortOrder:     int(relation.GetSortOrder()),
		MetadataJson:  relation.GetMetadataJson(),
		CreatedAt:     relation.GetCreatedAt(),
		UpdatedAt:     relation.GetUpdatedAt(),
	}
}
