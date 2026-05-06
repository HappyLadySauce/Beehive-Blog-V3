package repo

import "gorm.io/gorm"

type Store struct {
	db *gorm.DB

	Assets     *AssetRepository
	Categories *FileCategoryRepository
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db:         db,
		Assets:     &AssetRepository{db: db},
		Categories: &FileCategoryRepository{db: db},
	}
}

func (s *Store) DB() *gorm.DB {
	return s.db
}

func (s *Store) WithTx(tx *gorm.DB) *Store {
	return NewStore(tx)
}
