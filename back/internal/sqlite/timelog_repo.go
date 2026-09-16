package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"capuchin/internal/domain"
)

// TimelogRepository SQLite реализация репозитория записи времени.
type TimelogRepository struct {
	db *sqlx.DB
}

// NewTimelogRepository возвращает новый TimelogRepository.
func NewTimelogRepository(db *sqlx.DB) *TimelogRepository {
	return &TimelogRepository{
		db: db,
	}
}

// Create создаёт запись времени.
func (*TimelogRepository) Create(ctx context.Context, db SQLExecutor, timelog domain.Timelog) error {
	q := `-- timelogs::Create
		INSERT INTO "timelogs"
		("id", "user_id", "project_id", "date", "time_start", "time_end", "duration_seconds", "billable_rate", "billable_amount", "comment", "created_at", "task_id")
		VALUES
		(:id, :user_id, :project_id, :date, :time_start, :time_end, :duration_seconds, :billable_rate, :billable_amount, :comment, :created_at, :task_id)`

	if _, err := db.NamedExecContext(ctx, q, timelogModelFromDomain(timelog)); err != nil {
		return fmt.Errorf("insert timelog: %v", err)
	}

	return nil
}

// Update обновляет запись времени.
func (*TimelogRepository) Update(ctx context.Context, db SQLExecutor, timelog domain.Timelog) error {
	q := `-- timelogs::Update
		UPDATE "timelogs" SET
			"project_id" = :project_id,
			"date" = :date,
			"time_start" = :time_start,
			"time_end" = :time_end,
			"duration_seconds" = :duration_seconds,
			"billable_rate" = :billable_rate,
			"billable_amount" = :billable_amount,
			"comment" = :comment,
			"updated_at" = :updated_at,
			"task_id" = :task_id
		WHERE "user_id" = :user_id
			AND "id" = :id`

	if _, err := db.NamedExecContext(ctx, q, timelogModelFromDomain(timelog)); err != nil {
		return fmt.Errorf("update timelog: %v", err)
	}

	return nil
}

// Delete удаляет запись времени.
func (*TimelogRepository) Delete(ctx context.Context, db SQLExecutor, userID string, id string) error {
	q := `-- timelogs::Delete
		DELETE FROM "timelogs"
		WHERE "user_id" = :user_id
			AND "id" = :id`

	m := map[string]any{
		"user_id": userID,
		"id":      id,
	}

	if _, err := db.NamedExecContext(ctx, q, m); err != nil {
		return fmt.Errorf("delete timelog: %v", err)
	}

	return nil
}

