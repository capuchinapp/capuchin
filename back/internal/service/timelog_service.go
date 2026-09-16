package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"capuchin/internal/domain"
	"capuchin/internal/sqlite"
)

const (
	hoursPerDay = 24
)

// TimelogService реализует взаимодействие с журналом.
type TimelogService struct {
	timelogRepo  TimelogRepository
	clientRepo   ClientRepository
	projectRepo  ProjectRepository
	taskRepo     TaskRepository
	auditLogRepo AuditLogRepository

	txManager TransactionManager

	idFunc  func() string
	nowFunc func() time.Time
}

// NewTimelogService возвращает новый TimelogService.
func NewTimelogService(
	timelogRepo TimelogRepository,
	clientRepo ClientRepository,
	projectRepo ProjectRepository,
	taskRepo TaskRepository,
	auditLogRepo AuditLogRepository,
	txManager TransactionManager,
) *TimelogService {
	return &TimelogService{
		timelogRepo:  timelogRepo,
		clientRepo:   clientRepo,
		projectRepo:  projectRepo,
		taskRepo:     taskRepo,
		auditLogRepo: auditLogRepo,

		txManager: txManager,

		nowFunc: time.Now,
		idFunc:  domain.NewID,
	}
}

// Create создаёт запись в журнале.
func (s *TimelogService) Create(ctx context.Context, timelog domain.Timelog) (domain.Timelog, error) {
	project, err := s.projectRepo.FindByID(ctx, timelog.UserID, timelog.ProjectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Timelog{}, domain.ErrProjectNotFound
		}

		return domain.Timelog{}, fmt.Errorf("find project: %v", err)
	}
	if project.ArchivedAt != nil {
		return domain.Timelog{}, domain.ErrProjectArchived
	}

	client, err := s.clientRepo.FindByID(ctx, timelog.UserID, project.ClientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Timelog{}, domain.ErrClientNotFound
		}

		return domain.Timelog{}, fmt.Errorf("find client: %v", err)
	}
	if client.ArchivedAt != nil {
		return domain.Timelog{}, domain.ErrClientArchived
	}

	if timelog.TaskID != nil {
		task, err := s.taskRepo.FindByID(ctx, timelog.UserID, *timelog.TaskID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return domain.Timelog{}, domain.ErrTaskNotFound
			}

			return domain.Timelog{}, fmt.Errorf("find task: %v", err)
		}
		if task.ArchivedAt != nil {
			return domain.Timelog{}, domain.ErrTaskArchived
		}
		if task.CompletedAt != nil {
			return domain.Timelog{}, domain.ErrTaskCompleted
		}

		timelog.TaskName = &task.Name
		timelog.TaskCompletedAt = task.CompletedAt
	}

	if err := s.stopRunningRecord(ctx, timelog.UserID, timelog.Date, timelog.TimeStart); err != nil {
		if domain.IsFailedPreconditionError(err) {
			return domain.Timelog{}, err
		}

		return domain.Timelog{}, fmt.Errorf("stop running record: %v", err)
	}

	timelog.ClientID = client.ID
	timelog.ClientName = client.Name
	timelog.ProjectID = project.ID
	timelog.ProjectName = project.Name

	tl, err := s.insertRecord(ctx, timelog)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return domain.Timelog{}, domain.ErrAlreadyExists
		}
		if domain.IsFailedPreconditionError(err) {
			return domain.Timelog{}, err
		}

		return domain.Timelog{}, fmt.Errorf("insert record: %v", err)
	}

	return tl, nil
}

