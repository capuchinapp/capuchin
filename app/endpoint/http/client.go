package http

import (
	"capuchin/app/entity"
	"capuchin/app/repository"
	"capuchin/app/util"
	"capuchin/app/util/nullable"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClientHandler struct {
	validate *validator.Validate
	repo     repository.ClientRepository
}

func RegisterClientEndpoints(r fiber.Router, validate *validator.Validate, cr repository.ClientRepository) {
	h := ClientHandler{
		validate: validate,
		repo:     cr,
	}

	g := r.Group("/clients")
	g.Get("/", h.Index)
	g.Post("/", h.Create)
	g.Get("/:uuid", h.Get)
	g.Patch("/:uuid", h.Update)
	g.Post("/:uuid/archive", h.Archive)
	g.Post("/:uuid/unarchive", h.Unarchive)
}

func (h *ClientHandler) Index(c *fiber.Ctx) error {
	cs, err := h.repo.FindAll()
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.JSON(cs)
}

type clientCreateInput struct {
	Name         string              `json:"name" validate:"required"`
	BillableRate uint64              `json:"billableRate" validate:"min=0"`
	Comment      nullable.NullString `json:"comment"`
}

func (h *ClientHandler) Create(c *fiber.Ctx) error {
	var inp *clientCreateInput = new(clientCreateInput)

	if err := c.BodyParser(inp); err != nil {
		return util.ErrTrace("BodyParser", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return util.ErrTrace("validateStruct", err)
	}

	id := uuid.New()

	client := entity.Client{
		UUID:         id,
		Name:         inp.Name,
		BillableRate: inp.BillableRate,
		Comment:      inp.Comment,
		CreatedAt:    util.TimeNowUtc(),
	}

	_, err := h.repo.Create(&client)
	if err != nil {
		return util.ErrTrace("Create", err)
	}

	cl, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.Status(fiber.StatusCreated).JSON(cl)
}

func (h *ClientHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	cl, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.JSON(cl)
}

type clientUpdateInput struct {
	Name         string              `json:"name" validate:"required"`
	BillableRate uint64              `json:"billableRate" validate:"min=0"`
	Comment      nullable.NullString `json:"comment"`
}

func (h *ClientHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	var inp *clientUpdateInput = new(clientUpdateInput)

	if err := c.BodyParser(inp); err != nil {
		return util.ErrTrace("BodyParser", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return util.ErrTrace("validateStruct", err)
	}

	cl, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	cl.Name = inp.Name
	cl.BillableRate = inp.BillableRate
	cl.Comment = inp.Comment
	cl.UpdatedAt = nullable.NewNullTime(util.TimeNowUtc(), true)

	_, err = h.repo.Update(&cl)
	if err != nil {
		return util.ErrTrace("Update", err)
	}

	return c.JSON(cl)
}

func (h *ClientHandler) Archive(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	cl, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	cl.ArchivedAt = nullable.NewNullTime(util.TimeNowUtc(), true)

	_, err = h.repo.Update(&cl)
	if err != nil {
		return util.ErrTrace("Update", err)
	}

	return c.JSON(cl)
}

func (h *ClientHandler) Unarchive(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	cl, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	cl.ArchivedAt = nullable.NewNullTime(util.TimeNowUtc(), false)

	_, err = h.repo.Update(&cl)
	if err != nil {
		return util.ErrTrace("Update", err)
	}

	return c.JSON(cl)
}
