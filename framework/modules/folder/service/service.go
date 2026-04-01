package service

import (
	"context"
	"fmt"
	"time"

	frameworkmodel "github.com/khiemnd777/noah_framework/modules/folder/model"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/folder/repository"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
)

type CacheConfig struct {
	ShortTTL time.Duration
	LongTTL  time.Duration
}

type Service struct {
	repo  *frameworkrepository.Repository
	cache CacheConfig
}

func New(repo *frameworkrepository.Repository, cache CacheConfig) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) Create(ctx context.Context, userID int, input frameworkmodel.Folder) (*frameworkmodel.Folder, error) {
	var result *frameworkmodel.Folder
	err := s.invalidateAfter([]string{folderListKey(userID), folderListFirstPageKey(userID)}, func() error {
		input.UserID = userID
		if input.Color == "" {
			input.Color = "#ffffff"
		}
		folder, err := s.repo.Create(ctx, input)
		result = folder
		return err
	})
	return result, err
}

func (s *Service) Get(ctx context.Context, id, userID int) (*frameworkmodel.Folder, error) {
	return frameworkcache.GetOrSet(folderKey(userID, id), s.cache.ShortTTL, func() (*frameworkmodel.Folder, error) {
		return s.repo.FindByID(ctx, id)
	})
}

func (s *Service) ListPaginated(ctx context.Context, userID, page, limit int) ([]*frameworkmodel.Folder, bool, error) {
	offset := (page - 1) * limit
	if page == 1 {
		return cachePaginated(folderListFirstPageKey(userID), s.cache.LongTTL, func() ([]*frameworkmodel.Folder, bool, error) {
			return s.repo.ListByUserPaginated(ctx, userID, limit, offset)
		})
	}
	return s.repo.ListByUserPaginated(ctx, userID, limit, offset)
}

func (s *Service) Update(ctx context.Context, id, userID int, input frameworkmodel.Folder) (*frameworkmodel.Folder, error) {
	var result *frameworkmodel.Folder
	err := s.invalidateAfter([]string{folderKey(userID, id), folderListKey(userID), folderListFirstPageKey(userID)}, func() error {
		folder, err := s.repo.Update(ctx, id, input)
		result = folder
		return err
	})
	return result, err
}

func (s *Service) Delete(ctx context.Context, id, userID int) error {
	return s.invalidateAfter([]string{folderKey(userID, id), folderListKey(userID), folderListFirstPageKey(userID)}, func() error {
		return s.repo.SoftDelete(ctx, id)
	})
}

func (s *Service) invalidateAfter(keys []string, fn func() error) error {
	if err := fn(); err != nil {
		return err
	}
	store, err := frameworkcache.DefaultStore()
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := store.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

type paginatedList[T any] struct {
	Items   []T  `json:"items"`
	HasMore bool `json:"hasMore"`
}

func cachePaginated[T any](key string, ttl time.Duration, fn func() ([]T, bool, error)) ([]T, bool, error) {
	result, err := frameworkcache.GetOrSet(key, ttl, func() (*paginatedList[T], error) {
		items, hasMore, err := fn()
		if err != nil {
			return nil, err
		}
		return &paginatedList[T]{Items: items, HasMore: hasMore}, nil
	})
	if err != nil {
		return nil, false, err
	}
	return result.Items, result.HasMore, nil
}

func folderKey(userID, id int) string {
	return fmt.Sprintf("user:%d:folder:id:%d", userID, id)
}

func folderListKey(userID int) string {
	return fmt.Sprintf("user:%d:folder:list", userID)
}

func folderListFirstPageKey(userID int) string {
	return fmt.Sprintf("user:%d:folder:list:first-page", userID)
}
