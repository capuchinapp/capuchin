package middleware

import (
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var ulidRegex = regexp.MustCompile(`[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}`)

// NewRequestLog создает новое промежуточное программное обеспечение для регистрации запросов.
func NewRequestLog(log Logger, metrics Metrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		metrics.HTTPRequestsCurrentInc()
		err := c.Next()
		metrics.HTTPRequestsCurrentDec()

		duration := time.Since(start)

		method := c.Method()
		path := cleanPath(c.Path())
		status := c.Response().StatusCode()

		metrics.HTTPRequestsTotalInc(method, path, strconv.Itoa(status))
		metrics.HTTPRequestDurationSecondsObserve(method, path, duration.Seconds())

		log.Info(
			"Request",
			zap.String("ip", c.IP()),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("protocol", c.Protocol()),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("user-agent", c.Get("User-Agent")),
		)

		return err
	}
}

func cleanPath(s string) string {
	return ulidRegex.ReplaceAllString(s, ":id")
}
