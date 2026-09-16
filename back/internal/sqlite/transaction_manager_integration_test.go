//go:build integration && sqlite

package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"capuchin/internal/domain"
)

func TestIntegration_TransactionManager_ExecuteInTransaction(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JNNS4YMMZ7RH4KTBJPJ4VPFN"
	)

	// Seed

	defer clientsDelete(t, ctx, db, userID)
	defer auditLogsDelete(t, ctx, db, userID)

	// Test

	clientStore := &ClientRepository{db: db}
	auditLogStore := &AuditLogRepository{}
	txManager := &TransactionManager{db: db}

	client := domain.Client{
		ID:           "01JK7GT8HCJ26AK3HK26KWJKCX",
		UserID:       userID,
		Name:         "Client 1",
		BillableRate: 1000,
		Comment:      "Comment 1",
		CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	}

	auditLog := domain.AuditLog{
		ActionAt:   time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		ActionType: domain.AuditLogActionTypeCreate,
		UserID:     userID,
		TableName:  domain.AuditLogTableNameClients,
		RecordID:   "01JN5FHFW67ZS5TRKQTE47Z5X0",
		OldValue:   nil,
		NewValue:   nil,
	}

	// Should rollback
	err := txManager.ExecuteInTransaction(ctx,
		func(ctx context.Context, tx SQLExecutor) error {
			return domain.ErrAlreadyExists
		},
		func(ctx context.Context, tx SQLExecutor) error {
			return auditLogStore.Create(ctx, tx, auditLog)
		})
	assert.ErrorIs(t, err, domain.ErrAlreadyExists)

	_, err = auditLogSelect(t, ctx, db, userID)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	// Should success
	err = txManager.ExecuteInTransaction(ctx,
		func(ctx context.Context, tx SQLExecutor) error {
			return clientStore.Create(ctx, tx, client)
		},
		func(ctx context.Context, tx SQLExecutor) error {
			return auditLogStore.Create(ctx, tx, auditLog)
		})
	assert.NoError(t, err)

	gotClients, err := clientsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	gotClients[0].CreatedAt = gotClients[0].CreatedAt.UTC()
	assert.Equal(t, client, gotClients[0])

	gotAuditLog, err := auditLogSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Equal(t, auditLogModelFromDomain(auditLog), gotAuditLog)
}
