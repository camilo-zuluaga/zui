package sync

import (
	"context"

	"github.com/camilo-zuluaga/zotero-tui/cache"
	"github.com/camilo-zuluaga/zotero-tui/zotero"
)

type SyncService struct {
	DB     *cache.Cache
	Client *zotero.ZoteroClient
}

func New(db *cache.Cache, z *zotero.ZoteroClient) *SyncService {
	return &SyncService{DB: db, Client: z}
}

func (s *SyncService) SyncCollections(ctx context.Context) ([]zotero.Collection, error) {
	cols, err := s.Client.FetchCollections(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.DB.UpsertCollections(cols); err != nil {
		return nil, err
	}

	return cols, nil
}
