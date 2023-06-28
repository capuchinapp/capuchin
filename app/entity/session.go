package entity

import (
	"capuchin/app/util/nullable"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	RowID     nullable.NullInt64 `json:"id" db:"rowid"`
	UUID      uuid.UUID          `json:"-" db:"uuid"`
	CheckedAt time.Time          `json:"checkedAt" db:"checked_at"`
	IsCurrent bool               `json:"isCurrent"`
}
