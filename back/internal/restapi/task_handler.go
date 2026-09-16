package restapi

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	internalvalidator "capuchin/internal/validator"
)

// TaskHandler представляет работу с задачами.
type TaskHandler struct {
	validate    *validator.Validate
	sanitizer   Sanitizer
	taskService TaskService
	taskRepo    TaskRepository

	logger *zap.Logger
}

// RegisterTaskEndpoints регистрирует конечные точки для работы с задачами.
func RegisterTaskEndpoints(
	router fiber.Router,
	validate *validator.Validate,
	sanitizer Sanitizer,
	taskService TaskService,
	taskRepo TaskRepository,
	logger *zap.Logger,
) {
	h := TaskHandler{
		validate:    validate,
		sanitizer:   sanitizer,
		taskService: taskService,
		taskRepo:    taskRepo,

		logger: logger.Named("task_handler"),
	}

	g := router.Group("/tasks")
	g.Get("/", h.Index)
	g.Post("/", h.Create)
	g.Get("/:id", h.Get)
	g.Patch("/:id", h.Update)
	g.Post("/:id/archive", h.Archive)
	g.Post("/:id/unarchive", h.Unarchive)
	g.Post("/:id/complete", h.Complete)
	g.Post("/:id/incomplete", h.Incomplete)
	g.Get("/:id/report", h.Report)
}

// Index обрабатывает GET /tasks.
func (h *TaskHandler) Index(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inp := &taskIndexFilter{}
	if err := c.QueryParser(inp); err != nil {
		return fmt.Errorf("query parser task index filter: %v", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	res, err := h.taskRepo.FindByFilter(c.Context(), userID, domain.TaskFilter{
		ProjectID:               inp.ProjectID,
		WithoutArchivedProjects: inp.WithoutArchivedProjects,
	})
	if err != nil {
		return fmt.Errorf("find tasks: %v", err)
	}

	return c.JSON(res)
}

// Create обрабатывает POST /tasks.
func (h *TaskHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inp := &taskCreateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser task create input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	res, err := h.taskService.Create(c.Context(), domain.Task{
		UserID:    userID,
		ProjectID: inp.ProjectID,
		Name:      h.sanitizer.Sanitize(inp.Name),
		Comment:   h.sanitizer.Sanitize(inp.Comment),
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return c.Status(fiber.StatusConflict).SendString("task with this name already exists")
		}
		if domain.IsFailedPreconditionError(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fmt.Errorf("create task: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Get обрабатывает GET /tasks/:id.
func (h *TaskHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	res, err := h.taskRepo.FindByID(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("find task: %v", err)
	}

	return c.JSON(res)
}

// Update обрабатывает PATCH /tasks/:id.
func (h *TaskHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	inp := &taskUpdateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser task update input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	res, err := h.taskService.Update(c.Context(), domain.Task{
		ID:        id,
		UserID:    userID,
		ProjectID: inp.ProjectID,
		Name:      h.sanitizer.Sanitize(inp.Name),
		Comment:   h.sanitizer.Sanitize(inp.Comment),
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}
		if domain.IsFailedPreconditionError(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fmt.Errorf("update task: %v", err)
	}

	return c.JSON(res)
}

// Archive обрабатывает POST /tasks/:id/archive.
func (h *TaskHandler) Archive(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	res, err := h.taskService.Archive(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}
		if errors.Is(err, domain.ErrImpossibleAction) {
			return c.Status(fiber.StatusConflict).SendString(domain.ErrImpossibleAction.Error())
		}

		return fmt.Errorf("archive task: %v", err)
	}

	return c.JSON(res)
}

// Unarchive обрабатывает POST /tasks/:id/unarchive.
func (h *TaskHandler) Unarchive(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	res, err := h.taskService.Unarchive(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("unarchive task: %v", err)
	}

	return c.JSON(res)
}

// Complete обрабатывает POST /tasks/:id/complete.
func (h *TaskHandler) Complete(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	res, err := h.taskService.Complete(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}
		if errors.Is(err, domain.ErrImpossibleAction) {
			return c.Status(fiber.StatusConflict).SendString(domain.ErrImpossibleAction.Error())
		}

		return fmt.Errorf("complete task: %v", err)
	}

	return c.JSON(res)
}

// Incomplete обрабатывает POST /tasks/:id/incomplete.
func (h *TaskHandler) Incomplete(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	res, err := h.taskService.Incomplete(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("incomplete task: %v", err)
	}

	return c.JSON(res)
}

// Report обрабатывает GET /tasks/:id/report.
func (h *TaskHandler) Report(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}

	res, err := h.taskService.Report(c.Context(), userID, id)
	if err != nil {
		return fmt.Errorf("report time task: %v", err)
	}

	return c.JSON(res)
}
