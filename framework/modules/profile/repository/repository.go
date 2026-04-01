package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	frameworkmodel "github.com/khiemnd777/noah_framework/modules/profile/model"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id int) (*frameworkmodel.User, error) {
	const query = `
		SELECT id, email, password, name, phone, active, deleted_at, avatar, provider, provider_id, ref_code, qr_code, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanUser(row)
}

func (r *Repository) UpdateByID(ctx context.Context, id int, name string, phone, email *string, avatar *string, refCode *string, qrCode *string) (*frameworkmodel.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	current, err := scanUser(tx.QueryRowContext(ctx, `
		SELECT id, email, password, name, phone, active, deleted_at, avatar, provider, provider_id, ref_code, qr_code, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id))
	if err != nil {
		return nil, err
	}

	if current.QrCode == nil && refCode != nil && qrCode != nil {
		current.RefCode = refCode
		current.QrCode = qrCode
	}
	current.Name = name
	current.Phone = stringValue(phone)
	current.Email = stringValue(email)
	current.Avatar = stringValue(avatar)

	row := tx.QueryRowContext(ctx, `
		UPDATE users
		SET name = $2, phone = $3, email = $4, avatar = $5, ref_code = $6, qr_code = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, password, name, phone, active, deleted_at, avatar, provider, provider_id, ref_code, qr_code, created_at, updated_at
	`, id, current.Name, nullableString(phone), nullableString(email), nullableString(avatar), nullableString(current.RefCode), nullableString(current.QrCode))
	updated, err := scanUser(row)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *Repository) ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error {
	var currentHash string
	if err := r.db.QueryRowContext(ctx, `SELECT password FROM users WHERE id = $1`, id).Scan(&currentHash); err != nil {
		return fmt.Errorf("failed to get current password hash: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(currentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `UPDATE users SET password = $2, updated_at = NOW() WHERE id = $1`, id, string(newHash))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func (r *Repository) CheckPhoneExists(ctx context.Context, userID int, phone string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE id <> $1 AND phone = $2 AND deleted_at IS NULL
		)
	`, userID, phone).Scan(&exists)
	return exists, err
}

func (r *Repository) CheckEmailExists(ctx context.Context, userID int, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE id <> $1 AND email = $2 AND deleted_at IS NULL
		)
	`, userID, email).Scan(&exists)
	return exists, err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `UPDATE users SET active = false, deleted_at = $2, updated_at = $2 WHERE id = $1`, id, time.Now())
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

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(row userScanner) (*frameworkmodel.User, error) {
	var user frameworkmodel.User
	var deletedAt sql.NullTime
	var phone sql.NullString
	var email sql.NullString
	var avatar sql.NullString
	var provider sql.NullString
	var providerID sql.NullString
	var refCode sql.NullString
	var qrCode sql.NullString

	if err := row.Scan(
		&user.ID,
		&email,
		&user.Password,
		&user.Name,
		&phone,
		&user.Active,
		&deletedAt,
		&avatar,
		&provider,
		&providerID,
		&refCode,
		&qrCode,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	user.Email = email.String
	user.Phone = phone.String
	user.Avatar = avatar.String
	user.Provider = provider.String
	user.ProviderID = providerID.String
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}
	if refCode.Valid {
		value := refCode.String
		user.RefCode = &value
	}
	if qrCode.Valid {
		value := qrCode.String
		user.QrCode = &value
	}
	return &user, nil
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}
	if *value == "" {
		return nil
	}
	return *value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
