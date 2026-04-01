package service

import (
	"context"
	"time"

	frameworkmodel "github.com/khiemnd777/noah_framework/modules/attribute/model"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/attribute/repository"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
)

type CacheConfig struct {
	ShortTTL time.Duration
	LongTTL  time.Duration
}

type AttributeService struct {
	repo  frameworkrepository.Repository
	cache CacheConfig
}

func New(repo frameworkrepository.Repository, cacheConfig CacheConfig) *AttributeService {
	return &AttributeService{repo: repo, cache: cacheConfig}
}

func (s *AttributeService) EnsureSchema(ctx context.Context) error {
	return s.repo.EnsureSchema(ctx)
}

func (s *AttributeService) CreateAttribute(ctx context.Context, userID int, name, attributeType string) (*frameworkmodel.Attribute, error) {
	var result *frameworkmodel.Attribute
	err := s.updateManyAndInvalidate([]string{
		attributeListKey(userID),
		attributeListFirstPageKey(userID),
	}, func() error {
		attr, err := s.repo.CreateAttribute(ctx, frameworkmodel.Attribute{
			UserID:        userID,
			AttributeName: name,
			AttributeType: attributeType,
		})
		result = attr
		return err
	})
	return result, err
}

func (s *AttributeService) GetAttribute(ctx context.Context, id, userID int) (*frameworkmodel.Attribute, error) {
	return frameworkcache.GetOrSet(attributeKey(userID, id), s.cache.ShortTTL, func() (*frameworkmodel.Attribute, error) {
		return s.repo.GetAttributeByID(ctx, id)
	})
}

func (s *AttributeService) ListAttributes(ctx context.Context, userID int) ([]*frameworkmodel.Attribute, error) {
	return cacheList(attributeListKey(userID), s.cache.LongTTL, func() ([]*frameworkmodel.Attribute, error) {
		return s.repo.ListAttributesByUser(ctx, userID)
	})
}

func (s *AttributeService) ListAttributesPaginated(ctx context.Context, userID, page, limit int) ([]*frameworkmodel.Attribute, bool, error) {
	offset := (page - 1) * limit
	if page == 1 {
		return cachePaginatedList(attributeListFirstPageKey(userID), s.cache.LongTTL, func() ([]*frameworkmodel.Attribute, bool, error) {
			return s.repo.ListAttributesByUserPaginated(ctx, userID, limit, offset)
		})
	}
	return s.repo.ListAttributesByUserPaginated(ctx, userID, limit, offset)
}

