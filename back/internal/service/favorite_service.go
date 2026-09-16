package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"capuchin/internal/domain"
	"capuchin/internal/sqlite"
)

// FavoriteService реализует взаимодействие с избранным.
type FavoriteService struct {
	favoriteRepo FavoriteRepository
	clientRepo   ClientRepository
	projectRepo  ProjectRepository
	taskRepo     TaskRepository
	auditLogRepo AuditLogRepository

	txManager TransactionManager

	idFunc  func() string
	nowFunc func() time.Time
}

// NewFavoriteService возвращает новый FavoriteService.
func NewFavoriteService(
	favoriteRepo FavoriteRepository,
	clientRepo ClientRepository,
	projectRepo ProjectRepository,
	taskRepo TaskRepository,
	auditLogRepo AuditLogRepository,
	txManager TransactionManager,
) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favoriteRepo,
		clientRepo:   clientRepo,
		projectRepo:  projectRepo,
		taskRepo:     taskRepo,
		auditLogRepo: auditLogRepo,

		txManager: txManager,

		idFunc:  domain.NewID,
		nowFunc: time.Now,
	}
}

// Create создаёт избранное.
func (s *FavoriteService) Create(ctx context.Context, favorite domain.Favorite) (domain.Favorite, error) { //nolint:gocognit // TODO: уменьшить сложность
	project, err := s.projectRepo.FindByID(ctx, favorite.UserID, favorite.ProjectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Favorite{}, domain.ErrProjectNotFound
		}

		return domain.Favorite{}, fmt.Errorf("find project: %v", err)
	}
	if project.ArchivedAt != nil {
		return domain.Favorite{}, domain.ErrProjectArchived
	}

	client, err := s.clientRepo.FindByID(ctx, favorite.UserID, project.ClientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Favorite{}, domain.ErrClientNotFound
		}

		return domain.Favorite{}, fmt.Errorf("find client: %v", err)
	}
	if client.ArchivedAt != nil {
		return domain.Favorite{}, domain.ErrClientArchived
	}

	favorite.ID = s.idFunc()
	favorite.ClientID = client.ID
	favorite.ClientName = client.Name
	favorite.ProjectID = project.ID
	favorite.ProjectName = project.Name

	if favorite.TaskID != nil {
		task, err := s.taskRepo.FindByID(ctx, favorite.UserID, *favorite.TaskID)
		if err != nil && errors.Is(err, domain.ErrNotFound) {
			return domain.Favorite{}, domain.ErrTaskNotFound
		}
		if err != nil {
			return domain.Favorite{}, fmt.Errorf("find task: %v", err)
		}
		if task.ArchivedAt != nil {
			return domain.Favorite{}, domain.ErrTaskArchived
		}
		if task.CompletedAt != nil {
			return domain.Favorite{}, domain.ErrTaskCompleted
		}

		favorite.TaskName = &task.Name
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.favoriteRepo.Create(ctx, tx, favorite)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   s.nowFunc(),
			ActionType: domain.AuditLogActionTypeCreate,
			UserID:     favorite.UserID,
			TableName:  domain.AuditLogTableNameFavorites,
			RecordID:   favorite.ID,
			OldValue:   nil,
			NewValue:   func() *string { v := domain.StructToJSON(favorite); return &v }(),
		})
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return domain.Favorite{}, domain.ErrAlreadyExists
		}

		return domain.Favorite{}, fmt.Errorf("create favorite: %v", err)
	}

	return favorite, nil
}

// Update обновляет избранное.
func (s *FavoriteService) Update( //nolint:gocognit,gocyclo,cyclop // TODO: уменьшить сложность
	ctx context.Context,
	favorite domain.Favorite,
) (domain.Favorite, error) {
	exist, err := s.favoriteRepo.FindByID(ctx, favorite.UserID, favorite.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Favorite{}, domain.ErrNotFound
		}

		return domain.Favorite{}, fmt.Errorf("find favorite: %v", err)
	}

	project, err := s.projectRepo.FindByID(ctx, favorite.UserID, favorite.ProjectID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Favorite{}, domain.ErrProjectNotFound
		}

		return domain.Favorite{}, fmt.Errorf("find project: %v", err)
	}
	if project.ArchivedAt != nil {
		return domain.Favorite{}, domain.ErrProjectArchived
	}

	client, err := s.clientRepo.FindByID(ctx, favorite.UserID, project.ClientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Favorite{}, domain.ErrClientNotFound
		}

		return domain.Favorite{}, fmt.Errorf("find client: %v", err)
	}
	if client.ArchivedAt != nil {
		return domain.Favorite{}, domain.ErrClientArchived
	}

	newFavorite := domain.Favorite{
		ID:           exist.ID,
		UserID:       exist.UserID,
		Name:         favorite.Name,
		ProjectID:    project.ID,
		ProjectName:  project.Name,
		ClientID:     client.ID,
		ClientName:   client.Name,
		BillableRate: favorite.BillableRate,
		Comment:      favorite.Comment,
	}

	if favorite.TaskID != nil {
		task, err := s.taskRepo.FindByID(ctx, favorite.UserID, *favorite.TaskID)
		if err != nil && errors.Is(err, domain.ErrNotFound) {
			return domain.Favorite{}, domain.ErrTaskNotFound
		}
		if err != nil {
			return domain.Favorite{}, fmt.Errorf("find task: %v", err)
		}
		if task.ArchivedAt != nil {
			return domain.Favorite{}, domain.ErrTaskArchived
		}
		if task.CompletedAt != nil {
			return domain.Favorite{}, domain.ErrTaskCompleted
		}

		newFavorite.TaskID = favorite.TaskID
		newFavorite.TaskName = &task.Name
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.favoriteRepo.Update(ctx, tx, newFavorite)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   s.nowFunc(),
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newFavorite.UserID,
			TableName:  domain.AuditLogTableNameFavorites,
			RecordID:   newFavorite.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newFavorite); return &v }(),
		})
	})
	if err != nil {
		return domain.Favorite{}, fmt.Errorf("update favorite: %v", err)
	}

	return newFavorite, nil
}

// Delete удаляет избранное.
func (s *FavoriteService) Delete(ctx context.Context, userID string, id string) error { //nolint:dupl // Это лишь кажется дубликатом
	exist, err := s.favoriteRepo.FindByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("find favorite: %v", err)
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.favoriteRepo.Delete(ctx, tx, exist.UserID, exist.ID)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   s.nowFunc(),
			ActionType: domain.AuditLogActionTypeDelete,
			UserID:     exist.UserID,
			TableName:  domain.AuditLogTableNameFavorites,
			RecordID:   exist.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   nil,
		})
	})
	if err != nil {
		return fmt.Errorf("delete favorite: %v", err)
	}

	return nil
}
