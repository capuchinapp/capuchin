package entity

import (
	"capuchin/app/util/nullable"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	UUID         uuid.UUID           `json:"uuid" db:"uuid"`
	ClientUUID   uuid.UUID           `json:"clientUUID" db:"client_uuid"`
	ClientName   string              `json:"clientName" db:"client_name"`
	Name         string              `json:"name" db:"name"`
	BillableRate uint64              `json:"billableRate" db:"billable_rate"`
	Comment      nullable.NullString `json:"comment" db:"comment"`
	CreatedAt    time.Time           `json:"createdAt" db:"created_at"`
	UpdatedAt    nullable.NullTime   `json:"updatedAt" db:"updated_at"`
	ArchivedAt   nullable.NullTime   `json:"archivedAt" db:"archived_at"`
}

type ProjectFilter struct {
	WithoutArchivedClients bool `json:"clientArchivedAt"`
}
