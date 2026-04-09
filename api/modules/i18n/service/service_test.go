package service

import (
	"context"
	"strings"
	"testing"

	"github.com/khiemnd777/noah_api/modules/i18n/model"
	"github.com/khiemnd777/noah_api/shared/utils/table"
)

type repoStub struct {
	getByIDResult       *model.LanguageDTO
	upsertImportedInput []*model.LanguageResourceDTO
}

func (r *repoStub) List(ctx context.Context, query table.TableQuery) (table.TableListResult[model.LanguageDTO], error) {
	return table.TableListResult[model.LanguageDTO]{}, nil
}

func (r *repoStub) GetByID(ctx context.Context, id int) (*model.LanguageDTO, error) {
	return r.getByIDResult, nil
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
