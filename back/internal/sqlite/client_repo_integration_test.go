//go:build integration && sqlite

package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"capuchin/internal/domain"
)

func TestIntegration_ClientRepository_Create(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JK7GRFMM086BYZWBYJBPS156"
	)

	// Seed

	defer clientsDelete(t, ctx, db, userID)

	// Test

	store := &ClientRepository{db: db}

	// Create
	client := domain.Client{
		ID:           "01JK7GT8HCJ26AK3HK26KWJKCX",
		UserID:       userID,
		Name:         "Client 1",
		BillableRate: 1000,
		Comment:      "Comment 1",
		CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	}

	err := store.Create(ctx, db, client)
	assert.NoError(t, err)

	got, err := clientsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Equal(t, client, got[0])

	// Create fail - already exists
	err = store.Create(ctx, db, domain.Client{
		ID:           "01JK7GWDGYSKAJ2CGS2K601R6M",
		UserID:       userID,
		Name:         "Client 1",
		BillableRate: 1000,
		Comment:      "Comment 1",
		CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	})
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func TestIntegration_ClientRepository_Update(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JK7J444PC9Z7ZPJPKP29887J"
	)

	// Seed

	clients := []domain.Client{
		{
			ID:           "01JK7J4819BTYMXGSV03CNEF6C",
			UserID:       userID,
			Name:         "Client 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           "01JK7J4BS58CEJZ2CWGR13YTV8",
			UserID:       userID,
			Name:         "Client 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
			CreatedAt:    time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
	}

	clientsInsert(t, ctx, db, clients)
	defer clientsDelete(t, ctx, db, userID)

	// Test

	store := &ClientRepository{db: db}

	clients[0].Name = "Client 111"
	clients[0].BillableRate = 10000
	clients[0].Comment = "Comment 111"
	clients[0].UpdatedAt = func() *time.Time { t := time.Date(2025, 2, 3, 9, 15, 25, 0, time.UTC); return &t }()
	clients[0].ArchivedAt = func() *time.Time { t := time.Date(2025, 2, 3, 9, 16, 25, 0, time.UTC); return &t }()

	err := store.Update(ctx, db, clients[0])
	assert.NoError(t, err)

	got, err := clientsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, clients, got)
}

func TestIntegration_ClientRepository_FindAll(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID1 = "01JK7HQYAX3NPFE7YK73EHHCME"
		userID2 = "01JK7HTBS2P9P129CZZXNSAJ0C"
	)

	// Seed

	clients := []domain.Client{
		{
			ID:           "01JK7HQQ2YRDC1SGSDF26S51WT",
			UserID:       userID1,
			Name:         "Client 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           "01JK7HR35PCC855JFW5EFJG201",
			UserID:       userID1,
			Name:         "Client 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
			CreatedAt:    time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
		{
			ID:           "01JK7HTMYQNKWWQVZERZ8ZPCDC",
			UserID:       userID2,
			Name:         "Client 3",
			BillableRate: 3000,
			Comment:      "Comment 3",
			CreatedAt:    time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
	}

	clientsInsert(t, ctx, db, clients)
	defer clientsDelete(t, ctx, db, userID1)
	defer clientsDelete(t, ctx, db, userID2)

	// Test

	store := &ClientRepository{db: db}

	got, err := store.FindAll(ctx, userID1)
	assert.NoError(t, err)
	assert.ElementsMatch(t, clients[:2], got)
}

func TestIntegration_ClientRepository_FindByID(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID    = "01JK7GZ8YWD5WK72F54JF4YNS7"
		clientID1 = "01JK7GZGYWBMX9RNEAYTSN7FFC"
		clientID2 = "01JK7GZMRB1SZNQWGJ6BQFSN3B"
	)

	// Seed

	clients := []domain.Client{
		{
			ID:           clientID1,
			UserID:       userID,
			Name:         "Client 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           clientID2,
			UserID:       userID,
			Name:         "Client 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
			CreatedAt:    time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
	}

	clientsInsert(t, ctx, db, clients)
	defer clientsDelete(t, ctx, db, userID)

	// Test

	store := &ClientRepository{db: db}

	// Found
	got, err := store.FindByID(ctx, userID, clientID2)
	assert.NoError(t, err)
	assert.Equal(t, clients[1], got)

	// Not found
	_, err = store.FindByID(ctx, userID, "01JK7H1YZ9W9YFNV196Q09TQWW")
	assert.Equal(t, domain.ErrNotFound, err)
}

func clientsInsert(t *testing.T, ctx context.Context, db *sqlx.DB, clients []domain.Client) {
	t.Helper()

	q := `INSERT INTO "clients"
		("id", "user_id", "name", "billable_rate", "comment", "created_at", "updated_at", "archived_at")
		VALUES
		(:id, :user_id, :name, :billable_rate, :comment, :created_at, :updated_at, :archived_at)`

	for _, c := range clients {
		_, err := db.NamedExecContext(ctx, q, clientModelFromDomain(c))
		assert.NoError(t, err)
	}
}
