package sqlite

import (
	"capuchin/app/entity"
	"capuchin/app/util"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

func (r ProjectRepository) Create(p *entity.Project) (int64, error) {
	q := `INSERT INTO projects
		(uuid, client_uuid, name, billable_rate, comment, created_at, updated_at, archived_at)
		VALUES
		(:uuid, :client_uuid, :name, :billable_rate, :comment, :created_at, NULL, NULL)`
	res, err := r.db.NamedExec(q, p)
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

func (r ProjectRepository) Update(p *entity.Project) (int64, error) {
	q := `UPDATE projects SET
		client_uuid = :client_uuid,
		name = :name,
		billable_rate = :billable_rate,
		comment = :comment,
		updated_at = :updated_at,
		archived_at = :archived_at
		WHERE uuid = :uuid`
	res, err := r.db.NamedExec(q, p)
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

func (r ProjectRepository) FindAll(f entity.ProjectFilter) ([]entity.Project, error) {
	pl := []entity.Project{}

	woac := ""
	if f.WithoutArchivedClients {
		woac = "AND cl.archived_at IS NULL"
	}

	q := "SELECT p.*, cl.name client_name FROM projects p INNER JOIN clients cl ON cl.uuid = p.client_uuid %s ORDER BY (p.archived_at IS NULL) DESC, p.archived_at DESC, p.name ASC"
	if err := r.db.Select(&pl, fmt.Sprintf(q, woac)); err != nil {
		return nil, util.ErrTrace("Select", err)
	}

	return pl, nil
}

func (r ProjectRepository) FindById(uuid uuid.UUID) (entity.Project, error) {
	p := entity.Project{}

	q := `SELECT
		p.*,
		cl.name client_name
		FROM projects p
		INNER JOIN clients cl
			ON cl.uuid = p.client_uuid
		WHERE p.uuid = ?
		LIMIT 1`
	if err := r.db.Get(&p, q, uuid); err != nil {
		return entity.Project{}, util.ErrTrace("Get", err)
	}

	return p, nil
}

func (r ProjectRepository) FindByClientId(clientUUID uuid.UUID) ([]entity.Project, error) {
	pl := []entity.Project{}

	q := `SELECT
		p.*,
		cl.name client_name
		FROM projects p
		INNER JOIN clients cl
			ON cl.uuid = p.client_uuid
		WHERE p.client_uuid = ?
		ORDER BY (p.archived_at IS NULL) DESC, p.archived_at DESC, p.name ASC`
	if err := r.db.Select(&pl, q, clientUUID); err != nil {
		return nil, util.ErrTrace("Select", err)
	}

	return pl, nil
}
