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

func TestIntegration_TimelogRepository_Create(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKFD6PM2QEE2Q7YHYFNW8GXB"
	)

	// Seed

	defer timelogsDelete(t, ctx, db, userID)

	// Test

	store := &TimelogRepository{db: db}

	timelog := domain.Timelog{
		ID:              "01JKFD85PPXEXH2BRCCAXZP87W",
		UserID:          userID,
		ProjectID:       "01JKFD8B58RR8D0MQAXQE6J1VV",
		TaskID:          func() *string { s := "01JKFD8F0ZBG3SGJ4WZXGGPGH2"; return &s }(),
		Date:            "2025-02-03",
		TimeStart:       "09:13:25",
		TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
		DurationSeconds: 3600,
		BillableRate:    1000,
		BillableAmount:  1000,
		Comment:         "Comment 1",
		CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
	}

	err := store.Create(ctx, db, timelog)
	assert.NoError(t, err)

	got, err := timelogsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.Equal(t, timelog, got[0])
}

func TestIntegration_TimelogRepository_Update(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKFDETZ5RGR692SXB11KJMYG"
	)

	// Seed

	timelogs := []domain.Timelog{
		{
			ID:              "01JKFDEYH8H5ZCJ8TD4935DBZ7",
			UserID:          userID,
			ProjectID:       "01JKFDF1WF5Y9PJ0KHRJ50P2BF",
			TaskID:          func() *string { s := "01JKFDF59Y9H172FKX7C0P6KE9"; return &s }(),
			Date:            "2025-02-03",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Comment 1",
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
	}

	timelogsInsert(t, ctx, db, timelogs)
	defer timelogsDelete(t, ctx, db, userID)

	// Test

	store := &TimelogRepository{db: db}

	timelogs[0].TaskID = nil
	timelogs[0].Date = "2025-03-03"
	timelogs[0].TimeStart = "10:13:25"
	timelogs[0].TimeEnd = func() *string { s := "12:13:25"; return &s }()
	timelogs[0].DurationSeconds = 7200
	timelogs[0].BillableRate = 1000
	timelogs[0].BillableAmount = 2000
	timelogs[0].Comment = "Comment 2"
	timelogs[0].UpdatedAt = func() *time.Time { t := time.Date(2025, 2, 3, 9, 15, 25, 0, time.UTC); return &t }()

	err := store.Update(ctx, db, timelogs[0])
	assert.NoError(t, err)

	got, err := timelogsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, timelogs, got)
}

func TestIntegration_TimelogRepository_Delete(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID = "01JKFDPF160VF1MGQX2VYPFK9T"
	)

	// Seed

	timelogs := []domain.Timelog{
		{
			ID:              "01JKFDPK3PDC104P1PAZY25YVT",
			UserID:          userID,
			ProjectID:       "01JKFDPPM2JT081NGHH04FY5Y2",
			TaskID:          func() *string { s := "01JKFDPT6ZDYMQDHDTB9348YDY"; return &s }(),
			Date:            "2025-02-03",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Comment 1",
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKFDQ1NZ8MPSV9CMBGYGWJFW",
			UserID:          userID,
			ProjectID:       "01JKFDPPM2JT081NGHH04FY5Y2",
			Date:            "2025-02-03",
			TimeStart:       "09:13:25",
			DurationSeconds: 0,
			BillableRate:    1000,
			BillableAmount:  0,
			Comment:         "Comment 2",
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
	}

	timelogsInsert(t, ctx, db, timelogs)
	defer timelogsDelete(t, ctx, db, userID)

	// Test

	store := &TimelogRepository{db: db}

	err := store.Delete(ctx, db, userID, "01JKFDPK3PDC104P1PAZY25YVT")
	assert.NoError(t, err)

	got, err := timelogsSelect(t, ctx, db, userID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, timelogs[1:], got)
}

