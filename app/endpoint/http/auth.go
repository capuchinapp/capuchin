package http

import (
	"capuchin/app/entity"
	"capuchin/app/repository"
	"capuchin/app/util"
	"capuchin/app/util/cookiemanager"
	"capuchin/app/util/nullable"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	validate *validator.Validate
	repo     repository.SessionRepository
	cm       cookiemanager.CookieManager
	pwd      string
}

func RegisterAuthEndpoints(r fiber.Router, validate *validator.Validate, sr repository.SessionRepository, cm cookiemanager.CookieManager, pwd string) {
	h := AuthHandler{
		validate: validate,
		repo:     sr,
		cm:       cm,
		pwd:      pwd,
	}

	g := r.Group("/auth")
	g.Post("/login", h.Login)
	g.Post("/logout", h.Logout)
}

type authLoginInput struct {
	Password string `json:"password" validate:"min=8"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var inp *authLoginInput = new(authLoginInput)

	if err := c.BodyParser(inp); err != nil {
		return util.ErrTrace("BodyParser", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return util.ErrTrace("validateStruct", err)
	}

	if inp.Password != h.pwd {
		return c.Status(fiber.StatusBadRequest).SendString("Wrong password")
	}

	id := uuid.New()

	session := entity.Session{
		RowID:     nullable.NewNullInt64(0, false),
		UUID:      id,
		CheckedAt: util.TimeNowUtc(),
	}

	if _, err := h.repo.Create(&session); err != nil {
		return util.ErrTrace("Create", err)
	}

	if _, err := h.repo.FindByUuid(id); err != nil {
		return util.ErrSql2Fiber(err)
	}

	h.cm.Create(c, id.String())

	return c.SendStatus(fiber.StatusCreated)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	cid := h.cm.Get(c)
	if cid == "" {
		return c.SendStatus(fiber.StatusOK)
	}

	id, err := uuid.Parse(cid)
	if err != nil {
		return util.ErrTrace("Parse", err)
	}

	_, _ = h.repo.DeleteByUuid(id)
	h.cm.Clear(c)

	return c.SendStatus(fiber.StatusOK)
}
