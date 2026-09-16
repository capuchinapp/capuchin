package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"capuchin/internal/domain"
	"capuchin/internal/sqlite"
)

// TaskService реализует взаимодействие с задачами.
type TaskService struct {
	taskRepo     TaskRepository
	clientRepo   ClientRepository
	projectRepo  ProjectRepository
	timelogRepo  TimelogRepository
	auditLogRepo AuditLogRepository

	txManager TransactionManager

	idFunc  func() string
	nowFunc func() time.Time
}

// NewTaskService возвращает новый TaskService.
func NewTaskService(
	taskRepo TaskRepository,
	clientRepo ClientRepository,
	projectRepo ProjectRepository,
	timelogRepo TimelogRepository,
	auditLogRepo AuditLogRepository,
	txManager TransactionManager,
) *TaskService {
	return &TaskService{
		taskRepo:     taskRepo,
		clientRepo:   clientRepo,
		projectRepo:  projectRepo,
		timelogRepo:  timelogRepo,
		auditLogRepo: auditLogRepo,

		txManager: txManager,

		idFunc:  domain.NewID,
		nowFunc: time.Now,
	}
}

// Create создаёт задачу.
func (s *TaskService) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	project, err := s.projectRepo.FindByID(ctx, task.UserID, task.ProjectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrProjectNotFound
		}

		return domain.Task{}, fmt.Errorf("find project: %v", err)
	}
	if project.ArchivedAt != nil {
		return domain.Task{}, domain.ErrProjectArchived
	}

	client, err := s.clientRepo.FindByID(ctx, task.UserID, project.ClientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrClientNotFound
		}

		return domain.Task{}, fmt.Errorf("find client: %v", err)
	}
	if client.ArchivedAt != nil {
		return domain.Task{}, domain.ErrClientArchived
	}

	task.ID = s.idFunc()
	task.ClientID = client.ID
	task.ClientName = client.Name
	task.ProjectID = project.ID
	task.ProjectName = project.Name
	task.CreatedAt = s.nowFunc()

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.taskRepo.Create(ctx, tx, task)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   task.CreatedAt,
			ActionType: domain.AuditLogActionTypeCreate,
			UserID:     task.UserID,
			TableName:  domain.AuditLogTableNameTasks,
			RecordID:   task.ID,
			OldValue:   nil,
			NewValue:   func() *string { v := domain.StructToJSON(task); return &v }(),
		})
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return domain.Task{}, domain.ErrAlreadyExists
		}

		return domain.Task{}, fmt.Errorf("create task: %v", err)
	}

	return task, nil
}

// Update обновляет задачу.
func (s *TaskService) Update(ctx context.Context, task domain.Task) (domain.Task, error) {
	exist, err := s.taskRepo.FindByID(ctx, task.UserID, task.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrNotFound
		}

		return domain.Task{}, fmt.Errorf("find task: %v", err)
	}

	project, err := s.projectRepo.FindByID(ctx, task.UserID, task.ProjectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrProjectNotFound
		}

		return domain.Task{}, fmt.Errorf("find project: %v", err)
	}
	if project.ArchivedAt != nil {
		return domain.Task{}, domain.ErrProjectArchived
	}

	client, err := s.clientRepo.FindByID(ctx, task.UserID, project.ClientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrClientNotFound
		}

		return domain.Task{}, fmt.Errorf("find client: %v", err)
	}
	if client.ArchivedAt != nil {
		return domain.Task{}, domain.ErrClientArchived
	}

	newTask := domain.Task{
		ID:          exist.ID,
		UserID:      exist.UserID,
		ProjectID:   project.ID,
		ProjectName: project.Name,
		ClientID:    client.ID,
		ClientName:  client.Name,
		Name:        task.Name,
		Comment:     task.Comment,
		CreatedAt:   exist.CreatedAt,
		UpdatedAt:   func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:  exist.ArchivedAt,
		CompletedAt: exist.CompletedAt,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.taskRepo.Update(ctx, tx, newTask)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newTask.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newTask.UserID,
			TableName:  domain.AuditLogTableNameTasks,
			RecordID:   newTask.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newTask); return &v }(),
		})
	})
	if err != nil {
		return domain.Task{}, fmt.Errorf("update task: %v", err)
	}

	return newTask, nil
}

