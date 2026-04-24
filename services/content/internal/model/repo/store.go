package repo

import "gorm.io/gorm"

type Store struct {
	db *gorm.DB

	Items       *ItemRepository
	Revisions   *RevisionRepository
	Tags        *TagRepository
	ContentTags *ContentTagRepository
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db:          db,
		Items:       &ItemRepository{db: db},
		Revisions:   &RevisionRepository{db: db},
		Tags:        &TagRepository{db: db},
		ContentTags: &ContentTagRepository{db: db},
	}
}

func (s *Store) DB() *gorm.DB {
	return s.db
}

func (s *Store) WithTx(tx *gorm.DB) *Store {
	return NewStore(tx)
}
