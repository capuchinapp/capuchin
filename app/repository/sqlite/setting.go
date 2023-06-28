package sqlite

import (
	"capuchin/app/entity"
	"capuchin/app/util"
	"errors"

	"github.com/jmoiron/sqlx"
)

type SettingRepository struct {
	db *sqlx.DB
}

func NewSettingRepository(db *sqlx.DB) *SettingRepository {
	return &SettingRepository{
		db: db,
	}
}

func (r SettingRepository) Update(s *entity.Setting) (int64, error) {
	res, err := r.db.NamedExec("UPDATE settings SET value = :value WHERE key = :key", s)
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

func (r SettingRepository) FindAll() ([]entity.Setting, error) {
	sl := []entity.Setting{}

	if err := r.db.Select(&sl, "SELECT * FROM settings"); err != nil {
		return nil, util.ErrTrace("Select", err)
	}

	return sl, nil
}

func (r SettingRepository) FindByKey(key string) (entity.Setting, error) {
	s := entity.Setting{}

	if err := r.db.Get(&s, "SELECT * FROM settings WHERE key = ? LIMIT 1", key); err != nil {
		return entity.Setting{}, util.ErrTrace("Get", err)
	}

	return s, nil
}
