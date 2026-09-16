package sqlite

import (
	"database/sql"

	"capuchin/internal/domain"
)

// favoriteModel представляет элемент избранного.
type favoriteModel struct {
	ID           string         `db:"id"`
	UserID       string         `db:"user_id"`
	Name         string         `db:"name"`
	ProjectID    string         `db:"project_id"`
	ProjectName  string         `db:"project_name"`
	ClientID     string         `db:"client_id"`
	ClientName   string         `db:"client_name"`
	TaskID       sql.NullString `db:"task_id"`
	TaskName     sql.NullString `db:"task_name"`
	BillableRate int32          `db:"billable_rate"`
	Comment      sql.NullString `db:"comment"`
}

// favoriteModelFromDomain преобразует доменный favorite в sqlite.
func favoriteModelFromDomain(src domain.Favorite) favoriteModel {
	return favoriteModel{
		ID:           src.ID,
		UserID:       src.UserID,
		Name:         src.Name,
		ProjectID:    src.ProjectID,
		ProjectName:  src.ProjectName,
		ClientID:     src.ClientID,
		ClientName:   src.ClientName,
		TaskID:       nullableStringFromDomain(src.TaskID),
		TaskName:     nullableStringFromDomain(src.TaskName),
		BillableRate: src.BillableRate,
		Comment:      commentFromDomain(src.Comment),
	}
}

// favoriteModelToDomain преобразует sqlite favorite в доменный.
func favoriteModelToDomain(src favoriteModel) domain.Favorite {
	return domain.Favorite{
		ID:           src.ID,
		UserID:       src.UserID,
		Name:         src.Name,
		ProjectID:    src.ProjectID,
		ProjectName:  src.ProjectName,
		ClientID:     src.ClientID,
		ClientName:   src.ClientName,
		TaskID:       nullableStringToDomain(src.TaskID),
		TaskName:     nullableStringToDomain(src.TaskName),
		BillableRate: src.BillableRate,
		Comment:      commentToDomain(src.Comment),
	}
}

// favoriteModelsToDomains преобразует sqlite favorites в доменные.
func favoriteModelsToDomains(src []favoriteModel) []domain.Favorite {
	dst := make([]domain.Favorite, len(src))

	for i := range src {
		dst[i] = favoriteModelToDomain(src[i])
	}

	return dst
}
