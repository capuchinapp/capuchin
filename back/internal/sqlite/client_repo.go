package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"capuchin/internal/domain"
)

// ClientRepository SQLite реализация репозитория клиента.
type ClientRepository struct {
	db *sqlx.DB
}

// NewClientRepository возвращает новый ClientRepository.
func NewClientRepository(db *sqlx.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

// Create создаёт клиента.
func (*ClientRepository) Create(ctx context.Context, db SQLExecutor, client domain.Client) error {
	q := `-- clients::Create
		INSERT INTO "clients"
		("id", "user_id", "name", "billable_rate", "comment", "created_at")
		VALUES
		(:id, :user_id, :name, :billable_rate, :comment, :created_at)`

	if _, err := db.NamedExecContext(ctx, q, clientModelFromDomain(client)); err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("insert client: %v", err)
	}

	return nil
}

// Update обновляет клиента.
func (*ClientRepository) Update(ctx context.Context, db SQLExecutor, client domain.Client) error {
	q := `-- clients::Update
		UPDATE "clients" SET
			"name" = :name,
			"billable_rate" = :billable_rate,
			"comment" = :comment,
			"updated_at" = :updated_at,
			"archived_at" = :archived_at
		WHERE "user_id" = :user_id
			AND "id" = :id`

	if _, err := db.NamedExecContext(ctx, q, clientModelFromDomain(client)); err != nil {
		return fmt.Errorf("update client: %v", err)
	}

	return nil
}

// FindAll возвращает список клиентов.
func (r *ClientRepository) FindAll(ctx context.Context, userID string) ([]domain.Client, error) {
	cl := []clientModel{}

	q := `-- clients::FindAll
		SELECT
			"id",
			"user_id",
			"name",
			"billable_rate",
			"comment",
			"created_at",
			"updated_at",
			"archived_at"
		FROM "clients"
		WHERE "user_id" = ?
		ORDER BY ("archived_at" IS NULL) DESC,
			"archived_at" DESC,
			"name" ASC`

	if err := r.db.SelectContext(ctx, &cl, q, userID); err != nil {
		return nil, fmt.Errorf("select clients: %v", err)
	}

	return clientModelsToDomains(cl), nil
}

// FindByID возвращает клиента по ID.
func (r *ClientRepository) FindByID(ctx context.Context, userID string, id string) (domain.Client, error) {
	c := clientModel{}

	q := `-- clients::FindByID
		SELECT
			"id",
			"user_id",
			"name",
			"billable_rate",
			"comment",
			"created_at",
			"updated_at",
			"archived_at"
		FROM "clients"
		WHERE "user_id" = ?
			AND "id" = ?
		LIMIT 1`

	if err := r.db.GetContext(ctx, &c, q, userID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Client{}, domain.ErrNotFound
		}

		return domain.Client{}, fmt.Errorf("select client by id: %v", err)
	}

	return clientModelToDomain(c), nil
}
