package domain

import (
	"time"
)

// Timelog представляет запись в журнале.
type Timelog struct {
	ID              string     `json:"id"`
	UserID          string     `json:"-"`
	ProjectID       string     `json:"projectId"`
	ProjectName     string     `json:"projectName"`
	ClientID        string     `json:"clientId"`
	ClientName      string     `json:"clientName"`
	Date            string     `json:"date"`
	TimeStart       string     `json:"timeStart"`
	TimeEnd         *string    `json:"timeEnd"`
	DurationSeconds uint32     `json:"durationSeconds"`
	BillableRate    int32      `json:"billableRate"`
	BillableAmount  int64      `json:"billableAmount"`
	Comment         string     `json:"comment"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	TaskID          *string    `json:"taskId"`
	TaskName        *string    `json:"taskName"`
	TaskCompletedAt *time.Time `json:"taskCompletedAt"`
}

// TimelogFilter представляет фильтр для записей в журнале.
type TimelogFilter struct {
	DateFrom  string
	DateTo    string
	ClientID  string
	ProjectID string
}
