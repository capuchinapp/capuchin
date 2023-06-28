package repository

import (
	"capuchin/app/entity"
)

type SettingRepository interface {
	Update(s *entity.Setting) (int64, error)
	FindAll() ([]entity.Setting, error)
	FindByKey(key string) (entity.Setting, error)
}
