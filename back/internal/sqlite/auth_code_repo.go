package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"capuchin/internal/domain"
)

// AuthCodeRepository SQLite реализация репозитория кода аутентификации.
type AuthCodeRepository struct {
	db *sqlx.DB
}

// NewAuthCodeRepository возвращает новый AuthCodeRepository.
func NewAuthCodeRepository(db *sqlx.DB) *AuthCodeRepository {
	return &AuthCodeRepository{
		db: db,
	}
}

// Create создаёт код авторизации.
func (r *AuthCodeRepository) Create(ctx context.Context, authCode domain.AuthCode) error {
	q := `-- auth_codes::Create
		INSERT INTO "auth_codes"
		("id", "user_id", "user_email", "code", "created_at")
		VALUES
		(:id, :user_id, :user_email, :code, :created_at)`

	if _, err := r.db.NamedExecContext(ctx, q, authCodeModelFromDomain(authCode)); err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("insert auth code: %v", err)
	}

	return nil
}

// Delete удаляет код авторизации.
func (r *AuthCodeRepository) Delete(ctx context.Context, id string) error {
	q := `-- auth_codes::Delete
		DELETE FROM "auth_codes"
		WHERE "id" = :id`

	m := map[string]any{
		"id": id,
	}

	if _, err := r.db.NamedExecContext(ctx, q, m); err != nil {
		return fmt.Errorf("delete auth code: %v", err)
	}

	return nil
}

// Find возвращает код авторизации по ID.
func (r *AuthCodeRepository) Find(ctx context.Context, email string, code string) (domain.AuthCode, error) {
	c := authCodeModel{}

	q := `-- auth_codes::Find
		SELECT
			"id",
			"user_id",
			"user_email",
			"code",
			"created_at"
		FROM "auth_codes"
		WHERE "user_email" = ?
			AND "code" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &c, q, email, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.AuthCode{}, domain.ErrNotFound
		}

		return domain.AuthCode{}, fmt.Errorf("select auth code: %v", err)
	}

	return authCodeModelToDomain(c), nil
}

// CancelOld удаляет старые коды авторизации.
func (r *AuthCodeRepository) CancelOld(ctx context.Context, olderThan time.Time) error {
	q := `-- auth_codes::CancelOld
		DELETE FROM "auth_codes"
		WHERE "created_at" < :dt`

	m := map[string]any{
		"dt": timeToUnix(olderThan),
	}

	if _, err := r.db.NamedExecContext(ctx, q, m); err != nil {
		return fmt.Errorf("delete old auth codes: %v", err)
	}

	return nil
}