// Update обновляет запись в журнале.
func (s *TimelogService) Update(ctx context.Context, timelog domain.Timelog) (domain.Timelog, error) {
	if timelog.TimeEnd == nil {
		return domain.Timelog{}, domain.NewFailedPreconditionError("cannot edit running timelog") //nolint:wrapcheck // Всё в порядке
	}

	exist, err := s.timelogRepo.FindByID(ctx, timelog.UserID, timelog.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Timelog{}, domain.ErrNotFound
		}

		return domain.Timelog{}, fmt.Errorf("find timelog: %v", err)
	}

	project, err := s.projectRepo.FindByID(ctx, timelog.UserID, timelog.ProjectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Timelog{}, domain.ErrProjectNotFound
		}

		return domain.Timelog{}, fmt.Errorf("find project: %v", err)
	}
	if project.ArchivedAt != nil {
		return domain.Timelog{}, domain.ErrProjectArchived
	}

	client, err := s.clientRepo.FindByID(ctx, timelog.UserID, project.ClientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Timelog{}, domain.ErrClientNotFound
		}

		return domain.Timelog{}, fmt.Errorf("find client: %v", err)
	}
	if client.ArchivedAt != nil {
		return domain.Timelog{}, domain.ErrClientArchived
	}

	dur, err := getDuration(timelog.TimeStart, *timelog.TimeEnd)
	if err != nil {
		if domain.IsFailedPreconditionError(err) {
			return domain.Timelog{}, err
		}

		return domain.Timelog{}, fmt.Errorf("get duration: %v", err)
	}

	newTimelog := domain.Timelog{
		ID:              exist.ID,
		UserID:          exist.UserID,
		ProjectID:       project.ID,
		ProjectName:     project.Name,
		ClientID:        client.ID,
		ClientName:      client.Name,
		Date:            timelog.Date,
		TimeStart:       timelog.TimeStart,
		TimeEnd:         timelog.TimeEnd,
		DurationSeconds: uint32(dur.Seconds()),
		BillableRate:    timelog.BillableRate,
		BillableAmount:  int64(dur.Hours() * float64(timelog.BillableRate)),
		Comment:         timelog.Comment,
		CreatedAt:       exist.CreatedAt,
		UpdatedAt:       func() *time.Time { t := s.nowFunc(); return &t }(),
	}

	if timelog.TaskID != nil {
		task, err := s.taskRepo.FindByID(ctx, timelog.UserID, *timelog.TaskID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return domain.Timelog{}, domain.ErrTaskNotFound
			}

			return domain.Timelog{}, fmt.Errorf("find task: %v", err)
		}
		if task.ArchivedAt != nil {
			return domain.Timelog{}, domain.ErrTaskArchived
		}
		if task.CompletedAt != nil {
			return domain.Timelog{}, domain.ErrTaskCompleted
		}

		newTimelog.TaskID = timelog.TaskID
		newTimelog.TaskName = &task.Name
		newTimelog.TaskCompletedAt = task.CompletedAt
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.timelogRepo.Update(ctx, tx, newTimelog)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newTimelog.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newTimelog.UserID,
			TableName:  domain.AuditLogTableNameTimelogs,
			RecordID:   newTimelog.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newTimelog); return &v }(),
		})
	})
	if err != nil {
		return domain.Timelog{}, fmt.Errorf("update timelog: %v", err)
	}

	return newTimelog, nil
}

// Delete удаляет запись времени.
func (s *TimelogService) Delete(ctx context.Context, userID string, id string) error { //nolint:dupl // Это лишь кажется дубликатом
	exist, err := s.timelogRepo.FindByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("find timelog: %v", err)
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.timelogRepo.Delete(ctx, tx, exist.UserID, exist.ID)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   s.nowFunc(),
			ActionType: domain.AuditLogActionTypeDelete,
			UserID:     exist.UserID,
			TableName:  domain.AuditLogTableNameTimelogs,
			RecordID:   exist.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   nil,
		})
	})
	if err != nil {
		return fmt.Errorf("delete timelog: %v", err)
	}

	return nil
}

// Stop останавливает запись в журнале.
func (s *TimelogService) Stop(ctx context.Context, userID string, timelogID string, date string, timeEnd string) (domain.Timelog, error) {
	tl, err := s.timelogRepo.FindByID(ctx, userID, timelogID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Timelog{}, domain.ErrNotFound
		}

		return domain.Timelog{}, fmt.Errorf("find timelog: %v", err)
	}

	if tl.TimeEnd != nil {
		// Не будем останавливать уже остановленную запись
		return tl, nil
	}

	tl, err = s.stopRecord(ctx, tl, date, timeEnd)
	if err != nil {
		if domain.IsFailedPreconditionError(err) {
			return domain.Timelog{}, err
		}

		return domain.Timelog{}, fmt.Errorf("stop record: %v", err)
	}

	return tl, nil
}

