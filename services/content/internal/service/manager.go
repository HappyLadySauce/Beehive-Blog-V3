package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"gorm.io/gorm"
)

const (
	StatusDraft     = "draft"
	StatusReview    = "review"
	StatusPublished = "published"
	StatusArchived  = "archived"

	VisibilityPublic  = "public"
	VisibilityMember  = "member"
	VisibilityPrivate = "private"

	AIAccessAllowed = "allowed"
	AIAccessDenied  = "denied"

	EditorHuman = "human"
)

type Manager struct {
	CreateContent          *CreateContentService
	UpdateContent          *UpdateContentService
	GetContent             *GetContentService
	ListStudioContents     *ListStudioContentsService
	ArchiveContent         *ArchiveContentService
	ListContentRevisions   *ListContentRevisionsService
	GetContentRevision     *GetContentRevisionService
	CreateTag              *CreateTagService
	UpdateTag              *UpdateTagService
	DeleteTag              *DeleteTagService
	ListTags               *ListTagsService
	ListPublicContents     *ListPublicContentsService
	GetPublicContentBySlug *GetPublicContentBySlugService
	Ping                   *PingService
}

type Dependencies struct {
	Config         config.Config
	Store          *repo.Store
	Clock          func() time.Time
	CheckReadiness func(ctx context.Context) error
}

type Actor struct {
	UserID    int64
	SessionID int64
	Role      string
}

func NewManager(deps Dependencies) *Manager {
	if deps.Clock == nil {
		deps.Clock = func() time.Time { return time.Now().UTC() }
	}
	return &Manager{
		CreateContent:          &CreateContentService{deps: deps},
		UpdateContent:          &UpdateContentService{deps: deps},
		GetContent:             &GetContentService{deps: deps},
		ListStudioContents:     &ListStudioContentsService{deps: deps},
		ArchiveContent:         &ArchiveContentService{deps: deps},
		ListContentRevisions:   &ListContentRevisionsService{deps: deps},
		GetContentRevision:     &GetContentRevisionService{deps: deps},
		CreateTag:              &CreateTagService{deps: deps},
		UpdateTag:              &UpdateTagService{deps: deps},
		DeleteTag:              &DeleteTagService{deps: deps},
		ListTags:               &ListTagsService{deps: deps},
		ListPublicContents:     &ListPublicContentsService{deps: deps},
		GetPublicContentBySlug: &GetPublicContentBySlugService{deps: deps},
		Ping:                   &PingService{deps: deps},
	}
}

func requireActor(actor Actor) error {
	if actor.UserID <= 0 {
		return errs.New(errs.CodeContentAccessForbidden, "content access forbidden")
	}
	return nil
}

func withTransaction(ctx context.Context, store *repo.Store, fn func(txStore *repo.Store) error) error {
	if store == nil || store.DB() == nil {
		return errs.New(errs.CodeContentInternal, "content store is not initialized")
	}
	return store.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(store.WithTx(tx))
	})
}

func pageValues(page, pageSize int32) (int, int) {
	p := int(page)
	if p <= 0 {
		p = 1
	}
	ps := int(pageSize)
	if ps <= 0 {
		ps = 20
	}
	if ps > 100 {
		ps = 100
	}
	return p, ps
}

func stringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func parseID(value string, code errs.Code, message string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || id <= 0 {
		return 0, errs.New(code, message)
	}
	return id, nil
}

func contentTypeName(value pb.ContentType) (string, error) {
	switch value {
	case pb.ContentType_CONTENT_TYPE_ARTICLE:
		return "article", nil
	case pb.ContentType_CONTENT_TYPE_NOTE:
		return "note", nil
	case pb.ContentType_CONTENT_TYPE_PROJECT:
		return "project", nil
	case pb.ContentType_CONTENT_TYPE_EXPERIENCE:
		return "experience", nil
	case pb.ContentType_CONTENT_TYPE_TIMELINE_EVENT:
		return "timeline_event", nil
	case pb.ContentType_CONTENT_TYPE_INSIGHT:
		return "insight", nil
	case pb.ContentType_CONTENT_TYPE_PORTFOLIO:
		return "portfolio", nil
	case pb.ContentType_CONTENT_TYPE_PAGE:
		return "page", nil
	default:
		return "", errs.New(errs.CodeContentInvalidType, "invalid content type")
	}
}

