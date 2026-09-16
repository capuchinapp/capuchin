package sqlite

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// uniqueViolationError возвращает реальную ошибку SQLite о нарушении
// уникального ограничения.
func uniqueViolationError(t *testing.T, ctx context.Context) error {
	t.Helper()

	db, err := NewDatabase(ctx, filepath.Join(t.TempDir(), "db.sqlite"))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	_, err = db.ExecContext(ctx, `CREATE TABLE t (id TEXT PRIMARY KEY, email TEXT UNIQUE)`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `INSERT INTO t (id, email) VALUES ('1', 'a@b.c')`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `INSERT INTO t (id, email) VALUES ('2', 'a@b.c')`)
	require.Error(t, err)

	return fmt.Errorf("insert duplicate: %w", err)
}

func TestIsUniqueViolation(t *testing.T) {
	ctx := context.Background()

	uniqErr := uniqueViolationError(t, ctx)

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "нарушение уникальности SQLite",
			err:  uniqErr,
			want: true,
		},
		{
			name: "обёрнутая ошибка нарушения уникальности",
			err:  fmt.Errorf("insert: %w", uniqErr),
			want: true,
		},
		{
			name: "посторонняя ошибка",
			err:  errors.New("unexpected error"),
			want: false,
		},
		{
			name: "nil ошибка",
			err:  nil,
			want: false,
		},
		{
			name: "ошибка с другим кодом SQLite",
			err:  errors.New("constraint failed: NOT NULL constraint failed"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isUniqueViolation(tt.err))
		})
	}
}
