package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap/zapcore"

	"capuchin/internal/domain"
)

// SessionRepo представляет репозиторий сессий.
type SessionRepo interface {
	FindByCID(ctx context.Context, cid string) (domain.Session, error)
}

// SessionCache представляет кеш сессий.
type SessionCache interface {
	Add(key string, value domain.Session) (evicted bool)
	Get(key string) (value domain.Session, ok bool)
}

// SessionUpdater представляет инструмент обновления данных о сессии.
type SessionUpdater interface {
	UpdateTime(sessionID string)
}

// CookieManager представляет менеджер файлов cookie.
type CookieManager interface {
	Get(c *fiber.Ctx) string
}

// Logger представляет регистратор логов.
type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
}

// Metrics представляет регистратор метрик.
type Metrics interface {
	HTTPRequestsCurrentInc()
	HTTPRequestsCurrentDec()
	HTTPRequestsTotalInc(method string, path string, status string)
	HTTPRequestDurationSecondsObserve(method string, path string, elapsedSeconds float64)
}
