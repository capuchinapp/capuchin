package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"capuchin/internal/domain"
)

// SettingRepository SQLite реализация репозитория настроек.
type SettingRepository struct {
	db *sqlx.DB
}

// NewSettingRepository возвращает новый SettingRepository.
func NewSettingRepository(db *sqlx.DB) *SettingRepository {
	return &SettingRepository{
		db: db,
	}
}

// InsertOrUpdate создаёт или обновляет настройку.
func (r *SettingRepository) InsertOrUpdate(ctx context.Context, setting domain.Setting) error {
	q := `-- settings::InsertOrUpdate
		INSERT INTO "settings"
		("user_id", "key", "value")
		VALUES
		(:user_id, :key, :value)
		ON CONFLICT ("user_id", "key")
		DO UPDATE SET
			value = EXCLUDED.value;`

	if _, err := r.db.NamedExecContext(ctx, q, settingModelFromDomain(setting)); err != nil {
		return fmt.Errorf("update setting: %v", err)
	}

	return nil
}

// FindAll возвращает список настроек.
func (r *SettingRepository) FindAll(ctx context.Context, userID string) ([]domain.Setting, error) {
	sl := []settingModel{}

	q := `-- settings::FindAll
		SELECT
			"user_id",
			"key",
			"value"
		FROM "settings"
		WHERE "user_id" = ?`

	if err := r.db.SelectContext(ctx, &sl, q, userID); err != nil {
		return nil, fmt.Errorf("select all settings: %v", err)
	}

	return settingModelsToDomains(sl), nil
}

// FindByKey возвращает настройку по ключу.
func (r *SettingRepository) FindByKey(ctx context.Context, userID string, key string) (domain.Setting, error) {
	s := settingModel{}

	q := `-- settings::FindByKey
		SELECT
			"user_id",
			"key",
			"value"
		FROM "settings"
		WHERE "user_id" = ?
			AND "key" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &s, q, userID, key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Setting{}, domain.ErrNotFound
		}

		return domain.Setting{}, fmt.Errorf("select setting by key: %v", err)
	}

	return settingModelToDomain(s), nil
}
