package sqlite

import (
	"database/sql"

	"capuchin/internal/domain"
)

// taskModel представляет задачу.
type taskModel struct {
	ID          string         `db:"id"`
	UserID      string         `db:"user_id"`
	ProjectID   string         `db:"project_id"`
	ProjectName string         `db:"project_name"`
	ClientID    string         `db:"client_id"`
	ClientName  string         `db:"client_name"`
	Name        string         `db:"name"`
	Comment     sql.NullString `db:"comment"`
	CreatedAt   int64          `db:"created_at"`
	UpdatedAt   sql.NullInt64  `db:"updated_at"`
	ArchivedAt  sql.NullInt64  `db:"archived_at"`
	CompletedAt sql.NullInt64  `db:"completed_at"`
}

// taskModelFromDomain преобразует задачу из домена в sqlite.
func taskModelFromDomain(src domain.Task) taskModel {
	return taskModel{
		ID:          src.ID,
		UserID:      src.UserID,
		ProjectID:   src.ProjectID,
		ProjectName: src.ProjectName,
		ClientID:    src.ClientID,
		ClientName:  src.ClientName,
		Name:        src.Name,
		Comment:     commentFromDomain(src.Comment),
		CreatedAt:   timeToUnix(src.CreatedAt),
		UpdatedAt:   nullableUnixFromDomain(src.UpdatedAt),
		ArchivedAt:  nullableUnixFromDomain(src.ArchivedAt),
		CompletedAt: nullableUnixFromDomain(src.CompletedAt),
	}
}

// taskModelToDomain преобразует задачу из sqlite в домен.
func taskModelToDomain(src taskModel) domain.Task {
	return domain.Task{
		ID:          src.ID,
		UserID:      src.UserID,
		ProjectID:   src.ProjectID,
		ProjectName: src.ProjectName,
		ClientID:    src.ClientID,
		ClientName:  src.ClientName,
		Name:        src.Name,
		Comment:     commentToDomain(src.Comment),
		CreatedAt:   unixToTime(src.CreatedAt),
		UpdatedAt:   nullableUnixToDomain(src.UpdatedAt),
		ArchivedAt:  nullableUnixToDomain(src.ArchivedAt),
		CompletedAt: nullableUnixToDomain(src.CompletedAt),
	}
}

// taskModelsToDomains преобразует задачи из sqlite в домен.
func taskModelsToDomains(src []taskModel) []domain.Task {
	dst := make([]domain.Task, len(src))

	for i, s := range src {
		dst[i] = taskModelToDomain(s)
	}

	return dst
}
