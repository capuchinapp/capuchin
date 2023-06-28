package sqlite

import (
	"capuchin/app/entity"
	"capuchin/app/util"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TimelogRepository struct {
	db *sqlx.DB
}

func NewTimelogRepository(db *sqlx.DB) *TimelogRepository {
	return &TimelogRepository{
		db: db,
	}
}

func (r TimelogRepository) Create(t *entity.Timelog) (int64, error) {
	q := `INSERT INTO timelogs
		(uuid, project_uuid, date, time_start, time_end, duration_seconds, billable_rate, billable_amount, comment, created_at)
		VALUES
		(:uuid, :project_uuid, :date, :time_start, :time_end, :duration_seconds, :billable_rate, :billable_amount, :comment, :created_at)`
	res, err := r.db.NamedExec(q, t)
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

func (r TimelogRepository) Update(t *entity.Timelog) (int64, error) {
	q := `UPDATE timelogs SET
		project_uuid = :project_uuid,
		date = :date,
		time_start = :time_start,
		time_end = :time_end,
		duration_seconds = :duration_seconds,
		billable_rate = :billable_rate,
		billable_amount = :billable_amount,
		comment = :comment,
		updated_at = :updated_at
		WHERE uuid = :uuid`
	res, err := r.db.NamedExec(q, t)
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

func (r TimelogRepository) Delete(t *entity.Timelog) (int64, error) {
	q := `UPDATE timelogs SET
		deleted_at = :deleted_at
		WHERE uuid = :uuid`
	res, err := r.db.NamedExec(q, t)
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

func (r TimelogRepository) FindById(uuid uuid.UUID) (entity.Timelog, error) {
	t := entity.Timelog{}

	q := `SELECT
			tl.*,
			pr.uuid project_uuid,
			pr.name project_name,
			cl.uuid client_uuid,
			cl.name client_name
		FROM timelogs tl
		INNER JOIN projects pr
			ON pr.uuid = tl.project_uuid
		INNER JOIN clients cl
			ON cl.uuid = pr.client_uuid
		WHERE tl.uuid = ?
			AND tl.deleted_at IS NULL
		LIMIT 1`

	if err := r.db.Get(&t, q, uuid); err != nil {
		return entity.Timelog{}, util.ErrTrace("Get", err)
	}

	return t, nil
}

func (r TimelogRepository) FindByPeriod(dateFrom string, dateTo string) ([]entity.Timelog, error) {
	tl := []entity.Timelog{}

	q := `SELECT
			tl.*,
			pr.uuid project_uuid,
			pr.name project_name,
			cl.uuid client_uuid,
			cl.name client_name
		FROM timelogs tl
		INNER JOIN projects pr
			ON pr.uuid = tl.project_uuid
		INNER JOIN clients cl
			ON cl.uuid = pr.client_uuid
		WHERE tl.date >= ? AND tl.date <= ?
			AND tl.deleted_at IS NULL
		ORDER BY tl.date DESC, tl.time_start DESC`

	if err := r.db.Select(&tl, q, dateFrom, dateTo); err != nil {
		return nil, util.ErrTrace("Select", err)
	}

	return tl, nil
}

func (r TimelogRepository) FindByPeriodAndClientId(dateFrom string, dateTo string, clientUuid uuid.UUID) ([]entity.Timelog, error) {
	tl := []entity.Timelog{}

	q := `SELECT
			tl.*,
			pr.uuid project_uuid,
			pr.name project_name,
			cl.uuid client_uuid,
			cl.name client_name
		FROM timelogs tl
		INNER JOIN projects pr
			ON pr.uuid = tl.project_uuid
		INNER JOIN clients cl
			ON cl.uuid = pr.client_uuid
		WHERE tl.date >= ? AND tl.date <= ?
			AND pr.client_uuid = ?
			AND tl.deleted_at IS NULL
		ORDER BY tl.date DESC, tl.time_start DESC`

	if err := r.db.Select(&tl, q, dateFrom, dateTo, clientUuid); err != nil {
		return nil, util.ErrTrace("Select", err)
	}

	return tl, nil
}

func (r TimelogRepository) FindByPeriodAndProjectId(dateFrom string, dateTo string, projectUuid uuid.UUID) ([]entity.Timelog, error) {
	tl := []entity.Timelog{}

	q := `SELECT
			tl.*,
			pr.uuid project_uuid,
			pr.name project_name,
			cl.uuid client_uuid,
			cl.name client_name
		FROM timelogs tl
		INNER JOIN projects pr
			ON pr.uuid = tl.project_uuid
		INNER JOIN clients cl
			ON cl.uuid = pr.client_uuid
		WHERE tl.date >= ? AND tl.date <= ?
			AND tl.project_uuid = ?
			AND tl.deleted_at IS NULL
		ORDER BY tl.date DESC, tl.time_start DESC`

	if err := r.db.Select(&tl, q, dateFrom, dateTo, projectUuid); err != nil {
		return nil, util.ErrTrace("Select", err)
	}

	return tl, nil
}

func (r TimelogRepository) FindRunning() (entity.Timelog, error) {
	t := entity.Timelog{}

	if err := r.db.Get(&t, "SELECT * FROM timelogs WHERE time_end IS NULL AND deleted_at IS NULL LIMIT 1"); err != nil {
		return entity.Timelog{}, util.ErrTrace("Get", err)
	}

	return t, nil
}
