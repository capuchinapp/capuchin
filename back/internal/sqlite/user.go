package sqlite

import (
	"database/sql"

	"capuchin/internal/domain"
)

// userModel представляет пользователя.
type userModel struct {
	ID           string         `db:"id"`
	Email        string         `db:"email"`
	ActivateCode sql.NullString `db:"activate_code"`
	CreatedAt    int64          `db:"created_at"`
	UpdatedAt    sql.NullInt64  `db:"updated_at"`
}

// userModelFromDomain преобразует пользователя домена в sqlite.
func userModelFromDomain(src domain.User) userModel {
	return userModel{
		ID:           src.ID,
		Email:        src.Email,
		ActivateCode: nullableStringFromDomain(src.ActivateCode),
		CreatedAt:    timeToUnix(src.CreatedAt),
		UpdatedAt:    nullableUnixFromDomain(src.UpdatedAt),
	}
}

// userModelToDomain преобразует sqlite пользователя в пользователя домена.
func userModelToDomain(src userModel) domain.User {
	return domain.User{
		ID:           src.ID,
		Email:        src.Email,
		ActivateCode: nullableStringToDomain(src.ActivateCode),
		CreatedAt:    unixToTime(src.CreatedAt),
		UpdatedAt:    nullableUnixToDomain(src.UpdatedAt),
	}
}
