package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	frameworkmodel "github.com/khiemnd777/noah_framework/modules/attribute/model"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

type Repository interface {
	EnsureSchema(ctx context.Context) error
	CreateAttribute(ctx context.Context, data frameworkmodel.Attribute) (*frameworkmodel.Attribute, error)
	GetAttributeByID(ctx context.Context, id int) (*frameworkmodel.Attribute, error)
	ListAttributesByUser(ctx context.Context, userID int) ([]*frameworkmodel.Attribute, error)
	ListAttributesByUserPaginated(ctx context.Context, userID, limit, offset int) ([]*frameworkmodel.Attribute, bool, error)
	UpdateAttribute(ctx context.Context, id int, name, attributeType string) (*frameworkmodel.Attribute, error)
	SoftDeleteAttribute(ctx context.Context, id int) error
	CreateOption(ctx context.Context, option frameworkmodel.AttributeOption) (*frameworkmodel.AttributeOption, error)
	UpdateOption(ctx context.Context, id int, value string, order int) (*frameworkmodel.AttributeOption, error)
	SoftDeleteOption(ctx context.Context, id int) error
	ListOptionsByAttribute(ctx context.Context, attributeID int) ([]*frameworkmodel.AttributeOption, error)
	GetOptionByID(ctx context.Context, id int) (*frameworkmodel.AttributeOption, error)
	BatchUpdateDisplayOrder(ctx context.Context, orders []frameworkmodel.OptionOrder) error
}

type SQLRepository struct {
	db *sql.DB
}

func New(client frameworkdb.Client) (*SQLRepository, error) {
	sqlDB, err := frameworkruntime.SQLDB(client)
	if err != nil {
		return nil, err
	}
	return &SQLRepository{db: sqlDB}, nil
}