func (s *TimelogService) FindLastN(ctx context.Context, userID string, n int) ([]domain.Timelog, error) {
	// Возьмем больше записей чем нужно для уникальности
	const multiplier = 10
	n *= multiplier

	timelogsOverdraft, err := s.timelogRepo.FindLastN(ctx, userID, n)
	if err != nil {
		return []domain.Timelog{}, fmt.Errorf("find last timelogs: %v", err)
	}

	// Найдем уникальные записи по полям project_id, task_id, comment
	uniqueKeys := make(map[string]struct{})
	timelogs := make([]domain.Timelog, 0, n)
	for _, tl := range timelogsOverdraft {
		var taskID string
		if tl.TaskID != nil {
			if tl.TaskCompletedAt != nil {
				continue
			}

			taskID = *tl.TaskID
		}

		key := fmt.Sprintf("%s_%s_%s", tl.ProjectID, taskID, tl.Comment)
		if _, ok := uniqueKeys[key]; ok {
			continue
		}

		uniqueKeys[key] = struct{}{}

		timelogs = append(timelogs, tl)
		if len(timelogs) >= n {
			break
		}
	}

	return timelogs, nil
}

// insertRecord добавляет новую запись в журнал.
func (s *TimelogService) insertRecord(ctx context.Context, tl domain.Timelog) (domain.Timelog, error) {
	if tl.UserID == "" {
		return domain.Timelog{}, ErrUserIDRequired
	}

	tl.ID = s.idFunc()
	tl.CreatedAt = s.nowFunc()

	if tl.TimeEnd != nil {
		dur, err := getDuration(tl.TimeStart, *tl.TimeEnd)
		if err != nil {
			if domain.IsFailedPreconditionError(err) {
				return domain.Timelog{}, err
			}

			return domain.Timelog{}, fmt.Errorf("get duration: %v", err)
		}

		tl.DurationSeconds = uint32(dur.Seconds())
		tl.BillableAmount = int64(dur.Hours() * float64(tl.BillableRate))
	}

	err := s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.timelogRepo.Create(ctx, tx, tl)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   tl.CreatedAt,
			ActionType: domain.AuditLogActionTypeCreate,
			UserID:     tl.UserID,
			TableName:  domain.AuditLogTableNameTimelogs,
			RecordID:   tl.ID,
			OldValue:   nil,
			NewValue:   func() *string { v := domain.StructToJSON(tl); return &v }(),
		})
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return domain.Timelog{}, domain.ErrAlreadyExists
		}

		return domain.Timelog{}, fmt.Errorf("create timelog: %v", err)
	}

	return tl, nil
}

