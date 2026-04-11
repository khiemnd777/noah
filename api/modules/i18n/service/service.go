package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
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
	ListActive(ctx context.Context) ([]*model.LanguageOptionDTO, error)
	GetCurrentUserLanguage(ctx context.Context, userID int) (*model.CurrentUserLanguageDTO, error)
	UpdateCurrentUserLanguage(ctx context.Context, userID int, code string) (*model.CurrentUserLanguageDTO, error)
	GetAdminResourcesByCode(ctx context.Context, code string) (*model.AdminResourcesDTO, error)
	Create(ctx context.Context, input model.LanguageDTO) (*model.LanguageDTO, error)
	Update(ctx context.Context, id int, input model.LanguageDTO) (*model.LanguageDTO, error)
	Delete(ctx context.Context, id int) error
	ImportXML(ctx context.Context, id int, raw []byte) (*model.LanguageDTO, error)
	ExportXML(ctx context.Context, id int) (*model.LanguageXMLDocument, error)
	SyncFromDirectory(ctx context.Context, dir string) error
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

func (s *languageService) ListActive(ctx context.Context) ([]*model.LanguageOptionDTO, error) {
	type boxed = []*model.LanguageOptionDTO
	ptr, err := cache.Get(activeLanguagesKey(), cache.TTLMedium, func() (*boxed, error) {
		res, err := s.repo.ListActive(ctx)
		if err != nil {
			return nil, err
		}
		return &res, nil
	})
	if err != nil {
		return nil, err
	}
	return *ptr, nil
}

func (s *languageService) GetCurrentUserLanguage(ctx context.Context, userID int) (*model.CurrentUserLanguageDTO, error) {
	selectedCode, err := s.repo.GetUserLanguagePreference(ctx, userID)
	if err != nil {
		return nil, err
	}
	resolved, err := s.resolveLanguage(ctx, valueOrString(selectedCode))
	if err != nil {
		return nil, err
	}
	return &model.CurrentUserLanguageDTO{
		SelectedCode:  selectedCode,
		EffectiveCode: resolved.Code,
		Language:      toLanguageOption(resolved),
	}, nil
}

func (s *languageService) UpdateCurrentUserLanguage(ctx context.Context, userID int, code string) (*model.CurrentUserLanguageDTO, error) {
	normalized := strings.TrimSpace(strings.ToLower(code))
	if normalized == "" {
		if err := s.repo.SetUserLanguagePreference(ctx, userID, nil); err != nil {
			return nil, err
		}
		cache.InvalidateKeys(userLanguagePreferenceKey(userID))
		return s.GetCurrentUserLanguage(ctx, userID)
	}

	resolved, err := s.repo.GetByCode(ctx, normalized)
	if err != nil {
		return nil, fmt.Errorf("language not found or inactive")
	}
	if err := s.repo.SetUserLanguagePreference(ctx, userID, &resolved.Code); err != nil {
		return nil, err
	}
	cache.InvalidateKeys(userLanguagePreferenceKey(userID))
	return &model.CurrentUserLanguageDTO{
		SelectedCode:  &resolved.Code,
		EffectiveCode: resolved.Code,
		Language:      toLanguageOption(resolved),
	}, nil
}

func (s *languageService) GetAdminResourcesByCode(ctx context.Context, code string) (*model.AdminResourcesDTO, error) {
	resolved, err := s.resolveLanguage(ctx, code)
	if err != nil {
		return nil, err
	}

	resources := make(map[string]string, len(resolved.Resources))
	defaultResources, err := s.resolveDefaultResources(ctx, resolved.Code)
	if err != nil {
		return nil, err
	}
	for key, value := range defaultResources {
		resources[key] = value
	}
	for _, item := range resolved.Resources {
		resources[item.Key] = item.Value
	}

	return &model.AdminResourcesDTO{
		RequestedCode: strings.TrimSpace(strings.ToLower(code)),
		EffectiveCode: resolved.Code,
		Language:      toLanguageOption(resolved),
		Resources:     resources,
	}, nil
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
	doc, err := parseLanguageXML(raw)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.UpsertImportedResources(ctx, id, doc.Resources)
	if err != nil {
		return nil, err
	}
	invalidateLanguageCaches(id)
	return res, nil
}

