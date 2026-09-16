package service

import (
	"context"
	"time"

	"capuchin/internal/domain"
	"capuchin/internal/mailer"
	"capuchin/internal/sqlite"
)

// MailService представляет сервис отправки писем.
type MailService interface {
	Send(ctx context.Context, recipients []string, subject string, tpl mailer.Template, content any) error
}

// TransactionManager SQLite реализация менеджера транзакций.
type TransactionManager interface {
	ExecuteInTransaction(ctx context.Context, opFuncs ...func(ctx context.Context, tx sqlite.SQLExecutor) error) error
}

// SessionRepository представляет репозиторий сессии.
type SessionRepository interface {
	Create(ctx context.Context, session domain.Session) error
	FindAll(ctx context.Context, userID string) ([]domain.Session, error)
	CancelOld(ctx context.Context, olderThan time.Time) error
}

// ClientRepository представляет репозиторий клиента.
type ClientRepository interface {
	Create(ctx context.Context, db sqlite.SQLExecutor, client domain.Client) error
	Update(ctx context.Context, db sqlite.SQLExecutor, client domain.Client) error
	FindByID(ctx context.Context, userID string, id string) (domain.Client, error)
}

// ProjectRepository представляет репозиторий проекта.
type ProjectRepository interface {
	Create(ctx context.Context, db sqlite.SQLExecutor, project domain.Project) error
	Update(ctx context.Context, db sqlite.SQLExecutor, project domain.Project) error
	FindByID(ctx context.Context, userID string, id string) (domain.Project, error)
}

// TimelogRepository представляет репозиторий записи времени.
type TimelogRepository interface {
	Create(ctx context.Context, db sqlite.SQLExecutor, timelog domain.Timelog) error
	Update(ctx context.Context, db sqlite.SQLExecutor, timelog domain.Timelog) error
	Delete(ctx context.Context, db sqlite.SQLExecutor, userID string, id string) error
	FindByID(ctx context.Context, userID string, id string) (domain.Timelog, error)
	FindRunning(ctx context.Context, userID string) (domain.Timelog, error)
	FindLastN(ctx context.Context, userID string, n int) ([]domain.Timelog, error)
	TaskReport(ctx context.Context, userID string, taskID string) (domain.TaskReport, error)
}

// UserRepository представляет репозиторий пользователя.
type UserRepository interface {
	Create(ctx context.Context, db sqlite.SQLExecutor, user domain.User) error
	CancelActivateCode(ctx context.Context, user domain.User) error
	DeleteNotActivated(ctx context.Context, olderThan time.Time) error
	FindByID(ctx context.Context, id string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

// AuthCodeRepository представляет репозиторий кода авторизации.
type AuthCodeRepository interface {
	Create(ctx context.Context, authCode domain.AuthCode) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, email string, code string) (domain.AuthCode, error)
	CancelOld(ctx context.Context, olderThan time.Time) error
}

// TaskRepository представляет репозиторий задач.
type TaskRepository interface {
	Create(ctx context.Context, db sqlite.SQLExecutor, task domain.Task) error
	Update(ctx context.Context, db sqlite.SQLExecutor, task domain.Task) error
	FindByID(ctx context.Context, userID string, id string) (domain.Task, error)
}

// FavoriteRepository представляет репозиторий избранного.
type FavoriteRepository interface {
	Create(ctx context.Context, db sqlite.SQLExecutor, favorite domain.Favorite) error
	Update(ctx context.Context, db sqlite.SQLExecutor, favorite domain.Favorite) error
	Delete(ctx context.Context, db sqlite.SQLExecutor, userID string, id string) error
	FindByID(ctx context.Context, userID string, id string) (domain.Favorite, error)
}

// AuditLogRepository представляет репозиторий лога действий.
type AuditLogRepository interface {
	Create(ctx context.Context, db sqlite.SQLExecutor, auditLog domain.AuditLog) error
}
