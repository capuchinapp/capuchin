package sqlite

import (
	"database/sql"

	"capuchin/internal/domain"
)

// projectModel представляет проект.
type projectModel struct {
	ID           string         `db:"id"`
	UserID       string         `db:"user_id"`
	ClientID     string         `db:"client_id"`
	ClientName   string         `db:"client_name"`
	Name         string         `db:"name"`
	BillableRate int32          `db:"billable_rate"`
	Comment      sql.NullString `db:"comment"`
	CreatedAt    int64          `db:"created_at"`
	UpdatedAt    sql.NullInt64  `db:"updated_at"`
	ArchivedAt   sql.NullInt64  `db:"archived_at"`
}

// projectModelFromDomain преобразует проект домена в sqlite.
func projectModelFromDomain(src domain.Project) projectModel {
	return projectModel{
		ID:           src.ID,
		UserID:       src.UserID,
		ClientID:     src.ClientID,
		ClientName:   src.ClientName,
		Name:         src.Name,
		BillableRate: src.BillableRate,
		Comment:      commentFromDomain(src.Comment),
		CreatedAt:    timeToUnix(src.CreatedAt),
		UpdatedAt:    nullableUnixFromDomain(src.UpdatedAt),
		ArchivedAt:   nullableUnixFromDomain(src.ArchivedAt),
	}
}

// projectModelToDomain преобразует sqlite в проект домена.
func projectModelToDomain(src projectModel) domain.Project {
	return domain.Project{
		ID:           src.ID,
		UserID:       src.UserID,
		ClientID:     src.ClientID,
		ClientName:   src.ClientName,
		Name:         src.Name,
		BillableRate: src.BillableRate,
		Comment:      commentToDomain(src.Comment),
		CreatedAt:    unixToTime(src.CreatedAt),
		UpdatedAt:    nullableUnixToDomain(src.UpdatedAt),
		ArchivedAt:   nullableUnixToDomain(src.ArchivedAt),
	}
}

// projectModelsToDomains преобразует sqlite проекты в проекты домена.
func projectModelsToDomains(src []projectModel) []domain.Project {
	dst := make([]domain.Project, len(src))

	for i, s := range src {
		dst[i] = projectModelToDomain(s)
	}

	return dst
}
