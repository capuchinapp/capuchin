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

func TestIntegration_ProjectRepository_Create(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JK82GCDNYKX92PJZQCDZ2NSR"
	)

	// Seed

	defer projectsDelete(t, ctx, db, userID)

	// Test

	store := &ProjectRepository{db: db}

	// Create
	project := domain.Project{
		ID:           "01JK82GJDY5Z4A7HA2XWRHQDVQ",
		UserID:       userID,
		ClientID:     "01JK82HDZHA0437PV5372T3EWZ",
		Name:         "Project 1",
		BillableRate: 1000,
		Comment:      "Comment 1",
		CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	}

	err := store.Create(ctx, db, project)
	assert.NoError(t, err)

	got, err := projectsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Equal(t, project, got[0])

	// Create fail - already exists
	err = store.Create(ctx, db, domain.Project{
		ID:           "01JK82GJDY5Z4A7HA2XWRHQDVQ",
		UserID:       userID,
		ClientID:     "01JK82HDZHA0437PV5372T3EWZ",
		Name:         "Project 1",
		BillableRate: 1000,
		Comment:      "Comment 1",
		CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
	})
	assert.Equal(t, domain.ErrAlreadyExists, err)
}

func TestIntegration_ProjectRepository_Update(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JK82J9Q9XDXQXZYM2377S1WC"
	)

	// Seed

	projects := []domain.Project{
		{
			ID:           "01JK82JDK2RPXK91VXBGSBDPPX",
			UserID:       userID,
			ClientID:     "01JK82K4AMAR65XWTZKZ0Y4DWR",
			Name:         "Project 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           "01JK82JHNT51FPCM5AYC16AM5K",
			UserID:       userID,
			ClientID:     "01JK82K89TX1K5NV7WY3AGADXN",
			Name:         "Project 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
			CreatedAt:    time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
	}

	projectsInsert(t, ctx, db, projects)
	defer projectsDelete(t, ctx, db, userID)

	// Test

	store := &ProjectRepository{db: db}

	projects[0].Name = "Project 111"
	projects[0].BillableRate = 10000
	projects[0].Comment = "Comment 111"
	projects[0].UpdatedAt = func() *time.Time { t := time.Date(2025, 2, 3, 9, 15, 25, 0, time.UTC); return &t }()
	projects[0].ArchivedAt = func() *time.Time { t := time.Date(2025, 2, 3, 9, 16, 25, 0, time.UTC); return &t }()

	err := store.Update(ctx, db, projects[0])
	assert.NoError(t, err)

	got, err := projectsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, projects, got)
}

func TestIntegration_ProjectRepository_FindByFilter(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID    = "01JK82NEXBSKC3DC4FSW2MJEH9"
		clientID1 = "01JK82PKCA1W418CX5XGDDT37A"
		clientID2 = "01JK82PT4Q773P7FG6JQBVX0ZA"
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
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
			ArchivedAt:   func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
		},
	}

	clientsInsert(t, ctx, db, clients)
	defer clientsDelete(t, ctx, db, userID)

	projects := []domain.Project{
		{
			ID:           "01JK82NQ218YANNADPKE0481VW",
			UserID:       userID,
			ClientID:     clientID1,
			ClientName:   "Client 1",
			Name:         "Project 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           "01JK82NTQ6GPW1RJ8SSGYQ8Z9W",
			UserID:       userID,
			ClientID:     clientID1,
			ClientName:   "Client 1",
			Name:         "Project 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
			CreatedAt:    time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
		{
			ID:           "01JK82NYC7TJFPR30P1Y0E965Y",
			UserID:       userID,
			ClientID:     clientID2,
			ClientName:   "Client 2",
			Name:         "Project 3",
			BillableRate: 3000,
			Comment:      "Comment 3",
			CreatedAt:    time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
	}

	projectsInsert(t, ctx, db, projects)
	defer projectsDelete(t, ctx, db, userID)

	// Test

	store := &ProjectRepository{db: db}

	// Find by filter without archived clients
	got, err := store.FindByFilter(ctx, userID, domain.ProjectFilter{
		WithoutArchivedClients: true,
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, projects[:2], got)

	// Find by filter by client ID
	got, err = store.FindByFilter(ctx, userID, domain.ProjectFilter{
		ClientID: clientID2,
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, projects[2:], got)

	// Find without filter
	got, err = store.FindByFilter(ctx, userID, domain.ProjectFilter{})
	assert.NoError(t, err)
	assert.ElementsMatch(t, projects, got)
}

func TestIntegration_ProjectRepository_FindByID(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID     = "01JK82W0EY2NSM4R64ND9R8YNA"
		clientID   = "01JK839TE840RB2KW76PMHPG0P"
		projectID1 = "01JK82W3ZCZYJNX5MRXV23HTT7"
		projectID2 = "01JK82W872QKJ2XNB7SSAQ84QH"
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
			ClientName:   "Client 1",
			Name:         "Project 1",
			BillableRate: 1000,
			Comment:      "Comment 1",
			CreatedAt:    time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			ID:           projectID2,
			UserID:       userID,
			ClientID:     clientID,
			ClientName:   "Client 1",
			Name:         "Project 2",
			BillableRate: 2000,
			Comment:      "Comment 2",
			CreatedAt:    time.Date(2025, 2, 3, 9, 12, 25, 0, time.UTC),
		},
	}

	projectsInsert(t, ctx, db, projects)
	defer projectsDelete(t, ctx, db, userID)

	// Test

	store := &ProjectRepository{db: db}

	// Found
	got, err := store.FindByID(ctx, userID, projectID2)
	assert.NoError(t, err)
	assert.Equal(t, projects[1], got)

	// Not found
	_, err = store.FindByID(ctx, userID, "01JK82WP5SRR8GW7PRHYN80XMW")
	assert.Equal(t, domain.ErrNotFound, err)
}

func projectsInsert(t *testing.T, ctx context.Context, db *sqlx.DB, projects []domain.Project) {
	t.Helper()

	q := `INSERT INTO "projects"
		("id", "user_id", "client_id", "name", "comment", "billable_rate", "created_at", "updated_at", "archived_at")
		VALUES
		(:id, :user_id, :client_id, :name, :comment, :billable_rate, :created_at, :updated_at, :archived_at)`

	for _, c := range projects {
		_, err := db.NamedExecContext(ctx, q, projectModelFromDomain(c))
		assert.NoError(t, err)
	}
}

func projectsSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) ([]domain.Project, error) {
	t.Helper()

	ms := []projectModel{}

	q := `SELECT
			"id",
			"user_id",
			"client_id",
			"name",
			"billable_rate",
			"comment",
			"created_at",
			"updated_at",
			"archived_at"
		FROM "projects"
		WHERE "user_id" = ?`

	if err := db.SelectContext(ctx, &ms, q, userID); err != nil {
		return nil, fmt.Errorf("select projects by id: %v", err)
	}

	return projectModelsToDomains(ms), nil
}

func projectsDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "projects" WHERE "user_id" = :user_id`

	p := map[string]any{
		"user_id": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}
