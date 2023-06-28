package sqlite

import (
	"capuchin/app/entity"
	"capuchin/app/util"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ClientRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (r ClientRepository) Create(c *entity.Client) (int64, error) {
	q := `INSERT INTO clients
		(uuid, name, billable_rate, comment, created_at, updated_at, archived_at)
		VALUES
		(:uuid, :name, :billable_rate, :comment, :created_at, NULL, NULL)`
	res, err := r.db.NamedExec(q, c)
	if err != nil {
		return 0, util.ErrTrace("NamedExec", err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return 0, util.ErrTrace("RowsAffected", err)
	}

	if cnt == 0 {
		return 0, errors.New("The query affected 0 rows")
	}

	return cnt, nil
}

func (r ClientRepository) Update(c *entity.Client) (int64, error) {
	q := `UPDATE clients SET
		name = :name,
		billable_rate = :billable_rate,
		comment = :comment,
		updated_at = :updated_at,
		archived_at = :archived_at
		WHERE uuid = :uuid`
	res, err := r.db.NamedExec(q, c)
	if err != nil {
		return 0, util.ErrTrace("NamedExec", err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return 0, util.ErrTrace("RowsAffected", err)
	}

	if cnt == 0 {
		return 0, errors.New("The query affected 0 rows")
	}

	return cnt, nil
}

func (r ClientRepository) FindAll() ([]entity.Client, error) {
	cl := []entity.Client{}

	q := `SELECT * FROM clients
		ORDER BY (archived_at IS NULL) DESC, archived_at DESC, name ASC`
	if err := r.db.Select(&cl, q); err != nil {
		return nil, util.ErrTrace("Select", err)
	}

	return cl, nil
}

func (r ClientRepository) FindById(uuid uuid.UUID) (entity.Client, error) {
	c := entity.Client{}

	if err := r.db.Get(&c, "SELECT * FROM clients WHERE uuid = ? LIMIT 1", uuid); err != nil {
		return entity.Client{}, util.ErrTrace("Get", err)
	}

	return c, nil
}
