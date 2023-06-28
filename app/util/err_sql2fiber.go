package util

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrSql2Fiber(e error) error {
	switch {
	case errors.Is(e, sql.ErrNoRows):
		return fiber.ErrNotFound
	default:
		return e
	}
}
