package sqlite

import (
	"context"
	"database/sql"
)

// SQLExecutor представляет интерфейс выполнения запроса как в транзакции так и без.
type SQLExecutor interface {
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}
