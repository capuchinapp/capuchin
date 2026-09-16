package restapi

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"capuchin/internal/domain"
)

// SessionHandler представляет работу с сессиями.
type SessionHandler struct {
	sessionService SessionService
	sessionRepo    SessionRepository
	cookieManager  CookieManager

	nowFunc func() time.Time
}

// RegisterSessionEndpoints регистрирует конечные точки для сессий.
func RegisterSessionEndpoints(
	router fiber.Router,
	sessionService SessionService,
	sessionRepo SessionRepository,
	cookieManager CookieManager,
) {
	h := SessionHandler{
		sessionService: sessionService,
		sessionRepo:    sessionRepo,
		cookieManager:  cookieManager,

		nowFunc: time.Now,
	}

	g := router.Group("/sessions")
	g.Get("/", h.Index)
	g.Delete("/:id", h.Delete)
}

// Index обрабатывает GET /sessions.
func (h *SessionHandler) Index(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	outdatedTime := h.nowFunc().AddDate(0, 0, -1*h.cookieManager.GetExpiresDays())
	ss, err := h.sessionService.List(c.Context(), userID, outdatedTime)
	if err != nil {
		return fmt.Errorf("list sessions: %v", err)
	}

	cookieID := h.cookieManager.Get(c)
	for k, v := range ss {
		if v.CookieID == cookieID {
			ss[k].IsCurrent = true
		}
	}

	return c.JSON(ss)
}

// Delete обрабатывает DELETE /sessions/:id.
func (h *SessionHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	s, err := h.sessionRepo.FindByID(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.SendStatus(fiber.StatusOK)
		}

		return fmt.Errorf("find by row id: %v", err)
	}

	cookieID := h.cookieManager.Get(c)
	if s.CookieID == cookieID {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := h.sessionRepo.DeleteByCID(c.Context(), s.CookieID); err != nil {
		return fmt.Errorf("delete by id: %v", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
