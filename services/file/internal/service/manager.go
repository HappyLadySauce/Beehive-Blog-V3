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
	ObjectStorage  storage.ObjectStorage
	CheckReadiness func(context.Context) error
}

type Manager struct {
	conf           config.Config
	store          *repo.Store
	objectStorage  storage.ObjectStorage
	checkReadiness func(context.Context) error
}

func NewManager(deps Dependencies) *Manager {
	return &Manager{
		conf:           deps.Config,
		store:          deps.Store,
		objectStorage:  deps.ObjectStorage,
		checkReadiness: deps.CheckReadiness,
	}
}

func (m *Manager) Ping(ctx context.Context) error {
	if m == nil || m.store == nil || m.objectStorage == nil {
		return serviceNotInitialized()
	}
	if m.checkReadiness != nil {
		if err := m.checkReadiness(ctx); err != nil {
			return dependencyUnavailable(err)
		}
	}
	return nil
}
