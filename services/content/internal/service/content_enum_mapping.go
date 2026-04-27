package service

import (
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

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

func relationTypeName(value pb.ContentRelationType, allowUnspecified bool) (string, error) {
	switch value {
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_UNSPECIFIED:
		if allowUnspecified {
			return "", nil
		}
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_BELONGS_TO:
		return RelationBelongsTo, nil
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO:
		return RelationRelatedTo, nil
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_DERIVED_FROM:
		return RelationDerivedFrom, nil
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_REFERENCES:
		return RelationReferences, nil
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_PART_OF:
		return RelationPartOf, nil
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_DEPENDS_ON:
		return RelationDependsOn, nil
	case pb.ContentRelationType_CONTENT_RELATION_TYPE_TIMELINE_OF:
		return RelationTimelineOf, nil
	}
	return "", errs.New(errs.CodeContentInvalidArgument, "invalid relation type")
}

func relationTypeProto(value string) pb.ContentRelationType {
	switch value {
	case RelationBelongsTo:
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_BELONGS_TO
	case RelationRelatedTo:
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO
	case RelationDerivedFrom:
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_DERIVED_FROM
	case RelationReferences:
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_REFERENCES
	case RelationPartOf:
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_PART_OF
	case RelationDependsOn:
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_DEPENDS_ON
	case RelationTimelineOf:
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_TIMELINE_OF
	default:
		return pb.ContentRelationType_CONTENT_RELATION_TYPE_UNSPECIFIED
	}
}
