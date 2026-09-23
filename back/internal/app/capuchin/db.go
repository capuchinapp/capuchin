package capuchin

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"capuchin/internal/sqlite"
)

// newDB создает соединение с базой данных SQLite.
func newDB(ctx context.Context, conf *Configuration) (*sqlx.DB, error) {
	db, err := sqlite.NewDatabase(ctx, conf.Sqlite.DBPath)
	if err != nil {
		return nil, fmt.Errorf("create database connection: %v", err)
	}

	return db, nil
}
