package restapi

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"capuchin/internal/domain"
)

// CookieManager представляет менеджер файлов cookie.
type CookieManager interface {
	Create(c *fiber.Ctx, v string)
	Get(c *fiber.Ctx) string
	GetExpiresDays() int
}

// Sanitizer представляет санитайзер для фильтрации входящих данных.
type Sanitizer interface {
	Sanitize(input string) string
}

// AuthService представляет сервис аутентификации.
type AuthService interface {
	Register(ctx context.Context, ip string, email string) error
	Activate(ctx context.Context, userID string, code string) (sessionID string, err error)
	Login(ctx context.Context, email string) error
	ApplyCode(ctx context.Context, email string, code string) (sessionID string, err error)
}

// SessionService представляет сервис сессий.
type SessionService interface {
	List(ctx context.Context, userID string, outdatedTime time.Time) ([]domain.Session, error)
}

// ClientService представляет сервис клиентов.
type ClientService interface {
	Create(ctx context.Context, client domain.Client) (domain.Client, error)
	Update(ctx context.Context, client domain.Client) (domain.Client, error)
	Archive(ctx context.Context, userID string, clientID string) (domain.Client, error)
	Unarchive(ctx context.Context, userID string, clientID string) (domain.Client, error)
}

// ProjectService представляет сервис проектов.
type ProjectService interface {
	Create(ctx context.Context, project domain.Project) (domain.Project, error)
	Update(ctx context.Context, project domain.Project) (domain.Project, error)
	Archive(ctx context.Context, userID string, projectID string) (domain.Project, error)
	Unarchive(ctx context.Context, userID string, projectID string) (domain.Project, error)
}

// TimelogService представляет сервис записи времени.
type TimelogService interface {
	Create(ctx context.Context, timelog domain.Timelog) (domain.Timelog, error)
	Update(ctx context.Context, timelog domain.Timelog) (domain.Timelog, error)
	Delete(ctx context.Context, userID string, id string) error
	Stop(ctx context.Context, userID string, timelogID string, date string, timeEnd string) (domain.Timelog, error)
	FindLastN(ctx context.Context, userID string, n int) ([]domain.Timelog, error)
}

// TaskService представляет сервис задач.
type TaskService interface {
	Create(ctx context.Context, task domain.Task) (domain.Task, error)
	Update(ctx context.Context, task domain.Task) (domain.Task, error)
	Archive(ctx context.Context, userID string, taskID string) (domain.Task, error)
	Unarchive(ctx context.Context, userID string, taskID string) (domain.Task, error)
	Complete(ctx context.Context, userID string, taskID string) (domain.Task, error)
	Incomplete(ctx context.Context, userID string, taskID string) (domain.Task, error)
	Report(ctx context.Context, userID string, taskID string) (domain.TaskReport, error)
}

// FavoriteService представляет сервис избранного.
type FavoriteService interface {
	Create(ctx context.Context, favorite domain.Favorite) (domain.Favorite, error)
	Update(ctx context.Context, favorite domain.Favorite) (domain.Favorite, error)
	Delete(ctx context.Context, userID string, id string) error
}

// SessionRepository представляет репозиторий сессий.
type SessionRepository interface {
	DeleteByCID(ctx context.Context, cid string) error
	FindByCID(ctx context.Context, cid string) (domain.Session, error)
	FindByID(ctx context.Context, userID string, id string) (domain.Session, error)
}

// ClientRepository представляет репозиторий клиентов.
type ClientRepository interface {
	FindAll(ctx context.Context, userID string) ([]domain.Client, error)
	FindByID(ctx context.Context, userID string, id string) (domain.Client, error)
}

// ProjectRepository представляет репозиторий проектов.
type ProjectRepository interface {
	FindByFilter(ctx context.Context, userID string, filter domain.ProjectFilter) ([]domain.Project, error)
	FindByID(ctx context.Context, userID string, id string) (domain.Project, error)
}

// TimelogRepository представляет репозиторий записи времени.
type TimelogRepository interface {
	FindByID(ctx context.Context, userID string, id string) (domain.Timelog, error)
	FindByFilter(ctx context.Context, userID string, filter domain.TimelogFilter) ([]domain.Timelog, error)
	FindRunning(ctx context.Context, userID string) (domain.Timelog, error)
}

// SettingRepository представляет репозиторий настроек.
type SettingRepository interface {
	InsertOrUpdate(ctx context.Context, setting domain.Setting) error
	FindAll(ctx context.Context, userID string) ([]domain.Setting, error)
	FindByKey(ctx context.Context, userID string, key string) (domain.Setting, error)
}

// TaskRepository представляет репозиторий задач.
type TaskRepository interface {
	FindByFilter(ctx context.Context, userID string, filter domain.TaskFilter) ([]domain.Task, error)
	FindByID(ctx context.Context, userID string, id string) (domain.Task, error)
}

// FavoriteRepository представляет репозиторий избранного.
type FavoriteRepository interface {
	FindAll(ctx context.Context, userID string) ([]domain.Favorite, error)
	FindByID(ctx context.Context, userID string, id string) (domain.Favorite, error)
}
