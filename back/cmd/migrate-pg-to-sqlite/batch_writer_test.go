//go:build migration

package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openSqlite(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", "file::memory:?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	db.SetMaxOpenConns(1)

	_, err = db.ExecContext(context.Background(), `CREATE TABLE users (id TEXT PRIMARY KEY, email TEXT NOT NULL, created_at INTEGER NOT NULL)`)
	require.NoError(t, err)

	return db
}

func TestVerifyEmpty(t *testing.T) {
	tests := []struct {
		name     string
		insert   string
		expError bool
	}{
		{name: "empty tables pass", insert: ""},
		{name: "goose_db_version ignored", insert: `CREATE TABLE goose_db_version (id INTEGER NOT NULL PRIMARY KEY)`},
		{
			name:     "non-empty app table fails",
			insert:   `INSERT INTO users (id, email, created_at) VALUES ('u1', 'a@b.c', 1)`,
			expError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := openSqlite(t)
			ctx := context.Background()

			if tt.insert != "" {
				_, err := db.ExecContext(ctx, tt.insert)
				require.NoError(t, err)
			}

			err := verifyEmpty(ctx, db)
			if tt.expError {
				assert.ErrorContains(t, err, "not empty")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBatchWriter(t *testing.T) {
	db := openSqlite(t)
	ctx := context.Background()

	bw, err := newBatchWriter(ctx, db, "users", []string{"id", "email", "created_at"})
	require.NoError(t, err)

	rows := batchSize*2 + 7 // several full batches plus a remainder

	for i := 0; i < rows; i++ {
		err := bw.add(
			fmt.Sprintf("id-%06d", i),
			"user@example.com",
			int64(i),
		)
		require.NoError(t, err)
	}

	require.NoError(t, bw.close())
	assert.Equal(t, int64(rows), bw.rowsCount)
	assert.Equal(t, 3, bw.batches) // 2 full batches + 1 remainder

	var count int64

	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(rows), count)
}

func TestBatchWriterRollbackOnError(t *testing.T) {
	db := openSqlite(t)
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `CREATE TABLE users2 (id TEXT PRIMARY KEY, email TEXT NOT NULL)`)
	require.NoError(t, err)

	bw, err := newBatchWriter(ctx, db, "users2", []string{"id", "email"})
	require.NoError(t, err)

	err = bw.add("ok", "valid@example.com")
	require.NoError(t, err)

	// Дубликат id спровоцирует ошибку на flush при close()
	err = bw.add("ok", "dup@example.com")
	require.NoError(t, err)

	err = bw.close()
	assert.Error(t, err)

	var count int64

	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users2`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}