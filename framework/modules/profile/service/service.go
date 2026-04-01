package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	frameworkmodel "github.com/khiemnd777/noah_framework/modules/profile/model"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/profile/repository"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
)

type CacheConfig struct {
	MediumTTL time.Duration
}

type Service struct {
	repo  *frameworkrepository.Repository
	cache CacheConfig
}

func New(repo *frameworkrepository.Repository, cache CacheConfig) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) GetProfile(ctx context.Context, userID int) (*frameworkmodel.User, error) {
	return frameworkcache.GetOrSet(profileKey(userID), s.cache.MediumTTL, func() (*frameworkmodel.User, error) {
		return s.repo.GetByID(ctx, userID)
	})
}

func (s *Service) CheckEmailExists(ctx context.Context, userID int, email string) (bool, error) {
	return s.repo.CheckEmailExists(ctx, userID, email)
}

func (s *Service) CheckPhoneExists(ctx context.Context, userID int, phone string) (bool, error) {
	return s.repo.CheckPhoneExists(ctx, userID, phone)
}

func (s *Service) UpdateProfile(ctx context.Context, userID int, name, avatar string, phone, email *string) (*frameworkmodel.User, error) {
	if email != nil && *email != "" {
		if exists, _ := s.repo.CheckEmailExists(ctx, userID, *email); exists {
			return nil, frameworkmodel.ErrEmailExists
		}
	}
	if isPhone(deref(phone)) {
		if exists, _ := s.repo.CheckPhoneExists(ctx, userID, *phone); exists {
			return nil, frameworkmodel.ErrPhoneExists
		}
	}

	var updated *frameworkmodel.User
	err := s.invalidateAfter([]string{
		userKey(userID),
		profileKey(userID),
		fmt.Sprintf("user:%d:bank_qr_code", userID),
		fmt.Sprintf("user:%d:qr_code", userID),
	}, func() error {
		refCode := frameworkmodel.NewReferenceCode()
		qrCode := frameworkmodel.NewQRCode(refCode)
		var err error
		updated, err = s.repo.UpdateByID(ctx, userID, name, phone, email, &avatar, &refCode, &qrCode)
		return err
	})
	return updated, err
}

func (s *Service) ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error {
	return s.repo.ChangePassword(ctx, userID, currentPassword, newPassword)
}

func (s *Service) Delete(ctx context.Context, userID int) error {
	return s.invalidateAfter([]string{
		userKey(userID),
		profileKey(userID),
		fmt.Sprintf("user:%d:*", userID),
	}, func() error {
		return s.repo.Delete(ctx, userID)
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
		if strings.Contains(key, "*") {
			if err := store.DeleteByPattern(key); err != nil {
				return err
			}
			continue
		}
		if err := store.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

func isPhone(input string) bool {
	re := regexp.MustCompile(`^(?:\+84|84|0)\d{9,10}$`)
	value := strings.ReplaceAll(input, " ", "")
	value = strings.ReplaceAll(value, ".", "")
	value = strings.ReplaceAll(value, "-", "")
	return re.MatchString(value)
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func userKey(userID int) string {
	return fmt.Sprintf("user:%d", userID)
}

func profileKey(userID int) string {
	return fmt.Sprintf("profile:%d", userID)
}
