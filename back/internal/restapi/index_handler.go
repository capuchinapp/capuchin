package restapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"capuchin/internal/domain"
)

// IndexHandler представляет работу с главной страницей.
type IndexHandler struct {
	appVersion    string
	cookieManager CookieManager
	sessionRepo   SessionRepository
	timelogRepo   TimelogRepository

	logger *zap.Logger
}

// RegisterIndexEndpoints регистрирует конечную точку для главной страницы.
func RegisterIndexEndpoints(
	router fiber.Router,
	appVersion string,
	cookieManager CookieManager,
	sessionRepo SessionRepository,
	timelogRepo TimelogRepository,
	logger *zap.Logger,
) {
	h := IndexHandler{
		appVersion:    appVersion,
		cookieManager: cookieManager,
		sessionRepo:   sessionRepo,
		timelogRepo:   timelogRepo,

		logger: logger.Named("index_handler"),
	}

	router.Get("/", h.Index)
}

// Index обрабатывает GET /.
func (h *IndexHandler) Index(c *fiber.Ctx) error {
	session, isAuth := h.checkAuth(c)

	var runningTimelogDatetime *string
	if isAuth {
		runningTimelogDatetime = h.runningTimelogDatetime(c.Context(), session.UserID)
	}

	return c.JSON(index{
		Name:                   "Capuchin API",
		IsAuth:                 isAuth,
		AppVersion:             h.appVersion,
		RunningTimelogDatetime: runningTimelogDatetime,
	})
}

func (h *IndexHandler) checkAuth(c *fiber.Ctx) (session domain.Session, ok bool) {
	cid := h.cookieManager.Get(c)
	if !domain.IsID(cid) {
		return domain.Session{}, false
	}

	session, err := h.sessionRepo.FindByCID(c.Context(), cid)
	if err == nil && session.CookieID == cid {
		return session, true
	}

	return domain.Session{}, false
}

func (h *IndexHandler) runningTimelogDatetime(ctx context.Context, userID string) *string {
	runningTimelog, err := h.timelogRepo.FindRunning(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}

		h.logger.Warn("failed to find running timelog", zap.Error(err))

		return nil
	}

	dt := fmt.Sprintf("%s %s", runningTimelog.Date, runningTimelog.TimeStart)

	return &dt
}
