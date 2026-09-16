package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"capuchin/internal/domain"
	"capuchin/internal/sqlite"
)

// ProjectService реализует взаимодействие с проектами.
type ProjectService struct {
	projectRepo  ProjectRepository
	clientRepo   ClientRepository
	auditLogRepo AuditLogRepository

	txManager TransactionManager

	idFunc  func() string
	nowFunc func() time.Time
}

// NewProjectService возвращает новый ProjectService.
func NewProjectService(
	projectRepo ProjectRepository,
	clientRepo ClientRepository,
	auditLogRepo AuditLogRepository,
	txManager TransactionManager,
) *ProjectService {
	return &ProjectService{
		projectRepo:  projectRepo,
		clientRepo:   clientRepo,
		auditLogRepo: auditLogRepo,

		txManager: txManager,

		idFunc:  domain.NewID,
		nowFunc: time.Now,
	}
}

// Create создаёт проект.
func (s *ProjectService) Create(ctx context.Context, project domain.Project) (domain.Project, error) {
	client, err := s.clientRepo.FindByID(ctx, project.UserID, project.ClientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Project{}, domain.ErrClientNotFound
		}

		return domain.Project{}, fmt.Errorf("find client: %v", err)
	}
	if client.ArchivedAt != nil {
		return domain.Project{}, domain.ErrClientArchived
	}

	project.ID = s.idFunc()
	project.ClientName = client.Name
	project.CreatedAt = s.nowFunc()

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.projectRepo.Create(ctx, tx, project)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   project.CreatedAt,
			ActionType: domain.AuditLogActionTypeCreate,
			UserID:     project.UserID,
			TableName:  domain.AuditLogTableNameProjects,
			RecordID:   project.ID,
			OldValue:   nil,
			NewValue:   func() *string { v := domain.StructToJSON(project); return &v }(),
		})
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return domain.Project{}, domain.ErrAlreadyExists
		}

		return domain.Project{}, fmt.Errorf("create project: %v", err)
	}

	return project, nil
}

// Update обновляет проект.
func (s *ProjectService) Update(ctx context.Context, project domain.Project) (domain.Project, error) {
	exist, err := s.projectRepo.FindByID(ctx, project.UserID, project.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Project{}, domain.ErrNotFound
		}

		return domain.Project{}, fmt.Errorf("find project: %v", err)
	}

	client, err := s.clientRepo.FindByID(ctx, project.UserID, project.ClientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Project{}, domain.ErrClientNotFound
		}

		return domain.Project{}, fmt.Errorf("find client: %v", err)
	}
	if client.ArchivedAt != nil {
		return domain.Project{}, domain.ErrClientArchived
	}

	newProject := domain.Project{
		ID:           exist.ID,
		UserID:       exist.UserID,
		ClientID:     client.ID,
		ClientName:   client.Name,
		Name:         project.Name,
		BillableRate: project.BillableRate,
		Comment:      project.Comment,
		CreatedAt:    exist.CreatedAt,
		UpdatedAt:    func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:   exist.ArchivedAt,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.projectRepo.Update(ctx, tx, newProject)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newProject.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newProject.UserID,
			TableName:  domain.AuditLogTableNameProjects,
			RecordID:   newProject.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newProject); return &v }(),
		})
	})
	if err != nil {
		return domain.Project{}, fmt.Errorf("update project: %v", err)
	}

	return newProject, nil
}

// Archive архивирует проект.
func (s *ProjectService) Archive(ctx context.Context, userID string, projectID string) (domain.Project, error) {
	exist, err := s.projectRepo.FindByID(ctx, userID, projectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Project{}, domain.ErrNotFound
		}

		return domain.Project{}, fmt.Errorf("find project: %v", err)
	}

	newProject := domain.Project{
		ID:           exist.ID,
		UserID:       exist.UserID,
		ClientID:     exist.ClientID,
		ClientName:   exist.ClientName,
		Name:         exist.Name,
		BillableRate: exist.BillableRate,
		Comment:      exist.Comment,
		CreatedAt:    exist.CreatedAt,
		UpdatedAt:    func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:   func() *time.Time { t := s.nowFunc(); return &t }(),
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.projectRepo.Update(ctx, tx, newProject)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newProject.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newProject.UserID,
			TableName:  domain.AuditLogTableNameProjects,
			RecordID:   newProject.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newProject); return &v }(),
		})
	})
	if err != nil {
		return domain.Project{}, fmt.Errorf("archive project: %v", err)
	}

	return newProject, nil
}

// Unarchive разархивирует проект.
func (s *ProjectService) Unarchive(ctx context.Context, userID string, projectID string) (domain.Project, error) {
	exist, err := s.projectRepo.FindByID(ctx, userID, projectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Project{}, domain.ErrNotFound
		}

		return domain.Project{}, fmt.Errorf("find project: %v", err)
	}

	newProject := domain.Project{
		ID:           exist.ID,
		UserID:       exist.UserID,
		ClientID:     exist.ClientID,
		ClientName:   exist.ClientName,
		Name:         exist.Name,
		BillableRate: exist.BillableRate,
		Comment:      exist.Comment,
		CreatedAt:    exist.CreatedAt,
		UpdatedAt:    func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:   nil,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.projectRepo.Update(ctx, tx, newProject)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newProject.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newProject.UserID,
			TableName:  domain.AuditLogTableNameProjects,
			RecordID:   newProject.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newProject); return &v }(),
		})
	})
	if err != nil {
		return domain.Project{}, fmt.Errorf("unarchive project: %v", err)
	}

	return newProject, nil
}
