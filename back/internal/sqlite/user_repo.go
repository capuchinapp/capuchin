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

// UserRepository SQLite реализация репозитория пользователя.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository возвращает новый UserRepository.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create создаёт пользователя.
func (*UserRepository) Create(ctx context.Context, db SQLExecutor, user domain.User) error {
	q := `-- users::Create
		INSERT INTO "users"
		("id", "email", "activate_code", "created_at")
		VALUES
		(:id, :email, :activate_code, :created_at)`

	if _, err := db.NamedExecContext(ctx, q, userModelFromDomain(user)); err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("insert user: %v", err)
	}

	return nil
}

// CancelActivateCode отменяет код активации.
func (r *UserRepository) CancelActivateCode(ctx context.Context, user domain.User) error {
	q := `-- users::CancelActivateCode
		UPDATE "users" SET
			"activate_code" = NULL
		WHERE "id" = :id`

	if _, err := r.db.NamedExecContext(ctx, q, userModelFromDomain(user)); err != nil {
		return fmt.Errorf("update checked_at: %v", err)
	}

	return nil
}

// DeleteNotActivated удаляет неактивированных пользователей.
func (r *UserRepository) DeleteNotActivated(ctx context.Context, olderThan time.Time) error {
	q := `-- users::DeleteNotActivated
		DELETE FROM "users"
		WHERE "created_at" < :dt
			AND "activate_code" IS NOT NULL`

	m := map[string]any{
		"dt": timeToUnix(olderThan),
	}

	if _, err := r.db.NamedExecContext(ctx, q, m); err != nil {
		return fmt.Errorf("delete users not activated: %v", err)
	}

	return nil
}

// FindByID возвращает пользователя по ID.
func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	u := userModel{}

	q := `-- users::FindByID
		SELECT
			"id",
			"email",
			"activate_code",
			"created_at",
			"updated_at"
		FROM "users"
		WHERE "id" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &u, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("select user by id: %v", err)
	}

	return userModelToDomain(u), nil
}

// FindByEmail возвращает пользователя по email.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u := userModel{}

	q := `-- users::FindByEmail
		SELECT
			"id",
			"email",
			"activate_code",
			"created_at",
			"updated_at"
		FROM "users"
		WHERE "email" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &u, q, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("select user by email: %v", err)
	}

	return userModelToDomain(u), nil
}
