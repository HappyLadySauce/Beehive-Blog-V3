package service

import (
	"errors"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
)

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
	case repo.ConstraintContentRelation:
		return errs.Wrap(err, errs.CodeContentRelationAlreadyExists, "content relation already exists")
	}
	return errs.Wrap(err, errs.CodeContentInternal, "content internal error")
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
