package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func New(filepath string) (*sqlx.DB, error) {
	db := sqlx.MustOpen("sqlite", fmt.Sprintf("file:%s?cache=shared&mode=rwc", filepath))

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("❌ [db] db.Ping: %v\n", err)
	}

	fmt.Printf("📒 Database loaded from file %s\n", filepath)

	if _, err := prepareDb(db); err != nil {
		return nil, fmt.Errorf("❌ [db] prepareDb: %v\n", err)
	}

	fmt.Println("📒 Tables have been prepared")

	return db, nil
}

func prepareDb(db *sqlx.DB) (sql.Result, error) {
	return db.Exec(`CREATE TABLE IF NOT EXISTS clients (
		uuid TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		billable_rate INTEGER NOT NULL,
		comment TEXT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NULL,
		archived_at DATETIME NULL
	) WITHOUT ROWID;

	CREATE TABLE IF NOT EXISTS projects (
		uuid TEXT NOT NULL PRIMARY KEY,
		client_uuid TEXT NOT NULL,
		name TEXT NOT NULL UNIQUE,
		billable_rate INTEGER NOT NULL,
		comment TEXT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NULL,
		archived_at DATETIME NULL
	) WITHOUT ROWID;
	CREATE INDEX IF NOT EXISTS client_uuid ON projects (client_uuid);

	CREATE TABLE IF NOT EXISTS timelogs (
		uuid TEXT NOT NULL PRIMARY KEY,
		project_uuid TEXT NOT NULL,
		date TEXT NOT NULL,
		time_start TEXT NOT NULL,
		time_end TEXT NULL,
		duration_seconds INTEGER NOT NULL,
		billable_rate INTEGER NOT NULL,
		billable_amount INTEGER NOT NULL,
		comment TEXT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NULL,
		deleted_at DATETIME NULL
	) WITHOUT ROWID;
	CREATE INDEX IF NOT EXISTS project_uuid ON timelogs (project_uuid);

	CREATE TABLE IF NOT EXISTS settings (
		key TEXT NOT NULL PRIMARY KEY,
		value TEXT NOT NULL
	) WITHOUT ROWID;
	INSERT OR IGNORE INTO settings(key, value) VALUES('dateFormat', '');

	CREATE TABLE IF NOT EXISTS sessions (
		uuid TEXT NOT NULL PRIMARY KEY,
		checked_at DATETIME NOT NULL
	);`)
}
