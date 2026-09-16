package sqlite

import (
	"errors"

	sqlitedriver "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// isUniqueViolation возвращает true, если ошибка вызвана нарушением
// уникального ограничения SQLite (SQLITE_CONSTRAINT_UNIQUE).
func isUniqueViolation(err error) bool {
	var e *sqlitedriver.Error

	return errors.As(err, &e) && e.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}
