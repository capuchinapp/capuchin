package repository

import (
	"capuchin/app/entity"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(p *entity.Project) (int64, error)
	Update(p *entity.Project) (int64, error)
	FindAll(f entity.ProjectFilter) ([]entity.Project, error)
	FindById(uuid uuid.UUID) (entity.Project, error)
	FindByClientId(clientUUID uuid.UUID) ([]entity.Project, error)
}
