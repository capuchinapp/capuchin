package http

import (
	"capuchin/app/repository"
	"capuchin/app/util"
	"capuchin/app/util/cookiemanager"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SessionHandler struct {
	validate  *validator.Validate
	repo      repository.SessionRepository
	cm        cookiemanager.CookieManager
	debugMode *bool
}

func RegisterSessionEndpoints(r fiber.Router, validate *validator.Validate, sr repository.SessionRepository, cm cookiemanager.CookieManager, debugMode *bool) {
	h := SessionHandler{
		validate:  validate,
		repo:      sr,
		cm:        cm,
		debugMode: debugMode,
	}

	g := r.Group("/sessions")
	g.Get("/", h.Index)
	g.Delete("/:rowid", h.Delete)
}

func (h *SessionHandler) Index(c *fiber.Ctx) error {
	outdatedDateTime := util.TimeNowUtc().AddDate(0, 0, -1*int(h.cm.ExpiresDays)).String()
	if _, err := h.repo.CancelOld(outdatedDateTime); err != nil {
		return util.ErrTrace("CancelOld", err)
	}

	ss, err := h.repo.FindAll()
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	if *h.debugMode {
		ss[0].IsCurrent = true
	} else {
		cookieID := h.cm.Get(c)

		for k, v := range ss {
			if v.UUID.String() == cookieID {
				ss[k].IsCurrent = true
			}
		}
	}

	return c.JSON(ss)
}

func (h *SessionHandler) Delete(c *fiber.Ctx) error {
	rowid, err := strconv.Atoi((c.Params("rowid", "0")))
	if err != nil {
		return util.ErrTrace("Atoi", err)
	}

	s, err := h.repo.FindByRowId(int64(rowid))
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	if !*h.debugMode {
		cookieID := h.cm.Get(c)
		if s.UUID.String() == cookieID {
			return c.SendStatus(fiber.StatusForbidden)
		}
	}

	if _, err := h.repo.DeleteByRowId(s.RowID.Int64); err != nil {
		return util.ErrTrace("DeleteByRowId", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
