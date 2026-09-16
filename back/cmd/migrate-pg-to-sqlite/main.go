//go:build migration

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func main() {
	pgDSN := flag.String("pg-dsn", "", "PostgreSQL DSN (source)")
	sqlitePath := flag.String("sqlite-path", "", "SQLite file path (target)")
	dryRun := flag.Bool("dry-run", false, "read source and count without writing")

	flag.Parse()

	if *pgDSN == "" || *sqlitePath == "" {
		log.Fatal("flags -pg-dsn and -sqlite-path are required")
	}

	if err := run(*pgDSN, *sqlitePath, *dryRun); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}

func run(pgDSN, sqlitePath string, dryRun bool) error {
	ctx := context.Background()

	pgDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		return fmt.Errorf("open postgres: %w", err)
	}
	defer pgDB.Close()

	sqDB, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_txlock=immediate", sqlitePath))
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}
	defer sqDB.Close()

	if err := sqDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping sqlite: %w", err)
	}

	if err := verifyEmpty(ctx, sqDB); err != nil {
		return err
	}

	if dryRun {
		log.Println("dry-run mode: reading and counting source only, no writes")

		return countAll(ctx, pgDB)
	}

	log.Println("migrating tables in dependency order")

	tables := []struct {
		name string
		fn   func(ctx context.Context, pgDB, sqDB *sql.DB) error
	}{
		{"users", migrateUsers},
		{"clients", migrateClients},
		{"projects", migrateProjects},
		{"tasks", migrateTasks},
		{"timelogs", migrateTimelogs},
		{"favorites", migrateFavorites},
		{"settings", migrateSettings},
		{"sessions", migrateSessions},
		{"auth_codes", migrateAuthCodes},
		{"audit_log", migrateAuditLog},
	}

	for _, t := range tables {
		err = t.fn(ctx, pgDB, sqDB)
		if err != nil {
			return fmt.Errorf("migrate %s: %w", t.name, err)
		}
	}

	log.Println("migration completed successfully")

	return nil
}

// verifyEmpty проверяет, что все таблицы приёмника пусты (свежая схема после goose).
func verifyEmpty(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, `SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name != 'goose_db_version'`)
	if err != nil {
		return fmt.Errorf("list sqlite tables: %w", err)
	}
	defer rows.Close()

	var tables []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return err
		}

		tables = append(tables, name)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for _, table := range tables {
		var count int64

		err := db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %q", table)).Scan(&count)
		if err != nil {
			return fmt.Errorf("count %s: %w", table, err)
		}

		if count > 0 {
			return fmt.Errorf("target sqlite database is not empty: table %q contains %d rows; expected a fresh schema from goose migrations", table, count)
		}
	}

	return nil
}

func countRows(ctx context.Context, db *sql.DB, table string) (int64, error) {
	var count int64

	err := db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %q", table)).Scan(&count)

	return count, err
}

func countAll(ctx context.Context, pgDB *sql.DB) error {
	tables := []string{"users", "clients", "projects", "tasks", "timelogs", "favorites", "settings", "sessions", "auth_codes", "audit_log"}

	for _, table := range tables {
		count, err := countRows(ctx, pgDB, table)
		if err != nil {
			return fmt.Errorf("count %s: %w", table, err)
		}

		log.Printf("  %s: %d rows\n", table, count)
	}

	return nil
}

// batchWriter записывает строки в одном INSERT-батче внутри одной транзакции.
type batchWriter struct {
	ctx       context.Context
	tx        *sql.Tx
	table     string
	columns   []string
	pending   []any
	rowsCount int64
	batches   int
}

func newBatchWriter(ctx context.Context, db *sql.DB, table string, columns []string) (*batchWriter, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &batchWriter{ctx: ctx, tx: tx, table: table, columns: columns}, nil
}

// add накапливает строку и сбрасывает батч при достижении batchSize.
func (b *batchWriter) add(row ...any) error {
	b.pending = append(b.pending, row...)
	b.rowsCount++

	if len(b.pending) >= batchSize*len(b.columns) {
		return b.flush()
	}

	return nil
}

// flush выполняет многострочный INSERT с накопленными строками, не завершая транзакцию.
func (b *batchWriter) flush() error {
	if len(b.pending) == 0 {
		return nil
	}

	rowsInBatch := len(b.pending) / len(b.columns)
	values := make([]string, 0, rowsInBatch)
	args := b.pending

	for i := 0; i < rowsInBatch; i++ {
		values = append(values, "("+strings.TrimSuffix(strings.Repeat("?,", len(b.columns)), ",")+")")
	}

	query := fmt.Sprintf(
		"INSERT INTO %q (%s) VALUES %s",
		b.table,
		strings.Join(b.columns, ", "),
		strings.Join(values, ", "),
	)

	_, err := b.tx.ExecContext(b.ctx, query, args...)
	if err != nil {
		return err
	}

	b.pending = nil
	b.batches++

	return nil
}

// close сбрасывает остаток батча, завершает транзакцию и возвращает число записанных строк.
func (b *batchWriter) close() error {
	if err := b.flush(); err != nil {
		b.tx.Rollback() //nolint:errcheck

		return err
	}

	if err := b.tx.Commit(); err != nil {
		return err
	}

	return nil
}