// stopRecord останавливает запись в журнале.
func (s *TimelogService) stopRecord(ctx context.Context, tl domain.Timelog, date string, timeEnd string) (domain.Timelog, error) {
	const (
		timeStartDay = "00:00:00"
		timeEndDay   = "23:59:59"
	)

	var newTimeEnd *string
	if date == tl.Date {
		newTimeEnd = &timeEnd
	} else {
		// Закрываем запись до конца дня
		newTimeEnd = func() *string {
			s := timeEndDay
			return &s
		}()
	}

	dur, err := getDuration(tl.TimeStart, *newTimeEnd)
	if err != nil {
		if domain.IsFailedPreconditionError(err) {
			return domain.Timelog{}, err
		}

		return domain.Timelog{}, fmt.Errorf("get duration: %v", err)
	}

	newTimelog := domain.Timelog{
		ID:              tl.ID,
		UserID:          tl.UserID,
		ProjectID:       tl.ProjectID,
		ProjectName:     tl.ProjectName,
		ClientID:        tl.ClientID,
		ClientName:      tl.ClientName,
		Date:            tl.Date,
		TimeStart:       tl.TimeStart,
		TimeEnd:         newTimeEnd,
		DurationSeconds: uint32(dur.Seconds()),
		BillableRate:    tl.BillableRate,
		BillableAmount:  int64(dur.Hours() * float64(tl.BillableRate)),
		Comment:         tl.Comment,
		CreatedAt:       tl.CreatedAt,
		UpdatedAt:       func() *time.Time { t := s.nowFunc(); return &t }(),
		TaskID:          tl.TaskID,
		TaskName:        tl.TaskName,
		TaskCompletedAt: tl.TaskCompletedAt,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.timelogRepo.Update(ctx, tx, newTimelog)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newTimelog.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newTimelog.UserID,
			TableName:  domain.AuditLogTableNameTimelogs,
			RecordID:   newTimelog.ID,
			OldValue:   func() *string { v := domain.StructToJSON(tl); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newTimelog); return &v }(),
		})
	})
	if err != nil {
		return domain.Timelog{}, fmt.Errorf("update timelog: %v", err)
	}

	if date == tl.Date {
		return newTimelog, nil
	}

	startDate, err := time.Parse("2006-01-02", tl.Date)
	if err != nil {
		return domain.Timelog{}, fmt.Errorf("start date parse: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return domain.Timelog{}, fmt.Errorf("end date parse: %v", err)
	}

	// Найти количество дней между двумя датами
	diff := int(endDate.Sub(startDate).Hours() / hoursPerDay)
	if diff < 0 {
		return domain.Timelog{}, domain.ErrCurrentDateLessThanStartDate
	}

	curDate := startDate

	// Проходим все дни, кроме последнего
	if diff > 1 {
		diff--
		for i := 0; i < diff; i++ {
			curDate = curDate.AddDate(0, 0, 1)

			_, err := s.insertRecord(ctx, domain.Timelog{
				UserID:       tl.UserID,
				ProjectID:    tl.ProjectID,
				TaskID:       tl.TaskID,
				Date:         curDate.Format("2006-01-02"),
				TimeStart:    timeStartDay,
				TimeEnd:      func() *string { s := timeEndDay; return &s }(),
				BillableRate: tl.BillableRate,
				Comment:      tl.Comment,
			})
			if err != nil {
				if domain.IsFailedPreconditionError(err) {
					return domain.Timelog{}, err
				}

				return domain.Timelog{}, fmt.Errorf("insert record: %v", err)
			}
		}
	}

	// Добавление последнего дня
	curDate = curDate.AddDate(0, 0, 1)

	tl, err = s.insertRecord(ctx, domain.Timelog{
		UserID:          tl.UserID,
		ClientID:        tl.ClientID,
		ClientName:      tl.ClientName,
		ProjectID:       tl.ProjectID,
		ProjectName:     tl.ProjectName,
		TaskID:          tl.TaskID,
		TaskName:        tl.TaskName,
		TaskCompletedAt: tl.TaskCompletedAt,
		Date:            curDate.Format("2006-01-02"),
		TimeStart:       timeStartDay,
		TimeEnd:         &timeEnd,
		BillableRate:    tl.BillableRate,
		Comment:         tl.Comment,
	})
	if err != nil {
		if domain.IsFailedPreconditionError(err) {
			return domain.Timelog{}, err
		}

		return domain.Timelog{}, fmt.Errorf("insert record: %v", err)
	}

	return tl, nil
}

// stopRunningRecord останавливает запущенную запись в журнале.
func (s *TimelogService) stopRunningRecord(ctx context.Context, userID string, date string, timeEnd string) error {
	tl, err := s.timelogRepo.FindRunning(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("find running timelog: %v", err)
	}

	_, err = s.stopRecord(ctx, tl, date, timeEnd)
	if err != nil {
		if domain.IsFailedPreconditionError(err) {
			return err
		}

		return fmt.Errorf("stop record: %v", err)
	}

	return nil
}

// getDuration возвращает длительность задачи.
func getDuration(timeStart string, timeEnd string) (time.Duration, error) {
	from, err := time.Parse("15:04:05", timeStart)
	if err != nil {
		return 0, fmt.Errorf("from parse: %v", err)
	}

	to, err := time.Parse("15:04:05", timeEnd)
	if err != nil {
		return 0, fmt.Errorf("to parse: %v", err)
	}

	diff := to.Sub(from)
	if diff < 0 {
		return 0, domain.ErrTimeEndLessThanTimeStart
	}

	return diff, nil
}
