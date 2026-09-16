package restapi

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	internalvalidator "capuchin/internal/validator"
)

// SettingHandler представляет работу с настройками.
type SettingHandler struct {
	validate    *validator.Validate
	sanitizer   Sanitizer
	settingRepo SettingRepository

	logger *zap.Logger
}

// RegisterSettingEndpoints регистрирует конечные точки для работы с настройками.
func RegisterSettingEndpoints(
	router fiber.Router,
	validate *validator.Validate,
	sanitizer Sanitizer,
	settingRepo SettingRepository,
	logger *zap.Logger,
) {
	h := SettingHandler{
		validate:    validate,
		sanitizer:   sanitizer,
		settingRepo: settingRepo,

		logger: logger.Named("setting_handler"),
	}

	g := router.Group("/settings")
	g.Get("/", h.Index)
	g.Put("/", h.Update)
}

// Index обрабатывает GET /settings.
func (h *SettingHandler) Index(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	ss, err := h.settingRepo.FindAll(c.Context(), userID)
	if err != nil {
		return fmt.Errorf("find settings: %v", err)
	}

	return c.JSON(ss)
}

// Update обрабатывает PUT /settings.
func (h *SettingHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals(LocalsUserIDKey).(string)
	if !ok {
		return ErrGetUserID
	}

	inp := &settingUpdateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser setting update input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	if err := h.settingRepo.InsertOrUpdate(c.Context(), domain.Setting{
		UserID: userID,
		Key:    "dateFormat",
		Value:  h.sanitizer.Sanitize(inp.DateFormat),
	}); err != nil {
		return fmt.Errorf("insert or update setting dateFormat: %v", err)
	}

	if err := h.settingRepo.InsertOrUpdate(c.Context(), domain.Setting{
		UserID: userID,
		Key:    "workingDays",
		Value:  h.sanitizer.Sanitize(inp.WorkingDays),
	}); err != nil {
		return fmt.Errorf("insert or update setting workingDays: %v", err)
	}

	return h.Index(c)
}
