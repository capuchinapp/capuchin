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

type ProjectHandler struct {
	validate *validator.Validate
	repo     repository.ProjectRepository
}

func RegisterProjectEndpoints(r fiber.Router, validate *validator.Validate, pr repository.ProjectRepository) {
	h := ProjectHandler{
		validate: validate,
		repo:     pr,
	}

	g := r.Group("/projects")
	g.Get("/", h.Index)
	g.Post("/", h.Create)
	g.Get("/:uuid", h.Get)
	g.Patch("/:uuid", h.Update)
	g.Post("/:uuid/archive", h.Archive)
	g.Post("/:uuid/unarchive", h.Unarchive)
}

type projectIndexFilter struct {
	ClientUUID             string `query:"client_uuid"`
	WithoutArchivedClients bool   `query:"filter_archived_clients"`
}

func (h *ProjectHandler) Index(c *fiber.Ctx) error {
	pif := new(projectIndexFilter)
	if err := c.QueryParser(pif); err != nil {
		return util.ErrTrace("QueryParser", err)
	}

	if pif.ClientUUID != "" {
		id, err := uuid.Parse(pif.ClientUUID)
		if err != nil {
			return util.ErrTrace("uuid.Parse", err)
		}

		ps, err := h.repo.FindByClientId(id)
		if err != nil {
			return util.ErrSql2Fiber(err)
		}

		return c.JSON(ps)
	}

	ps, err := h.repo.FindAll(entity.ProjectFilter{
		WithoutArchivedClients: pif.WithoutArchivedClients,
	})
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.JSON(ps)
}

type projectCreateInput struct {
	ClientUUID   uuid.UUID           `json:"clientUUID" validate:"required"`
	Name         string              `json:"name" validate:"required"`
	BillableRate uint64              `json:"billableRate" validate:"min=0"`
	Comment      nullable.NullString `json:"comment"`
}

func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	inp := new(projectCreateInput)
	if err := c.BodyParser(inp); err != nil {
		return util.ErrTrace("BodyParser", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return util.ErrTrace("validateStruct", err)
	}

	id := uuid.New()

	project := entity.Project{
		UUID:         id,
		ClientUUID:   inp.ClientUUID,
		Name:         inp.Name,
		BillableRate: inp.BillableRate,
		Comment:      inp.Comment,
		CreatedAt:    util.TimeNowUtc(),
	}

	_, err := h.repo.Create(&project)
	if err != nil {
		return util.ErrTrace("Create", err)
	}

	p, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.Status(fiber.StatusCreated).JSON(p)
}

func (h *ProjectHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	p, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.JSON(p)
}

type projectUpdateInput struct {
	ClientUUID   uuid.UUID           `json:"clientUUID" validate:"required"`
	Name         string              `json:"name" validate:"required"`
	BillableRate uint64              `json:"billableRate" validate:"min=0"`
	Comment      nullable.NullString `json:"comment"`
}

func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	inp := new(projectUpdateInput)
	if err := c.BodyParser(inp); err != nil {
		return util.ErrTrace("BodyParser", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return util.ErrTrace("validateStruct", err)
	}

	p, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	p.ClientUUID = inp.ClientUUID
	p.Name = inp.Name
	p.BillableRate = inp.BillableRate
	p.Comment = inp.Comment
	p.UpdatedAt = nullable.NewNullTime(util.TimeNowUtc(), true)

	_, err = h.repo.Update(&p)
	if err != nil {
		return util.ErrTrace("Update", err)
	}

	return c.JSON(p)
}

func (h *ProjectHandler) Archive(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	p, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	p.ArchivedAt = nullable.NewNullTime(util.TimeNowUtc(), true)

	_, err = h.repo.Update(&p)
	if err != nil {
		return util.ErrTrace("Update", err)
	}

	return c.JSON(p)
}

func (h *ProjectHandler) Unarchive(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	p, err := h.repo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	p.ArchivedAt = nullable.NewNullTime(util.TimeNowUtc(), false)

	_, err = h.repo.Update(&p)
	if err != nil {
		return util.ErrTrace("Update", err)
	}

	return c.JSON(p)
}
