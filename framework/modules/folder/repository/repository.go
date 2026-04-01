package repository

import (
	"context"
	"database/sql"
	"errors"

	frameworkmodel "github.com/khiemnd777/noah_framework/modules/folder/model"
)

var ErrNotFound = errors.New("folder not found")

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, input frameworkmodel.Folder) (*frameworkmodel.Folder, error) {
	const query = `
		INSERT INTO folders (user_id, folder_name, color, shared, parent_id, deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, false, NOW(), NOW())
		RETURNING id, user_id, folder_name, color, shared, parent_id, deleted, created_at, updated_at
	`
	row := r.db.QueryRowContext(ctx, query, input.UserID, input.FolderName, input.Color, input.Shared, input.ParentID)
	return scanFolderRow(row)
}

func (r *Repository) FindByID(ctx context.Context, id int) (*frameworkmodel.Folder, error) {
	const query = `
		SELECT id, user_id, folder_name, color, shared, parent_id, deleted, created_at, updated_at
		FROM folders
		WHERE id = $1 AND deleted = false
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanFolderRow(row)
}

func (r *Repository) ListByUser(ctx context.Context, userID int) ([]*frameworkmodel.Folder, error) {
	const query = `
		SELECT id, user_id, folder_name, color, shared, parent_id, deleted, created_at, updated_at
		FROM folders
		WHERE user_id = $1 AND deleted = false
		ORDER BY updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFolders(rows)
}

func (r *Repository) ListByUserPaginated(ctx context.Context, userID, limit, offset int) ([]*frameworkmodel.Folder, bool, error) {
	const query = `
		SELECT id, user_id, folder_name, color, shared, parent_id, deleted, created_at, updated_at
		FROM folders
		WHERE user_id = $1 AND deleted = false
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit+1, offset)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	items, err := scanFolders(rows)
	if err != nil {
		return nil, false, err
	}
	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	return items, hasMore, nil
}

func (r *Repository) Update(ctx context.Context, id int, input frameworkmodel.Folder) (*frameworkmodel.Folder, error) {
	const query = `
		UPDATE folders
		SET folder_name = $2, color = $3, shared = $4, parent_id = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, folder_name, color, shared, parent_id, deleted, created_at, updated_at
	`
	row := r.db.QueryRowContext(ctx, query, id, input.FolderName, input.Color, input.Shared, input.ParentID)
	return scanFolderRow(row)
}

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	const query = `UPDATE folders SET deleted = true, updated_at = NOW() WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

type folderScanner interface {
	Scan(dest ...any) error
}

func scanFolderRow(row folderScanner) (*frameworkmodel.Folder, error) {
	var folder frameworkmodel.Folder
	var parentID sql.NullInt64
	if err := row.Scan(
		&folder.ID,
		&folder.UserID,
		&folder.FolderName,
		&folder.Color,
		&folder.Shared,
		&parentID,
		&folder.Deleted,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if parentID.Valid {
		value := int(parentID.Int64)
		folder.ParentID = &value
	}
	return &folder, nil
}

func scanFolders(rows *sql.Rows) ([]*frameworkmodel.Folder, error) {
	var items []*frameworkmodel.Folder
	for rows.Next() {
		item, err := scanFolderRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