func contentTypeProto(value string) pb.ContentType {
	switch value {
	case "article":
		return pb.ContentType_CONTENT_TYPE_ARTICLE
	case "note":
		return pb.ContentType_CONTENT_TYPE_NOTE
	case "project":
		return pb.ContentType_CONTENT_TYPE_PROJECT
	case "experience":
		return pb.ContentType_CONTENT_TYPE_EXPERIENCE
	case "timeline_event":
		return pb.ContentType_CONTENT_TYPE_TIMELINE_EVENT
	case "insight":
		return pb.ContentType_CONTENT_TYPE_INSIGHT
	case "portfolio":
		return pb.ContentType_CONTENT_TYPE_PORTFOLIO
	case "page":
		return pb.ContentType_CONTENT_TYPE_PAGE
	default:
		return pb.ContentType_CONTENT_TYPE_UNSPECIFIED
	}
}

func statusName(value pb.ContentStatus, defaultValue string) (string, error) {
	switch value {
	case pb.ContentStatus_CONTENT_STATUS_UNSPECIFIED:
		return defaultValue, nil
	case pb.ContentStatus_CONTENT_STATUS_DRAFT:
		return StatusDraft, nil
	case pb.ContentStatus_CONTENT_STATUS_REVIEW:
		return StatusReview, nil
	case pb.ContentStatus_CONTENT_STATUS_PUBLISHED:
		return StatusPublished, nil
	case pb.ContentStatus_CONTENT_STATUS_ARCHIVED:
		return StatusArchived, nil
	default:
		return "", errs.New(errs.CodeContentInvalidStatus, "invalid content status")
	}
}

func statusProto(value string) pb.ContentStatus {
	switch value {
	case StatusDraft:
		return pb.ContentStatus_CONTENT_STATUS_DRAFT
	case StatusReview:
		return pb.ContentStatus_CONTENT_STATUS_REVIEW
	case StatusPublished:
		return pb.ContentStatus_CONTENT_STATUS_PUBLISHED
	case StatusArchived:
		return pb.ContentStatus_CONTENT_STATUS_ARCHIVED
	default:
		return pb.ContentStatus_CONTENT_STATUS_UNSPECIFIED
	}
}

func visibilityName(value pb.ContentVisibility, defaultValue string) (string, error) {
	switch value {
	case pb.ContentVisibility_CONTENT_VISIBILITY_UNSPECIFIED:
		return defaultValue, nil
	case pb.ContentVisibility_CONTENT_VISIBILITY_PUBLIC:
		return VisibilityPublic, nil
	case pb.ContentVisibility_CONTENT_VISIBILITY_MEMBER:
		return VisibilityMember, nil
	case pb.ContentVisibility_CONTENT_VISIBILITY_PRIVATE:
		return VisibilityPrivate, nil
	default:
		return "", errs.New(errs.CodeContentInvalidVisibility, "invalid visibility")
	}
}

func visibilityProto(value string) pb.ContentVisibility {
	switch value {
	case VisibilityPublic:
		return pb.ContentVisibility_CONTENT_VISIBILITY_PUBLIC
	case VisibilityMember:
		return pb.ContentVisibility_CONTENT_VISIBILITY_MEMBER
	case VisibilityPrivate:
		return pb.ContentVisibility_CONTENT_VISIBILITY_PRIVATE
	default:
		return pb.ContentVisibility_CONTENT_VISIBILITY_UNSPECIFIED
	}
}

func aiAccessName(value pb.AIAccess, defaultValue string) (string, error) {
	switch value {
	case pb.AIAccess_AI_ACCESS_UNSPECIFIED:
		return defaultValue, nil
	case pb.AIAccess_AI_ACCESS_ALLOWED:
		return AIAccessAllowed, nil
	case pb.AIAccess_AI_ACCESS_DENIED:
		return AIAccessDenied, nil
	default:
		return "", errs.New(errs.CodeContentInvalidAIAccess, "invalid ai access")
	}
}

