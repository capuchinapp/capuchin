package sqlite

import "capuchin/internal/domain"

// authCodeModel представляет код авторизации.
type authCodeModel struct {
	ID        string `db:"id"`
	UserID    string `db:"user_id"`
	UserEmail string `db:"user_email"`
	Code      string `db:"code"`
	CreatedAt int64  `db:"created_at"`
}

// authCodeModelFromDomain преобразует код аутентификации домена в код аутентификации SQLite.
func authCodeModelFromDomain(src domain.AuthCode) authCodeModel {
	return authCodeModel{
		ID:        src.ID,
		UserID:    src.UserID,
		UserEmail: src.UserEmail,
		Code:      src.Code,
		CreatedAt: timeToUnix(src.CreatedAt),
	}
}

// authCodeModelToDomain преобразует код аутентификации SQLite в код аутентификации домена.
func authCodeModelToDomain(src authCodeModel) domain.AuthCode {
	return domain.AuthCode{
		ID:        src.ID,
		UserID:    src.UserID,
		UserEmail: src.UserEmail,
		Code:      src.Code,
		CreatedAt: unixToTime(src.CreatedAt),
	}
}
