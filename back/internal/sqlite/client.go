package sqlite

import (
	"database/sql"

	"capuchin/internal/domain"
)

// clientModel представляет клиента.
type clientModel struct {
	ID           string         `db:"id"`
	UserID       string         `db:"user_id"`
	Name         string         `db:"name"`
	BillableRate int32          `db:"billable_rate"`
	Comment      sql.NullString `db:"comment"`
	CreatedAt    int64          `db:"created_at"`
	UpdatedAt    sql.NullInt64  `db:"updated_at"`
	ArchivedAt   sql.NullInt64  `db:"archived_at"`
}

// clientModelFromDomain преобразует клиента домена в sqlite.
func clientModelFromDomain(src domain.Client) clientModel {
	return clientModel{
		ID:           src.ID,
		UserID:       src.UserID,
		Name:         src.Name,
		BillableRate: src.BillableRate,
		Comment:      commentFromDomain(src.Comment),
		CreatedAt:    timeToUnix(src.CreatedAt),
		UpdatedAt:    nullableUnixFromDomain(src.UpdatedAt),
		ArchivedAt:   nullableUnixFromDomain(src.ArchivedAt),
	}
}

// clientModelToDomain преобразует sqlite в клиента домена.
func clientModelToDomain(src clientModel) domain.Client {
	return domain.Client{
		ID:           src.ID,
		UserID:       src.UserID,
		Name:         src.Name,
		BillableRate: src.BillableRate,
		Comment:      commentToDomain(src.Comment),
		CreatedAt:    unixToTime(src.CreatedAt),
		UpdatedAt:    nullableUnixToDomain(src.UpdatedAt),
		ArchivedAt:   nullableUnixToDomain(src.ArchivedAt),
	}
}

// clientModelsToDomain преобразует sqlite clients в клиентов домена.
func clientModelsToDomains(src []clientModel) []domain.Client {
	dst := make([]domain.Client, len(src))

	for i := range src {
		dst[i] = clientModelToDomain(src[i])
	}

	return dst
}
