package restapi

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	internalvalidator "capuchin/internal/validator"
)

const (
	defaultAuthLimiterMax        = 3
	defaultAuthLimiterExpiration = 30 * time.Second
)

// AuthHandler представляет работу авторизации.
type AuthHandler struct {
	validate      *validator.Validate
	authService   AuthService
	cookieManager CookieManager
	sessionRepo   SessionRepository

	logger *zap.Logger

	idFunc   func() string
	nowFunc  func() time.Time
	randFunc func(l int) string
}

// RegisterAuthEndpoints регистрирует конечные точки для авторизации.
func RegisterAuthEndpoints(
	router fiber.Router,
	validate *validator.Validate,
	authService AuthService,
	cookieManager CookieManager,
	sessionRepo SessionRepository,
	logger *zap.Logger,
) {
	h := AuthHandler{
		validate:      validate,
		authService:   authService,
		cookieManager: cookieManager,
		sessionRepo:   sessionRepo,

		logger: logger.Named("auth_handler"),

		idFunc:   domain.NewID,
		nowFunc:  time.Now,
		randFunc: domain.NewString,
	}

	g := router.Group("/auth")

	g.Use(limiter.New(limiter.Config{
		Max:        defaultAuthLimiterMax,
		Expiration: defaultAuthLimiterExpiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	g.Post("/register", h.Register)
	g.Post("/activate", h.Activate)
	g.Post("/login", h.Login)
	g.Post("/apply_code", h.ApplyCode)
	g.Post("/logout", h.Logout)
}

// Register обрабатывает POST /auth/register.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	inp := &authRegisterInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser auth register input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	ip := c.IP()

	err := h.authService.Register(c.Context(), ip, inp.Email)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return c.Status(fiber.StatusConflict).SendString("user with this email already exists")
		}

		return fmt.Errorf("register: %v", err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

// Activate обрабатывает POST /auth/activate.
func (h *AuthHandler) Activate(c *fiber.Ctx) error {
	inp := &authActivateInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser auth activate input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	id, err := h.authService.Activate(c.Context(), inp.UserID, inp.Code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}

		return fmt.Errorf("activate user: %v", err)
	}

	h.cookieManager.Create(c, id)

	return c.SendStatus(fiber.StatusCreated)
}

// Login обрабатывает POST /auth/login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	inp := &authLoginInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser auth login input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	err := h.authService.Login(c.Context(), inp.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fiber.ErrNotFound
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			return c.SendStatus(fiber.StatusConflict)
		}
		if errors.Is(err, domain.ErrUserNotActivated) {
			return c.Status(fiber.StatusConflict).SendString(domain.ErrUserNotActivated.Error())
		}

		return fmt.Errorf("login: %v", err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// ApplyCode обрабатывает POST /auth/apply_code.
func (h *AuthHandler) ApplyCode(c *fiber.Ctx) error {
	inp := &authApplyCodeInput{}
	if err := c.BodyParser(inp); err != nil {
		h.logger.Error("body parser auth apply code input", zap.Error(err))

		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.validate.Struct(inp); err != nil {
		return internalvalidator.ParseError(err)
	}

	id, err := h.authService.ApplyCode(c.Context(), inp.Email, inp.Code)
	if err != nil {
		return fmt.Errorf("apply code: %v", err)
	}

	h.cookieManager.Create(c, id)

	return c.SendStatus(fiber.StatusCreated)
}

// Logout обрабатывает POST /auth/logout.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	cid := h.cookieManager.Get(c)
	if cid == "" {
		return c.SendStatus(fiber.StatusOK)
	}

	if !domain.IsID(cid) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid session id")
	}

	if err := h.sessionRepo.DeleteByCID(c.Context(), cid); err != nil {
		return fmt.Errorf("delete by id: %v", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
