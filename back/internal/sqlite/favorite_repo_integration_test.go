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

func TestIntegration_FavoriteRepository_Create(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JK9ZXWJ02ZBQ2ND2R5RD8QFH"
	)

	// Seed

	defer favoritesDelete(t, ctx, db, userID)

	// Test

	store := &FavoriteRepository{db: db}

	// Create
	favorite := domain.Favorite{
		ID:           "01JKA0EFFVPAK3YY720FTQ74TC",
		UserID:       userID,
		ProjectID:    "01JKA0EMM8Z2YVJ7KVPCGTS7S9",
		TaskID:       func() *string { s := "01JKA0G0QSX9PEPK2N5KGKMR9N"; return &s }(),
		Name:         "Favorite 1",
		BillableRate: 1000,
		Comment:      "Comment 1",
	}

	err := store.Create(ctx, db, favorite)
	assert.NoError(t, err)

	got, err := favoritesSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Equal(t, favorite, got[0])

	// Create fail - already exists
	err = store.Create(ctx, db, domain.Favorite{
		ID:           "01JKA0EFFVPAK3YY720FTQ74TC",
		UserID:       userID,
		ProjectID:    "01JKA0EMM8Z2YVJ7KVPCGTS7S9",
		TaskID:       func() *string { s := "01JKA0G0QSX9PEPK2N5KGKMR9N"; return &s }(),
		Name:         "Favorite 1",
		BillableRate: 1000,
		Comment:      "Comment 1",
	})
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func TestIntegration_FavoriteRepository_Update(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKA0JWTPJNR0M6JC7GJMQC4R"
	)

	// Seed

	favorites := []domain.Favorite{
		{
			ID:           "01JKA0M99W8A9XW7MZJ2G32260",
			UserID:       userID,
			ProjectID:    "01JKA0N5X65YVFRERJ2251MQYB",
			TaskID:       func() *string { s := "01JKA0NET5S9Q90ERDSH25SPG3"; return &s }(),
			Name:         "Favorite 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
		},
		{
			ID:           "01JKA0N0V246H6KRSGFVPMWVJP",
			UserID:       userID,
			ProjectID:    "01JKA0N5X65YVFRERJ2251MQYB",
			TaskID:       func() *string { s := "01JKA0NJJPDVT2DDS1T6YBKXC1"; return &s }(),
			Name:         "Favorite 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
		},
	}

	favoritesInsert(t, ctx, db, favorites)
	defer favoritesDelete(t, ctx, db, userID)

	// Test

	store := &FavoriteRepository{db: db}

	favorites[0].Name = "Favorite 111"
	favorites[0].BillableRate = 10000
	favorites[0].Comment = "Comment 111"

	err := store.Update(ctx, db, favorites[0])
	assert.NoError(t, err)

	got, err := favoritesSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, favorites, got)
}

func TestIntegration_FavoriteRepository_Delete(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKA1FN1AAPX877HW3WX7RKFP"
	)

	// Seed

	favorites := []domain.Favorite{
		{
			ID:           "01JKA1G3SFVW36GKAPDAHN0MQQ",
			UserID:       userID,
			ProjectID:    "01JKA1FVZD9NQ1JTGNZ14CV1H3",
			TaskID:       func() *string { s := "01JKA1GB7SHAYJHGWWJR71ZMTQ"; return &s }(),
			Name:         "Favorite 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
		},
		{
			ID:           "01JKA1G7F7BZBD4SRXAT4ZWGPT",
			UserID:       userID,
			ProjectID:    "01JKA1FVZD9NQ1JTGNZ14CV1H3",
			TaskID:       func() *string { s := "01JKA1GF16QDRWYNGPZBGMVM68"; return &s }(),
			Name:         "Favorite 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
		},
	}

	favoritesInsert(t, ctx, db, favorites)
	defer favoritesDelete(t, ctx, db, userID)

	// Test

	store := &FavoriteRepository{db: db}

	err := store.Delete(ctx, db, userID, "01JKA1G3SFVW36GKAPDAHN0MQQ")
	assert.NoError(t, err)

	got, err := favoritesSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, favorites[1:], got)
}

