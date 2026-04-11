package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/khiemnd777/noah_api/modules/i18n/model"
	"github.com/khiemnd777/noah_api/shared/utils/table"
)

type repoStub struct {
	getByIDResult       *model.LanguageDTO
	getByCodeResult     *model.LanguageDTO
	getDefaultResult    *model.LanguageDTO
	listActiveResult    []*model.LanguageOptionDTO
	preferenceResult    *string
	updatedPreference   *string
	upsertImportedInput []*model.LanguageResourceDTO
	getByCodeErr        error
	syncImportedInput   *model.LanguageDTO
	syncImportedResult  *model.LanguageDTO
	syncImportedErr     error
}

func (r *repoStub) List(ctx context.Context, query table.TableQuery) (table.TableListResult[model.LanguageDTO], error) {
	return table.TableListResult[model.LanguageDTO]{}, nil
}

func (r *repoStub) GetByID(ctx context.Context, id int) (*model.LanguageDTO, error) {
	return r.getByIDResult, nil
}

func (r *repoStub) GetByCode(ctx context.Context, code string) (*model.LanguageDTO, error) {
	if r.getByCodeErr != nil {
		return nil, r.getByCodeErr
	}
	return r.getByCodeResult, nil
}

func (r *repoStub) GetByCodeForSync(ctx context.Context, code string) (*model.LanguageDTO, error) {
	return r.getByCodeResult, nil
}

func (r *repoStub) ListActive(ctx context.Context) ([]*model.LanguageOptionDTO, error) {
	return r.listActiveResult, nil
}

func (r *repoStub) GetDefaultOrFirstActive(ctx context.Context) (*model.LanguageDTO, error) {
	return r.getDefaultResult, nil
}

func (r *repoStub) GetUserLanguagePreference(ctx context.Context, userID int) (*string, error) {
	return r.preferenceResult, nil
}

func (r *repoStub) SetUserLanguagePreference(ctx context.Context, userID int, code *string) error {
	r.updatedPreference = code
	return nil
}

func (r *repoStub) Create(ctx context.Context, input model.LanguageDTO) (*model.LanguageDTO, error) {
	return &input, nil
}

func (r *repoStub) Update(ctx context.Context, id int, input model.LanguageDTO) (*model.LanguageDTO, error) {
	return &input, nil
}

func (r *repoStub) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *repoStub) UpsertImportedResources(ctx context.Context, id int, resources []*model.LanguageResourceDTO) (*model.LanguageDTO, error) {
	r.upsertImportedInput = resources
	return &model.LanguageDTO{ID: id, Resources: resources}, nil
}

func (r *repoStub) SyncImportedLanguage(ctx context.Context, input model.LanguageDTO) (*model.LanguageDTO, error) {
	if r.syncImportedErr != nil {
		return nil, r.syncImportedErr
	}
	copied := input
	r.syncImportedInput = &copied
	if r.syncImportedResult != nil {
		return r.syncImportedResult, nil
	}
	return &model.LanguageDTO{ID: 1, Code: input.Code, Name: input.Name, NativeName: input.NativeName, Resources: input.Resources}, nil
}

