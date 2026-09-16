package sqlite

import (
	"database/sql"

	"capuchin/internal/domain"
)

// timelogModel представляет запись в журнале.
type timelogModel struct {
	ID              string         `db:"id"`
	UserID          string         `db:"user_id"`
	ProjectID       string         `db:"project_id"`
	ProjectName     string         `db:"project_name"`
	ClientID        string         `db:"client_id"`
	ClientName      string         `db:"client_name"`
	Date            string         `db:"date"`
	TimeStart       string         `db:"time_start"`
	TimeEnd         sql.NullString `db:"time_end"`
	DurationSeconds uint32         `db:"duration_seconds"`
	BillableRate    int32          `db:"billable_rate"`
	BillableAmount  int64          `db:"billable_amount"`
	Comment         sql.NullString `db:"comment"`
	CreatedAt       int64          `db:"created_at"`
	UpdatedAt       sql.NullInt64  `db:"updated_at"`
	TaskID          sql.NullString `db:"task_id"`
	TaskName        sql.NullString `db:"task_name"`
	TaskCompletedAt sql.NullInt64  `db:"task_completed_at"`
}

// timelogModelFromDomain преобразует доменный timelog в sqlite.
func timelogModelFromDomain(src domain.Timelog) timelogModel {
	return timelogModel{
		ID:              src.ID,
		UserID:          src.UserID,
		ProjectID:       src.ProjectID,
		ProjectName:     src.ProjectName,
		ClientID:        src.ClientID,
		ClientName:      src.ClientName,
		Date:            src.Date,
		TimeStart:       src.TimeStart,
		TimeEnd:         nullableStringFromDomain(src.TimeEnd),
		DurationSeconds: src.DurationSeconds,
		BillableRate:    src.BillableRate,
		BillableAmount:  src.BillableAmount,
		Comment:         commentFromDomain(src.Comment),
		CreatedAt:       timeToUnix(src.CreatedAt),
		UpdatedAt:       nullableUnixFromDomain(src.UpdatedAt),
		TaskID:          nullableStringFromDomain(src.TaskID),
		TaskName:        nullableStringFromDomain(src.TaskName),
		TaskCompletedAt: nullableUnixFromDomain(src.TaskCompletedAt),
	}
}

// timelogModelToDomain преобразует sqlite timelog в доменный.
func timelogModelToDomain(src timelogModel) domain.Timelog {
	return domain.Timelog{
		ID:              src.ID,
		UserID:          src.UserID,
		ProjectID:       src.ProjectID,
		ProjectName:     src.ProjectName,
		ClientID:        src.ClientID,
		ClientName:      src.ClientName,
		Date:            src.Date,
		TimeStart:       src.TimeStart,
		TimeEnd:         nullableStringToDomain(src.TimeEnd),
		DurationSeconds: src.DurationSeconds,
		BillableRate:    src.BillableRate,
		BillableAmount:  src.BillableAmount,
		Comment:         commentToDomain(src.Comment),
		CreatedAt:       unixToTime(src.CreatedAt),
		UpdatedAt:       nullableUnixToDomain(src.UpdatedAt),
		TaskID:          nullableStringToDomain(src.TaskID),
		TaskName:        nullableStringToDomain(src.TaskName),
		TaskCompletedAt: nullableUnixToDomain(src.TaskCompletedAt),
	}
}

// timelogModelsToDomains преобразует sqlite timelogs в доменные.
func timelogModelsToDomains(src []timelogModel) []domain.Timelog {
	dst := make([]domain.Timelog, len(src))

	for i := range src {
		dst[i] = timelogModelToDomain(src[i])
	}

	return dst
}

// taskReportModel представляет запись в журнале.
type taskReportModel struct {
	DurationSeconds  uint32 `db:"duration_seconds"`
	BillableAmount   int64  `db:"billable_amount"`
	UniqueDatesCount uint32 `db:"unique_dates_count"`
}

// taskReportModelToDomain преобразует sqlite task report в доменный.
func taskReportModelToDomain(src taskReportModel) domain.TaskReport {
	return domain.TaskReport{
		DurationSeconds:  src.DurationSeconds,
		BillableAmount:   src.BillableAmount,
		UniqueDatesCount: src.UniqueDatesCount,
	}
}
