package sqlite

import (
	"capuchin/app/entity"
	"capuchin/app/util"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r SessionRepository) Create(s *entity.Session) (int64, error) {
	q := `INSERT INTO sessions
		(uuid, checked_at)
		VALUES
		(:uuid, :checked_at)`
	res, err := r.db.NamedExec(q, s)
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

func (r SessionRepository) DeleteByUuid(uuid uuid.UUID) (int64, error) {
	m := map[string]interface{}{"uuid": uuid}
	res, err := r.db.NamedExec("DELETE FROM sessions WHERE uuid = :uuid", m)
	if err != nil {
		return 0, util.ErrTrace("NamedExec", err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return 0, util.ErrTrace("RowsAffected", err)
	}

	if cnt == 0 {
		cnt = 1
	}

	return cnt, nil
}

func (r SessionRepository) DeleteByRowId(i int64) (int64, error) {
	m := map[string]interface{}{"rowid": i}
	res, err := r.db.NamedExec("DELETE FROM sessions WHERE rowid = :rowid", m)
	if err != nil {
		return 0, util.ErrTrace("NamedExec", err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return 0, util.ErrTrace("RowsAffected", err)
	}

	if cnt == 0 {
		cnt = 1
	}

	return cnt, nil
}

func (r SessionRepository) FindAll() ([]entity.Session, error) {
	sl := []entity.Session{}

	if err := r.db.Select(&sl, "SELECT rowid, * FROM sessions ORDER BY checked_at DESC"); err != nil {
		return nil, util.ErrTrace("Select", err)
	}

	return sl, nil
}

func (r SessionRepository) FindByUuid(uuid uuid.UUID) (entity.Session, error) {
	s := entity.Session{}

	if err := r.db.Get(&s, "SELECT rowid, * FROM sessions WHERE uuid = ? LIMIT 1", uuid); err != nil {
		return entity.Session{}, util.ErrTrace("Get", err)
	}

	return s, nil
}

func (r SessionRepository) FindByRowId(i int64) (entity.Session, error) {
	s := entity.Session{}

	if err := r.db.Get(&s, "SELECT rowid, * FROM sessions WHERE rowid = ? LIMIT 1", i); err != nil {
		return entity.Session{}, util.ErrTrace("Get", err)
	}

	return s, nil
}

func (r SessionRepository) UpdateCheckedAt(s *entity.Session) (int64, error) {
	q := `UPDATE sessions SET
		checked_at = :checked_at
		WHERE uuid = :uuid`
	res, err := r.db.NamedExec(q, s)
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

func (r SessionRepository) CancelOld(dt string) (int64, error) {
	m := map[string]interface{}{"dt": dt}
	res, err := r.db.NamedExec("DELETE FROM sessions WHERE checked_at < :dt", m)
	if err != nil {
		return 0, util.ErrTrace("NamedExec", err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return 0, util.ErrTrace("RowsAffected", err)
	}

	return cnt, nil
}