func (r *SQLRepository) EnsureSchema(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS attributes (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			attribute_name TEXT NOT NULL,
			attribute_type TEXT NOT NULL DEFAULT 'text',
			deleted BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_attributes_user_deleted_updated
			ON attributes (user_id, deleted, updated_at DESC)`,
		`CREATE TABLE IF NOT EXISTS attribute_options (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			attribute_id BIGINT NOT NULL,
			option_value TEXT NOT NULL,
			display_order BIGINT NOT NULL DEFAULT 0,
			deleted BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_attribute_options_attribute_deleted_order
			ON attribute_options (attribute_id, deleted, display_order ASC)`,
	}

	for _, stmt := range statements {
		if _, err := r.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}

func (r *SQLRepository) CreateAttribute(ctx context.Context, data frameworkmodel.Attribute) (*frameworkmodel.Attribute, error) {
	row := r.db.QueryRowContext(
		ctx,
		`INSERT INTO attributes (user_id, attribute_name, attribute_type, deleted)
		 VALUES ($1, $2, $3, FALSE)
		 RETURNING id, user_id, attribute_name, attribute_type, deleted, created_at, updated_at`,
		data.UserID,
		data.AttributeName,
		data.AttributeType,
	)
	return scanAttribute(row)
}

func (r *SQLRepository) GetAttributeByID(ctx context.Context, id int) (*frameworkmodel.Attribute, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, attribute_name, attribute_type, deleted, created_at, updated_at
		 FROM attributes
		 WHERE id = $1`,
		id,
	)
	return scanAttribute(row)
}

func (r *SQLRepository) ListAttributesByUser(ctx context.Context, userID int) ([]*frameworkmodel.Attribute, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, attribute_name, attribute_type, deleted, created_at, updated_at
		 FROM attributes
		 WHERE user_id = $1 AND deleted = FALSE
		 ORDER BY updated_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collectAttributes(rows)
}

func (r *SQLRepository) ListAttributesByUserPaginated(ctx context.Context, userID, limit, offset int) ([]*frameworkmodel.Attribute, bool, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, attribute_name, attribute_type, deleted, created_at, updated_at
		 FROM attributes
		 WHERE user_id = $1 AND deleted = FALSE
		 ORDER BY updated_at DESC
		 LIMIT $2 OFFSET $3`,
		userID,
		limit+1,
		offset,
	)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	items, err := collectAttributes(rows)
	if err != nil {
		return nil, false, err
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	return items, hasMore, nil
}

func (r *SQLRepository) UpdateAttribute(ctx context.Context, id int, name, attributeType string) (*frameworkmodel.Attribute, error) {
	row := r.db.QueryRowContext(
		ctx,
		`UPDATE attributes
		 SET attribute_name = $2, attribute_type = $3, updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, user_id, attribute_name, attribute_type, deleted, created_at, updated_at`,
		id,
		name,
		attributeType,
	)
	return scanAttribute(row)
}

func (r *SQLRepository) SoftDeleteAttribute(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE attributes
		 SET deleted = TRUE, updated_at = NOW()
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	return ensureAffected(result)
}

func (r *SQLRepository) CreateOption(ctx context.Context, option frameworkmodel.AttributeOption) (*frameworkmodel.AttributeOption, error) {
	row := r.db.QueryRowContext(
		ctx,
		`INSERT INTO attribute_options (user_id, attribute_id, option_value, display_order, deleted)
		 VALUES ($1, $2, $3, $4, FALSE)
		 RETURNING id, user_id, attribute_id, option_value, display_order, deleted, created_at`,
		option.UserID,
		option.AttributeID,
		option.OptionValue,
		option.DisplayOrder,
	)
	return scanOption(row)
}

func (r *SQLRepository) UpdateOption(ctx context.Context, id int, value string, order int) (*frameworkmodel.AttributeOption, error) {
	row := r.db.QueryRowContext(
		ctx,
		`UPDATE attribute_options
		 SET option_value = $2, display_order = $3
		 WHERE id = $1
		 RETURNING id, user_id, attribute_id, option_value, display_order, deleted, created_at`,
		id,
		value,
		order,
	)
	return scanOption(row)
}

func (r *SQLRepository) SoftDeleteOption(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE attribute_options
		 SET deleted = TRUE
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	return ensureAffected(result)
}

func (r *SQLRepository) ListOptionsByAttribute(ctx context.Context, attributeID int) ([]*frameworkmodel.AttributeOption, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, attribute_id, option_value, display_order, deleted, created_at
		 FROM attribute_options
		 WHERE attribute_id = $1 AND deleted = FALSE
		 ORDER BY display_order ASC`,
		attributeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collectOptions(rows)
}

func (r *SQLRepository) GetOptionByID(ctx context.Context, id int) (*frameworkmodel.AttributeOption, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, attribute_id, option_value, display_order, deleted, created_at
		 FROM attribute_options
		 WHERE id = $1`,
		id,
	)
	return scanOption(row)
}

func (r *SQLRepository) BatchUpdateDisplayOrder(ctx context.Context, orders []frameworkmodel.OptionOrder) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if _, err := tx.ExecContext(
			ctx,
			`UPDATE attribute_options
			 SET display_order = $2
			 WHERE id = $1`,
			order.OptionID,
			order.DisplayOrder,
		); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

type attributeScanner interface {
	Scan(dest ...any) error
}

func scanAttribute(scanner attributeScanner) (*frameworkmodel.Attribute, error) {
	var item frameworkmodel.Attribute
	if err := scanner.Scan(
		&item.ID,
		&item.UserID,
		&item.AttributeName,
		&item.AttributeType,
		&item.Deleted,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &item, nil
}

func scanOption(scanner attributeScanner) (*frameworkmodel.AttributeOption, error) {
	var item frameworkmodel.AttributeOption
	if err := scanner.Scan(
		&item.ID,
		&item.UserID,
		&item.AttributeID,
		&item.OptionValue,
		&item.DisplayOrder,
		&item.Deleted,
		&item.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &item, nil
}

func collectAttributes(rows *sql.Rows) ([]*frameworkmodel.Attribute, error) {
	items := make([]*frameworkmodel.Attribute, 0)
	for rows.Next() {
		item, err := scanAttribute(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func collectOptions(rows *sql.Rows) ([]*frameworkmodel.AttributeOption, error) {
	items := make([]*frameworkmodel.AttributeOption, 0)
	for rows.Next() {
		item, err := scanOption(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func ensureAffected(result sql.Result) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func WrapRepositoryError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", operation, err)
}
