package repository

import (
	"context"
	"strings"

	"github.com/khiemnd777/noah_api/modules/i18n/model"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated/language"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated/languageresource"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated/user"
	"github.com/khiemnd777/noah_api/shared/mapper"
	"github.com/khiemnd777/noah_api/shared/utils/table"
)

type LanguageRepository interface {
	List(ctx context.Context, query table.TableQuery) (table.TableListResult[model.LanguageDTO], error)
	GetByID(ctx context.Context, id int) (*model.LanguageDTO, error)
	GetByCode(ctx context.Context, code string) (*model.LanguageDTO, error)
	GetByCodeForSync(ctx context.Context, code string) (*model.LanguageDTO, error)
	ListActive(ctx context.Context) ([]*model.LanguageOptionDTO, error)
	GetDefaultOrFirstActive(ctx context.Context) (*model.LanguageDTO, error)
	GetUserLanguagePreference(ctx context.Context, userID int) (*string, error)
	SetUserLanguagePreference(ctx context.Context, userID int, code *string) error
	Create(ctx context.Context, input model.LanguageDTO) (*model.LanguageDTO, error)
	Update(ctx context.Context, id int, input model.LanguageDTO) (*model.LanguageDTO, error)
	Delete(ctx context.Context, id int) error
	UpsertImportedResources(ctx context.Context, id int, resources []*model.LanguageResourceDTO) (*model.LanguageDTO, error)
	SyncImportedLanguage(ctx context.Context, input model.LanguageDTO) (*model.LanguageDTO, error)
}

type languageRepository struct {
	db *generated.Client
}

func NewLanguageRepository(db *generated.Client) LanguageRepository {
	return &languageRepository{db: db}
}

func (r *languageRepository) List(ctx context.Context, query table.TableQuery) (table.TableListResult[model.LanguageDTO], error) {
	return table.TableListV2(
		ctx,
		r.db.Language.Query().Where(language.Deleted(false)),
		query,
		language.Table,
		language.FieldID,
		language.FieldID,
		func(q *generated.LanguageQuery) *generated.LanguageQuery {
			return q
		},
		func(src []*generated.Language) []*model.LanguageDTO {
			return mapper.MapListAs[*generated.Language, *model.LanguageDTO](src)
		},
	)
}

