package restapi

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	internalvalidator "capuchin/internal/validator"
)

// TimelogHandler представляет работу с записями в журнале.
type TimelogHandler struct {
	validate       *validator.Validate
	sanitizer      Sanitizer
	timelogService TimelogService
	timelogRepo    TimelogRepository

	logger *zap.Logger
}

// RegisterTimelogEndpoints регистрирует конечные точки для работы с журналом.
func RegisterTimelogEndpoints(
	router fiber.Router,
	validate *validator.Validate,
	sanitizer Sanitizer,
	timelogService TimelogService,
	timelogRepo TimelogRepository,
	logger *zap.Logger,
) {
	h := TimelogHandler{
		validate:       validate,
		sanitizer:      sanitizer,
		timelogService: timelogService,
		timelogRepo:    timelogRepo,

		logger: logger.Named("timelog_handler"),
	}

	g := router.Group("/timelogs")
	g.Get("/", h.Index)
	g.Post("/", h.Create)
	g.Get("/:id", h.Get)
	g.Patch("/:id", h.Update)
	g.Post("/:id/stop", h.Stop)
	g.Delete("/:id", h.Delete)
	g.Get("/last/:n", h.LastN)
}

// Index обрабатывает GET /timelogs.
func (h *TimelogHandler) Index(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inp := &timelogIndexFilter{}
	if err := c.QueryParser(inp); err != nil {
		return fmt.Errorf("query parser timelog index filter: %v", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	tls, err := h.timelogRepo.FindByFilter(c.Context(), userID, domain.TimelogFilter{
		DateFrom:  inp.DateFrom,
		DateTo:    inp.DateTo,
		ClientID:  inp.ClientID,
		ProjectID: inp.ProjectID,
	})
	if err != nil {
		return fmt.Errorf("find timelogs: %v", err)
	}

	return c.JSON(tls)
}

// Create обрабатывает POST /timelogs.
func (h *TimelogHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inp := &timelogCreateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser timelog create input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	tl, err := h.timelogService.Create(c.Context(), domain.Timelog{
		UserID:       userID,
		ProjectID:    inp.ProjectID,
		Date:         inp.Date,
		TimeStart:    inp.TimeStart,
		TimeEnd:      inp.TimeEnd,
		BillableRate: inp.BillableRate,
		Comment:      h.sanitizer.Sanitize(inp.Comment),
		TaskID:       inp.TaskID,
	})
	if err != nil {
		if domain.IsFailedPreconditionError(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fmt.Errorf("create timelog: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(tl)
}

// Get обрабатывает GET /timelogs/:id.
func (h *TimelogHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid timelog id")
	}

	tl, err := h.timelogRepo.FindByID(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("find timelog: %v", err)
	}

	return c.JSON(tl)
}

// Update обрабатывает PUT /timelogs/:id.
func (h *TimelogHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid timelog id")
	}

	inp := &timelogUpdateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser timelog update input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	tl, err := h.timelogService.Update(c.Context(), domain.Timelog{
		ID:           id,
		UserID:       userID,
		ProjectID:    inp.ProjectID,
		Date:         inp.Date,
		TimeStart:    inp.TimeStart,
		TimeEnd:      &inp.TimeEnd,
		BillableRate: inp.BillableRate,
		Comment:      h.sanitizer.Sanitize(inp.Comment),
		TaskID:       inp.TaskID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}
		if domain.IsFailedPreconditionError(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fmt.Errorf("update timelog: %v", err)
	}

	return c.JSON(tl)
}

// Delete обрабатывает DELETE /timelogs/:id.
func (h *TimelogHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid timelog id")
	}

	if err := h.timelogService.Delete(c.Context(), userID, id); err != nil {
		return fmt.Errorf("delete timelog: %v", err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// Stop обрабатывает POST /timelogs/:id/stop.
func (h *TimelogHandler) Stop(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid timelog id")
	}

	inp := &timelogStopInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser timelog stop input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	tl, err := h.timelogService.Stop(c.Context(), userID, id, inp.Date, inp.TimeEnd)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}
		if domain.IsFailedPreconditionError(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fmt.Errorf("stop timelog: %v", err)
	}

	return c.JSON(tl)
}

// LastN обрабатывает GET /timelogs/last/:n.
func (h *TimelogHandler) LastN(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inputN := c.Params("n", "")
	n, err := strconv.Atoi(inputN)
	if err != nil {
		h.logger.Error("invalid n",
			zap.String("input_n", inputN),
			zap.Error(err),
		)

		return fiber.NewError(fiber.StatusBadRequest, "invalid n parameter")
	}

	if n <= 0 {
		h.logger.Error("n must be greater than zero", zap.Int("n", n))

		return fiber.NewError(fiber.StatusBadRequest, "n must be greater than zero")
	}

	tls, err := h.timelogService.FindLastN(c.Context(), userID, n)
	if err != nil {
		return fmt.Errorf("find last timelogs: %v", err)
	}

	return c.JSON(tls)
}
