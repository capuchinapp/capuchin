package cookiemanager

import (
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CookieManager представляет менеджер файлов cookie.
type CookieManager struct {
	name        string
	expiresDays int
	secure      bool
	domain      string

	re *regexp.Regexp

	nowFunc func() time.Time
}

// Новый CookieManager.
func New(name string, expiresDays int, secure bool, domain string) *CookieManager {
	return &CookieManager{
		name:        name,
		expiresDays: expiresDays,
		secure:      secure,
		domain:      domain,

		re: regexp.MustCompile(`(?m)` + name + `=(\w{26})`),

		nowFunc: time.Now,
	}
}

// Создать файл cookie.
func (cm *CookieManager) Create(c *fiber.Ctx, v string) {
	cookie := &fiber.Cookie{
		Name:     cm.name,
		Value:    v,
		Path:     "/",
		Domain:   cm.domain,
		Expires:  cm.nowFunc().UTC().AddDate(0, 0, cm.expiresDays),
		Secure:   cm.secure,
		HTTPOnly: true,
		SameSite: "Strict",
	}

	c.Cookie(cookie)
}

// Получить cookie.
func (cm *CookieManager) Get(c *fiber.Ctx) string {
	cookie := string(c.Request().Header.Peek("Cookie"))
	matches := cm.re.FindAllStringSubmatch(cookie, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return matches[0][1]
	}

	return ""
}

// GetExpiresDays возвращает количество дней до истечения срока действия файлов cookie.
func (cm *CookieManager) GetExpiresDays() int {
	return cm.expiresDays
}
