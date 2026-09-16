//go:build integration && sqlite

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"capuchin/internal/domain"
)

func auditLogSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) (auditLogModel, error) {
	t.Helper()

	m := auditLogModel{}

	q := `SELECT
			"action_at",
			"action_type",
			"user_id",
			"table_name",
			"record_id",
			"old_value",
			"new_value"
		FROM "audit_log"
		WHERE "user_id" = ?
		LIMIT 1`

	if err := db.GetContext(ctx, &m, q, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auditLogModel{}, domain.ErrNotFound
		}

		return auditLogModel{}, fmt.Errorf("select audit log: %v", err)
	}

	return m, nil
}

func auditLogsDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "audit_log" WHERE "user_id" = :user_id`

	p := map[string]any{
		"user_id": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}

func clientsSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) ([]domain.Client, error) {
	t.Helper()

	ms := []clientModel{}

	q := `SELECT
			"id",
			"user_id",
			"name",
			"billable_rate",
			"comment",
			"created_at",
			"updated_at",
			"archived_at"
		FROM "clients"
		WHERE "user_id" = ?`

	if err := db.SelectContext(ctx, &ms, q, userID); err != nil {
		return nil, fmt.Errorf("select clients by id: %v", err)
	}

	return clientModelsToDomains(ms), nil
}

func clientsDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "clients" WHERE "user_id" = :user_id`

	p := map[string]any{
		"user_id": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}