func (s *languageService) SyncFromDirectory(ctx context.Context, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read language xml dir failed: %w", err)
	}

	fileNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".xml" {
			continue
		}
		fileNames = append(fileNames, entry.Name())
	}
	sort.Strings(fileNames)

	seenCodes := make(map[string]string, len(fileNames))
	for _, name := range fileNames {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("read language xml %q failed: %w", name, err)
		}

		doc, err := parseLanguageXML(raw)
		if err != nil {
			return fmt.Errorf("parse language xml %q failed: %w", name, err)
		}

		if previous, ok := seenCodes[doc.Code]; ok {
			return fmt.Errorf("duplicate language code %q in %q and %q", doc.Code, previous, name)
		}
		seenCodes[doc.Code] = name

		res, err := s.repo.SyncImportedLanguage(ctx, model.LanguageDTO{
			Code:       doc.Code,
			Name:       doc.Name,
			NativeName: doc.NativeName,
			Active:     true,
			Resources:  doc.Resources,
		})
		if err != nil {
			return fmt.Errorf("sync language %q from %q failed: %w", doc.Code, name, err)
		}
		invalidateLanguageCaches(res.ID)
	}

	return nil
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

func parseLanguageXML(raw []byte) (*model.LanguageDTO, error) {
	var doc model.LanguageXMLDocument
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("invalid xml: %w", err)
	}

	normalized := normalizeLanguage(model.LanguageDTO{
		Code:       doc.Code,
		Name:       doc.Name,
		NativeName: doc.NativeName,
		Resources:  doc.Resources,
	})

	if strings.TrimSpace(normalized.Code) == "" {
		return nil, fmt.Errorf("code is required")
	}
	if strings.TrimSpace(normalized.Name) == "" {
		return nil, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(normalized.NativeName) == "" {
		return nil, fmt.Errorf("native_name is required")
	}

	seen := make(map[string]struct{}, len(normalized.Resources))
	for _, item := range normalized.Resources {
		if err := validateResource(item); err != nil {
			return nil, err
		}
		if _, ok := seen[item.Key]; ok {
			return nil, fmt.Errorf("duplicate resource key: %s", item.Key)
		}
		seen[item.Key] = struct{}{}
	}

	return &normalized, nil
}

func languageDetailKey(id int) string {
	return fmt.Sprintf("languages:detail:%d", id)
}

func invalidateLanguageCaches(id int) {
	cache.InvalidateKeys(
		"languages:list:*",
		activeLanguagesKey(),
		languageDetailKey(id),
	)
}

func valueOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func valueOrString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func activeLanguagesKey() string {
	return "languages:active"
}

func userLanguagePreferenceKey(userID int) string {
	return fmt.Sprintf("languages:user:%d:preference", userID)
}

func toLanguageOption(item *model.LanguageDTO) *model.LanguageOptionDTO {
	if item == nil {
		return nil
	}
	return &model.LanguageOptionDTO{
		Code:       item.Code,
		Name:       item.Name,
		NativeName: item.NativeName,
		IsDefault:  item.IsDefault,
	}
}

func (s *languageService) resolveLanguage(ctx context.Context, code string) (*model.LanguageDTO, error) {
	normalized := strings.TrimSpace(strings.ToLower(code))
	if normalized != "" {
		if language, err := s.repo.GetByCode(ctx, normalized); err == nil {
			return language, nil
		}
	}
	return s.repo.GetDefaultOrFirstActive(ctx)
}

func (s *languageService) resolveDefaultResources(ctx context.Context, effectiveCode string) (map[string]string, error) {
	fallback, err := s.repo.GetDefaultOrFirstActive(ctx)
	if err != nil {
		return nil, err
	}
	if fallback.Code == effectiveCode {
		return map[string]string{}, nil
	}

	res := make(map[string]string, len(fallback.Resources))
	for _, item := range fallback.Resources {
		res[item.Key] = item.Value
	}
	return res, nil
}