// FindByID возвращает запись времени по ID.
func (r *TimelogRepository) FindByID(ctx context.Context, userID string, id string) (domain.Timelog, error) {
	t := timelogModel{}

	q := `-- timelogs::FindByID
		SELECT
			"tl"."id",
			"tl"."user_id",
			"tl"."date",
			"tl"."time_start",
			"tl"."time_end",
			"tl"."duration_seconds",
			"tl"."billable_rate",
			"tl"."billable_amount",
			"tl"."comment",
			"tl"."created_at",
			"tl"."updated_at",
			"pr"."id" "project_id",
			"pr"."name" "project_name",
			"cl"."id" "client_id",
			"cl"."name" "client_name",
			"tsk"."id" "task_id",
			"tsk"."name" "task_name",
			"tsk"."completed_at" "task_completed_at"
		FROM "timelogs" "tl"
		INNER JOIN "projects" "pr"
			ON "pr"."id" = "tl"."project_id"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "pr"."client_id"
		LEFT JOIN "tasks" "tsk"
			ON "tsk"."id" = "tl"."task_id"
		WHERE "tl"."user_id" = ?
			AND "tl"."id" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &t, q, userID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Timelog{}, domain.ErrNotFound
		}

		return domain.Timelog{}, fmt.Errorf("select timelog by id: %v", err)
	}

	return timelogModelToDomain(t), nil
}

// FindByFilter возвращает записи времени по фильтру.
func (r *TimelogRepository) FindByFilter(ctx context.Context, userID string, filter domain.TimelogFilter) ([]domain.Timelog, error) {
	tl := []timelogModel{}

	args := map[string]any{
		"user_id":   userID,
		"date_from": filter.DateFrom,
		"date_to":   filter.DateTo,
	}

	client := ""
	if filter.ClientID != "" {
		client = `AND "pr"."client_id" = :client_id`
		args["client_id"] = filter.ClientID
	}

	project := ""
	if filter.ProjectID != "" {
		project = `AND "tl"."project_id" = :project_id`
		args["project_id"] = filter.ProjectID
	}

	q := `-- timelogs::FindByFilter
		SELECT
			"tl"."id",
			"tl"."user_id",
			"tl"."date",
			"tl"."time_start",
			"tl"."time_end",
			"tl"."duration_seconds",
			"tl"."billable_rate",
			"tl"."billable_amount",
			"tl"."comment",
			"tl"."created_at",
			"tl"."updated_at",
			"pr"."id" "project_id",
			"pr"."name" "project_name",
			"cl"."id" "client_id",
			"cl"."name" "client_name",
			"tsk"."id" "task_id",
			"tsk"."name" "task_name",
			"tsk"."completed_at" "task_completed_at"
		FROM "timelogs" "tl"
		INNER JOIN "projects" "pr"
			ON "pr"."id" = "tl"."project_id"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "pr"."client_id"
		LEFT JOIN "tasks" "tsk"
			ON "tsk"."id" = "tl"."task_id"
		WHERE "tl"."user_id" = :user_id
			AND "tl"."date" >= :date_from AND "tl"."date" <= :date_to
			` + client + `
			` + project + `
		ORDER BY "tl"."date" DESC,
			"tl"."time_start" DESC`

	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("prepare named: %v", err)
	}

	if err = stmt.SelectContext(ctx, &tl, args); err != nil {
		return nil, fmt.Errorf("select timelogs: %v", err)
	}

	return timelogModelsToDomains(tl), nil
}

// FindRunning возвращает запущенную запись времени.
func (r *TimelogRepository) FindRunning(ctx context.Context, userID string) (domain.Timelog, error) {
	t := timelogModel{}

	q := `-- timelogs::FindRunning
		SELECT
			"tl"."id",
			"tl"."user_id",
			"tl"."date",
			"tl"."time_start",
			"tl"."time_end",
			"tl"."duration_seconds",
			"tl"."billable_rate",
			"tl"."billable_amount",
			"tl"."comment",
			"tl"."created_at",
			"tl"."updated_at",
			"pr"."id" "project_id",
			"pr"."name" "project_name",
			"cl"."id" "client_id",
			"cl"."name" "client_name",
			"tsk"."id" "task_id",
			"tsk"."name" "task_name",
			"tsk"."completed_at" "task_completed_at"
		FROM "timelogs" "tl"
		INNER JOIN "projects" "pr"
			ON "pr"."id" = "tl"."project_id"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "pr"."client_id"
		LEFT JOIN "tasks" "tsk"
			ON "tsk"."id" = "tl"."task_id"
		WHERE "tl"."user_id" = ?
			AND "tl"."time_end" IS NULL
		LIMIT 1`

	if err := r.db.GetContext(ctx, &t, q, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Timelog{}, domain.ErrNotFound
		}

		return domain.Timelog{}, fmt.Errorf("select running timelog: %v", err)
	}

	return timelogModelToDomain(t), nil
}

// FindLastN возвращает последние N записей времени.
func (r *TimelogRepository) FindLastN(ctx context.Context, userID string, n int) ([]domain.Timelog, error) {
	tl := []timelogModel{}

	args := map[string]any{
		"user_id": userID,
		"limit":   n,
	}

	q := `-- timelogs::FindLastN
		SELECT
			"tl"."id",
			"tl"."user_id",
			"tl"."date",
			"tl"."time_start",
			"tl"."time_end",
			"tl"."duration_seconds",
			"tl"."billable_rate",
			"tl"."billable_amount",
			"tl"."comment",
			"tl"."created_at",
			"tl"."updated_at",
			"pr"."id" "project_id",
			"pr"."name" "project_name",
			"cl"."id" "client_id",
			"cl"."name" "client_name",
			"tsk"."id" "task_id",
			"tsk"."name" "task_name",
			"tsk"."completed_at" "task_completed_at"
		FROM "timelogs" "tl"
		INNER JOIN "projects" "pr"
			ON "pr"."id" = "tl"."project_id"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "pr"."client_id"
		LEFT JOIN "tasks" "tsk"
			ON "tsk"."id" = "tl"."task_id"
		WHERE "tl"."user_id" = :user_id
		ORDER BY "tl"."date" DESC,
			"tl"."time_start" DESC
		LIMIT :limit`

	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("prepare named: %v", err)
	}

	if err = stmt.SelectContext(ctx, &tl, args); err != nil {
		return nil, fmt.Errorf("select timelogs: %v", err)
	}

	return timelogModelsToDomains(tl), nil
}

// TaskReport возвращает отчёт по задаче.
func (r *TimelogRepository) TaskReport(ctx context.Context, userID string, taskID string) (domain.TaskReport, error) {
	m := taskReportModel{}

	q := `-- timelogs::TaskReport
		SELECT
			SUM("duration_seconds") "duration_seconds",
			SUM("billable_amount") "billable_amount",
    		COUNT(DISTINCT "date") AS "unique_dates_count"
		FROM "timelogs"
		WHERE "user_id" = ?
			AND "task_id" = ?
			AND "time_end" IS NOT NULL
		GROUP BY "task_id"
		LIMIT 1`

	if err := r.db.GetContext(ctx, &m, q, userID, taskID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TaskReport{}, nil
		}

		return domain.TaskReport{}, fmt.Errorf("select task report by task id: %v", err)
	}

	return taskReportModelToDomain(m), nil
}