func TestImportXMLRejectsDuplicateKeys(t *testing.T) {
	svc := NewLanguageService(&repoStub{})

	_, err := svc.ImportXML(context.Background(), 1, []byte(`
<language code="vi" name="Vietnamese" native_name="Tiếng Việt">
  <resources>
    <resource key="admin.settings.title">A</resource>
    <resource key="admin.settings.title">B</resource>
  </resources>
</language>`))
	if err == nil {
		t.Fatalf("expected duplicate key validation error")
	}
	if !strings.Contains(err.Error(), "duplicate resource key") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateRejectsInvalidResourceKey(t *testing.T) {
	svc := NewLanguageService(&repoStub{})

	_, err := svc.Create(context.Background(), model.LanguageDTO{
		Code:       "vi",
		Name:       "Vietnamese",
		NativeName: "Tiếng Việt",
		Active:     true,
		Resources: []*model.LanguageResourceDTO{
			{Key: "settings.title", Value: "Cài đặt"},
		},
	})
	if err == nil {
		t.Fatalf("expected invalid resource key validation error")
	}
	if !strings.Contains(err.Error(), "admin.{module}.*") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetCurrentUserLanguageFallsBackToDefault(t *testing.T) {
	selected := "fr"
	svc := NewLanguageService(&repoStub{
		preferenceResult: &selected,
		getDefaultResult: &model.LanguageDTO{
			Code:       "vi",
			Name:       "Vietnamese",
			NativeName: "Tiếng Việt",
			IsDefault:  true,
		},
		getByCodeErr: context.DeadlineExceeded,
	})

	res, err := svc.GetCurrentUserLanguage(context.Background(), 99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.EffectiveCode != "vi" {
		t.Fatalf("expected fallback effective code vi, got %q", res.EffectiveCode)
	}
	if res.SelectedCode == nil || *res.SelectedCode != "fr" {
		t.Fatalf("expected original selected code to be preserved")
	}
}

func TestUpdateCurrentUserLanguageRejectsInactiveLanguage(t *testing.T) {
	svc := NewLanguageService(&repoStub{
		getByCodeErr: context.DeadlineExceeded,
	})

	_, err := svc.UpdateCurrentUserLanguage(context.Background(), 5, "fr")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "language not found or inactive") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAdminResourcesByCodeMergesDefaultFallback(t *testing.T) {
	svc := NewLanguageService(&repoStub{
		getByCodeResult: &model.LanguageDTO{
			Code:       "en",
			Name:       "English",
			NativeName: "English",
			Resources: []*model.LanguageResourceDTO{
				{Key: "admin.settings.title", Value: "Settings"},
			},
		},
		getDefaultResult: &model.LanguageDTO{
			Code:       "vi",
			Name:       "Vietnamese",
			NativeName: "Tiếng Việt",
			IsDefault:  true,
			Resources: []*model.LanguageResourceDTO{
				{Key: "admin.settings.title", Value: "Cài đặt"},
				{Key: "admin.settings.save", Value: "Lưu"},
			},
		},
	})

	res, err := svc.GetAdminResourcesByCode(context.Background(), "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Resources["admin.settings.title"] != "Settings" {
		t.Fatalf("expected selected language to override default value")
	}
	if res.Resources["admin.settings.save"] != "Lưu" {
		t.Fatalf("expected missing key to fall back to default language")
	}
}

func TestSyncFromDirectoryImportsLanguageDocument(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "en.xml"), []byte(`
<language code="EN" name="English" native_name="English">
  <resources>
    <resource key="admin.settings.page_title">Settings</resource>
  </resources>
</language>`), 0o644); err != nil {
		t.Fatalf("write xml: %v", err)
	}

	repo := &repoStub{
		syncImportedResult: &model.LanguageDTO{ID: 7, Code: "en"},
	}
	svc := NewLanguageService(repo)

	if err := svc.SyncFromDirectory(context.Background(), dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.syncImportedInput == nil {
		t.Fatalf("expected sync input to be captured")
	}
	if repo.syncImportedInput.Code != "en" {
		t.Fatalf("expected normalized code en, got %q", repo.syncImportedInput.Code)
	}
	if repo.syncImportedInput.Active != true {
		t.Fatalf("expected synced languages to default active")
	}
	if len(repo.syncImportedInput.Resources) != 1 || repo.syncImportedInput.Resources[0].Key != "admin.settings.page_title" {
		t.Fatalf("expected parsed resources to be forwarded")
	}
}

func TestSyncFromDirectoryRejectsDuplicateLanguageCodes(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"a.xml": `<language code="en" name="English" native_name="English"><resources></resources></language>`,
		"b.xml": `<language code="EN" name="English 2" native_name="English"><resources></resources></language>`,
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatalf("write xml %s: %v", name, err)
		}
	}

	svc := NewLanguageService(&repoStub{})
	err := svc.SyncFromDirectory(context.Background(), dir)
	if err == nil {
		t.Fatalf("expected duplicate code error")
	}
	if !strings.Contains(err.Error(), `duplicate language code "en"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSyncFromDirectoryRejectsMalformedXML(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "broken.xml"), []byte(`<language code="en"`), 0o644); err != nil {
		t.Fatalf("write xml: %v", err)
	}

	svc := NewLanguageService(&repoStub{})
	err := svc.SyncFromDirectory(context.Background(), dir)
	if err == nil {
		t.Fatalf("expected parse error")
	}
	if !strings.Contains(err.Error(), `parse language xml "broken.xml" failed`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSyncFromDirectoryPropagatesRepositorySyncFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "en.xml"), []byte(`
<language code="en" name="English" native_name="English">
  <resources>
    <resource key="admin.settings.page_title">Settings</resource>
  </resources>
</language>`), 0o644); err != nil {
		t.Fatalf("write xml: %v", err)
	}

	svc := NewLanguageService(&repoStub{syncImportedErr: fmt.Errorf("db unavailable")})
	err := svc.SyncFromDirectory(context.Background(), dir)
	if err == nil {
		t.Fatalf("expected sync failure")
	}
	if !strings.Contains(err.Error(), `sync language "en" from "en.xml" failed`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
