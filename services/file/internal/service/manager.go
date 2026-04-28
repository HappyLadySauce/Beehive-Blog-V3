package service

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/storage"
)

type Dependencies struct {
	Config         config.Config
	Store          *repo.Store
	Storage        storage.ObjectStorage
	CheckReadiness func(context.Context) error
}

type Manager struct {
	conf           config.Config
	store          *repo.Store
	storage        storage.ObjectStorage
	checkReadiness func(context.Context) error
}

func NewManager(deps Dependencies) *Manager {
	return &Manager{
		conf:           deps.Config,
		store:          deps.Store,
		storage:        deps.Storage,
		checkReadiness: deps.CheckReadiness,
	}
}

func (m *Manager) Ping(ctx context.Context) error {
	if m == nil || m.store == nil || m.storage == nil {
		return serviceNotInitialized()
	}
	if m.checkReadiness != nil {
		if err := m.checkReadiness(ctx); err != nil {
			return dependencyUnavailable(err)
		}
	}
	if err := m.storage.Health(ctx); err != nil {
		return dependencyUnavailable(err)
	}
	return nil
}
