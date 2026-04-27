package content

import (
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

func BuildCreateRequest(req *types.ContentCreateReq) (*pb.CreateContentRequest, error) {
	contentType, err := ContentTypeToProto(req.Type)
	if err != nil {
		return nil, err
	}
	visibility, err := VisibilityToProtoDefault(req.Visibility)
	if err != nil {
		return nil, err
	}
	aiAccess, err := AIAccessToProtoDefault(req.AiAccess)
	if err != nil {
		return nil, err
	}
	sourceType, err := SourceTypeToProtoDefault(req.SourceType)
	if err != nil {
		return nil, err
	}
	return &pb.CreateContentRequest{
		Type:           contentType,
		Title:          req.Title,
		Slug:           req.Slug,
		Summary:        req.Summary,
		BodyMarkdown:   req.BodyMarkdown,
		BodyJson:       req.BodyJson,
		CoverImageUrl:  req.CoverImageUrl,
		Visibility:     visibility,
		AiAccess:       aiAccess,
		SourceType:     sourceType,
		CommentEnabled: req.CommentEnabled,
		IsFeatured:     req.IsFeatured,
		SortOrder:      int32(req.SortOrder),
		TagIds:         req.TagIds,
		ChangeSummary:  req.ChangeSummary,
	}, nil
}

func BuildUpdateRequest(req *types.ContentUpdateReq) (*pb.UpdateContentRequest, error) {
	contentType, err := ContentTypeToProto(req.Type)
	if err != nil {
		return nil, err
	}
	status, err := StatusToProto(req.Status)
	if err != nil {
		return nil, err
	}
	visibility, err := VisibilityToProto(req.Visibility)
	if err != nil {
		return nil, err
	}
	aiAccess, err := AIAccessToProto(req.AiAccess)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateContentRequest{
		ContentId:      req.ContentId,
		Type:           contentType,
		Title:          req.Title,
		Slug:           req.Slug,
		Summary:        req.Summary,
		BodyMarkdown:   req.BodyMarkdown,
		BodyJson:       req.BodyJson,
		CoverImageUrl:  req.CoverImageUrl,
		Status:         status,
		Visibility:     visibility,
		AiAccess:       aiAccess,
		CommentEnabled: req.CommentEnabled,
		IsFeatured:     req.IsFeatured,
		SortOrder:      int32(req.SortOrder),
		TagIds:         req.TagIds,
		ChangeSummary:  req.ChangeSummary,
	}, nil
}

func BuildListRequest(req *types.ContentListReq) (*pb.ListStudioContentsRequest, error) {
	contentType, err := OptionalContentTypeToProto(req.Type)
	if err != nil {
		return nil, err
	}
	status, err := OptionalStatusToProto(req.Status)
	if err != nil {
		return nil, err
	}
	visibility, err := OptionalVisibilityToProto(req.Visibility)
	if err != nil {
		return nil, err
	}
	return &pb.ListStudioContentsRequest{
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		Type:       contentType,
		Status:     status,
		Visibility: visibility,
		Keyword:    req.Keyword,
	}, nil
}

func BuildPublicListRequest(req *types.PublicContentListReq) (*pb.ListPublicContentsRequest, error) {
	contentType, err := OptionalContentTypeToProto(req.Type)
	if err != nil {
		return nil, err
	}
	return &pb.ListPublicContentsRequest{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Type:     contentType,
		Keyword:  req.Keyword,
	}, nil
}

func BuildCreateRelationRequest(req *types.ContentRelationCreateReq) (*pb.CreateContentRelationRequest, error) {
	relationType, err := RelationTypeToProto(req.RelationType)
	if err != nil {
		return nil, err
	}
	return &pb.CreateContentRelationRequest{
		ContentId:    req.ContentId,
		ToContentId:  req.ToContentId,
		RelationType: relationType,
		Weight:       int32(req.Weight),
		SortOrder:    int32(req.SortOrder),
		MetadataJson: req.MetadataJson,
	}, nil
}

func BuildListRelationsRequest(req *types.ContentRelationListReq) (*pb.ListContentRelationsRequest, error) {
	relationType, err := OptionalRelationTypeToProto(req.RelationType)
	if err != nil {
		return nil, err
	}
	return &pb.ListContentRelationsRequest{
		ContentId:    req.ContentId,
		Page:         int32(req.Page),
		PageSize:     int32(req.PageSize),
		RelationType: relationType,
	}, nil
}
