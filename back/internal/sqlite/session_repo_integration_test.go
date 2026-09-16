//go:build integration && sqlite

package sqlite

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"capuchin/internal/domain"
)

func TestIntegration_SessionRepository_Create(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKF7JCSTRTDJ5CQS2CW0BPX5"
	)

	// Seed

	defer sessionsDelete(t, ctx, db, userID)

	// Test

	store := &SessionRepository{db: db}

	session := domain.Session{
		ID:        "01JKF7MQ0P9B2EQDEJCBQE3GFQ",
		CookieID:  "01JKF7N4V2QYJ7WFXM2BMC3GY4",
		UserID:    userID,
		CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	}

	err := store.Create(ctx, session)
	assert.NoError(t, err)

	got, err := sessionsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Equal(t, session, got[0])
}

func TestIntegration_SessionRepository_DeleteByCID(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKF7XPP2X8QKY20ZV3WXT6VQ"
	)

	// Seed

	sessions := []domain.Session{
		{
			ID:        "01JKF7Y7YFTTM9YBK2789TK6P9",
			CookieID:  "01JKF7YBBKDEZ6D0KA33YHEYPY",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JKF7YEXX35SEGB3JCAWHD6EN",
			CookieID:  "01JKF7YJAZ4ZQ96BZ42TQKK53Q",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	sessionsInsert(t, ctx, db, sessions)
	defer sessionsDelete(t, ctx, db, userID)

	// Test

	store := &SessionRepository{db: db}

	err := store.DeleteByCID(ctx, "01JKF7YBBKDEZ6D0KA33YHEYPY")
	assert.NoError(t, err)

	got, err := sessionsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.ElementsMatch(t, sessions[1:], got)
}

func TestIntegration_SessionRepository_FindAll(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID1 = "01JKF8A88TF6G869NSGDEVQEAX"
		userID2 = "01JKF8ABX6MJ646VD4K5C0BYAF"
	)

	// Seed

	sessions := []domain.Session{
		{
			ID:        "01JKF8BW8PX4ADFTWT5CTG4V3R",
			CookieID:  "01JKF8C5F9BMT7BA1THMAGF8YP",
			UserID:    userID1,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JKF8C94Z06Z73P61KBV55JCE",
			CookieID:  "01JKF8CCSD9NGGNBX8JMJJGG90",
			UserID:    userID1,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JKF8CGFBN4MXRJWF829GKDNY",
			CookieID:  "01JKF8CMAETARNQ8A1XR443B03",
			UserID:    userID2,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	sessionsInsert(t, ctx, db, sessions)
	defer sessionsDelete(t, ctx, db, userID1)
	defer sessionsDelete(t, ctx, db, userID2)

	// Test

	store := &SessionRepository{db: db}

	// Found
	got, err := store.FindAll(ctx, userID1)
	assert.NoError(t, err)
	assert.ElementsMatch(t, sessions[:2], got)

	// Not found
	got0, err := store.FindAll(ctx, "01JKF8DCRTC7GWTK2DNQM9K0A7")
	assert.Len(t, got0, 0)
}

func TestIntegration_SessionRepository_FindByCID(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKF8MAJEJ8W0ZEBA6X79PNXD"
	)

	// Seed

	sessions := []domain.Session{
		{
			ID:        "01JKF8N4XCHGFHHJ9KW4TEDFNT",
			CookieID:  "01JKF8N8SBE80PR89K68YMHX4H",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JKF8NCE6NX54KZQB3BKD6Z4Y",
			CookieID:  "01JKF8NFZRVNJB5RD2N7BTTQPQ",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	sessionsInsert(t, ctx, db, sessions)
	defer sessionsDelete(t, ctx, db, userID)

	// Test

	store := &SessionRepository{db: db}

	// Found
	got, err := store.FindByCID(ctx, "01JKF8N8SBE80PR89K68YMHX4H")
	assert.NoError(t, err)
	assert.Equal(t, sessions[0], got)

	// Not found
	_, err = store.FindByCID(ctx, "01JKF8P6JXJH1TEP959WH0EWC7")
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestIntegration_SessionRepository_FindByID(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKF8RB09JRDSY18EKECEMPKT"
	)

	// Seed

	sessions := []domain.Session{
		{
			ID:        "01JKF8REWHAWT1D7P82H7T6ZG9",
			CookieID:  "01JKF8RNG6JJJ8HP8FT44PTWP0",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JKF8RSPS6XGFTJ3R94MZWRRR",
			CookieID:  "01JKF8RX61PPBFG5ZM61CEEG44",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	sessionsInsert(t, ctx, db, sessions)
	defer sessionsDelete(t, ctx, db, userID)

	// Test

	store := &SessionRepository{db: db}

	// Found
	got, err := store.FindByID(ctx, userID, "01JKF8REWHAWT1D7P82H7T6ZG9")
	assert.NoError(t, err)
	assert.Equal(t, sessions[0], got)

	// Not found
	_, err = store.FindByID(ctx, userID, "01JKF8S14KMC7WANA9HREW8YCX")
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestIntegration_SessionRepository_UpdateCheckedAt(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKF95QTF6B0AF4TVR1NB3Q1G"
	)

	// Seed

	sessions := []domain.Session{
		{
			ID:        "01JKF95VFWBC6C7J8405XG9ZBY",
			CookieID:  "01JKF95ZEQ3YGDEZGEJSYM0JGS",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JKF963A0AEJSZZ141VW48GXE",
			CookieID:  "01JKF96704ZKF0E2RE24RJKZ2N",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
	}

	sessionsInsert(t, ctx, db, sessions)
	defer sessionsDelete(t, ctx, db, userID)

	// Test

	store := &SessionRepository{db: db}

	sessions[0].CheckedAt = time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC)

	err := store.UpdateCheckedAt(ctx, sessions[0].ID, sessions[0].CheckedAt)
	assert.NoError(t, err)

	got, err := sessionsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, sessions, got)
}

func TestIntegration_SessionRepository_CancelOld(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKF8WHJ09HXB76N3YKMB7D1F"
	)

	// Seed

	sessions := []domain.Session{
		{
			ID:        "01JKF8WNG907441MYCP7DYRWJ1",
			CookieID:  "01JKF8WRSNTZE13YFBSND27RFX",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JKF8WWCQKF9CV7YKB8Y59R7Q",
			CookieID:  "01JKF8X2B1HG179KSADRZTWAS4",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
		{
			ID:        "01JKF8X5VKTGA5KPMEXDWRQNCQ",
			CookieID:  "01JKF8X9X9D42PS46VX54YNTMX",
			UserID:    userID,
			CheckedAt: time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
	}

	sessionsInsert(t, ctx, db, sessions)
	defer sessionsDelete(t, ctx, db, userID)

	// Test

	store := &SessionRepository{db: db}

	err := store.CancelOld(ctx, time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC))
	assert.NoError(t, err)

	got, err := sessionsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, sessions[1:], got)
}

func sessionsInsert(t *testing.T, ctx context.Context, db *sqlx.DB, sessions []domain.Session) {
	t.Helper()

	q := `INSERT INTO "sessions"
		("id", "cid", "uid", "checked_at")
		VALUES
		(:id, :cid, :uid, :checked_at)`

	for _, c := range sessions {
		_, err := db.NamedExecContext(ctx, q, sessionModelFromDomain(c))
		assert.NoError(t, err)
	}
}

func sessionsSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) ([]domain.Session, error) {
	t.Helper()

	ms := []sessionModel{}

	q := `SELECT
			"id",
			"cid",
			"uid",
			"checked_at"
		FROM "sessions"
		WHERE "uid" = ?`

	if err := db.SelectContext(ctx, &ms, q, userID); err != nil {
		return nil, fmt.Errorf("select sessions by id: %v", err)
	}

	return sessionModelsToDomains(ms), nil
}

func sessionsDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "sessions" WHERE "uid" = :uid`

	p := map[string]any{
		"uid": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}
