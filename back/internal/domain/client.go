package domain

import (
	"time"
)

// Client представляет клиента.
type Client struct {
	ID           string     `json:"id"`
	UserID       string     `json:"-"`
	Name         string     `json:"name"`
	BillableRate int32      `json:"billableRate"`
	Comment      string     `json:"comment"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	ArchivedAt   *time.Time `json:"archivedAt"`
}
