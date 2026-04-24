package content

import (
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
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

func ContentTypeToProto(value string) (pb.ContentType, error) {
	switch normalize(value) {
	case "article":
		return pb.ContentType_CONTENT_TYPE_ARTICLE, nil
	case "note":
		return pb.ContentType_CONTENT_TYPE_NOTE, nil
	case "project":
		return pb.ContentType_CONTENT_TYPE_PROJECT, nil
	case "experience":
		return pb.ContentType_CONTENT_TYPE_EXPERIENCE, nil
	case "timeline_event":
		return pb.ContentType_CONTENT_TYPE_TIMELINE_EVENT, nil
	case "insight":
		return pb.ContentType_CONTENT_TYPE_INSIGHT, nil
	case "portfolio":
		return pb.ContentType_CONTENT_TYPE_PORTFOLIO, nil
	case "page":
		return pb.ContentType_CONTENT_TYPE_PAGE, nil
	default:
		return pb.ContentType_CONTENT_TYPE_UNSPECIFIED, errs.New(errs.CodeContentInvalidType, "invalid content type")
	}
}

func OptionalContentTypeToProto(value string) (pb.ContentType, error) {
	if normalize(value) == "" {
		return pb.ContentType_CONTENT_TYPE_UNSPECIFIED, nil
	}
	return ContentTypeToProto(value)
}

func ContentTypeToString(value pb.ContentType) string {
	switch value {
	case pb.ContentType_CONTENT_TYPE_ARTICLE:
		return "article"
	case pb.ContentType_CONTENT_TYPE_NOTE:
		return "note"
	case pb.ContentType_CONTENT_TYPE_PROJECT:
		return "project"
	case pb.ContentType_CONTENT_TYPE_EXPERIENCE:
		return "experience"
	case pb.ContentType_CONTENT_TYPE_TIMELINE_EVENT:
		return "timeline_event"
	case pb.ContentType_CONTENT_TYPE_INSIGHT:
		return "insight"
	case pb.ContentType_CONTENT_TYPE_PORTFOLIO:
		return "portfolio"
	case pb.ContentType_CONTENT_TYPE_PAGE:
		return "page"
	default:
		return ""
	}
}

func StatusToProto(value string) (pb.ContentStatus, error) {
	switch normalize(value) {
	case "draft":
		return pb.ContentStatus_CONTENT_STATUS_DRAFT, nil
	case "review":
		return pb.ContentStatus_CONTENT_STATUS_REVIEW, nil
	case "published":
		return pb.ContentStatus_CONTENT_STATUS_PUBLISHED, nil
	case "archived":
		return pb.ContentStatus_CONTENT_STATUS_ARCHIVED, nil
	default:
		return pb.ContentStatus_CONTENT_STATUS_UNSPECIFIED, errs.New(errs.CodeContentInvalidStatus, "invalid content status")
	}
}

func OptionalStatusToProto(value string) (pb.ContentStatus, error) {
	if normalize(value) == "" {
		return pb.ContentStatus_CONTENT_STATUS_UNSPECIFIED, nil
	}
	return StatusToProto(value)
}

func StatusToString(value pb.ContentStatus) string {
	switch value {
	case pb.ContentStatus_CONTENT_STATUS_DRAFT:
		return "draft"
	case pb.ContentStatus_CONTENT_STATUS_REVIEW:
		return "review"
	case pb.ContentStatus_CONTENT_STATUS_PUBLISHED:
		return "published"
	case pb.ContentStatus_CONTENT_STATUS_ARCHIVED:
		return "archived"
	default:
		return ""
	}
}

func VisibilityToProto(value string) (pb.ContentVisibility, error) {
	switch normalize(value) {
	case "public":
		return pb.ContentVisibility_CONTENT_VISIBILITY_PUBLIC, nil
	case "member":
		return pb.ContentVisibility_CONTENT_VISIBILITY_MEMBER, nil
	case "private":
		return pb.ContentVisibility_CONTENT_VISIBILITY_PRIVATE, nil
	default:
		return pb.ContentVisibility_CONTENT_VISIBILITY_UNSPECIFIED, errs.New(errs.CodeContentInvalidVisibility, "invalid visibility")
	}
}

func VisibilityToProtoDefault(value string) (pb.ContentVisibility, error) {
	if normalize(value) == "" {
		return pb.ContentVisibility_CONTENT_VISIBILITY_UNSPECIFIED, nil
	}
	return VisibilityToProto(value)
}

func OptionalVisibilityToProto(value string) (pb.ContentVisibility, error) {
	return VisibilityToProtoDefault(value)
}

func VisibilityToString(value pb.ContentVisibility) string {
	switch value {
	case pb.ContentVisibility_CONTENT_VISIBILITY_PUBLIC:
		return "public"
	case pb.ContentVisibility_CONTENT_VISIBILITY_MEMBER:
		return "member"
	case pb.ContentVisibility_CONTENT_VISIBILITY_PRIVATE:
		return "private"
	default:
		return ""
	}
}

func AIAccessToProto(value string) (pb.AIAccess, error) {
	switch normalize(value) {
	case "allowed":
		return pb.AIAccess_AI_ACCESS_ALLOWED, nil
	case "denied":
		return pb.AIAccess_AI_ACCESS_DENIED, nil
	default:
		return pb.AIAccess_AI_ACCESS_UNSPECIFIED, errs.New(errs.CodeContentInvalidAIAccess, "invalid ai access")
	}
}

func AIAccessToProtoDefault(value string) (pb.AIAccess, error) {
	if normalize(value) == "" {
		return pb.AIAccess_AI_ACCESS_UNSPECIFIED, nil
	}
	return AIAccessToProto(value)
}

func AIAccessToString(value pb.AIAccess) string {
	switch value {
	case pb.AIAccess_AI_ACCESS_ALLOWED:
		return "allowed"
	case pb.AIAccess_AI_ACCESS_DENIED:
		return "denied"
	default:
		return ""
	}
}

func SourceTypeToProtoDefault(value string) (pb.SourceType, error) {
	if normalize(value) == "" {
		return pb.SourceType_SOURCE_TYPE_UNSPECIFIED, nil
	}
	switch normalize(value) {
	case "manual":
		return pb.SourceType_SOURCE_TYPE_MANUAL, nil
	case "import_v1":
		return pb.SourceType_SOURCE_TYPE_IMPORT_V1, nil
	case "import_markdown":
		return pb.SourceType_SOURCE_TYPE_IMPORT_MARKDOWN, nil
	case "agent_generated":
		return pb.SourceType_SOURCE_TYPE_AGENT_GENERATED, nil
	case "agent_assisted":
		return pb.SourceType_SOURCE_TYPE_AGENT_ASSISTED, nil
	default:
		return pb.SourceType_SOURCE_TYPE_UNSPECIFIED, errs.New(errs.CodeContentInvalidArgument, "invalid source type")
	}
}

func SourceTypeToString(value pb.SourceType) string {
	switch value {
	case pb.SourceType_SOURCE_TYPE_MANUAL:
		return "manual"
	case pb.SourceType_SOURCE_TYPE_IMPORT_V1:
		return "import_v1"
	case pb.SourceType_SOURCE_TYPE_IMPORT_MARKDOWN:
		return "import_markdown"
	case pb.SourceType_SOURCE_TYPE_AGENT_GENERATED:
		return "agent_generated"
	case pb.SourceType_SOURCE_TYPE_AGENT_ASSISTED:
		return "agent_assisted"
	default:
		return ""
	}
}

func EditorTypeToString(value pb.EditorType) string {
	switch value {
	case pb.EditorType_EDITOR_TYPE_HUMAN:
		return "human"
	case pb.EditorType_EDITOR_TYPE_AGENT:
		return "agent"
	case pb.EditorType_EDITOR_TYPE_SYSTEM:
		return "system"
	default:
		return ""
	}
}

func RelationTypeToProto(value string) (pb.ContentRelationType, error) {
	switch normalize(value) {
	case "belongs_to":
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_BELONGS_TO, nil
	case "related_to":
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO, nil
	case "derived_from":
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_DERIVED_FROM, nil
	case "references":
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_REFERENCES, nil
	case "part_of":
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_PART_OF, nil
	case "depends_on":
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_DEPENDS_ON, nil
	case "timeline_of":
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_TIMELINE_OF, nil
	default:
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_UNSPECIFIED, errs.New(errs.CodeContentInvalidArgument, "invalid relation type")
	}
}

func OptionalRelationTypeToProto(value string) (pb.ContentRelationType, error) {
	if normalize(value) == "" {
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_UNSPECIFIED, nil
	}
	return RelationTypeToProto(value)
}

func RelationTypeToString(value pb.ContentRelationType) string {
	switch value {
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_BELONGS_TO:
		return "belongs_to"
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO:
		return "related_to"
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_DERIVED_FROM:
		return "derived_from"
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_REFERENCES:
		return "references"
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_PART_OF:
		return "part_of"
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_DEPENDS_ON:
		return "depends_on"
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_TIMELINE_OF:
		return "timeline_of"
	default:
		return ""
	}
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
