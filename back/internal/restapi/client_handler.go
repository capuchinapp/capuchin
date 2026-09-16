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

// ClientHandler представляет работу с клиентами.
type ClientHandler struct {
	validate      *validator.Validate
	sanitizer     Sanitizer
	clientService ClientService
	clientRepo    ClientRepository

	logger *zap.Logger
}

// RegisterClientEndpoints регистрирует конечные точки для работы с клиентами.
func RegisterClientEndpoints(
	router fiber.Router,
	validate *validator.Validate,
	sanitizer Sanitizer,
	clientService ClientService,
	clientRepo ClientRepository,
	logger *zap.Logger,
) {
	h := ClientHandler{
		validate:      validate,
		sanitizer:     sanitizer,
		clientService: clientService,
		clientRepo:    clientRepo,

		logger: logger.Named("client_handler"),
	}

	g := router.Group("/clients")
	g.Get("/", h.Index)
	g.Post("/", h.Create)
	g.Get("/:id", h.Get)
	g.Patch("/:id", h.Update)
	g.Post("/:id/archive", h.Archive)
	g.Post("/:id/unarchive", h.Unarchive)
}

// Index обрабатывает GET /clients.
func (h *ClientHandler) Index(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	cs, err := h.clientRepo.FindAll(c.Context(), userID)
	if err != nil {
		return fmt.Errorf("find clients: %v", err)
	}

	return c.JSON(cs)
}

// Create обрабатывает POST /clients.
func (h *ClientHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inp := &clientCreateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser client create input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	cl, err := h.clientService.Create(c.Context(), domain.Client{
		UserID:       userID,
		Name:         h.sanitizer.Sanitize(inp.Name),
		BillableRate: inp.BillableRate,
		Comment:      h.sanitizer.Sanitize(inp.Comment),
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return c.Status(fiber.StatusConflict).SendString("client with this name already exists")
		}

		return fmt.Errorf("create client: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(cl)
}

// Get обрабатывает GET /clients/:id.
func (h *ClientHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid client id")
	}

	cl, err := h.clientRepo.FindByID(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("find client: %v", err)
	}

	return c.JSON(cl)
}

// Update обрабатывает PATCH /clients/:id.
func (h *ClientHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid client id")
	}

	inp := &clientUpdateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser client update input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	cl, err := h.clientService.Update(c.Context(), domain.Client{
		ID:           id,
		UserID:       userID,
		Name:         h.sanitizer.Sanitize(inp.Name),
		BillableRate: inp.BillableRate,
		Comment:      h.sanitizer.Sanitize(inp.Comment),
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("update client: %v", err)
	}

	return c.JSON(cl)
}

// Archive обрабатывает POST /clients/:id/archive.
func (h *ClientHandler) Archive(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid client id")
	}

	cl, err := h.clientService.Archive(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("archive client: %v", err)
	}

	return c.JSON(cl)
}

// Unarchive обрабатывает POST /clients/:id/unarchive.
func (h *ClientHandler) Unarchive(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid client id")
	}

	cl, err := h.clientService.Unarchive(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("unarchive client: %v", err)
	}

	return c.JSON(cl)
}
