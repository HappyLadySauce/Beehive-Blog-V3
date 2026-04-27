package service

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
)

type Manager struct {
	CreateContent          *CreateContentService
	UpdateContent          *UpdateContentService
	GetContent             *GetContentService
	ListStudioContents     *ListStudioContentsService
	ArchiveContent         *ArchiveContentService
	ListContentRevisions   *ListContentRevisionsService
	GetContentRevision     *GetContentRevisionService
	CreateContentRelation  *CreateContentRelationService
	DeleteContentRelation  *DeleteContentRelationService
	ListContentRelations   *ListContentRelationsService
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
		CreateContentRelation:  &CreateContentRelationService{deps: deps},
		DeleteContentRelation:  &DeleteContentRelationService{deps: deps},
		ListContentRelations:   &ListContentRelationsService{deps: deps},
		CreateTag:              &CreateTagService{deps: deps},
		UpdateTag:              &UpdateTagService{deps: deps},
		DeleteTag:              &DeleteTagService{deps: deps},
		ListTags:               &ListTagsService{deps: deps},
		ListPublicContents:     &ListPublicContentsService{deps: deps},
		GetPublicContentBySlug: &GetPublicContentBySlugService{deps: deps},
		Ping:                   &PingService{deps: deps},
	}
}