func TestIntegration_TimelogRepository_FindByID(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID    = "01JKFDYREN5FGDK5ZWJT0KH07X"
		clientID  = "01JKFDYZQY22YMY087KSE93PE0"
		projectID = "01JKFDZ3BWX3Q8RBY0REHKGKBM"
		taskID    = "01JKFDZ6QB28XX2BHWB2G5655S"
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
			ID:          taskID,
			UserID:      userID,
			ClientID:    clientID,
			ProjectID:   projectID,
			Name:        "Task 1",
			Comment:     "Comment 1",
			CreatedAt:   time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
			CompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
		},
	}

	tasksInsert(t, ctx, db, tasks)
	defer tasksDelete(t, ctx, db, userID)

	timelogs := []domain.Timelog{
		{
			ID:              "01JKFE14SHVS5BJ7RJ3R804MA4",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID,
			ProjectName:     "Project 1",
			TaskID:          func() *string { s := taskID; return &s }(),
			TaskName:        func() *string { s := "Task 1"; return &s }(),
			TaskCompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
			Date:            "2025-02-03",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Comment 1",
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKFE18QZZVPTH4VK5WM4C12S",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID,
			ProjectName:     "Project 1",
			Date:            "2025-02-03",
			TimeStart:       "09:13:25",
			DurationSeconds: 0,
			BillableRate:    1000,
			BillableAmount:  0,
			Comment:         "Comment 2",
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
	}

	timelogsInsert(t, ctx, db, timelogs)
	defer timelogsDelete(t, ctx, db, userID)

	// Test

	store := &TimelogRepository{db: db}

	// Found - without task
	got, err := store.FindByID(ctx, userID, "01JKFE18QZZVPTH4VK5WM4C12S")
	assert.NoError(t, err)
	assert.Equal(t, timelogs[1], got)

	// Found - with task
	got, err = store.FindByID(ctx, userID, "01JKFE14SHVS5BJ7RJ3R804MA4")
	assert.NoError(t, err)
	assert.Equal(t, timelogs[0], got)

	// Not found
	_, err = store.FindByID(ctx, userID, "01JKFE2CEPJA0WPEAG41N6VBG4")
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestIntegration_TimelogRepository_FindByFilter(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID     = "01JKG06W88400FYHGRY69ZSWR4"
		clientID   = "01JKG07VCWW9E035S7CGAF13G6"
		projectID1 = "01JKG0836GTG8T9ER2PXYE1XE6"
		projectID2 = "01JKG0872VB291BVGK2QX78W78"
		taskID     = "01JKG08KS7ZDVY9J3NNZGEGV58"
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
		},
	}

	projectsInsert(t, ctx, db, projects)
	defer projectsDelete(t, ctx, db, userID)

	tasks := []domain.Task{
		{
			ID:          taskID,
			UserID:      userID,
			ClientID:    clientID,
			ProjectID:   projectID2,
			Name:        "Task 1",
			Comment:     "Comment 1",
			CreatedAt:   time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
			CompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
		},
	}

	tasksInsert(t, ctx, db, tasks)
	defer tasksDelete(t, ctx, db, userID)

	timelogs := []domain.Timelog{
		// 2025-02-03
		{
			ID:              "01JKG315TT345Y0PNDPKS9YYZR",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID2,
			ProjectName:     "Project 2",
			Date:            "2025-02-03",
			TimeStart:       "11:13:25",
			DurationSeconds: 0,
			BillableRate:    1000,
			BillableAmount:  0,
			Comment:         "Timelog 6",
			CreatedAt:       time.Date(2025, 2, 3, 11, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKG311WGSYJ3T3ACJH0PJTNG",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID1,
			ProjectName:     "Project 1",
			TaskID:          func() *string { s := taskID; return &s }(),
			TaskName:        func() *string { s := "Task 1"; return &s }(),
			TaskCompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
			Date:            "2025-02-03",
			TimeStart:       "10:13:25",
			TimeEnd:         func() *string { s := "11:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 5",
			CreatedAt:       time.Date(2025, 2, 3, 10, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKG30XWT1GWZTB7PPWTSR1AQ",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID1,
			ProjectName:     "Project 1",
			Date:            "2025-02-03",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 4",
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
		// 2025-02-02
		{
			ID:              "01JKG30SP6S0R37F2XQN14NZQR",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID2,
			ProjectName:     "Project 2",
			Date:            "2025-02-02",
			TimeStart:       "11:13:25",
			TimeEnd:         func() *string { s := "12:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 3",
			CreatedAt:       time.Date(2025, 2, 2, 11, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKG30N6QMP37B0ZPJC6SG9ZX",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID1,
			ProjectName:     "Project 1",
			TaskID:          func() *string { s := taskID; return &s }(),
			TaskName:        func() *string { s := "Task 1"; return &s }(),
			TaskCompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
			Date:            "2025-02-02",
			TimeStart:       "10:13:25",
			TimeEnd:         func() *string { s := "11:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 2",
			CreatedAt:       time.Date(2025, 2, 2, 10, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKG30FHX44SSQHX282SR23ET",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID1,
			ProjectName:     "Project 1",
			Date:            "2025-02-02",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 1",
			CreatedAt:       time.Date(2025, 2, 2, 9, 13, 25, 0, time.UTC),
		},
	}

	timelogsInsert(t, ctx, db, timelogs)
	defer timelogsDelete(t, ctx, db, userID)

	// Test

	store := &TimelogRepository{db: db}

	var (
		got []domain.Timelog
		err error
	)

	// Find by filter by date
	got, err = store.FindByFilter(ctx, userID, domain.TimelogFilter{
		DateFrom: "2025-02-02",
		DateTo:   "2025-02-03",
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, timelogs, got)

	// Find by filter by date and client ID
	got, err = store.FindByFilter(ctx, userID, domain.TimelogFilter{
		DateFrom: "2025-02-03",
		DateTo:   "2025-02-03",
		ClientID: clientID,
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, timelogs[:3], got)

	// Find by filter by date and project ID
	got, err = store.FindByFilter(ctx, userID, domain.TimelogFilter{
		DateFrom:  "2025-02-03",
		DateTo:    "2025-02-03",
		ProjectID: projectID1,
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, timelogs[1:3], got)
}

func TestIntegration_TimelogRepository_FindRunning(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID    = "01JKGHQ2KRRDFN85WTCVDBEDFX"
		clientID  = "01JKGHQD8W5PXMD5D7BHHC7VK3"
		projectID = "01JKGHQH9SS4MHQ20C0K64AJ72"
		taskID    = "01JKGHQRFX0QGESRAJM30WVJX5"
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
			ID:          taskID,
			UserID:      userID,
			ClientID:    clientID,
			ProjectID:   projectID,
			Name:        "Task 1",
			Comment:     "Comment 1",
			CreatedAt:   time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
			CompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
		},
	}

	tasksInsert(t, ctx, db, tasks)
	defer tasksDelete(t, ctx, db, userID)

	timelogs := []domain.Timelog{
		{
			ID:              "01JKGJ7H07NFWYFSB8FDDRDM5B",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID,
			ProjectName:     "Project 1",
			Date:            "2025-02-03",
			TimeStart:       "11:13:25",
			DurationSeconds: 0,
			BillableRate:    1000,
			BillableAmount:  0,
			Comment:         "Timelog 3",
			CreatedAt:       time.Date(2025, 2, 3, 11, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKGJ7NE8X5CSKE4YYWJXRAXY",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID,
			ProjectName:     "Project 1",
			TaskID:          func() *string { s := taskID; return &s }(),
			TaskName:        func() *string { s := "Task 1"; return &s }(),
			TaskCompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
			Date:            "2025-02-03",
			TimeStart:       "10:13:25",
			TimeEnd:         func() *string { s := "11:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 2",
			CreatedAt:       time.Date(2025, 2, 3, 10, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKGJ7SCBRWBPYYV7GEVXTW9N",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID,
			ProjectName:     "Project 1",
			Date:            "2025-02-03",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog `",
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
	}

	timelogsInsert(t, ctx, db, timelogs)
	defer timelogsDelete(t, ctx, db, userID)

	// Test

	store := &TimelogRepository{db: db}

	got, err := store.FindRunning(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, timelogs[0], got)
}

func TestIntegration_TimelogRepository_FindLastN(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID     = "01JKGJ5MFGG1XCZKPAP9PA8PQ2"
		clientID   = "01JKGJ5R8P8WNJH1KBKH3T9B0P"
		projectID1 = "01JKGJ5VQRHRJRTR81TJ7197WF"
		projectID2 = "01JKGJ5Z5CBA9A9P0FTKR5S0YN"
		taskID     = "01JKGJ62K5Y9YY22F5CDJJPYG7"
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
		},
	}

	projectsInsert(t, ctx, db, projects)
	defer projectsDelete(t, ctx, db, userID)

	tasks := []domain.Task{
		{
			ID:          taskID,
			UserID:      userID,
			ClientID:    clientID,
			ProjectID:   projectID2,
			Name:        "Task 1",
			Comment:     "Comment 1",
			CreatedAt:   time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
			CompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
		},
	}

	tasksInsert(t, ctx, db, tasks)
	defer tasksDelete(t, ctx, db, userID)

	timelogs := []domain.Timelog{
		// 2025-02-03
		{
			ID:              "01JKGJ67YH9SD21WWZPYX3D5VQ",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID2,
			ProjectName:     "Project 2",
			Date:            "2025-02-03",
			TimeStart:       "11:13:25",
			DurationSeconds: 0,
			BillableRate:    1000,
			BillableAmount:  0,
			Comment:         "Timelog 6",
			CreatedAt:       time.Date(2025, 2, 3, 11, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKGJ6BJ1TXN023Q0BNMHK3T1",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID1,
			ProjectName:     "Project 1",
			TaskID:          func() *string { s := taskID; return &s }(),
			TaskName:        func() *string { s := "Task 1"; return &s }(),
			TaskCompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
			Date:            "2025-02-03",
			TimeStart:       "10:13:25",
			TimeEnd:         func() *string { s := "11:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 5",
			CreatedAt:       time.Date(2025, 2, 3, 10, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKGJ6FVSVRGNG2ZWBWR57474",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID1,
			ProjectName:     "Project 1",
			Date:            "2025-02-03",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 4",
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
		// 2025-02-02
		{
			ID:              "01JKGJ6M6KXWRN5Q0C4K9WECA0",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID2,
			ProjectName:     "Project 2",
			Date:            "2025-02-02",
			TimeStart:       "11:13:25",
			TimeEnd:         func() *string { s := "12:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 3",
			CreatedAt:       time.Date(2025, 2, 2, 11, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKGJ6RWQENZZ05ZDKC5N52CA",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID1,
			ProjectName:     "Project 1",
			TaskID:          func() *string { s := taskID; return &s }(),
			TaskName:        func() *string { s := "Task 1"; return &s }(),
			TaskCompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC); return &t }(),
			Date:            "2025-02-02",
			TimeStart:       "10:13:25",
			TimeEnd:         func() *string { s := "11:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 2",
			CreatedAt:       time.Date(2025, 2, 2, 10, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JKGJ6WJVJXM0TA919H7CJ0C3",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID1,
			ProjectName:     "Project 1",
			Date:            "2025-02-02",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			Comment:         "Timelog 1",
			CreatedAt:       time.Date(2025, 2, 2, 9, 13, 25, 0, time.UTC),
		},
	}

	timelogsInsert(t, ctx, db, timelogs)
	defer timelogsDelete(t, ctx, db, userID)

	// Test

	store := &TimelogRepository{db: db}

	var (
		got []domain.Timelog
		err error
	)

	got, err = store.FindLastN(ctx, userID, 4)
	assert.NoError(t, err)
	assert.ElementsMatch(t, timelogs[:4], got)
}

func TestIntegration_TimelogRepository_TaskReport(t *testing.T) {
	ctx := t.Context()

	db := dbConnect(t, ctx)
	defer db.Close()

	const (
		userID    = "01JMGT4B57ZVSG0RN9C9QWA3DV"
		clientID  = "01JMGT4F04N6VJW9XP7VNNKVWG"
		projectID = "01JMGT4JM98E0KDMEAW6QWTAQ0"
		taskID    = "01JMGT4P13AHJY5MFVJZ8HF0K4"
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

	timelogs := []domain.Timelog{
		{
			ID:              "01JMGT7V45Z4GXDF843SPDJW4B",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID,
			ProjectName:     "Project 1",
			TaskID:          func() *string { s := taskID; return &s }(),
			TaskName:        func() *string { s := "Task 1"; return &s }(),
			Date:            "2025-02-03",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JMGT7YWSK8DX3D79NPN3VHPM",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID,
			ProjectName:     "Project 1",
			TaskID:          func() *string { s := taskID; return &s }(),
			TaskName:        func() *string { s := "Task 1"; return &s }(),
			Date:            "2025-02-04",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JMGT82DF57ES0YP09YNRRABX",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID,
			ProjectName:     "Project 1",
			Date:            "2025-02-05",
			TimeStart:       "09:13:25",
			TimeEnd:         func() *string { s := "10:13:25"; return &s }(),
			DurationSeconds: 3600,
			BillableRate:    1000,
			BillableAmount:  1000,
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
		{
			ID:              "01JMGT86GDRTWHKFA0DQTJJ6R1",
			UserID:          userID,
			ClientID:        clientID,
			ClientName:      "Client 1",
			ProjectID:       projectID,
			ProjectName:     "Project 1",
			TaskID:          func() *string { s := taskID; return &s }(),
			TaskName:        func() *string { s := "Task 1"; return &s }(),
			Date:            "2025-02-06",
			TimeStart:       "09:13:25",
			DurationSeconds: 0,
			BillableRate:    1000,
			BillableAmount:  0,
			CreatedAt:       time.Date(2025, 2, 3, 9, 13, 25, 0, time.UTC),
		},
	}

	timelogsInsert(t, ctx, db, timelogs)
	defer timelogsDelete(t, ctx, db, userID)

	// Test

	store := &TimelogRepository{db: db}

	// Found
	got, err := store.TaskReport(ctx, userID, taskID)
	assert.NoError(t, err)
	assert.Equal(t, domain.TaskReport{
		DurationSeconds:  7200,
		BillableAmount:   2000,
		UniqueDatesCount: 2,
	}, got)

	// Not found
	got, err = store.TaskReport(ctx, userID, "01JMGT86GDRTWHKFA0DQTJJ6R1")
	assert.NoError(t, err)
	assert.Equal(t, domain.TaskReport{}, got)
}

func timelogsInsert(t *testing.T, ctx context.Context, db *sqlx.DB, timelogs []domain.Timelog) {
	t.Helper()

	q := `INSERT INTO "timelogs"
		("id", "user_id", "project_id", "task_id", "date", "time_start", "time_end", "duration_seconds", "billable_rate", "billable_amount", "comment", "created_at", "updated_at")
		VALUES
		(:id, :user_id, :project_id, :task_id, :date, :time_start, :time_end, :duration_seconds, :billable_rate, :billable_amount, :comment, :created_at, :updated_at)`

	for _, c := range timelogs {
		_, err := db.NamedExecContext(ctx, q, timelogModelFromDomain(c))
		assert.NoError(t, err)
	}
}

func timelogsSelect(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) ([]domain.Timelog, error) {
	t.Helper()

	ms := []timelogModel{}

	q := `SELECT
			"id",
			"user_id",
			"project_id",
			"task_id",
			"date",
			"time_start",
			"time_end",
			"duration_seconds",
			"billable_rate",
			"billable_amount",
			"comment",
			"created_at",
			"updated_at"
		FROM "timelogs"
		WHERE "user_id" = ?`

	if err := db.SelectContext(ctx, &ms, q, userID); err != nil {
		return nil, fmt.Errorf("select timelogs by id: %v", err)
	}

	return timelogModelsToDomains(ms), nil
}

func timelogsDelete(t *testing.T, ctx context.Context, db *sqlx.DB, userID string) {
	t.Helper()

	q := `DELETE FROM "timelogs" WHERE "user_id" = :user_id`

	p := map[string]any{
		"user_id": userID,
	}

	_, err := db.NamedExecContext(ctx, q, p)
	assert.NoError(t, err)
}
