package entity

import (
	"capuchin/app/util/nullable"
	"time"

	"github.com/google/uuid"
)

type Timelog struct {
	UUID            uuid.UUID           `json:"uuid" db:"uuid"`
	ProjectUUID     uuid.UUID           `json:"projectUUID" db:"project_uuid"`
	ProjectName     string              `json:"projectName" db:"project_name"`
	ClientUUID      uuid.UUID           `json:"clientUUID" db:"client_uuid"`
	ClientName      string              `json:"clientName" db:"client_name"`
	Date            string              `json:"date" db:"date"`
	TimeStart       string              `json:"timeStart" db:"time_start"`
	TimeEnd         nullable.NullString `json:"timeEnd" db:"time_end"`
	DurationSeconds uint32              `json:"durationSeconds" db:"duration_seconds"`
	BillableRate    uint64              `json:"billableRate" db:"billable_rate"`
	BillableAmount  uint64              `json:"billableAmount" db:"billable_amount"`
	Comment         nullable.NullString `json:"comment" db:"comment"`
	CreatedAt       time.Time           `json:"createdAt" db:"created_at"`
	UpdatedAt       nullable.NullTime   `json:"updatedAt" db:"updated_at"`
	DeletedAt       nullable.NullTime   `json:"deletedAt" db:"deleted_at"`
}
