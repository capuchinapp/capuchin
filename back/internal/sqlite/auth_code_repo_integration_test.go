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

func TestIntegration_AuthCodeRepository_Create(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JK4ZW3PN6S5M8XS9TXG2TNFE"
	)

	// Seed

	defer authCodesDelete(t, ctx, db, userID)

	// Test

	store := &AuthCodeRepository{db: db}

	// Create
	authCode := domain.AuthCode{
		ID:        "01JK501C8YGR596NJQTQA90W4N",
		UserID:    userID,
		UserEmail: "user@localhost.ru",
		Code:      "123456",
		CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	}

	err := store.Create(ctx, authCode)
	assert.NoError(t, err)

	got := authCodesSelect(t, ctx, db, userID)[0]
	assert.Equal(t, authCode, got)

	// Create fail - already exists
	err = store.Create(ctx, domain.AuthCode{
		ID:        "01JK501C8YGR596NJQTQA90W4N",
		UserID:    userID,
		UserEmail: "user@localhost.ru",
		Code:      "123456",
		CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	})
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func TestIntegration_AuthCodeRepository_Delete(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JK4ZW7P5J3B5QPXSH0W1WVK6"
	)

	// Seed

	authCodes := []domain.AuthCode{
		{
			ID:        "01JK51EZKP177MW4S2Y0NMWZ91",
			UserID:    userID,
			UserEmail: "user@localhost.ru",
			Code:      "111111",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JK51F3NBGQ4RAP7RN5W16VXY",
			UserID:    userID,
			UserEmail: "user@localhost.ru",
			Code:      "222222",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	authCodesInsert(t, ctx, db, authCodes)
	defer authCodesDelete(t, ctx, db, userID)

	// Test

	store := &AuthCodeRepository{db: db}

	err := store.Delete(ctx, "01JK51EZKP177MW4S2Y0NMWZ91")
	assert.NoError(t, err)

	got := authCodesSelect(t, ctx, db, userID)
	assert.Len(t, got, 1)
	assert.Equal(t, authCodes[1], got[0])
}

func TestIntegration_AuthCodeRepository_Find(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID1 = "01JK51ZE3VQATX66GWJN7VB98E"
		userID2 = "01JK51ZJ72R7JVSXXYS36X2QC4"
	)

	// Seed

	authCodes := []domain.AuthCode{
		{
			ID:        "01JK51Y5RPNCF9PBMTAB9HQ3V6",
			UserID:    userID1,
			UserEmail: "user1@localhost.ru",
			Code:      "111111",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JK51YANBKBAYBM0VJEJM3KRT",
			UserID:    userID2,
			UserEmail: "user2@localhost.ru",
			Code:      "222222",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	authCodesInsert(t, ctx, db, authCodes)
	defer authCodesDelete(t, ctx, db, userID1)
	defer authCodesDelete(t, ctx, db, userID2)

	// Test

	store := &AuthCodeRepository{db: db}

	// Found
	got, err := store.Find(ctx, "user2@localhost.ru", "222222")
	assert.NoError(t, err)
	assert.Equal(t, authCodes[1], got)

	// Not found
	_, err = store.Find(ctx, "user2@localhost.ru", "111111")
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestIntegration_AuthCodeRepository_CancelOld(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JK52R09DJ2EMAK0BDC2SBJEW"
	)

	// Seed

	authCodes := []domain.AuthCode{
		{
			ID:        "01JK51Y5RPNCF9PBMTAB9HQ3V6",
			UserID:    userID,
			UserEmail: "user@localhost.ru",
			Code:      "111111",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JK51YANBKBAYBM0VJEJM3KRT",
			UserID:    userID,
			UserEmail: "user@localhost.ru",
			Code:      "222222",
			CreatedAt: time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
		{
			ID:        "01JK52S20P25GCWDGSM0PZHFYM",
			UserID:    userID,
			UserEmail: "user@localhost.ru",
			Code:      "333333",
			CreatedAt: time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
	}

	authCodesInsert(t, ctx, db, authCodes)
	defer authCodesDelete(t, ctx, db, userID)

	// Test

	store := &AuthCodeRepository{db: db}

	err := store.CancelOld(ctx, time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC))
	assert.NoError(t, err)

	got := authCodesSelect(t, ctx, db, userID)
	assert.Len(t, got, 1)
	assert.Equal(t, authCodes[2], got[0])
}

func authCodesInsert(t *testing.T, ctx context.Context, db *sqlx.DB, authCodes []domain.AuthCode) {
	t.Helper()

	q := `INSERT INTO "auth_codes"
		("id", "user_id", "user_email", "code", "created_at")
		VALUES
		(:id, :user_id, :user_email, :code, :created_at)`

	for _, c := range authCodes {
		_, err := db.NamedExecContext(ctx, q, authCodeModelFromDomain(c))
		assert.NoError(t, err)
	}
}

func authCodesSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) []domain.AuthCode {
	t.Helper()

	ms := []authCodeModel{}

	q := `SELECT
			"id",
			"user_id",
			"user_email",
			"code",
			"created_at"
		FROM "auth_codes"
		WHERE "user_id" = ?`

	if err := db.SelectContext(ctx, &ms, q, userID); err != nil {
		assert.NoError(t, err)

		return nil
	}

	res := make([]domain.AuthCode, len(ms))
	for i := range ms {
		res[i] = authCodeModelToDomain(ms[i])
	}

	return res
}

func authCodesDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "auth_codes" WHERE "user_id" = :user_id`

	p := map[string]any{
		"user_id": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}
