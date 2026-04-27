package content

import (
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

func TestEnumMapping(t *testing.T) {
	t.Parallel()

	contentType, err := ContentTypeToProto("timeline_event")
	if err != nil || contentType != pb.ContentType_CONTENT_TYPE_TIMELINE_EVENT {
		t.Fatalf("unexpected content type: %s %v", contentType, err)
	}
	status, err := StatusToProto("published")
	if err != nil || status != pb.ContentStatus_CONTENT_STATUS_PUBLISHED {
		t.Fatalf("unexpected status: %s %v", status, err)
	}
	visibility, err := VisibilityToProto("member")
	if err != nil || visibility != pb.ContentVisibility_CONTENT_VISIBILITY_MEMBER {
		t.Fatalf("unexpected visibility: %s %v", visibility, err)
	}
	aiAccess, err := AIAccessToProto("allowed")
	if err != nil || aiAccess != pb.AIAccess_AI_ACCESS_ALLOWED {
		t.Fatalf("unexpected ai access: %s %v", aiAccess, err)
	}
	sourceType, err := SourceTypeToProtoDefault("agent_assisted")
	if err != nil || sourceType != pb.SourceType_SOURCE_TYPE_AGENT_ASSISTED {
		t.Fatalf("unexpected source type: %s %v", sourceType, err)
	}
	relationType, err := RelationTypeToProto("related_to")
	if err != nil || relationType != pb.ContentRelationType_CONTENT_RELATION_TYPE_RELATED_TO {
		t.Fatalf("unexpected relation type: %s %v", relationType, err)
	}
	if EditorTypeToString(pb.EditorType_EDITOR_TYPE_SYSTEM) != "system" {
		t.Fatalf("unexpected editor type mapping")
	}
}
