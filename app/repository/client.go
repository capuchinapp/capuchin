package repository

import (
	"capuchin/app/entity"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Create(c *entity.Client) (int64, error)
	Update(c *entity.Client) (int64, error)
	FindAll() ([]entity.Client, error)
	FindById(uuid uuid.UUID) (entity.Client, error)
}
