package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/khiemnd777/noah_api/modules/i18n/model"
	"github.com/khiemnd777/noah_api/modules/i18n/repository"
	"github.com/khiemnd777/noah_api/shared/cache"
	"github.com/khiemnd777/noah_api/shared/utils/table"
)

var adminKeyPattern = regexp.MustCompile(`^admin\.[a-z0-9_]+\..+$`)

type LanguageService interface {
	List(ctx context.Context, query table.TableQuery) (table.TableListResult[model.LanguageDTO], error)
	GetByID(ctx context.Context, id int) (*model.LanguageDTO, error)
	Create(ctx context.Context, input model.LanguageDTO) (*model.LanguageDTO, error)
	Update(ctx context.Context, id int, input model.LanguageDTO) (*model.LanguageDTO, error)
	Delete(ctx context.Context, id int) error
	ImportXML(ctx context.Context, id int, raw []byte) (*model.LanguageDTO, error)
	ExportXML(ctx context.Context, id int) (*model.LanguageXMLDocument, error)
}

type languageService struct {
	repo repository.LanguageRepository
}

func NewLanguageService(repo repository.LanguageRepository) LanguageService {
	return &languageService{repo: repo}
}

func (s *languageService) List(ctx context.Context, query table.TableQuery) (table.TableListResult[model.LanguageDTO], error) {
	type boxed = table.TableListResult[model.LanguageDTO]
	key := fmt.Sprintf("languages:list:%d:%d:%s:%s", query.Page, query.Limit, valueOrEmpty(query.OrderBy), query.Direction)
	ptr, err := cache.Get(key, cache.TTLMedium, func() (*boxed, error) {
		res, err := s.repo.List(ctx, query)
		if err != nil {
			return nil, err
		}
		return &res, nil
	})
	if err != nil {
		var zero boxed
		return zero, err
	}
	return *ptr, nil
}

func (s *languageService) GetByID(ctx context.Context, id int) (*model.LanguageDTO, error) {
	return cache.Get(languageDetailKey(id), cache.TTLLong, func() (*model.LanguageDTO, error) {
		return s.repo.GetByID(ctx, id)
	})
}

func (s *languageService) Create(ctx context.Context, input model.LanguageDTO) (*model.LanguageDTO, error) {
	if err := s.validateLanguage(input, false); err != nil {
		return nil, err
	}
	res, err := s.repo.Create(ctx, normalizeLanguage(input))
	if err != nil {
		return nil, err
	}
	invalidateLanguageCaches(res.ID)
	return res, nil
}

func (s *languageService) Update(ctx context.Context, id int, input model.LanguageDTO) (*model.LanguageDTO, error) {
	if err := s.validateLanguage(input, true); err != nil {
		return nil, err
	}
	res, err := s.repo.Update(ctx, id, normalizeLanguage(input))
	if err != nil {
		return nil, err
	}
	invalidateLanguageCaches(id)
	return res, nil
}

func (s *languageService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	invalidateLanguageCaches(id)
	return nil
}

func (s *languageService) ImportXML(ctx context.Context, id int, raw []byte) (*model.LanguageDTO, error) {
	var doc model.LanguageXMLDocument
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("invalid xml: %w", err)
	}

	resources := make([]*model.LanguageResourceDTO, 0, len(doc.Resources))
	seen := make(map[string]struct{}, len(doc.Resources))
	for _, item := range doc.Resources {
		key := strings.TrimSpace(item.Key)
		if key == "" {
			return nil, fmt.Errorf("resource key is required")
		}
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("duplicate resource key: %s", key)
		}
		seen[key] = struct{}{}
		resource := &model.LanguageResourceDTO{
			Key:   key,
			Value: item.Value,
		}
		if err := validateResource(resource); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	res, err := s.repo.UpsertImportedResources(ctx, id, resources)
	if err != nil {
		return nil, err
	}
	invalidateLanguageCaches(id)
	return res, nil
}

func (s *languageService) ExportXML(ctx context.Context, id int) (*model.LanguageXMLDocument, error) {
	detail, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resources := slices.Clone(detail.Resources)
	slices.SortFunc(resources, func(a, b *model.LanguageResourceDTO) int {
		return strings.Compare(a.Key, b.Key)
	})
	return &model.LanguageXMLDocument{
		Code:       detail.Code,
		Name:       detail.Name,
		NativeName: detail.NativeName,
		Resources:  resources,
	}, nil
}

func (s *languageService) validateLanguage(input model.LanguageDTO, isUpdate bool) error {
	if strings.TrimSpace(input.Code) == "" {
		return fmt.Errorf("code is required")
	}
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(input.NativeName) == "" {
		return fmt.Errorf("native_name is required")
	}
	for _, item := range input.Resources {
		if err := validateResource(item); err != nil {
			return err
		}
	}
	if !isUpdate && input.Resources == nil {
		return nil
	}

	seen := make(map[string]struct{}, len(input.Resources))
	for _, item := range input.Resources {
		key := strings.TrimSpace(item.Key)
		if _, ok := seen[key]; ok {
			return fmt.Errorf("duplicate resource key: %s", key)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func validateResource(item *model.LanguageResourceDTO) error {
	if item == nil {
		return fmt.Errorf("resource is required")
	}
	if strings.TrimSpace(item.Key) == "" {
		return fmt.Errorf("resource key is required")
	}
	if !adminKeyPattern.MatchString(strings.TrimSpace(item.Key)) {
		return fmt.Errorf("resource key must match admin.{module}.*")
	}
	return nil
}

func normalizeLanguage(input model.LanguageDTO) model.LanguageDTO {
	input.Code = strings.TrimSpace(strings.ToLower(input.Code))
	input.Name = strings.TrimSpace(input.Name)
	input.NativeName = strings.TrimSpace(input.NativeName)
	if input.Resources == nil {
		return input
	}
	resources := make([]*model.LanguageResourceDTO, 0, len(input.Resources))
	for _, item := range input.Resources {
		if item == nil {
			continue
		}
		resources = append(resources, &model.LanguageResourceDTO{
			ID:         item.ID,
			LanguageID: item.LanguageID,
			Key:        strings.TrimSpace(item.Key),
			Value:      item.Value,
		})
	}
	input.Resources = resources
	return input
}

func languageDetailKey(id int) string {
	return fmt.Sprintf("languages:detail:%d", id)
}

func invalidateLanguageCaches(id int) {
	cache.InvalidateKeys(
		"languages:list:*",
		languageDetailKey(id),
	)
}

func valueOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
