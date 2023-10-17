package http

import (
	"capuchin/app/entity"
	"capuchin/app/repository"
	"capuchin/app/util"
	"capuchin/app/util/nullable"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TimelogHandler struct {
	validate *validator.Validate
	tlRepo   repository.TimelogRepository
	prRepo   repository.ProjectRepository
}

func RegisterTimelogEndpoints(r fiber.Router, validate *validator.Validate, tr repository.TimelogRepository, pr repository.ProjectRepository) {
	h := TimelogHandler{
		validate: validate,
		tlRepo:   tr,
		prRepo:   pr,
	}

	g := r.Group("/timelogs")
	g.Get("/", h.Index)
	g.Post("/", h.Create)
	g.Get("/:uuid", h.Get)
	g.Patch("/:uuid", h.Update)
	g.Patch("/:uuid/stop", h.Stop)
	g.Delete("/:uuid", h.Delete)
}

func (h *TimelogHandler) Index(c *fiber.Ctx) error {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	clUuid := c.Query("client_uuid")
	prUuid := c.Query("project_uuid")

	if dateFrom == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Field date_from is required")
	}

	if dateTo == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Field date_to is required")
	}

	if clUuid != "" {
		id, err := uuid.Parse(clUuid)
		if err != nil {
			return util.ErrTrace("uuid.Parse", err)
		}

		tl, err := h.tlRepo.FindByPeriodAndClientId(dateFrom, dateTo, id)
		if err != nil {
			return util.ErrSql2Fiber(err)
		}

		return c.JSON(tl)
	}

	if prUuid != "" {
		id, err := uuid.Parse(prUuid)
		if err != nil {
			return util.ErrTrace("uuid.Parse", err)
		}

		tl, err := h.tlRepo.FindByPeriodAndProjectId(dateFrom, dateTo, id)
		if err != nil {
			return util.ErrSql2Fiber(err)
		}

		return c.JSON(tl)
	}

	tl, err := h.tlRepo.FindByPeriod(dateFrom, dateTo)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.JSON(tl)
}

type timelogCreateInput struct {
	ProjectUUID  uuid.UUID           `json:"projectUUID" validate:"required"`
	Date         string              `json:"date" validate:"datetime=2006-01-02"`
	TimeStart    string              `json:"timeStart" validate:"datetime=15:04:05"`
	TimeEnd      nullable.NullString `json:"timeEnd"`
	BillableRate uint64              `json:"billableRate" validate:"min=0"`
	Comment      nullable.NullString `json:"comment"`
}

