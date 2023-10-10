package cookiemanager

import (
	"capuchin/app/util"

	"github.com/gofiber/fiber/v2"
)

type CookieManager struct {
	Key         string
	ExpiresDays uint
}

func New(key string, expiresDays uint) CookieManager {
	return CookieManager{
		Key:         key,
		ExpiresDays: expiresDays,
	}
}

func (cm *CookieManager) Create(c *fiber.Ctx, v string) {
	cookie := new(fiber.Cookie)
	cookie.Name = cm.Key
	cookie.Value = v
	cookie.Path = "/api/"
	cookie.Expires = util.TimeNowUtc().AddDate(0, 0, int(cm.ExpiresDays))
	cookie.Secure = true
	cookie.HTTPOnly = true
	cookie.SameSite = "Strict"

	c.Cookie(cookie)
}

func (cm *CookieManager) Get(c *fiber.Ctx) string {
	return c.Cookies(cm.Key)
}

func (cm *CookieManager) Clear(c *fiber.Ctx) {
	cookie := new(fiber.Cookie)
	cookie.Name = cm.Key
	cookie.Value = ""
    cookie.Path = "/api/"
	cookie.MaxAge = 0
	cookie.Secure = true
	cookie.HTTPOnly = true
	cookie.SameSite = "Strict"

	c.Cookie(cookie)
}
