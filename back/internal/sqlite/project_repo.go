package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"capuchin/internal/domain"
)

// ProjectRepository SQLite реализация репозитория проекта.
type ProjectRepository struct {
	db *sqlx.DB
}

// NewProjectRepository возвращает новый ProjectRepository.
func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

// Create создаёт проект.
func (*ProjectRepository) Create(ctx context.Context, db SQLExecutor, project domain.Project) error {
	q := `-- projects::Create
		INSERT INTO "projects"
		("id", "user_id", "client_id", "name", "comment", "billable_rate", "created_at")
		VALUES
		(:id, :user_id, :client_id, :name, :comment, :billable_rate, :created_at)`

	if _, err := db.NamedExecContext(ctx, q, projectModelFromDomain(project)); err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("insert project: %v", err)
	}

	return nil
}

// Update обновляет проект.
func (*ProjectRepository) Update(ctx context.Context, db SQLExecutor, project domain.Project) error {
	q := `-- projects::Update
		UPDATE "projects" SET
			"client_id" = :client_id,
			"name" = :name,
			"billable_rate" = :billable_rate,
			"comment" = :comment,
			"updated_at" = :updated_at,
			"archived_at" = :archived_at
		WHERE "user_id" = :user_id
			AND "id" = :id`

	if _, err := db.NamedExecContext(ctx, q, projectModelFromDomain(project)); err != nil {
		return fmt.Errorf("update project: %v", err)
	}

	return nil
}

// FindByFilter возвращает проекты по фильтру.
func (r *ProjectRepository) FindByFilter(ctx context.Context, userID string, filter domain.ProjectFilter) ([]domain.Project, error) {
	pl := []projectModel{}

	args := map[string]any{
		"user_id": userID,
	}

	withoutArchivedClients := ""
	if filter.WithoutArchivedClients {
		withoutArchivedClients = `AND "cl"."archived_at" IS NULL`
	}

	client := ""
	if filter.ClientID != "" {
		client = `AND "p"."client_id" = :client_id`
		args["client_id"] = filter.ClientID
	}

	q := `-- projects::FindByFilter
		SELECT
			"p"."id",
			"p"."user_id",
			"p"."client_id",
			"p"."name",
			"p"."billable_rate",
			"p"."comment",
			"p"."created_at",
			"p"."updated_at",
			"p"."archived_at",
			"cl"."name" "client_name"
		FROM "projects" "p"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "p"."client_id"
				` + withoutArchivedClients + `
		WHERE "p"."user_id" = :user_id
			` + client + `
		ORDER BY ("p"."archived_at" IS NULL) DESC,
			"p"."archived_at" DESC,
			"p"."name" ASC`

	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("prepare named: %v", err)
	}

	if err = stmt.SelectContext(ctx, &pl, args); err != nil {
		return nil, fmt.Errorf("select projects: %v", err)
	}

	return projectModelsToDomains(pl), nil
}

// FindByID возвращает проект по ID.
func (r *ProjectRepository) FindByID(ctx context.Context, userID string, id string) (domain.Project, error) {
	p := projectModel{}

	q := `-- projects::FindByID
		SELECT
			"p"."id",
			"p"."user_id",
			"p"."client_id",
			"p"."name",
			"p"."billable_rate",
			"p"."comment",
			"p"."created_at",
			"p"."updated_at",
			"p"."archived_at",
			"cl"."name" "client_name"
		FROM "projects" "p"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "p"."client_id"
		WHERE "p"."user_id" = ?
			AND "p"."id" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &p, q, userID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, domain.ErrNotFound
		}

		return domain.Project{}, fmt.Errorf("select project by id: %v", err)
	}

	return projectModelToDomain(p), nil
}