func (r *languageRepository) GetByID(ctx context.Context, id int) (*model.LanguageDTO, error) {
	entity, err := r.db.Language.Query().
		Where(language.IDEQ(id), language.Deleted(false)).
		WithResources(func(q *generated.LanguageResourceQuery) {
			q.Order(generated.Asc(languageresource.FieldKey))
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return toLanguageDTO(entity), nil
}

func (r *languageRepository) GetByCode(ctx context.Context, code string) (*model.LanguageDTO, error) {
	entity, err := r.db.Language.Query().
		Where(language.CodeEQ(strings.TrimSpace(strings.ToLower(code))), language.Deleted(false), language.Active(true)).
		WithResources(func(q *generated.LanguageResourceQuery) {
			q.Order(generated.Asc(languageresource.FieldKey))
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return toLanguageDTO(entity), nil
}

func (r *languageRepository) GetByCodeForSync(ctx context.Context, code string) (*model.LanguageDTO, error) {
	entity, err := r.db.Language.Query().
		Where(language.CodeEQ(strings.TrimSpace(strings.ToLower(code))), language.Deleted(false)).
		WithResources(func(q *generated.LanguageResourceQuery) {
			q.Order(generated.Asc(languageresource.FieldKey))
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return toLanguageDTO(entity), nil
}

func (r *languageRepository) ListActive(ctx context.Context) ([]*model.LanguageOptionDTO, error) {
	items, err := r.db.Language.Query().
		Where(language.Deleted(false), language.Active(true)).
		Order(generated.Desc(language.FieldIsDefault), generated.Asc(language.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*model.LanguageOptionDTO, 0, len(items))
	for _, item := range items {
		res = append(res, &model.LanguageOptionDTO{
			Code:       item.Code,
			Name:       item.Name,
			NativeName: item.NativeName,
			IsDefault:  item.IsDefault,
		})
	}
	return res, nil
}

func (r *languageRepository) GetDefaultOrFirstActive(ctx context.Context) (*model.LanguageDTO, error) {
	entity, err := r.db.Language.Query().
		Where(language.Deleted(false), language.Active(true)).
		Order(generated.Desc(language.FieldIsDefault), generated.Asc(language.FieldName)).
		WithResources(func(q *generated.LanguageResourceQuery) {
			q.Order(generated.Asc(languageresource.FieldKey))
		}).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toLanguageDTO(entity), nil
}

func (r *languageRepository) GetUserLanguagePreference(ctx context.Context, userID int) (*string, error) {
	current, err := r.db.User.Query().
		Where(user.ID(userID), user.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return current.AdminLanguageCode, nil
}

func (r *languageRepository) SetUserLanguagePreference(ctx context.Context, userID int, code *string) error {
	return r.db.User.UpdateOneID(userID).
		SetNillableAdminLanguageCode(code).
		Exec(ctx)
}

func (r *languageRepository) Create(ctx context.Context, input model.LanguageDTO) (*model.LanguageDTO, error) {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback(tx)

	if err := clearDefaultIfNeeded(ctx, tx.Client(), input.IsDefault, 0); err != nil {
		return nil, err
	}

	created, err := tx.Language.Create().
		SetCode(strings.TrimSpace(input.Code)).
		SetName(strings.TrimSpace(input.Name)).
		SetNativeName(strings.TrimSpace(input.NativeName)).
		SetIsDefault(input.IsDefault).
		SetActive(input.Active).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := syncResources(ctx, tx.Client(), created.ID, input.Resources); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, created.ID)
}

func (r *languageRepository) Update(ctx context.Context, id int, input model.LanguageDTO) (*model.LanguageDTO, error) {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback(tx)

	if err := clearDefaultIfNeeded(ctx, tx.Client(), input.IsDefault, id); err != nil {
		return nil, err
	}

	_, err = tx.Language.UpdateOneID(id).
		Where(language.Deleted(false)).
		SetCode(strings.TrimSpace(input.Code)).
		SetName(strings.TrimSpace(input.Name)).
		SetNativeName(strings.TrimSpace(input.NativeName)).
		SetIsDefault(input.IsDefault).
		SetActive(input.Active).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if input.Resources != nil {
		if err := syncResources(ctx, tx.Client(), id, input.Resources); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *languageRepository) Delete(ctx context.Context, id int) error {
	return r.db.Language.UpdateOneID(id).
		Where(language.Deleted(false)).
		SetDeleted(true).
		SetIsDefault(false).
		Exec(ctx)
}

func (r *languageRepository) UpsertImportedResources(ctx context.Context, id int, resources []*model.LanguageResourceDTO) (*model.LanguageDTO, error) {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback(tx)

	for _, item := range resources {
		trimmedKey := strings.TrimSpace(item.Key)
		trimmedValue := item.Value

		existing, err := tx.LanguageResource.Query().
			Where(
				languageresource.LanguageIDEQ(id),
				languageresource.KeyEQ(trimmedKey),
			).
			Only(ctx)
		if err == nil {
			if _, err := tx.LanguageResource.UpdateOneID(existing.ID).
				SetValue(trimmedValue).
				Save(ctx); err != nil {
				return nil, err
			}
			continue
		}
		if !generated.IsNotFound(err) {
			return nil, err
		}

		if _, err := tx.LanguageResource.Create().
			SetLanguageID(id).
			SetKey(trimmedKey).
			SetValue(trimmedValue).
			Save(ctx); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *languageRepository) SyncImportedLanguage(ctx context.Context, input model.LanguageDTO) (*model.LanguageDTO, error) {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback(tx)

	normalizedCode := strings.TrimSpace(strings.ToLower(input.Code))
	normalizedName := strings.TrimSpace(input.Name)
	normalizedNativeName := strings.TrimSpace(input.NativeName)

	existing, err := tx.Language.Query().
		Where(language.CodeEQ(normalizedCode), language.Deleted(false)).
		Only(ctx)
	if err != nil && !generated.IsNotFound(err) {
		return nil, err
	}

	var languageID int
	if generated.IsNotFound(err) {
		created, createErr := tx.Language.Create().
			SetCode(normalizedCode).
			SetName(normalizedName).
			SetNativeName(normalizedNativeName).
			SetActive(true).
			SetIsDefault(false).
			Save(ctx)
		if createErr != nil {
			return nil, createErr
		}
		languageID = created.ID
	} else {
		updated, updateErr := tx.Language.UpdateOneID(existing.ID).
			Where(language.Deleted(false)).
			SetName(normalizedName).
			SetNativeName(normalizedNativeName).
			Save(ctx)
		if updateErr != nil {
			return nil, updateErr
		}
		languageID = updated.ID
	}

	if err := syncResources(ctx, tx.Client(), languageID, input.Resources); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, languageID)
}

func clearDefaultIfNeeded(ctx context.Context, db *generated.Client, isDefault bool, exceptID int) error {
	if !isDefault {
		return nil
	}
	query := db.Language.Update().
		Where(language.Deleted(false), language.IsDefault(true))
	if exceptID > 0 {
		query = query.Where(language.IDNEQ(exceptID))
	}
	return query.SetIsDefault(false).Exec(ctx)
}

func syncResources(ctx context.Context, db *generated.Client, languageID int, items []*model.LanguageResourceDTO) error {
	existing, err := db.LanguageResource.Query().
		Where(languageresource.LanguageIDEQ(languageID)).
		All(ctx)
	if err != nil {
		return err
	}

	existingByKey := make(map[string]*generated.LanguageResource, len(existing))
	for _, item := range existing {
		existingByKey[item.Key] = item
	}

	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		key := strings.TrimSpace(item.Key)
		if key == "" {
			continue
		}
		seen[key] = struct{}{}
		if current, ok := existingByKey[key]; ok {
			if _, err := db.LanguageResource.UpdateOneID(current.ID).
				SetValue(item.Value).
				Save(ctx); err != nil {
				return err
			}
			continue
		}

		if _, err := db.LanguageResource.Create().
			SetLanguageID(languageID).
			SetKey(key).
			SetValue(item.Value).
			Save(ctx); err != nil {
			return err
		}
	}

	for _, item := range existing {
		if _, ok := seen[item.Key]; ok {
			continue
		}
		if err := db.LanguageResource.DeleteOneID(item.ID).Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func toLanguageDTO(entity *generated.Language) *model.LanguageDTO {
	res := mapper.MapAs[*generated.Language, *model.LanguageDTO](entity)
	if entity.Edges.Resources != nil {
		resources := mapper.MapListAs[*generated.LanguageResource, *model.LanguageResourceDTO](entity.Edges.Resources)
		res.Resources = resources
	}
	return res
}

func rollback(tx *generated.Tx) {
	_ = tx.Rollback()
}
