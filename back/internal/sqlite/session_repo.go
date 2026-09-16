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

// SessionRepository SQLite реализация репозитория сессий.
type SessionRepository struct {
	db *sqlx.DB
}

// NewSessionRepository возвращает новый SessionRepository.
func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

// Create создаёт сессию.
func (r *SessionRepository) Create(ctx context.Context, session domain.Session) error {
	q := `-- sessions::Create
		INSERT INTO "sessions"
		("id", "cid", "uid", "checked_at")
		VALUES
		(:id, :cid, :uid, :checked_at)`

	if _, err := r.db.NamedExecContext(ctx, q, sessionModelFromDomain(session)); err != nil {
		return fmt.Errorf("insert session: %v", err)
	}

	return nil
}

// DeleteByCID удаляет сессию по CID.
func (r *SessionRepository) DeleteByCID(ctx context.Context, cid string) error {
	q := `-- sessions::DeleteByCID
		DELETE FROM "sessions"
		WHERE "cid" = :cid`

	m := map[string]any{
		"cid": cid,
	}

	if _, err := r.db.NamedExecContext(ctx, q, m); err != nil {
		return fmt.Errorf("delete session by cid: %v", err)
	}

	return nil
}

// FindAll возвращает список сессий.
func (r *SessionRepository) FindAll(ctx context.Context, userID string) ([]domain.Session, error) {
	sl := []sessionModel{}

	q := `-- sessions::FindAll
		SELECT
			"id",
			"cid",
			"uid",
			"checked_at"
		FROM "sessions"
		WHERE "uid" = ?
		ORDER BY "checked_at" DESC`

	if err := r.db.SelectContext(ctx, &sl, q, userID); err != nil {
		return nil, fmt.Errorf("select all sessions: %v", err)
	}

	return sessionModelsToDomains(sl), nil
}

// FindByCID возвращает сессию по CID.
func (r *SessionRepository) FindByCID(ctx context.Context, cid string) (domain.Session, error) {
	s := sessionModel{}

	q := `-- sessions::FindByCID
		SELECT
			"id",
			"cid",
			"uid",
			"checked_at"
		FROM "sessions"
		WHERE "cid" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &s, q, cid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Session{}, domain.ErrNotFound
		}

		return domain.Session{}, fmt.Errorf("select session by cid: %v", err)
	}

	return sessionModelToDomain(s), nil
}

// FindByID возвращает сессию по ID.
func (r *SessionRepository) FindByID(ctx context.Context, userID string, id string) (domain.Session, error) {
	s := sessionModel{}

	q := `-- sessions::FindByID
		SELECT
			"id",
			"cid",
			"uid",
			"checked_at"
		FROM "sessions"
		WHERE "uid" = ?
			AND "id" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &s, q, userID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Session{}, domain.ErrNotFound
		}

		return domain.Session{}, fmt.Errorf("select session by id: %v", err)
	}

	return sessionModelToDomain(s), nil
}

// UpdateCheckedAt обновляет поле "checked_at" сессии.
func (r *SessionRepository) UpdateCheckedAt(ctx context.Context, sessionID string, checkedAt time.Time) error {
	q := `-- sessions::UpdateCheckedAt
		UPDATE "sessions" SET
			"checked_at" = :checked_at
		WHERE "id" = :id`

	m := map[string]any{
		"id":         sessionID,
		"checked_at": timeToUnix(checkedAt),
	}

	if _, err := r.db.NamedExecContext(ctx, q, m); err != nil {
		return fmt.Errorf("update checked_at: %v", err)
	}

	return nil
}

// CancelOld удаляет старые сессии.
func (r *SessionRepository) CancelOld(ctx context.Context, olderThan time.Time) error {
	q := `-- sessions::CancelOld
		DELETE FROM "sessions"
		WHERE "checked_at" < :dt`

	m := map[string]any{
		"dt": timeToUnix(olderThan),
	}

	if _, err := r.db.NamedExecContext(ctx, q, m); err != nil {
		return fmt.Errorf("delete old sessions: %v", err)
	}

	return nil
}
