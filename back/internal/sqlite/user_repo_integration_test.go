//go:build integration && sqlite

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"capuchin/internal/domain"
)

func TestIntegration_UserRepository_Create(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JK602X0JH004N6WGAEPTS1QG"
	)

	// Seed

	defer userDelete(t, ctx, db, userID)

	// Test

	store := &UserRepository{db: db}

	// Create
	user := domain.User{
		ID:        userID,
		Email:     "user@localhost.ru",
		CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	}

	err := store.Create(ctx, db, user)
	assert.NoError(t, err)

	got, err := userSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Equal(t, user, got)

	// Create fail - already exists
	err = store.Create(ctx, db, domain.User{
		ID:        userID,
		Email:     "user@localhost.ru",
		CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	})
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func TestIntegration_UserRepository_CancelActivateCode(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID1 = "01JK60KP9J8Z38ENSQEKD63JGZ"
		userID2 = "01JK60KT0XD58NKXR11R9T35NB"
	)

	// Seed

	users := []domain.User{
		{
			ID:           userID1,
			Email:        "user1@localhost.ru",
			ActivateCode: func() *string { s := "1111111111111111111111111111111111111111"; return &s }(),
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           userID2,
			Email:        "user2@localhost.ru",
			ActivateCode: func() *string { s := "2222222222222222222222222222222222222222"; return &s }(),
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	userInsert(t, ctx, db, users[0])
	userInsert(t, ctx, db, users[1])
	defer userDelete(t, ctx, db, userID1)
	defer userDelete(t, ctx, db, userID2)

	// Test

	store := &UserRepository{db: db}

	err := store.CancelActivateCode(ctx, users[0])
	assert.NoError(t, err)

	users[0].ActivateCode = nil
	got1, err := userSelect(t, ctx, db, userID1)
	assert.NoError(t, err)
	assert.Equal(t, users[0], got1)

	got2, err := userSelect(t, ctx, db, userID2)
	assert.NoError(t, err)
	assert.Equal(t, users[1], got2)
}

func TestIntegration_UserRepository_DeleteNotActivated(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID1 = "01JK61HHNJP0D7SJ0KZ76981Y7"
		userID2 = "01JK61HN7VX0JSPAETC5D8F8FR"
	)

	// Seed

	users := []domain.User{
		{
			ID:           userID1,
			Email:        "user1@localhost.ru",
			ActivateCode: func() *string { s := "1111111111111111111111111111111111111111"; return &s }(),
			CreatedAt:    time.Date(2025, 2, 3, 9, 10, 25, 0, time.UTC),
		},
		{
			ID:           userID2,
			Email:        "user2@localhost.ru",
			ActivateCode: func() *string { s := "2222222222222222222222222222222222222222"; return &s }(),
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	userInsert(t, ctx, db, users[0])
	userInsert(t, ctx, db, users[1])
	defer userDelete(t, ctx, db, userID1)
	defer userDelete(t, ctx, db, userID2)

	// Test

	store := &UserRepository{db: db}

	err := store.DeleteNotActivated(ctx, time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC))
	assert.NoError(t, err)

	_, err = userSelect(t, ctx, db, userID1)
	assert.Equal(t, err, domain.ErrNotFound)

	got, err := userSelect(t, ctx, db, userID2)
	assert.NoError(t, err)
	assert.Equal(t, users[1], got)
}

func TestIntegration_UserRepository_FindByID(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID1 = "01JK60D6RP4KCS1ZRDFSA3QCM7"
		userID2 = "01JK60DB4GP72EZRVSM9CK5F8W"
	)

	// Seed

	users := []domain.User{
		{
			ID:        userID1,
			Email:     "user1@localhost.ru",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        userID2,
			Email:     "user2@localhost.ru",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	userInsert(t, ctx, db, users[0])
	userInsert(t, ctx, db, users[1])
	defer userDelete(t, ctx, db, userID1)
	defer userDelete(t, ctx, db, userID2)

	// Test

	store := &UserRepository{db: db}

	// Found
	got, err := store.FindByID(ctx, userID2)
	assert.NoError(t, err)
	assert.Equal(t, users[1], got)

	// Not found
	_, err = store.FindByID(ctx, "01JK60DX7A2DN9SPEET0JQ6NXE")
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestIntegration_UserRepository_FindByEmail(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID1 = "01JK609DPKC9K4356X1SVM3Z42"
		userID2 = "01JK609HDHZCHN66MTSJJ1F3DR"
	)

	// Seed

	users := []domain.User{
		{
			ID:        userID1,
			Email:     "user1@localhost.ru",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        userID2,
			Email:     "user2@localhost.ru",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	userInsert(t, ctx, db, users[0])
	userInsert(t, ctx, db, users[1])
	defer userDelete(t, ctx, db, userID1)
	defer userDelete(t, ctx, db, userID2)

	// Test

	store := &UserRepository{db: db}

	// Found
	got, err := store.FindByEmail(ctx, "user2@localhost.ru")
	assert.NoError(t, err)
	assert.Equal(t, users[1], got)

	// Not found
	_, err = store.FindByEmail(ctx, "user3@localhost.ru")
	assert.Equal(t, domain.ErrNotFound, err)
}

func userInsert(t *testing.T, ctx context.Context, db *sqlx.DB, user domain.User) {
	t.Helper()

	q := `INSERT INTO "users"
		("id", "email", "activate_code", "created_at")
		VALUES
		(:id, :email, :activate_code, :created_at)`

	_, err := db.NamedExecContext(ctx, q, userModelFromDomain(user))
	assert.NoError(t, err)
}

func userSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) (domain.User, error) {
	t.Helper()

	m := userModel{}

	q := `SELECT
			"id",
			"email",
			"activate_code",
			"created_at",
			"updated_at"
		FROM "users"
		WHERE "id" = ?
		LIMIT 1`

	if err := db.GetContext(ctx, &m, q, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("select user by id: %v", err)
	}

	return userModelToDomain(m), nil
}

func userDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "users" WHERE "id" = :id`

	p := map[string]any{
		"id": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}
