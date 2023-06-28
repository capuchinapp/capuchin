package session

import (
	"capuchin/app/repository"
	"capuchin/app/util"
	"capuchin/app/util/cookiemanager"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Config struct {
	Paths             []string
	CookieManager     cookiemanager.CookieManager
	SessionRepository repository.SessionRepository
}

func New(cfg Config, debugMode *bool, authMode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		found := false
		for _, v := range cfg.Paths {
			if strings.HasPrefix(path, v) {
				found = true
			}
		}

		if !found || *debugMode || authMode == "0" {
			return c.Next()
		}

		id, err := uuid.Parse(cfg.CookieManager.Get(c))
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		s, err := cfg.SessionRepository.FindByUuid(id)
		if err != nil {
			cfg.CookieManager.Clear(c)

			return c.SendStatus(fiber.StatusUnauthorized)
		}

		s.CheckedAt = util.TimeNowUtc()

		_, err = cfg.SessionRepository.UpdateCheckedAt(&s)
		if err != nil {
			return util.ErrTrace("UpdateCheckedAt", err)
		}

		return c.Next()
	}
}
