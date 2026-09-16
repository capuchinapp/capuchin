package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"capuchin/internal/domain"
)

// TaskRepository SQLite реализация репозитория задач.
type TaskRepository struct {
	db *sqlx.DB
}

// NewTaskRepository возвращает новый TaskRepository.
func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

// Create создаёт задачу.
func (*TaskRepository) Create(ctx context.Context, db SQLExecutor, task domain.Task) error {
	q := `-- tasks::Create
		INSERT INTO "tasks"
		("id", "user_id", "project_id", "name", "comment", "created_at")
		VALUES
		(:id, :user_id, :project_id, :name, :comment, :created_at)`

	if _, err := db.NamedExecContext(ctx, q, taskModelFromDomain(task)); err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("insert task: %v", err)
	}

	return nil
}

// Update обновляет задачу.
func (*TaskRepository) Update(ctx context.Context, db SQLExecutor, task domain.Task) error {
	q := `-- tasks::Update
		UPDATE "tasks" SET
			"project_id" = :project_id,
			"name" = :name,
			"comment" = :comment,
			"updated_at" = :updated_at,
			"archived_at" = :archived_at,
			"completed_at" = :completed_at
		WHERE "user_id" = :user_id
			AND "id" = :id`

	if _, err := db.NamedExecContext(ctx, q, taskModelFromDomain(task)); err != nil {
		return fmt.Errorf("update task: %v", err)
	}

	return nil
}

// FindByFilter возвращает задачи по фильтру.
func (r *TaskRepository) FindByFilter(ctx context.Context, userID string, filter domain.TaskFilter) ([]domain.Task, error) {
	pl := []taskModel{}

	args := map[string]any{
		"user_id": userID,
	}

	withoutArchivedProjects := ""
	if filter.WithoutArchivedProjects {
		withoutArchivedProjects = `AND "p"."archived_at" IS NULL`
	}

	project := ""
	if filter.ProjectID != "" {
		project = `AND "t"."project_id" = :project_id`
		args["project_id"] = filter.ProjectID
	}

	q := `-- tasks::FindByFilter
		SELECT
			"t"."id",
			"t"."user_id",
			"t"."name",
			"t"."comment",
			"t"."created_at",
			"t"."updated_at",
			"t"."archived_at",
			"t"."completed_at",
			"p"."id" "project_id",
			"p"."name" "project_name",
			"cl"."id" "client_id",
			"cl"."name" "client_name"
		FROM "tasks" "t"
		INNER JOIN "projects" "p"
			ON "p"."id" = "t"."project_id"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "p"."client_id"
				` + withoutArchivedProjects + `
		WHERE "t"."user_id" = :user_id
			` + project + `
		ORDER BY ("t"."archived_at" IS NULL) DESC,
			"t"."archived_at" DESC,
			("t"."completed_at" IS NULL) DESC,
			"t"."completed_at" DESC,
			"t"."name" ASC` // Сортировка открытых задач по name, архивные и завершённые сортируются по дате действия, так удобнее.

	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("prepare named: %v", err)
	}

	if err = stmt.SelectContext(ctx, &pl, args); err != nil {
		return nil, fmt.Errorf("select tasks: %v", err)
	}

	return taskModelsToDomains(pl), nil
}

// FindByID возвращает задачу по ID.
func (r *TaskRepository) FindByID(ctx context.Context, userID string, id string) (domain.Task, error) {
	p := taskModel{}

	q := `-- tasks::FindByID
		SELECT
			"t"."id",
			"t"."user_id",
			"t"."name",
			"t"."comment",
			"t"."created_at",
			"t"."updated_at",
			"t"."archived_at",
			"t"."completed_at",
			"p"."id" "project_id",
			"p"."name" "project_name",
			"cl"."id" "client_id",
			"cl"."name" "client_name"
		FROM "tasks" "t"
		INNER JOIN "projects" "p"
			ON "p"."id" = "t"."project_id"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "p"."client_id"
		WHERE "t"."user_id" = ?
			AND "t"."id" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &p, q, userID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, domain.ErrNotFound
		}

		return domain.Task{}, fmt.Errorf("select task by id: %v", err)
	}

	return taskModelToDomain(p), nil
}
