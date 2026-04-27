package service

import (
	"context"
	"strconv"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

func itemDetail(ctx context.Context, store *repo.Store, item *entity.Item) (*pb.ContentDetail, error) {
	revision, err := store.Revisions.GetCurrent(ctx, item)
	if err != nil {
		return nil, mapRepoErr(err, errs.CodeContentRevisionNotFound, "content revision not found")
	}
	tags, err := store.ContentTags.ListTags(ctx, item.ID)
	if err != nil {
		return nil, mapRepoErr(err, errs.CodeContentInternal, "content internal error")
	}
	return toContentDetail(item, revision, tags), nil
}

func toContentSummary(item *entity.Item, tags []entity.Tag) *pb.ContentSummary {
	return &pb.ContentSummary{
		ContentId:     strconv.FormatInt(item.ID, 10),
		Type:          contentTypeProto(item.Type),
		Title:         item.Title,
		Slug:          item.Slug,
		Summary:       deref(item.Summary),
		CoverImageUrl: deref(item.CoverImageURL),
		Status:        statusProto(item.Status),
		Visibility:    visibilityProto(item.Visibility),
		AiAccess:      aiAccessProto(item.AIAccess),
		PublishedAt:   unixPtr(item.PublishedAt),
		ArchivedAt:    unixPtr(item.ArchivedAt),
		CreatedAt:     item.CreatedAt.Unix(),
		UpdatedAt:     item.UpdatedAt.Unix(),
		Tags:          toTags(tags),
	}
}

func toContentDetail(item *entity.Item, revision *entity.Revision, tags []entity.Tag) *pb.ContentDetail {
	return &pb.ContentDetail{
		ContentId:         strconv.FormatInt(item.ID, 10),
		Type:              contentTypeProto(item.Type),
		Title:             item.Title,
		Slug:              item.Slug,
		Summary:           deref(item.Summary),
		BodyMarkdown:      revision.BodyMarkdown,
		BodyJson:          deref(revision.BodyJSON),
		CoverImageUrl:     deref(item.CoverImageURL),
		Status:            statusProto(item.Status),
		Visibility:        visibilityProto(item.Visibility),
		AiAccess:          aiAccessProto(item.AIAccess),
		OwnerUserId:       strconv.FormatInt(item.OwnerUserID, 10),
		AuthorUserId:      strconv.FormatInt(item.AuthorUserID, 10),
		SourceType:        sourceProto(item.SourceType),
		CurrentRevisionId: strconv.FormatInt(derefInt64(item.CurrentRevisionID), 10),
		CommentEnabled:    item.CommentEnabled,
		IsFeatured:        item.IsFeatured,
		SortOrder:         item.SortOrder,
		PublishedAt:       unixPtr(item.PublishedAt),
		ArchivedAt:        unixPtr(item.ArchivedAt),
		CreatedAt:         item.CreatedAt.Unix(),
		UpdatedAt:         item.UpdatedAt.Unix(),
		Tags:              toTags(tags),
	}
}

func toRevisionSummary(revision entity.Revision) *pb.ContentRevisionSummary {
	return &pb.ContentRevisionSummary{
		RevisionId:    strconv.FormatInt(revision.ID, 10),
		ContentId:     strconv.FormatInt(revision.ContentID, 10),
		RevisionNo:    revision.RevisionNo,
		EditorType:    pb.EditorType_EDITOR_TYPE_HUMAN,
		ChangeSummary: deref(revision.ChangeSummary),
		SourceType:    sourceProto(revision.SourceType),
		CreatedAt:     revision.CreatedAt.Unix(),
	}
}

func toRevisionDetail(revision *entity.Revision) *pb.ContentRevisionDetail {
	return &pb.ContentRevisionDetail{
		RevisionId:          strconv.FormatInt(revision.ID, 10),
		ContentId:           strconv.FormatInt(revision.ContentID, 10),
		RevisionNo:          revision.RevisionNo,
		TitleSnapshot:       revision.TitleSnapshot,
		SummarySnapshot:     deref(revision.SummarySnapshot),
		BodyMarkdown:        revision.BodyMarkdown,
		BodyJson:            deref(revision.BodyJSON),
		EditorType:          pb.EditorType_EDITOR_TYPE_HUMAN,
		EditorUserId:        strconv.FormatInt(derefInt64(revision.EditorUserID), 10),
		EditorAgentClientId: strconv.FormatInt(derefInt64(revision.EditorAgentClientID), 10),
		ChangeSummary:       deref(revision.ChangeSummary),
		SourceType:          sourceProto(revision.SourceType),
		CreatedAt:           revision.CreatedAt.Unix(),
	}
}

func toRelation(relation *entity.Relation) *pb.ContentRelationView {
	if relation == nil {
		return nil
	}
	return &pb.ContentRelationView{
		RelationId:    strconv.FormatInt(relation.ID, 10),
		FromContentId: strconv.FormatInt(relation.FromContentID, 10),
		ToContentId:   strconv.FormatInt(relation.ToContentID, 10),
		RelationType:  relationTypeProto(relation.RelationType),
		Weight:        relation.Weight,
		SortOrder:     relation.SortOrder,
		MetadataJson:  deref(relation.MetadataJSON),
		CreatedAt:     relation.CreatedAt.Unix(),
		UpdatedAt:     relation.UpdatedAt.Unix(),
	}
}

func toRelations(relations []entity.Relation) []*pb.ContentRelationView {
	result := make([]*pb.ContentRelationView, 0, len(relations))
	for _, relation := range relations {
		result = append(result, toRelation(&relation))
	}
	return result
}

func toTags(tags []entity.Tag) []*pb.ContentTag {
	result := make([]*pb.ContentTag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, toTag(&tag))
	}
	return result
}

func toTag(tag *entity.Tag) *pb.ContentTag {
	return &pb.ContentTag{
		TagId:       strconv.FormatInt(tag.ID, 10),
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: deref(tag.Description),
		Color:       deref(tag.Color),
		CreatedAt:   tag.CreatedAt.Unix(),
		UpdatedAt:   tag.UpdatedAt.Unix(),
	}
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func derefInt64(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

func unixPtr(value *time.Time) int64 {
	if value == nil {
		return 0
	}
	return value.Unix()
}
