package http

import (
	"capuchin/app/repository"
	"capuchin/app/util/cookiemanager"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterIndexEndpoints(a *fiber.App, cm cookiemanager.CookieManager, ssr repository.SessionRepository, auth string, debugMode *bool) {
	a.Get("/api", func(c *fiber.Ctx) error {
		var isAuth bool

		if *debugMode {
			isAuth = true
		} else {
			isAuth = checkAuth(c, cm, ssr)
		}

		return c.JSON(fiber.Map{
			"app":      "Capuchin API",
			"authMode": auth == "1",
			"isAuth":   isAuth,
		})
	})
}

func checkAuth(c *fiber.Ctx, cm cookiemanager.CookieManager, ssr repository.SessionRepository) bool {
	id, err := uuid.Parse(cm.Get(c))
	if err != nil {
		return false
	}

	if _, err := ssr.FindByUuid(id); err != nil {
		cm.Clear(c)

		return false
	}

	return true
}
