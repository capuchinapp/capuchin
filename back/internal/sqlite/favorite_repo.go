package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"capuchin/internal/domain"
)

// FavoriteRepository SQLite реализация репозитория избранного.
type FavoriteRepository struct {
	db *sqlx.DB
}

// NewFavoriteRepository возвращает новый FavoriteRepository.
func NewFavoriteRepository(db *sqlx.DB) *FavoriteRepository {
	return &FavoriteRepository{
		db: db,
	}
}

// Create создаёт избранное.
func (*FavoriteRepository) Create(ctx context.Context, db SQLExecutor, favorite domain.Favorite) error {
	q := `-- favorites::Create
		INSERT INTO "favorites"
		("id", "user_id", "name", "project_id", "task_id", "billable_rate", "comment")
		VALUES
		(:id, :user_id, :name, :project_id, :task_id, :billable_rate, :comment)`

	if _, err := db.NamedExecContext(ctx, q, favoriteModelFromDomain(favorite)); err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("insert favorite: %v", err)
	}

	return nil
}

// Update обновляет избранное.
func (*FavoriteRepository) Update(ctx context.Context, db SQLExecutor, favorite domain.Favorite) error {
	q := `-- favorites::Update
		UPDATE "favorites" SET
			"name" = :name,
			"project_id" = :project_id,
			"task_id" = :task_id,
			"billable_rate" = :billable_rate,
			"comment" = :comment
		WHERE "user_id" = :user_id
			AND "id" = :id`

	if _, err := db.NamedExecContext(ctx, q, favoriteModelFromDomain(favorite)); err != nil {
		return fmt.Errorf("update favorite: %v", err)
	}

	return nil
}

// Delete удаляет избранное.
func (*FavoriteRepository) Delete(ctx context.Context, db SQLExecutor, userID string, id string) error {
	q := `-- favorites::Delete
		DELETE FROM "favorites"
		WHERE "user_id" = :user_id
			AND "id" = :id`

	m := map[string]any{
		"user_id": userID,
		"id":      id,
	}

	if _, err := db.NamedExecContext(ctx, q, m); err != nil {
		return fmt.Errorf("delete favorite: %v", err)
	}

	return nil
}

// FindByID возвращает избранное по ID.
func (r *FavoriteRepository) FindByID(ctx context.Context, userID string, id string) (domain.Favorite, error) {
	item := favoriteModel{}

	q := `-- favorites::FindByID
		SELECT
			"f"."id",
			"f"."user_id",
			"f"."name",
			"f"."billable_rate",
			"f"."comment",
			"pr"."id" "project_id",
			"pr"."name" "project_name",
			"cl"."id" "client_id",
			"cl"."name" "client_name",
			"tsk"."id" "task_id",
			"tsk"."name" "task_name"
		FROM "favorites" "f"
		INNER JOIN "projects" "pr"
			ON "pr"."id" = "f"."project_id"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "pr"."client_id"
		LEFT JOIN "tasks" "tsk"
			ON "tsk"."id" = "f"."task_id"
		WHERE "f"."user_id" = ?
			AND "f"."id" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &item, q, userID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Favorite{}, domain.ErrNotFound
		}

		return domain.Favorite{}, fmt.Errorf("select favorite by id: %v", err)
	}

	return favoriteModelToDomain(item), nil
}

// FindAll возвращает список избранного.
func (r *FavoriteRepository) FindAll(ctx context.Context, userID string) ([]domain.Favorite, error) {
	items := []favoriteModel{}

	args := map[string]any{
		"user_id": userID,
	}

	q := `-- favorites::FindAll
		SELECT
			"f"."id",
			"f"."user_id",
			"f"."name",
			"f"."billable_rate",
			"f"."comment",
			"pr"."id" "project_id",
			"pr"."name" "project_name",
			"cl"."id" "client_id",
			"cl"."name" "client_name",
			"tsk"."id" "task_id",
			"tsk"."name" "task_name"
		FROM "favorites" "f"
		INNER JOIN "projects" "pr"
			ON "pr"."id" = "f"."project_id"
		INNER JOIN "clients" "cl"
			ON "cl"."id" = "pr"."client_id"
		LEFT JOIN "tasks" "tsk"
			ON "tsk"."id" = "f"."task_id"
		WHERE "f"."user_id" = :user_id
		ORDER BY "f"."name"`

	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("prepare named: %v", err)
	}

	if err = stmt.SelectContext(ctx, &items, args); err != nil {
		return nil, fmt.Errorf("select favorites: %v", err)
	}

	return favoriteModelsToDomains(items), nil
}