// Archive архивирует задачу.
func (s *TaskService) Archive(ctx context.Context, userID string, taskID string) (domain.Task, error) {
	exist, err := s.taskRepo.FindByID(ctx, userID, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrNotFound
		}

		return domain.Task{}, fmt.Errorf("find task: %v", err)
	}

	if exist.CompletedAt != nil {
		return domain.Task{}, domain.ErrImpossibleAction
	}

	newTask := domain.Task{
		ID:          exist.ID,
		UserID:      exist.UserID,
		ProjectID:   exist.ProjectID,
		ProjectName: exist.ProjectName,
		ClientID:    exist.ClientID,
		ClientName:  exist.ClientName,
		Name:        exist.Name,
		Comment:     exist.Comment,
		CreatedAt:   exist.CreatedAt,
		UpdatedAt:   func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:  func() *time.Time { t := s.nowFunc(); return &t }(),
		CompletedAt: exist.CompletedAt,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.taskRepo.Update(ctx, tx, newTask)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newTask.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newTask.UserID,
			TableName:  domain.AuditLogTableNameTasks,
			RecordID:   newTask.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newTask); return &v }(),
		})
	})
	if err != nil {
		return domain.Task{}, fmt.Errorf("archive task: %v", err)
	}

	return newTask, nil
}

// Unarchive разархивирует задачу.
func (s *TaskService) Unarchive(ctx context.Context, userID string, taskID string) (domain.Task, error) {
	exist, err := s.taskRepo.FindByID(ctx, userID, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrNotFound
		}

		return domain.Task{}, fmt.Errorf("find task: %v", err)
	}

	newTask := domain.Task{
		ID:          exist.ID,
		UserID:      exist.UserID,
		ProjectID:   exist.ProjectID,
		ProjectName: exist.ProjectName,
		ClientID:    exist.ClientID,
		ClientName:  exist.ClientName,
		Name:        exist.Name,
		Comment:     exist.Comment,
		CreatedAt:   exist.CreatedAt,
		UpdatedAt:   func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:  nil,
		CompletedAt: exist.CompletedAt,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.taskRepo.Update(ctx, tx, newTask)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newTask.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newTask.UserID,
			TableName:  domain.AuditLogTableNameTasks,
			RecordID:   newTask.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newTask); return &v }(),
		})
	})
	if err != nil {
		return domain.Task{}, fmt.Errorf("unarchive task: %v", err)
	}

	return newTask, nil
}

// Complete выполняет задачу.
func (s *TaskService) Complete(ctx context.Context, userID string, taskID string) (domain.Task, error) {
	exist, err := s.taskRepo.FindByID(ctx, userID, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrNotFound
		}

		return domain.Task{}, fmt.Errorf("find task: %v", err)
	}

	if exist.ArchivedAt != nil {
		return domain.Task{}, domain.ErrImpossibleAction
	}

	newTask := domain.Task{
		ID:          exist.ID,
		UserID:      exist.UserID,
		ProjectID:   exist.ProjectID,
		ProjectName: exist.ProjectName,
		ClientID:    exist.ClientID,
		ClientName:  exist.ClientName,
		Name:        exist.Name,
		Comment:     exist.Comment,
		CreatedAt:   exist.CreatedAt,
		UpdatedAt:   func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:  exist.ArchivedAt,
		CompletedAt: func() *time.Time { t := s.nowFunc(); return &t }(),
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.taskRepo.Update(ctx, tx, newTask)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newTask.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newTask.UserID,
			TableName:  domain.AuditLogTableNameTasks,
			RecordID:   newTask.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newTask); return &v }(),
		})
	})
	if err != nil {
		return domain.Task{}, fmt.Errorf("complete task: %v", err)
	}

	return newTask, nil
}

// Incomplete отменяет выполнение задачи.
func (s *TaskService) Incomplete(ctx context.Context, userID string, taskID string) (domain.Task, error) {
	exist, err := s.taskRepo.FindByID(ctx, userID, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrNotFound
		}

		return domain.Task{}, fmt.Errorf("find task: %v", err)
	}

	newTask := domain.Task{
		ID:          exist.ID,
		UserID:      exist.UserID,
		ProjectID:   exist.ProjectID,
		ProjectName: exist.ProjectName,
		ClientID:    exist.ClientID,
		ClientName:  exist.ClientName,
		Name:        exist.Name,
		Comment:     exist.Comment,
		CreatedAt:   exist.CreatedAt,
		UpdatedAt:   func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:  exist.ArchivedAt,
		CompletedAt: nil,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.taskRepo.Update(ctx, tx, newTask)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newTask.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newTask.UserID,
			TableName:  domain.AuditLogTableNameTasks,
			RecordID:   newTask.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newTask); return &v }(),
		})
	})
	if err != nil {
		return domain.Task{}, fmt.Errorf("incomplete task: %v", err)
	}

	return newTask, nil
}

// Report возвращает отчёт по задаче.
func (s *TaskService) Report(ctx context.Context, userID string, taskID string) (domain.TaskReport, error) {
	report, err := s.timelogRepo.TaskReport(ctx, userID, taskID)
	if err != nil {
		return domain.TaskReport{}, fmt.Errorf("get task report time: %v", err)
	}

	return report, nil
}
