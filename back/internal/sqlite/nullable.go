package sqlite

import (
	"database/sql"
	"time"
)

// nullableStringToDomain преобразует строку в sql.NullString.
func nullableStringFromDomain(src *string) sql.NullString {
	if src != nil {
		return sql.NullString{String: *src, Valid: true}
	}

	return sql.NullString{Valid: false}
}

// nullableStringToDomain преобразует sql.NullString в строку.
func nullableStringToDomain(src sql.NullString) *string {
	if src.Valid {
		return &src.String
	}

	return nil
}

// timeToUnix конвертирует время домена в Unix-секунды (UTC) для хранения.
func timeToUnix(src time.Time) int64 {
	return src.UTC().Unix()
}

// unixToTime конвертирует Unix-секунды из хранилища во время домена (UTC).
func unixToTime(src int64) time.Time {
	return time.Unix(src, 0).UTC()
}

// nullableUnixFromDomain преобразует *time.Time в sql.NullInt64.
func nullableUnixFromDomain(src *time.Time) sql.NullInt64 {
	if src != nil {
		return sql.NullInt64{Int64: timeToUnix(*src), Valid: true}
	}

	return sql.NullInt64{Valid: false}
}

// nullableUnixToDomain преобразует sql.NullInt64 в *time.Time.
func nullableUnixToDomain(src sql.NullInt64) *time.Time {
	if src.Valid {
		t := unixToTime(src.Int64)

		return &t
	}

	return nil
}

// commentFromDomain преобразует комментарий в sql.NullString.
func commentFromDomain(src string) sql.NullString {
	if src != "" {
		return sql.NullString{String: src, Valid: true}
	}

	return sql.NullString{Valid: false}
}

// commentToDomain преобразует sql.NullString в комментарий.
func commentToDomain(src sql.NullString) string {
	if src.Valid {
		return src.String
	}

	return ""
}
