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

// ProjectHandler представляет работу с проектами.
type ProjectHandler struct {
	validate       *validator.Validate
	sanitizer      Sanitizer
	projectService ProjectService
	projectRepo    ProjectRepository

	logger *zap.Logger
}

// RegisterProjectEndpoints регистрирует конечные точки для работы с проектами.
func RegisterProjectEndpoints(
	router fiber.Router,
	validate *validator.Validate,
	sanitizer Sanitizer,
	projectService ProjectService,
	projectRepo ProjectRepository,
	logger *zap.Logger,
) {
	h := ProjectHandler{
		validate:       validate,
		sanitizer:      sanitizer,
		projectService: projectService,
		projectRepo:    projectRepo,

		logger: logger.Named("project_handler"),
	}

	g := router.Group("/projects")
	g.Get("/", h.Index)
	g.Post("/", h.Create)
	g.Get("/:id", h.Get)
	g.Patch("/:id", h.Update)
	g.Post("/:id/archive", h.Archive)
	g.Post("/:id/unarchive", h.Unarchive)
}

// Index обрабатывает GET /projects.
func (h *ProjectHandler) Index(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inp := &projectIndexFilter{}
	if err := c.QueryParser(inp); err != nil {
		return fmt.Errorf("query parser project index filter: %v", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	ps, err := h.projectRepo.FindByFilter(c.Context(), userID, domain.ProjectFilter{
		ClientID:               inp.ClientID,
		WithoutArchivedClients: inp.WithoutArchivedClients,
	})
	if err != nil {
		return fmt.Errorf("find projects: %v", err)
	}

	return c.JSON(ps)
}

// Create обрабатывает POST /projects.
func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inp := &projectCreateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser project create input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	p, err := h.projectService.Create(c.Context(), domain.Project{
		UserID:       userID,
		ClientID:     inp.ClientID,
		Name:         h.sanitizer.Sanitize(inp.Name),
		BillableRate: inp.BillableRate,
		Comment:      h.sanitizer.Sanitize(inp.Comment),
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return c.Status(fiber.StatusConflict).SendString("project with this name already exists")
		}
		if domain.IsFailedPreconditionError(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fmt.Errorf("create project: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(p)
}

// Get обрабатывает GET /projects/:id.
func (h *ProjectHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid project id")
	}

	p, err := h.projectRepo.FindByID(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("find project: %v", err)
	}

	return c.JSON(p)
}

// Update обрабатывает PATCH /projects/:id.
func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid project id")
	}

	inp := &projectUpdateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser project update input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	p, err := h.projectService.Update(c.Context(), domain.Project{
		ID:           id,
		UserID:       userID,
		ClientID:     inp.ClientID,
		Name:         h.sanitizer.Sanitize(inp.Name),
		BillableRate: inp.BillableRate,
		Comment:      h.sanitizer.Sanitize(inp.Comment),
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}
		if domain.IsFailedPreconditionError(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fmt.Errorf("update project: %v", err)
	}

	return c.JSON(p)
}

// Archive обрабатывает POST /projects/:id/archive.
func (h *ProjectHandler) Archive(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid project id")
	}

	p, err := h.projectService.Archive(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("archive project: %v", err)
	}

	return c.JSON(p)
}

// Unarchive обрабатывает POST /projects/:id/unarchive.
func (h *ProjectHandler) Unarchive(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid project id")
	}

	p, err := h.projectService.Unarchive(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("unarchive project: %v", err)
	}

	return c.JSON(p)
}