func aiAccessProto(value string) pb.AIAccess {
	if value == AIAccessAllowed {
		return pb.AIAccess_AI_ACCESS_ALLOWED
	}
	if value == AIAccessDenied {
		return pb.AIAccess_AI_ACCESS_DENIED
	}
	return pb.AIAccess_AI_ACCESS_UNSPECIFIED
}

func sourceName(value pb.SourceType, defaultValue string) (string, error) {
	switch value {
	case pb.SourceType_SOURCE_TYPE_UNSPECIFIED:
		return defaultValue, nil
	case pb.SourceType_SOURCE_TYPE_MANUAL:
		return "manual", nil
	case pb.SourceType_SOURCE_TYPE_IMPORT_V1:
		return "import_v1", nil
	case pb.SourceType_SOURCE_TYPE_IMPORT_MARKDOWN:
		return "import_markdown", nil
	case pb.SourceType_SOURCE_TYPE_AGENT_GENERATED:
		return "agent_generated", nil
	case pb.SourceType_SOURCE_TYPE_AGENT_ASSISTED:
		return "agent_assisted", nil
	default:
		return "", errs.New(errs.CodeContentInvalidArgument, "invalid source type")
	}
}

func sourceProto(value string) pb.SourceType {
	switch value {
	case "manual":
		return pb.SourceType_SOURCE_TYPE_MANUAL
	case "import_v1":
		return pb.SourceType_SOURCE_TYPE_IMPORT_V1
	case "import_markdown":
		return pb.SourceType_SOURCE_TYPE_IMPORT_MARKDOWN
	case "agent_generated":
		return pb.SourceType_SOURCE_TYPE_AGENT_GENERATED
	case "agent_assisted":
		return pb.SourceType_SOURCE_TYPE_AGENT_ASSISTED
	default:
		return pb.SourceType_SOURCE_TYPE_UNSPECIFIED
	}
}

func mapRepoErr(err error, notFoundCode errs.Code, notFoundMessage string) error {
	if err == nil {
		return nil
	}
	if repo.IsNotFound(err) {
		return errs.New(notFoundCode, notFoundMessage)
	}
	switch repo.UniqueConstraint(err) {
	case repo.ConstraintItemSlug:
		return errs.Wrap(err, errs.CodeContentSlugAlreadyExists, "slug already exists")
	case repo.ConstraintTagName, repo.ConstraintTagSlug:
		return errs.Wrap(err, errs.CodeContentTagAlreadyExists, "tag already exists")
	}
	return errs.Wrap(err, errs.CodeContentInternal, "content internal error")
}

func validateTitleSlug(title, slug string) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(slug) == "" {
		return errs.New(errs.CodeContentInvalidArgument, "title and slug are required")
	}
	return nil
}

func loadTags(ctx context.Context, store *repo.Store, rawIDs []string) ([]int64, error) {
	result := make([]int64, 0, len(rawIDs))
	seen := map[int64]struct{}{}
	for _, raw := range rawIDs {
		id, err := parseID(raw, errs.CodeContentInvalidArgument, "invalid tag id")
		if err != nil {
			return nil, err
		}
		if _, ok := seen[id]; ok {
			continue
		}
		if _, err := store.Tags.GetByID(ctx, id); err != nil {
			return nil, mapRepoErr(err, errs.CodeContentTagNotFound, "tag not found")
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result, nil
}

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

func internalErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, errs.E(errs.CodeContentInternal)) {
		return err
	}
	return errs.Wrap(err, errs.CodeContentInternal, "content internal error")
}

func serviceNotInitialized() error {
	return errs.New(errs.CodeContentInternal, "content service is not initialized")
}

func errInvalidArgument(message string) error {
	return errs.New(errs.CodeContentInvalidArgument, message)
}

func ensureStore(store *repo.Store) error {
	if store == nil {
		return serviceNotInitialized()
	}
	return nil
}
