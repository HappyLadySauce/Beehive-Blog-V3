package service

import (
	"context"
	"strconv"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type CreateContentRelationService struct{ deps Dependencies }
type DeleteContentRelationService struct{ deps Dependencies }
type ListContentRelationsService struct{ deps Dependencies }

func (s *CreateContentRelationService) Execute(ctx context.Context, actor Actor, req *pb.CreateContentRelationRequest) (*pb.CreateContentRelationResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errInvalidArgument("request is required")
	}
	fromContentID, err := parseID(req.GetContentId(), errs.CodeContentInvalidArgument, "invalid content id")
	if err != nil {
		return nil, err
	}
	toContentID, err := parseID(req.GetToContentId(), errs.CodeContentInvalidArgument, "invalid related content id")
	if err != nil {
		return nil, err
	}
	if fromContentID == toContentID {
		return nil, errs.New(errs.CodeContentInvalidArgument, "content relation cannot point to itself")
	}
	relationType, err := relationTypeName(req.GetRelationType(), false)
	if err != nil {
		return nil, err
	}
	metadataJSON, err := metadataJSONPtr(req.GetMetadataJson())
	if err != nil {
		return nil, err
	}

	var relation *entity.Relation
	err = withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		if _, err := txStore.Items.GetByID(ctx, fromContentID); err != nil {
			return mapRepoErr(err, errs.CodeContentNotFound, "content not found")
		}
		if _, err := txStore.Items.GetByID(ctx, toContentID); err != nil {
			return mapRepoErr(err, errs.CodeContentNotFound, "content not found")
		}
		relation = &entity.Relation{
			FromContentID: fromContentID,
			ToContentID:   toContentID,
			RelationType:  relationType,
			Weight:        req.GetWeight(),
			SortOrder:     req.GetSortOrder(),
			MetadataJSON:  metadataJSON,
		}
		if err := txStore.Relations.Create(ctx, relation); err != nil {
			return mapRepoErr(err, errs.CodeContentRelationNotFound, "content relation not found")
		}
		now := s.deps.Clock()
		payload := baseContentPayload(fromContentID, actor.UserID, now)
		payload["relation_id"] = strconv.FormatInt(relation.ID, 10)
		payload["operation"] = "create"
		payload["to_content_id"] = strconv.FormatInt(toContentID, 10)
		payload["relation_type"] = relationType
		if err := writeOutboxEvent(ctx, txStore, outboxEventInput{EventType: EventContentRelationChanged, ResourceID: fromContentID, Payload: payload, OccurredAt: now}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateContentRelationResponse{Relation: toRelation(relation)}, nil
}

func (s *DeleteContentRelationService) Execute(ctx context.Context, actor Actor, req *pb.DeleteContentRelationRequest) (*pb.DeleteContentRelationResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errInvalidArgument("request is required")
	}
	contentID, err := parseID(req.GetContentId(), errs.CodeContentInvalidArgument, "invalid content id")
	if err != nil {
		return nil, err
	}
	relationID, err := parseID(req.GetRelationId(), errs.CodeContentInvalidArgument, "invalid relation id")
	if err != nil {
		return nil, err
	}
	err = withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		relation, err := txStore.Relations.GetByIDForContent(ctx, contentID, relationID)
		if err != nil {
			return mapRepoErr(err, errs.CodeContentRelationNotFound, "content relation not found")
		}
		if err := txStore.Relations.Delete(ctx, relation); err != nil {
			return mapRepoErr(err, errs.CodeContentRelationNotFound, "content relation not found")
		}
		now := s.deps.Clock()
		payload := baseContentPayload(contentID, actor.UserID, now)
		payload["relation_id"] = strconv.FormatInt(relation.ID, 10)
		payload["operation"] = "delete"
		payload["to_content_id"] = strconv.FormatInt(relation.ToContentID, 10)
		payload["relation_type"] = relation.RelationType
		if err := writeOutboxEvent(ctx, txStore, outboxEventInput{EventType: EventContentRelationChanged, ResourceID: contentID, Payload: payload, OccurredAt: now}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteContentRelationResponse{Ok: true}, nil
}

func (s *ListContentRelationsService) Execute(ctx context.Context, actor Actor, req *pb.ListContentRelationsRequest) (*pb.ListContentRelationsResponse, error) {
	if err := requireActor(actor); err != nil {
		return nil, err
	}
	if err := ensureStore(s.deps.Store); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errInvalidArgument("request is required")
	}
	contentID, err := parseID(req.GetContentId(), errs.CodeContentInvalidArgument, "invalid content id")
	if err != nil {
		return nil, err
	}
	relationType, err := relationTypeName(req.GetRelationType(), true)
	if err != nil {
		return nil, err
	}
	if _, err := s.deps.Store.Items.GetByID(ctx, contentID); err != nil {
		return nil, mapRepoErr(err, errs.CodeContentNotFound, "content not found")
	}
	page, pageSize := pageValues(req.GetPage(), req.GetPageSize())
	relations, total, err := s.deps.Store.Relations.List(ctx, repo.RelationFilter{FromContentID: contentID, RelationType: relationType}, page, pageSize)
	if err != nil {
		return nil, internalErr(err)
	}
	return &pb.ListContentRelationsResponse{Items: toRelations(relations), Total: total, Page: int32(page), PageSize: int32(pageSize)}, nil
}
