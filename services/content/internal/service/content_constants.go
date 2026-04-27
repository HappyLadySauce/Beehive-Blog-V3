package service

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

	RelationBelongsTo   = "belongs_to"
	RelationRelatedTo   = "related_to"
	RelationDerivedFrom = "derived_from"
	RelationReferences  = "references"
	RelationPartOf      = "part_of"
	RelationDependsOn   = "depends_on"
	RelationTimelineOf  = "timeline_of"

	RoleAdmin = "admin"

	EditorHuman = "human"
)
