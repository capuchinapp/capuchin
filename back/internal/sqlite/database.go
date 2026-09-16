package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // blank-import: регистрация database/sql-драйвера "sqlite"
)

// buildDSN формирует DSN подключения к файлу SQLite с прагмами.
//
// Прагмы:
//   - busy_timeout(5000) — ждать снятия блокировки вместо мгновенного SQLITE_BUSY;
//   - journal_mode(WAL) — конкурентное чтение без блокировки писателя;
//   - synchronous(NORMAL) — баланс между безопасностью и скоростью в режиме WAL;
//   - _txlock=immediate — BeginTx/BeginTxx выполняет BEGIN IMMEDIATE.
func buildDSN(path string) string {
	return fmt.Sprintf(
		"file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_txlock=immediate",
		path,
	)
}

// NewDatabase возвращает соединение с файлом базы данных SQLite.
func NewDatabase(ctx context.Context, path string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "sqlite", buildDSN(path))
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %v", err)
	}

	return db, nil
}
