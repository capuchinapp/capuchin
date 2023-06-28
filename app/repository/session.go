package repository

import (
	"capuchin/app/entity"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(s *entity.Session) (int64, error)
	DeleteByUuid(uuid uuid.UUID) (int64, error)
	DeleteByRowId(i int64) (int64, error)
	FindAll() ([]entity.Session, error)
	FindByUuid(uuid uuid.UUID) (entity.Session, error)
	FindByRowId(i int64) (entity.Session, error)
	UpdateCheckedAt(s *entity.Session) (int64, error)
	CancelOld(dt string) (int64, error)
}
