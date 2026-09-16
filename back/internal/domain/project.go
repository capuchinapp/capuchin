package domain

import (
	"time"
)

// Project представляет проект.
type Project struct {
	ID           string     `json:"id"`
	UserID       string     `json:"-"`
	ClientID     string     `json:"clientId"`
	ClientName   string     `json:"clientName"`
	Name         string     `json:"name"`
	BillableRate int32      `json:"billableRate"`
	Comment      string     `json:"comment"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	ArchivedAt   *time.Time `json:"archivedAt"`
}

// ProjectFilter представляет фильтр проектов.
type ProjectFilter struct {
	ClientID               string
	WithoutArchivedClients bool
}