func timeToUnix(src time.Time) int64 {
	return src.UTC().Unix()
}

func nullableTimeToUnix(src sql.NullTime) sql.NullInt64 {
	if src.Valid {
		return sql.NullInt64{Int64: timeToUnix(src.Time), Valid: true}
	}

	return sql.NullInt64{Valid: false}
}

const batchSize = 1000

// --- users ---

type pgUser struct {
	ID           string
	Email        string
	ActivateCode sql.NullString
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
}

func migrateUsers(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT id, email, activate_code, created_at, updated_at FROM users`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "users", []string{"id", "email", "activate_code", "created_at", "updated_at"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var u pgUser

		if err := rows.Scan(&u.ID, &u.Email, &u.ActivateCode, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return err
		}

		if err := bw.add(u.ID, u.Email, u.ActivateCode, timeToUnix(u.CreatedAt), nullableTimeToUnix(u.UpdatedAt)); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "users", bw.rowsCount, bw.batches)
}

// --- clients ---

type pgClient struct {
	ID           string
	UserID       string
	Name         string
	BillableRate int32
	Comment      sql.NullString
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	ArchivedAt   sql.NullTime
}

func migrateClients(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT id, user_id, name, billable_rate, comment, created_at, updated_at, archived_at FROM clients`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "clients", []string{"id", "user_id", "name", "billable_rate", "comment", "created_at", "updated_at", "archived_at"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var c pgClient

		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.BillableRate, &c.Comment, &c.CreatedAt, &c.UpdatedAt, &c.ArchivedAt); err != nil {
			return err
		}

		if err := bw.add(c.ID, c.UserID, c.Name, c.BillableRate, c.Comment, timeToUnix(c.CreatedAt), nullableTimeToUnix(c.UpdatedAt), nullableTimeToUnix(c.ArchivedAt)); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "clients", bw.rowsCount, bw.batches)
}

// --- projects ---

type pgProject struct {
	ID           string
	UserID       string
	ClientID     string
	Name         string
	BillableRate int32
	Comment      sql.NullString
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	ArchivedAt   sql.NullTime
}

func migrateProjects(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT id, user_id, client_id, name, billable_rate, comment, created_at, updated_at, archived_at FROM projects`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "projects", []string{"id", "user_id", "client_id", "name", "billable_rate", "comment", "created_at", "updated_at", "archived_at"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var p pgProject

		if err := rows.Scan(&p.ID, &p.UserID, &p.ClientID, &p.Name, &p.BillableRate, &p.Comment, &p.CreatedAt, &p.UpdatedAt, &p.ArchivedAt); err != nil {
			return err
		}

		if err := bw.add(p.ID, p.UserID, p.ClientID, p.Name, p.BillableRate, p.Comment, timeToUnix(p.CreatedAt), nullableTimeToUnix(p.UpdatedAt), nullableTimeToUnix(p.ArchivedAt)); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "projects", bw.rowsCount, bw.batches)
}

// --- tasks ---

type pgTask struct {
	ID          string
	UserID      string
	ProjectID   string
	Name        string
	Comment     sql.NullString
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
	ArchivedAt  sql.NullTime
	CompletedAt sql.NullTime
}

func migrateTasks(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT id, user_id, project_id, name, comment, created_at, updated_at, archived_at, completed_at FROM tasks`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "tasks", []string{"id", "user_id", "project_id", "name", "comment", "created_at", "updated_at", "archived_at", "completed_at"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var t pgTask

		if err := rows.Scan(&t.ID, &t.UserID, &t.ProjectID, &t.Name, &t.Comment, &t.CreatedAt, &t.UpdatedAt, &t.ArchivedAt, &t.CompletedAt); err != nil {
			return err
		}

		if err := bw.add(t.ID, t.UserID, t.ProjectID, t.Name, t.Comment, timeToUnix(t.CreatedAt), nullableTimeToUnix(t.UpdatedAt), nullableTimeToUnix(t.ArchivedAt), nullableTimeToUnix(t.CompletedAt)); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "tasks", bw.rowsCount, bw.batches)
}

// --- timelogs ---

type pgTimelog struct {
	ID              string
	UserID          string
	ProjectID       string
	Date            string
	TimeStart       string
	TimeEnd         sql.NullString
	DurationSeconds uint32
	BillableRate    int32
	BillableAmount  int64
	Comment         sql.NullString
	CreatedAt       time.Time
	UpdatedAt       sql.NullTime
	TaskID          sql.NullString
}

