package middleware

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	"capuchin/internal/restapi"
)

// AuthenticationConfig это конфигурация промежуточного программного обеспечения аутентификации.
type AuthenticationConfig struct {
	CookieManager  CookieManager
	SessionRepo    SessionRepo
	SessionCache   SessionCache
	SessionUpdater SessionUpdater
	Logger         Logger
	NowFunc        func() time.Time
	Paths          []string
}

// NewAuthentication создает новое промежуточное программное обеспечение аутентификации.
func NewAuthentication(cfg AuthenticationConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		found := false
		for _, v := range cfg.Paths {
			if strings.HasPrefix(path, v) {
				found = true
			}
		}

		if !found {
			return c.Next()
		}

		cid := cfg.CookieManager.Get(c)
		if !domain.IsID(cid) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		ctx := c.Context()

		s, ok := getSession(ctx, cid, cfg.SessionCache, cfg.SessionRepo, cfg.Logger)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Обновляем время в кэше (реальное обновление в БД произойдет позже)
		cfg.SessionUpdater.UpdateTime(s.ID)

		c.Locals(restapi.LocalsUserIDKey, s.UserID)

		return c.Next()
	}
}

func getSession(
	ctx context.Context,
	cid string,
	sessionCache SessionCache,
	sessionRepo SessionRepo,
	log Logger,
) (domain.Session, bool) {
	s, ok := sessionCache.Get(cid)
	if !ok {
		var err error

		s, err = sessionRepo.FindByCID(ctx, cid)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return domain.Session{}, false
			}

			log.Error("find session by cid",
				zap.Error(err),
			)
			return domain.Session{}, false
		}

		sessionCache.Add(cid, s)
	}

	return s, true
}
