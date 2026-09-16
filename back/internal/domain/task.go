package domain

import (
	"time"
)

// Task представляет задачу.
type Task struct {
	ID          string     `json:"id"`
	UserID      string     `json:"-"`
	ProjectID   string     `json:"projectId"`
	ProjectName string     `json:"projectName"`
	ClientID    string     `json:"clientId"`
	ClientName  string     `json:"clientName"`
	Name        string     `json:"name"`
	Comment     string     `json:"comment"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

// TaskReport представляет отчёт по задаче.
type TaskReport struct {
	DurationSeconds  uint32 `json:"durationSeconds"`
	BillableAmount   int64  `json:"billableAmount"`
	UniqueDatesCount uint32 `json:"uniqueDatesCount"`
}

// TaskFilter представляет фильтр задач.
type TaskFilter struct {
	ProjectID               string
	WithoutArchivedProjects bool
}