func (h *TimelogHandler) Create(c *fiber.Ctx) error {
	var inp *timelogCreateInput = new(timelogCreateInput)

	if err := c.BodyParser(inp); err != nil {
		return util.ErrTrace("BodyParser", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return util.ErrTrace("validateStruct", err)
	}

	err := h.stopRunningRecord(inp.Date, inp.TimeStart)
	if err != nil {
		return util.ErrTrace("stopRunningRecord", err)
	}

	id, err := h.insertRecord(inp)
	if err != nil {
		return util.ErrTrace("insertRecord", err)
	}

	r, err := h.tlRepo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.Status(fiber.StatusCreated).JSON(r)
}

func (h *TimelogHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	tl, err := h.tlRepo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	return c.JSON(tl)
}

type timelogUpdateInput struct {
	ProjectUUID  uuid.UUID           `json:"projectUUID" validate:"required"`
	Date         string              `json:"date" validate:"datetime=2006-01-02"`
	TimeStart    string              `json:"timeStart" validate:"datetime=15:04:05"`
	TimeEnd      string              `json:"timeEnd" validate:"datetime=15:04:05"`
	BillableRate uint64              `json:"billableRate" validate:"min=0"`
	Comment      nullable.NullString `json:"comment"`
}

func (h *TimelogHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	var inp *timelogUpdateInput = new(timelogUpdateInput)

	err = c.BodyParser(inp)
	if err != nil {
		return util.ErrTrace("BodyParser", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return util.ErrTrace("validateStruct", err)
	}

	tl, err := h.tlRepo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	if inp.TimeEnd != "" {
		dur, err := h.getDuration(inp.TimeStart, inp.TimeEnd)
		if err != nil {
			return util.ErrTrace("inp.getDuration", err)
		}

		tl.TimeEnd = nullable.NewNullString(inp.TimeEnd, true)
		tl.DurationSeconds = uint32(dur.Seconds())
		tl.BillableRate = inp.BillableRate
		tl.BillableAmount = uint64(dur.Hours() * float64(inp.BillableRate))
	} else if tl.TimeEnd.Valid {
		dur, err := h.getDuration(tl.TimeStart, tl.TimeEnd.String)
		if err != nil {
			return util.ErrTrace("tl.getDuration", err)
		}

		tl.DurationSeconds = uint32(dur.Seconds())
		tl.BillableAmount = uint64(dur.Hours() * float64(tl.BillableRate))
	}

	tl.ProjectUUID = inp.ProjectUUID
	tl.Date = inp.Date
	tl.TimeStart = inp.TimeStart
	tl.Comment = inp.Comment
	tl.UpdatedAt = nullable.NewNullTime(util.TimeNowUtc(), true)

	_, err = h.tlRepo.Update(&tl)
	if err != nil {
		return util.ErrTrace("Update", err)
	}

	return c.JSON(tl)
}

func (h *TimelogHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	tl, err := h.tlRepo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	tl.DeletedAt = nullable.NewNullTime(util.TimeNowUtc(), true)

	_, err = h.tlRepo.Delete(&tl)
	if err != nil {
		return util.ErrTrace("Delete", err)
	}

	return c.JSON(tl)
}

type timelogStopInput struct {
	Date    string `json:"date" validate:"datetime=2006-01-02"`
	TimeEnd string `json:"timeEnd" validate:"datetime=15:04:05"`
}

func (h *TimelogHandler) Stop(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid", ""))
	if err != nil {
		return util.ErrTrace("uuid.Parse", err)
	}

	var inp *timelogStopInput = new(timelogStopInput)

	err = c.BodyParser(inp)
	if err != nil {
		return util.ErrTrace("BodyParser", err)
	}

	if err := h.validate.Struct(inp); err != nil {
		return util.ErrTrace("validateStruct", err)
	}

	tl, err := h.tlRepo.FindById(id)
	if err != nil {
		return util.ErrSql2Fiber(err)
	}

	if tl.TimeEnd.Valid {
		return c.JSON(tl)
	}

	r, err := h.stopRecord(inp, &tl)
	if err != nil {
		return util.ErrTrace("stopRecord", err)
	}

	return c.JSON(r)
}

func (h *TimelogHandler) insertRecord(inp *timelogCreateInput) (uuid.UUID, error) {
	id := uuid.New()

	tl := entity.Timelog{
		UUID:         id,
		ProjectUUID:  inp.ProjectUUID,
		Date:         inp.Date,
		TimeStart:    inp.TimeStart,
		TimeEnd:      inp.TimeEnd,
		BillableRate: inp.BillableRate,
		Comment:      inp.Comment,
		CreatedAt:    util.TimeNowUtc(),
	}

	if tl.TimeEnd.Valid {
		dur, err := h.getDuration(tl.TimeStart, tl.TimeEnd.String)
		if err != nil {
			return uuid.Nil, util.ErrTrace("getDuration", err)
		}

		tl.DurationSeconds = uint32(dur.Seconds())
		tl.BillableAmount = uint64(dur.Hours() * float64(tl.BillableRate))
	}

	_, err := h.tlRepo.Create(&tl)
	if err != nil {
		return uuid.Nil, util.ErrTrace("Create", err)
	}

	return id, nil
}

func (h *TimelogHandler) stopRecord(inp *timelogStopInput, tl *entity.Timelog) (*entity.Timelog, error) {
	if inp.Date == tl.Date {
		tl.TimeEnd = nullable.NewNullString(inp.TimeEnd, true)
	} else {
		// Закрываем запись концом дня
		tl.TimeEnd = nullable.NewNullString("23:59:59", true)
	}

	dur, err := h.getDuration(tl.TimeStart, tl.TimeEnd.String)
	if err != nil {
		return nil, util.ErrTrace("getDuration", err)
	}

	tl.DurationSeconds = uint32(dur.Seconds())
	tl.BillableAmount = uint64(dur.Hours() * float64(tl.BillableRate))
	tl.UpdatedAt = nullable.NewNullTime(util.TimeNowUtc(), true)

	_, err = h.tlRepo.Update(tl)
	if err != nil {
		return nil, util.ErrTrace("Update", err)
	}

	if inp.Date == tl.Date {
		return tl, nil
	}

	startDate, err := time.Parse("2006-01-02", tl.Date)
	if err != nil {
		return nil, util.ErrTrace("startDateParse", err)
	}

	endDate, err := time.Parse("2006-01-02", inp.Date)
	if err != nil {
		return nil, util.ErrTrace("endDateParse", err)
	}

	// Находим количество дней между двумя датами
	diff := int(endDate.Sub(startDate).Hours() / 24)
	if diff < 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "The current date cannot be less than the start date")
	}

	curDate := startDate

	// Перебираем все дни кроме последнего
	if diff > 1 {
		diff--
		for i := 0; i < diff; i++ {
			curDate = curDate.AddDate(0, 0, 1)

			_, err := h.insertRecord(&timelogCreateInput{
				ProjectUUID:  tl.ProjectUUID,
				Date:         curDate.Format("2006-01-02"),
				TimeStart:    "00:00:00",
				TimeEnd:      nullable.NewNullString("23:59:59", true),
				BillableRate: tl.BillableRate,
				Comment:      tl.Comment,
			})
			if err != nil {
				return nil, util.ErrTrace("insertRecord", err)
			}
		}
	}

	// Добавляем последний день
	curDate = curDate.AddDate(0, 0, 1)

	id, err := h.insertRecord(&timelogCreateInput{
		ProjectUUID:  tl.ProjectUUID,
		Date:         curDate.Format("2006-01-02"),
		TimeStart:    "00:00:00",
		TimeEnd:      nullable.NewNullString(inp.TimeEnd, true),
		BillableRate: tl.BillableRate,
		Comment:      tl.Comment,
	})
	if err != nil {
		return nil, util.ErrTrace("insertRecord", err)
	}

	r, err := h.tlRepo.FindById(id)
	if err != nil {
		return nil, util.ErrSql2Fiber(err)
	}

	return &r, nil
}

func (h *TimelogHandler) stopRunningRecord(date string, timeStart string) error {
	tl, err := h.tlRepo.FindRunning()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return util.ErrTrace("FindRunning", err)
	}

	inp := timelogStopInput{
		Date:    date,
		TimeEnd: timeStart,
	}

	_, err = h.stopRecord(&inp, &tl)
	if err != nil {
		return util.ErrTrace("stopRecord", err)
	}

	return nil
}

func (h *TimelogHandler) getDuration(timeStart string, timeEnd string) (time.Duration, error) {
	from, err := time.Parse("15:04:05", timeStart)
	if err != nil {
		return 0, util.ErrTrace("fromParse", err)
	}

	to, err := time.Parse("15:04:05", timeEnd)
	if err != nil {
		return 0, util.ErrTrace("toParse", err)
	}

	diff := to.Sub(from)
	if diff < 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "The timeEnd cannot be less than the timeStart")
	}

	return diff, nil
}
