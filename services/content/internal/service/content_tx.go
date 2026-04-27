package service

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	"gorm.io/gorm"
)

func withTransaction(ctx context.Context, store *repo.Store, fn func(txStore *repo.Store) error) error {
	if store == nil || store.DB() == nil {
		return errs.New(errs.CodeContentInternal, "content store is not initialized")
	}
	return store.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(store.WithTx(tx))
	})
}