func TestIntegration_FavoriteRepository_FindByID(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID    = "01JKA1RZNZ2029NYSD0PJ8T3WQ"
		clientID  = "01JKA1S3J6DY5H20YXE09TC201"
		projectID = "01JKA1S7CX0SYHDJCXH7DHDKPT"
		taskID    = "01JKA2JY438Y8QZPNNRCR6SFFF"
	)

	// Seed

	clients := []domain.Client{
		{
			ID:           clientID,
			UserID:       userID,
			Name:         "Client 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	clientsInsert(t, ctx, db, clients)
	defer clientsDelete(t, ctx, db, userID)

	projects := []domain.Project{
		{
			ID:           projectID,
			UserID:       userID,
			ClientID:     clientID,
			Name:         "Project 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	projectsInsert(t, ctx, db, projects)
	defer projectsDelete(t, ctx, db, userID)

	tasks := []domain.Task{
		{
			ID:        taskID,
			UserID:    userID,
			ClientID:  clientID,
			ProjectID: projectID,
			Name:      "Task 1",
			Comment:   "Comment 1",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	tasksInsert(t, ctx, db, tasks)
	defer tasksDelete(t, ctx, db, userID)

	favorites := []domain.Favorite{
		{
			ID:           "01JKA1SWN8PDZEW2YWZY12BZHM",
			UserID:       userID,
			ClientID:     clientID,
			ClientName:   "Client 1",
			ProjectID:    projectID,
			ProjectName:  "Project 1",
			Name:         "Favorite 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
		},
		{
			ID:           "01JKA1T0T9ZPW5E2HR3SSVT498",
			UserID:       userID,
			ClientID:     clientID,
			ClientName:   "Client 1",
			ProjectID:    projectID,
			ProjectName:  "Project 1",
			Name:         "Favorite 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
		},
		{
			ID:           "01JKA2N73WHMCMSEG12HHQZHRQ",
			UserID:       userID,
			ClientID:     clientID,
			ClientName:   "Client 1",
			ProjectID:    projectID,
			ProjectName:  "Project 1",
			TaskID:       func() *string { s := taskID; return &s }(),
			TaskName:     func() *string { s := "Task 1"; return &s }(),
			Name:         "Favorite 3",
			BillableRate: 3000,
			Comment:      "Comment 3",
		},
	}

	favoritesInsert(t, ctx, db, favorites)
	defer favoritesDelete(t, ctx, db, userID)

	// Test

	store := &FavoriteRepository{db: db}

	// Found - without task
	got, err := store.FindByID(ctx, userID, "01JKA1T0T9ZPW5E2HR3SSVT498")
	assert.NoError(t, err)
	assert.Equal(t, favorites[1], got)

	// Found - with task
	got, err = store.FindByID(ctx, userID, "01JKA2N73WHMCMSEG12HHQZHRQ")
	assert.NoError(t, err)
	assert.Equal(t, favorites[2], got)

	// Not found
	_, err = store.FindByID(ctx, userID, "01JKA1TG9GP8RN019FT90T3V86")
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestIntegration_FavoriteRepository_FindAll(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID1    = "01JKA2TXRRR6SC1Z5XMRTHFK78"
		userID2    = "01JKA2V1S2TWN8G1E8VE9NBTS9"
		clientID1  = "01JKA2VGCGYWXD3R429EDC7VNG"
		clientID2  = "01JKA2VMMA9H83G6SEMVT6BT6G"
		projectID1 = "01JKA2VRT92VB544AGC6KBPPNJ"
		projectID2 = "01JKA2VX09P7JE8D9CW7NE4SVM"
		taskID     = "01JKA2W0TNQ14HGYY7JKZFJRRB"
	)

	// Seed

	clients := []domain.Client{
		{
			ID:           clientID1,
			UserID:       userID1,
			Name:         "Client 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           clientID2,
			UserID:       userID2,
			Name:         "Client 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	clientsInsert(t, ctx, db, clients)
	defer clientsDelete(t, ctx, db, userID1)
	defer clientsDelete(t, ctx, db, userID2)

	projects := []domain.Project{
		{
			ID:           projectID1,
			UserID:       userID1,
			ClientID:     clientID1,
			Name:         "Project 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           projectID2,
			UserID:       userID2,
			ClientID:     clientID2,
			Name:         "Project 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	projectsInsert(t, ctx, db, projects)
	defer projectsDelete(t, ctx, db, userID1)
	defer projectsDelete(t, ctx, db, userID2)

	tasks := []domain.Task{
		{
			ID:        taskID,
			UserID:    userID2,
			ClientID:  clientID2,
			ProjectID: projectID2,
			Name:      "Task 1",
			Comment:   "Comment 1",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}

	tasksInsert(t, ctx, db, tasks)
	defer tasksDelete(t, ctx, db, userID2)

	favorites := []domain.Favorite{
		{
			ID:           "01JKA30EC2J5AVME94GDX0QFHV",
			UserID:       userID1,
			ClientID:     clientID1,
			ClientName:   "Client 1",
			ProjectID:    projectID1,
			ProjectName:  "Project 1",
			Name:         "Favorite 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
		},
		{
			ID:           "01JKA30Q08TR7ASM0T8888Y1QK",
			UserID:       userID1,
			ClientID:     clientID1,
			ClientName:   "Client 1",
			ProjectID:    projectID1,
			ProjectName:  "Project 1",
			Name:         "Favorite 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
		},
		{
			ID:           "01JKA30JA9QYSH90123XRAJ10B",
			UserID:       userID2,
			ClientID:     clientID2,
			ClientName:   "Client 2",
			ProjectID:    projectID2,
			ProjectName:  "Project 2",
			TaskID:       func() *string { s := taskID; return &s }(),
			TaskName:     func() *string { s := "Task 1"; return &s }(),
			Name:         "Favorite 3",
			BillableRate: 3000,
			Comment:      "Comment 3",
		},
	}

	favoritesInsert(t, ctx, db, favorites)
	defer favoritesDelete(t, ctx, db, userID1)
	defer favoritesDelete(t, ctx, db, userID2)

	// Test

	store := &FavoriteRepository{db: db}

	// Found - without tasks
	got, err := store.FindAll(ctx, userID1)
	assert.NoError(t, err)
	assert.Equal(t, favorites[:2], got)

	// Found - with tasks
	got, err = store.FindAll(ctx, userID2)
	assert.NoError(t, err)
	assert.Equal(t, favorites[2:], got)
}

func favoritesInsert(t *testing.T, ctx context.Context, db *sqlx.DB, favorites []domain.Favorite) {
	t.Helper()

	q := `INSERT INTO "favorites"
		("id", "user_id", "name", "project_id", "task_id", "billable_rate", "comment")
		VALUES
		(:id, :user_id, :name, :project_id, :task_id, :billable_rate, :comment)`

	for _, c := range favorites {
		_, err := db.NamedExecContext(ctx, q, favoriteModelFromDomain(c))
		assert.NoError(t, err)
	}
}

func favoritesSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) ([]domain.Favorite, error) {
	t.Helper()

	ms := []favoriteModel{}

	q := `SELECT
			"id",
			"user_id",
			"name",
			"project_id",
			"task_id",
			"billable_rate",
			"comment"
		FROM "favorites"
		WHERE "user_id" = ?`

	if err := db.SelectContext(ctx, &ms, q, userID); err != nil {
		return nil, fmt.Errorf("select favorites by id: %v", err)
	}

	return favoriteModelsToDomains(ms), nil
}

func favoritesDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "favorites" WHERE "user_id" = :user_id`

	p := map[string]any{
		"user_id": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}