func migrateTimelogs(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT id, user_id, project_id, date, time_start, time_end, duration_seconds, billable_rate, billable_amount, comment, created_at, updated_at, task_id FROM timelogs`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "timelogs", []string{"id", "user_id", "project_id", "date", "time_start", "time_end", "duration_seconds", "billable_rate", "billable_amount", "comment", "created_at", "updated_at", "task_id"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var tl pgTimelog

		if err := rows.Scan(&tl.ID, &tl.UserID, &tl.ProjectID, &tl.Date, &tl.TimeStart, &tl.TimeEnd, &tl.DurationSeconds, &tl.BillableRate, &tl.BillableAmount, &tl.Comment, &tl.CreatedAt, &tl.UpdatedAt, &tl.TaskID); err != nil {
			return err
		}

		if err := bw.add(tl.ID, tl.UserID, tl.ProjectID, tl.Date, tl.TimeStart, tl.TimeEnd, tl.DurationSeconds, tl.BillableRate, tl.BillableAmount, tl.Comment, timeToUnix(tl.CreatedAt), nullableTimeToUnix(tl.UpdatedAt), tl.TaskID); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "timelogs", bw.rowsCount, bw.batches)
}

// --- favorites ---

type pgFavorite struct {
	ID           string
	UserID       string
	Name         string
	ProjectID    string
	TaskID       sql.NullString
	BillableRate int32
	Comment      sql.NullString
}

func migrateFavorites(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT id, user_id, name, project_id, task_id, billable_rate, comment FROM favorites`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "favorites", []string{"id", "user_id", "name", "project_id", "task_id", "billable_rate", "comment"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var f pgFavorite

		if err := rows.Scan(&f.ID, &f.UserID, &f.Name, &f.ProjectID, &f.TaskID, &f.BillableRate, &f.Comment); err != nil {
			return err
		}

		if err := bw.add(f.ID, f.UserID, f.Name, f.ProjectID, f.TaskID, f.BillableRate, f.Comment); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "favorites", bw.rowsCount, bw.batches)
}

// --- settings ---

type pgSetting struct {
	UserID string
	Key    string
	Value  string
}

func migrateSettings(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT user_id, "key", value FROM settings`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "settings", []string{"user_id", "key", "value"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var s pgSetting

		if err := rows.Scan(&s.UserID, &s.Key, &s.Value); err != nil {
			return err
		}

		if err := bw.add(s.UserID, s.Key, s.Value); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "settings", bw.rowsCount, bw.batches)
}

// --- sessions ---

type pgSession struct {
	ID        string
	CookieID  string
	UserID    string
	CheckedAt time.Time
}

func migrateSessions(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT id, cid, uid, checked_at FROM sessions`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "sessions", []string{"id", "cid", "uid", "checked_at"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var s pgSession

		if err := rows.Scan(&s.ID, &s.CookieID, &s.UserID, &s.CheckedAt); err != nil {
			return err
		}

		if err := bw.add(s.ID, s.CookieID, s.UserID, timeToUnix(s.CheckedAt)); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "sessions", bw.rowsCount, bw.batches)
}

// --- auth_codes ---

type pgAuthCode struct {
	ID        string
	UserID    string
	UserEmail string
	Code      string
	CreatedAt time.Time
}

func migrateAuthCodes(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT id, user_id, user_email, code, created_at FROM auth_codes`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "auth_codes", []string{"id", "user_id", "user_email", "code", "created_at"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var ac pgAuthCode

		if err := rows.Scan(&ac.ID, &ac.UserID, &ac.UserEmail, &ac.Code, &ac.CreatedAt); err != nil {
			return err
		}

		if err := bw.add(ac.ID, ac.UserID, ac.UserEmail, ac.Code, timeToUnix(ac.CreatedAt)); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "auth_codes", bw.rowsCount, bw.batches)
}

// --- audit_log ---

type pgAuditLog struct {
	ID         int64
	ActionAt   time.Time
	ActionType string
	UserID     string
	TableName  string
	RecordID   string
	OldValue   sql.NullString
	NewValue   sql.NullString
}

func migrateAuditLog(ctx context.Context, pgDB, sqDB *sql.DB) error {
	rows, err := pgDB.QueryContext(ctx, `SELECT id, action_at, action_type, user_id, table_name, record_id, old_value, new_value FROM audit_log`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bw, err := newBatchWriter(ctx, sqDB, "audit_log", []string{"id", "action_at", "action_type", "user_id", "table_name", "record_id", "old_value", "new_value"})
	if err != nil {
		return err
	}

	for rows.Next() {
		var al pgAuditLog

		if err := rows.Scan(&al.ID, &al.ActionAt, &al.ActionType, &al.UserID, &al.TableName, &al.RecordID, &al.OldValue, &al.NewValue); err != nil {
			return err
		}

		if err := bw.add(al.ID, timeToUnix(al.ActionAt), al.ActionType, al.UserID, al.TableName, al.RecordID, al.OldValue, al.NewValue); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err := bw.close(); err != nil {
		return err
	}

	return verifyCount(ctx, pgDB, "audit_log", bw.rowsCount, bw.batches)
}

// verifyCount сверяет число записанных строк с COUNT(*) источника.
func verifyCount(ctx context.Context, pgDB *sql.DB, table string, written int64, batches int) error {
	pgCount, err := countRows(ctx, pgDB, table)
	if err != nil {
		return err
	}

	if pgCount != written {
		return fmt.Errorf("%s count mismatch: pg=%d migrated=%d", table, pgCount, written)
	}

	log.Printf("  %s: %d rows in %d batches\n", table, written, batches)

	return nil
}