func (s *AttributeService) UpdateAttribute(ctx context.Context, id, userID int, name, attributeType string) (*frameworkmodel.Attribute, error) {
	err := s.updateManyAndInvalidate([]string{
		attributeListFirstPageKey(userID),
		attributeKey(userID, id),
		attributeListKey(userID),
	}, func() error {
		_, err := s.repo.UpdateAttribute(ctx, id, name, attributeType)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetAttribute(ctx, id, userID)
}

func (s *AttributeService) DeleteAttribute(ctx context.Context, id, userID int) error {
	return s.updateManyAndInvalidate([]string{
		attributeListFirstPageKey(userID),
		attributeKey(userID, id),
		attributeListKey(userID),
	}, func() error {
		return s.repo.SoftDeleteAttribute(ctx, id)
	})
}

type AttributeOptionService struct {
	repo  frameworkrepository.Repository
	cache CacheConfig
}

func NewOptionService(repo frameworkrepository.Repository, cacheConfig CacheConfig) *AttributeOptionService {
	return &AttributeOptionService{repo: repo, cache: cacheConfig}
}

func (s *AttributeOptionService) ListOptions(ctx context.Context, userID, attributeID int) ([]*frameworkmodel.AttributeOption, error) {
	return cacheList(attributeOptionListKey(userID, attributeID), s.cache.LongTTL, func() ([]*frameworkmodel.AttributeOption, error) {
		return s.repo.ListOptionsByAttribute(ctx, attributeID)
	})
}

func (s *AttributeOptionService) CreateOption(ctx context.Context, userID, attributeID int, value string, order int) (*frameworkmodel.AttributeOption, error) {
	var result *frameworkmodel.AttributeOption
	err := s.updateManyAndInvalidate([]string{attributeOptionListKey(userID, attributeID)}, func() error {
		option, err := s.repo.CreateOption(ctx, frameworkmodel.AttributeOption{
			UserID:       userID,
			AttributeID:  attributeID,
			OptionValue:  value,
			DisplayOrder: order,
		})
		result = option
		return err
	})
	return result, err
}

func (s *AttributeOptionService) UpdateOption(ctx context.Context, userID, attributeID, optionID int, value string, order int) (*frameworkmodel.AttributeOption, error) {
	err := s.updateManyAndInvalidate([]string{
		attributeOptionListKey(userID, attributeID),
		attributeOptionKey(userID, attributeID, optionID),
	}, func() error {
		_, err := s.repo.UpdateOption(ctx, optionID, value, order)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetOption(ctx, userID, attributeID, optionID)
}

func (s *AttributeOptionService) GetOption(ctx context.Context, userID, attributeID, optionID int) (*frameworkmodel.AttributeOption, error) {
	return frameworkcache.GetOrSet(attributeOptionKey(userID, attributeID, optionID), s.cache.ShortTTL, func() (*frameworkmodel.AttributeOption, error) {
		return s.repo.GetOptionByID(ctx, optionID)
	})
}

func (s *AttributeOptionService) DeleteOption(ctx context.Context, userID, attributeID, optionID int) error {
	return s.updateManyAndInvalidate([]string{
		attributeOptionListKey(userID, attributeID),
		attributeOptionKey(userID, attributeID, optionID),
	}, func() error {
		return s.repo.SoftDeleteOption(ctx, optionID)
	})
}

func (s *AttributeOptionService) BatchUpdateDisplayOrder(ctx context.Context, userID, attributeID int, orders []frameworkmodel.OptionOrder) error {
	return s.updateManyAndInvalidate([]string{attributeOptionListKey(userID, attributeID)}, func() error {
		return s.repo.BatchUpdateDisplayOrder(ctx, orders)
	})
}

func (s *AttributeService) updateManyAndInvalidate(keys []string, updateFn func() error) error {
	if err := updateFn(); err != nil {
		return err
	}
	return invalidateKeys(keys...)
}

func (s *AttributeOptionService) updateManyAndInvalidate(keys []string, updateFn func() error) error {
	if err := updateFn(); err != nil {
		return err
	}
	return invalidateKeys(keys...)
}

type paginatedList[T any] struct {
	Items   []T  `json:"items"`
	HasMore bool `json:"hasMore"`
}

func cacheList[T any](key string, ttl time.Duration, fallback func() ([]T, error)) ([]T, error) {
	result, err := frameworkcache.GetOrSet(key, ttl, func() (*[]T, error) {
		list, err := fallback()
		if err != nil {
			return nil, err
		}
		return &list, nil
	})
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func cachePaginatedList[T any](key string, ttl time.Duration, fallback func() ([]T, bool, error)) ([]T, bool, error) {
	result, err := frameworkcache.GetOrSet(key, ttl, func() (*paginatedList[T], error) {
		items, hasMore, err := fallback()
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

func invalidateKeys(keys ...string) error {
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

func attributeKey(userID, id int) string {
	return "user:" + itoa(userID) + ":attribute:id:" + itoa(id)
}

func attributeListKey(userID int) string {
	return "user:" + itoa(userID) + ":attribute:list"
}

func attributeListFirstPageKey(userID int) string {
	return "user:" + itoa(userID) + ":attribute:list:first-page"
}

func attributeOptionListKey(userID, attributeID int) string {
	return "user:" + itoa(userID) + ":attribute:" + itoa(attributeID) + ":options:list"
}

func attributeOptionKey(userID, attributeID, optionID int) string {
	return "user:" + itoa(userID) + ":attribute:" + itoa(attributeID) + ":options:" + itoa(optionID)
}
