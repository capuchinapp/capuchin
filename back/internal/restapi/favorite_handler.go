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

// FavoriteHandler представляет работу с избранным.
type FavoriteHandler struct {
	validate        *validator.Validate
	sanitizer       Sanitizer
	favoriteService FavoriteService
	favoriteRepo    FavoriteRepository

	logger *zap.Logger
}

// RegisterFavoriteEndpoints регистрирует конечные точки для работы с избранным.
func RegisterFavoriteEndpoints(
	router fiber.Router,
	validate *validator.Validate,
	sanitizer Sanitizer,
	favoriteService FavoriteService,
	favoriteRepo FavoriteRepository,
	logger *zap.Logger,
) {
	h := FavoriteHandler{
		validate:        validate,
		sanitizer:       sanitizer,
		favoriteService: favoriteService,
		favoriteRepo:    favoriteRepo,

		logger: logger.Named("favorite_handler"),
	}

	g := router.Group("/favorites")
	g.Get("/", h.Index)
	g.Post("/", h.Create)
	g.Get("/:id", h.Get)
	g.Patch("/:id", h.Update)
	g.Delete("/:id", h.Delete)
}

// Index обрабатывает GET /favorites.
func (h *FavoriteHandler) Index(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	items, err := h.favoriteRepo.FindAll(c.Context(), userID)
	if err != nil {
		return fmt.Errorf("find favorites: %v", err)
	}

	return c.JSON(items)
}

// Create обрабатывает POST /favorites.
func (h *FavoriteHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inp := &favoriteInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser favorite create input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	item, err := h.favoriteService.Create(c.Context(), domain.Favorite{
		UserID:       userID,
		Name:         h.sanitizer.Sanitize(inp.Name),
		ProjectID:    inp.ProjectID,
		TaskID:       inp.TaskID,
		BillableRate: inp.BillableRate,
		Comment:      h.sanitizer.Sanitize(inp.Comment),
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return c.Status(fiber.StatusConflict).SendString("favorite with this name already exists")
		}
		if domain.IsFailedPreconditionError(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fmt.Errorf("create favorite: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

// Get обрабатывает GET /favorites/:id.
func (h *FavoriteHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid favorite id")
	}

	item, err := h.favoriteRepo.FindByID(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("find favorite: %v", err)
	}

	return c.JSON(item)
}

// Update обрабатывает PUT /favorites/:id.
func (h *FavoriteHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid favorite id")
	}

	inp := &favoriteInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser favorite update input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	item, err := h.favoriteService.Update(c.Context(), domain.Favorite{
		ID:           id,
		UserID:       userID,
		Name:         h.sanitizer.Sanitize(inp.Name),
		ProjectID:    inp.ProjectID,
		TaskID:       inp.TaskID,
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

		return fmt.Errorf("update favorite: %v", err)
	}

	return c.JSON(item)
}

// Delete обрабатывает DELETE /favorites/:id.
func (h *FavoriteHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	id := c.Params("id", "")
	if !domain.IsID(id) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid favorite id")
	}

	if err := h.favoriteService.Delete(c.Context(), userID, id); err != nil {
		return fmt.Errorf("delete favorite: %v", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
