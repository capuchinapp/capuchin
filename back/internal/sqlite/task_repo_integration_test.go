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

func TestIntegration_TaskRepository_Create(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKCTE9RHTCQJS0J64361YSDE"
	)

	// Seed

	defer tasksDelete(t, ctx, db, userID)

	// Test

	store := &TaskRepository{db: db}

	// Create
	task := domain.Task{
		ID:        "01JKCTHF648YWTNAGCWKK3BDNP",
		UserID:    userID,
		ProjectID: "01JKCTHJW25YKKYGSWWZ70RFDY",
		Name:      "Task 1",
		Comment:   "Comment 1",
		CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	}

	err := store.Create(ctx, db, task)
	assert.NoError(t, err)

	got, err := tasksSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Equal(t, task, got[0])

	// Create fail - already exists
	err = store.Create(ctx, db, domain.Task{
		ID:        "01JKCTHF648YWTNAGCWKK3BDNP",
		UserID:    userID,
		ProjectID: "01JKCTHJW25YKKYGSWWZ70RFDY",
		Name:      "Task 1",
		Comment:   "Comment 1",
		CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	})
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func TestIntegration_TaskRepository_Update(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKCTKYYZW53EM0CKV1T3ZR1R"
	)

	// Seed

	tasks := []domain.Task{
		{
			ID:        "01JKCTMCTN00E8HDQPBXHNT0WX",
			UserID:    userID,
			ProjectID: "01JKCTM4KK1RDSR7J0KV6JGRJ0",
			Name:      "Task 1",
			Comment:   "Comment 1",
			CreatedAt: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:        "01JKCTMG8FY3A928WV2M3VX5BS",
			UserID:    userID,
			ProjectID: "01JKCTM4KK1RDSR7J0KV6JGRJ0",
			Name:      "Task 2",
			Comment:   "Comment 2",
			CreatedAt: time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
	}

	tasksInsert(t, ctx, db, tasks)
	defer tasksDelete(t, ctx, db, userID)

	// Test

	store := &TaskRepository{db: db}

	tasks[0].Name = "Task 111"
	tasks[0].Comment = "Comment 111"
	tasks[0].UpdatedAt = func() *time.Time { t := time.Date(2025, 2, 3, 9, 15, 25, 0, time.UTC); return &t }()
	tasks[0].ArchivedAt = func() *time.Time { t := time.Date(2025, 2, 3, 9, 16, 25, 0, time.UTC); return &t }()
	tasks[0].CompletedAt = func() *time.Time { t := time.Date(2025, 2, 3, 9, 17, 25, 0, time.UTC); return &t }()

	err := store.Update(ctx, db, tasks[0])
	assert.NoError(t, err)

	got, err := tasksSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, tasks, got)
}

func TestIntegration_TaskRepository_FindByFilter(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID     = "01JKCV7Z4H7M83MTKEZSSJX6ZW"
		clientID   = "01JKCV82CREA17SAWG06KK2K3K"
		projectID1 = "01JKCV86TB15DHBQBA0A7AEJK3"
		projectID2 = "01JKCV8A32FX8GC3FZKC60SPKY"
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
			ID:           projectID1,
			UserID:       userID,
			ClientID:     clientID,
			Name:         "Project 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           projectID2,
			UserID:       userID,
			ClientID:     clientID,
			Name:         "Project 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
			ArchivedAt:   func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
		},
	}

	projectsInsert(t, ctx, db, projects)
	defer projectsDelete(t, ctx, db, userID)

	tasks := []domain.Task{
		{
			ID:          "01JKCV8GTDGR8F64ZYEBJ3JN3D",
			UserID:      userID,
			ClientID:    clientID,
			ClientName:  "Client 1",
			ProjectID:   projectID1,
			ProjectName: "Project 1",
			Name:        "Task 1",
			Comment:     "Comment 1",
			CreatedAt:   time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:          "01JKCV8M8MPCE3G7T4K7NHQPCS",
			UserID:      userID,
			ClientID:    clientID,
			ClientName:  "Client 1",
			ProjectID:   projectID2,
			ProjectName: "Project 2",
			Name:        "Task 2",
			Comment:     "Comment 2",
			CreatedAt:   time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
	}

	tasksInsert(t, ctx, db, tasks)
	defer tasksDelete(t, ctx, db, userID)

	// Test

	store := &TaskRepository{db: db}

	// Find by filter without archived projects
	got, err := store.FindByFilter(ctx, userID, domain.TaskFilter{
		WithoutArchivedProjects: true,
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, tasks[:1], got)

	// Find by filter by project ID
	got, err = store.FindByFilter(ctx, userID, domain.TaskFilter{
		ProjectID: projectID2,
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, tasks[1:], got)

	// Find without filter
	got, err = store.FindByFilter(ctx, userID, domain.TaskFilter{})
	assert.NoError(t, err)
	assert.ElementsMatch(t, tasks, got)
}

func TestIntegration_TaskRepository_FindByID(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID    = "01JKCV8S2HC58VMJM2XTE7NDZJ"
		clientID  = "01JKCV8W4ST6X5C5DKFB6QY4XF"
		projectID = "01JKCVAFJ1QWG4X073JJJGX8T4"
		taskID1   = "01JKCV8ZVJH5V0J29K3RAXJCK6"
		taskID2   = "01JKCV932WW603Q7PTAGRFRH3Y"
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
			ID:          taskID1,
			UserID:      userID,
			ClientID:    clientID,
			ClientName:  "Client 1",
			ProjectID:   projectID,
			ProjectName: "Project 1",
			Name:        "Task 1",
			Comment:     "Comment 1",
			CreatedAt:   time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:          taskID2,
			UserID:      userID,
			ClientID:    clientID,
			ClientName:  "Client 1",
			ProjectID:   projectID,
			ProjectName: "Project 1",
			Name:        "Task 2",
			Comment:     "Comment 2",
			CreatedAt:   time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
	}

	tasksInsert(t, ctx, db, tasks)
	defer tasksDelete(t, ctx, db, userID)

	// Test

	store := &TaskRepository{db: db}

	// Found
	got, err := store.FindByID(ctx, userID, taskID2)
	assert.NoError(t, err)
	assert.Equal(t, tasks[1], got)

	// Not found
	_, err = store.FindByID(ctx, userID, "01JKCV998CZWBC99PX5B0YZ8TD")
	assert.Equal(t, domain.ErrNotFound, err)
}

func tasksInsert(t *testing.T, ctx context.Context, db *sqlx.DB, tasks []domain.Task) {
	t.Helper()

	q := `INSERT INTO "tasks"
		("id", "user_id", "project_id", "name", "comment", "created_at", "updated_at", "archived_at", "completed_at")
		VALUES
		(:id, :user_id, :project_id, :name, :comment, :created_at, :updated_at, :archived_at, :completed_at)`

	for _, c := range tasks {
		_, err := db.NamedExecContext(ctx, q, taskModelFromDomain(c))
		assert.NoError(t, err)
	}
}

func tasksSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) ([]domain.Task, error) {
	t.Helper()

	ms := []taskModel{}

	q := `SELECT
			"id",
			"user_id",
			"project_id",
			"name",
			"comment",
			"created_at",
			"updated_at",
			"archived_at",
			"completed_at"
		FROM "tasks"
		WHERE "user_id" = ?`

	if err := db.SelectContext(ctx, &ms, q, userID); err != nil {
		return nil, fmt.Errorf("select tasks by id: %v", err)
	}

	return taskModelsToDomains(ms), nil
}

func tasksDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "tasks" WHERE "user_id" = :user_id`

	p := map[string]any{
		"user_id": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}
