package sqlite

import "capuchin/internal/domain"

// sessionModel представляет сессию.
type sessionModel struct {
	ID        string `db:"id"`
	CookieID  string `db:"cid"`
	UserID    string `db:"uid"`
	CheckedAt int64  `db:"checked_at"`
}

// sessionModelFromDomain преобразует доменную сессию в sqlite.
func sessionModelFromDomain(src domain.Session) sessionModel {
	return sessionModel{
		ID:        src.ID,
		CookieID:  src.CookieID,
		UserID:    src.UserID,
		CheckedAt: timeToUnix(src.CheckedAt),
	}
}

// sessionModelToDomain преобразует sqlite сессию в доменную сессию.
func sessionModelToDomain(src sessionModel) domain.Session {
	return domain.Session{
		ID:        src.ID,
		CookieID:  src.CookieID,
		UserID:    src.UserID,
		CheckedAt: unixToTime(src.CheckedAt),
	}
}

// sessionModelsToDomains преобразует sqlite сессии в доменные сессии.
func sessionModelsToDomains(src []sessionModel) []domain.Session {
	dst := make([]domain.Session, len(src))

	for i, s := range src {
		dst[i] = sessionModelToDomain(s)
	}

	return dst
}
