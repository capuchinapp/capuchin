package repository

import (
	"capuchin/app/entity"

	"github.com/google/uuid"
)

type TimelogRepository interface {
	Create(t *entity.Timelog) (int64, error)
	Update(t *entity.Timelog) (int64, error)
	Delete(t *entity.Timelog) (int64, error)
	FindById(uuid uuid.UUID) (entity.Timelog, error)
	FindByPeriod(dateFrom string, dateTo string) ([]entity.Timelog, error)
	FindByPeriodAndClientId(dateFrom string, dateTo string, clientUuid uuid.UUID) ([]entity.Timelog, error)
	FindByPeriodAndProjectId(dateFrom string, dateTo string, projectUuid uuid.UUID) ([]entity.Timelog, error)
	FindRunning() (entity.Timelog, error)
}
