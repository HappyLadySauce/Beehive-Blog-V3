package service

import (
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
)

func serviceNotInitialized() error {
	return errs.New(errs.CodeFileInternal, "file service is not initialized")
}

func dependencyUnavailable(err error) error {
	return errs.Wrap(err, errs.CodeFileDependencyUnavailable, "file dependency is unavailable")
}

func invalidArgument(message string) error {
	return errs.New(errs.CodeFileInvalidArgument, message)
}
