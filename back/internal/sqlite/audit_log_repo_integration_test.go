//go:build integration && sqlite

package sqlite

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"capuchin/internal/domain"
)

func TestIntegration_AuditLogRepository_Create(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JN5FG4N6WW19N5G7V03QK638"
	)

	// Seed

	defer auditLogsDelete(t, ctx, db, userID)

	// Test

	store := &AuditLogRepository{}

	auditLog := domain.AuditLog{
		ActionAt:   time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		ActionType: domain.AuditLogActionTypeCreate,
		UserID:     userID,
		TableName:  domain.AuditLogTableNameClients,
		RecordID:   "01JN5FHFW67ZS5TRKQTE47Z5X0",
		OldValue:   nil,
		NewValue:   nil,
	}

	err := store.Create(ctx, db, auditLog)
	assert.NoError(t, err)

	got, err := auditLogSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Equal(t, auditLogModelFromDomain(auditLog), got)
}
