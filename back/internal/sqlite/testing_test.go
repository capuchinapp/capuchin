//go:build integration && sqlite

package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

// migrationsDir возвращает абсолютный путь к каталогу миграций.
func migrationsDir() string { //nolint:unused // Используется в тестах
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd, "..", "..", "migrations")
}

// dbConnect возвращает соединение с временным файлом базы данных SQLite
// с применённой схемой из миграций.
func dbConnect(t *testing.T, ctx context.Context) *sqlx.DB { //nolint:unused // Используется в тестах
	t.Helper()

	path := filepath.Join(t.TempDir(), "capuchin.db")

	db, err := NewDatabase(ctx, path)
	assert.NoError(t, err)

	db.SetMaxOpenConns(1)

	assert.NoError(t, goose.SetDialect("sqlite3"))
	err = goose.Up(db.DB, migrationsDir())
	assert.NoError(t, err, "применение миграций")

	t.Cleanup(func() {
		_ = db.Close()
		_ = os.RemoveAll(filepath.Dir(path))
	})

	return db
}
