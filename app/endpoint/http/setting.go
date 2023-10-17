package http

import (
	"capuchin/app/repository"
	"capuchin/app/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SettingHandler struct {
	validate *validator.Validate
	repo     repository.SettingRepository
}

func RegisterSettingEndpoints(r fiber.Router, validate *validator.Validate, sr repository.SettingRepository) {
	h := SettingHandler{
		validate: validate,
		repo:     sr,
	}

	g := r.Group("/settings")
	g.Get("/", h.Index)
	g.Put("/", h.Update)
}

func (h *SettingHandler) Index(c *fiber.Ctx) error {
	ss, err := h.repo.FindAll()
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.JSON(ss)
}

type settingUpdateInput struct {
	DateFormat string `json:"dateFormat" validate:"len=10"`
}

func (h *SettingHandler) Update(c *fiber.Ctx) error {
	var inp *settingUpdateInput = new(settingUpdateInput)

	if err := c.BodyParser(inp); err != nil {
		return util.ErrTrace("BodyParser", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return util.ErrTrace("validateStruct", err)
	}

	err := h.updateItem("dateFormat", inp.DateFormat)
	if err != nil {
		return util.ErrTrace("updateItem", err)
	}

	return h.Index(c)
}

func (h *SettingHandler) updateItem(key string, value string) error {
	r, err := h.repo.FindByKey(key)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	r.Value = value

	_, err = h.repo.Update(&r)
	if err != nil {
		return util.ErrTrace("Update", err)
	}

	return nil
}
