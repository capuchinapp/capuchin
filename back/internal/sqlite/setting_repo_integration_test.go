//go:build integration && sqlite

package sqlite

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"capuchin/internal/domain"
)

func TestIntegration_SettingRepository_InsertOrUpdate(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JJNRSVZBWKYG2XD4W04P0ZJC"
	)

	// Seed

	defer settingsDelete(t, ctx, db, userID)

	// Test

	store := &SettingRepository{db: db}

	settings := []domain.Setting{
		{
			UserID: userID,
			Key:    "test_key_1",
			Value:  "test_value_1",
		},
		{
			UserID: userID,
			Key:    "test_key_2",
			Value:  "test_value_2",
		},
		{
			UserID: userID,
			Key:    "test_key_1",
			Value:  "test_value_11",
		},
	}

	for _, setting := range settings {
		err := store.InsertOrUpdate(ctx, setting)
		assert.NoError(t, err)
	}

	got := settingsSelect(t, ctx, db, userID)
	assert.ElementsMatch(t, []domain.Setting{
		{
			UserID: userID,
			Key:    "test_key_1",
			Value:  "test_value_11",
		},
		{
			UserID: userID,
			Key:    "test_key_2",
			Value:  "test_value_2",
		},
	}, got)
}

func TestIntegration_SettingRepository_FindAll(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JJNMGM04S465SRBHBR82J989"
	)

	// Seed

	setting := domain.Setting{
		UserID: userID,
		Key:    "test_key",
		Value:  "test_value",
	}

	settingsInsert(t, ctx, db, []domain.Setting{setting})
	defer settingsDelete(t, ctx, db, userID)

	// Test

	store := &SettingRepository{db: db}

	got, err := store.FindAll(ctx, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []domain.Setting{setting}, got)
}

func TestIntegration_SettingRepository_FindByKey(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JJNRM7AM1RN6ND1C4W9GDG10"
	)

	// Seed

	settings := []domain.Setting{
		{
			UserID: userID,
			Key:    "test_key_1",
			Value:  "test_value_1",
		},
		{
			UserID: userID,
			Key:    "test_key_2",
			Value:  "test_value_2",
		},
	}

	settingsInsert(t, ctx, db, settings)
	defer settingsDelete(t, ctx, db, userID)

	// Test

	store := &SettingRepository{db: db}

	got, err := store.FindByKey(ctx, userID, "test_key_1")
	assert.NoError(t, err)
	assert.Equal(t, settings[0], got)
}

func settingsInsert(t *testing.T, ctx context.Context, db *sqlx.DB, settings []domain.Setting) {
	t.Helper()

	q := `INSERT INTO "settings"
		("user_id", "key", "value")
		VALUES
		(:user_id, :key, :value)`

	for _, s := range settings {
		_, err := db.NamedExecContext(ctx, q, settingModelFromDomain(s))
		assert.NoError(t, err)
	}
}

func settingsSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) []domain.Setting {
	t.Helper()

	ms := []settingModel{}

	q := `SELECT
			"user_id",
			"key",
			"value"
		FROM "settings"
		WHERE "user_id" = ?`

	if err := db.SelectContext(ctx, &ms, q, userID); err != nil {
		assert.NoError(t, err)

		return nil
	}

	return settingModelsToDomains(ms)
}

func settingsDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "settings" WHERE "user_id" = :user_id`

	p := map[string]any{
		"user_id": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}
