package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"capuchin/internal/domain"
	"capuchin/internal/sqlite"
)

// ClientService реализует взаимодействие с клиентами.
type ClientService struct {
	clientRepo   ClientRepository
	auditLogRepo AuditLogRepository

	txManager TransactionManager

	idFunc  func() string
	nowFunc func() time.Time
}

// NewClientService возвращает новый ClientService.
func NewClientService(
	clientRepo ClientRepository,
	auditLogRepo AuditLogRepository,
	txManager TransactionManager,
) *ClientService {
	return &ClientService{
		clientRepo:   clientRepo,
		auditLogRepo: auditLogRepo,

		txManager: txManager,

		idFunc:  domain.NewID,
		nowFunc: time.Now,
	}
}

// Create создаёт клиента.
func (s *ClientService) Create(ctx context.Context, client domain.Client) (domain.Client, error) {
	if client.UserID == "" {
		return domain.Client{}, ErrUserIDRequired
	}

	client.ID = s.idFunc()
	client.CreatedAt = s.nowFunc()

	err := s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.clientRepo.Create(ctx, tx, client)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   client.CreatedAt,
			ActionType: domain.AuditLogActionTypeCreate,
			UserID:     client.UserID,
			TableName:  domain.AuditLogTableNameClients,
			RecordID:   client.ID,
			OldValue:   nil,
			NewValue:   func() *string { v := domain.StructToJSON(client); return &v }(),
		})
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return domain.Client{}, domain.ErrAlreadyExists
		}

		return domain.Client{}, fmt.Errorf("create client: %v", err)
	}

	return client, nil
}

// Update обновляет клиента.
func (s *ClientService) Update(ctx context.Context, client domain.Client) (domain.Client, error) {
	exist, err := s.clientRepo.FindByID(ctx, client.UserID, client.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Client{}, domain.ErrNotFound
		}

		return domain.Client{}, fmt.Errorf("find client: %v", err)
	}

	newClient := domain.Client{
		ID:           exist.ID,
		UserID:       exist.UserID,
		Name:         client.Name,
		BillableRate: client.BillableRate,
		Comment:      client.Comment,
		CreatedAt:    exist.CreatedAt,
		UpdatedAt:    func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:   exist.ArchivedAt,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.clientRepo.Update(ctx, tx, newClient)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newClient.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newClient.UserID,
			TableName:  domain.AuditLogTableNameClients,
			RecordID:   newClient.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newClient); return &v }(),
		})
	})
	if err != nil {
		return domain.Client{}, fmt.Errorf("update client: %v", err)
	}

	return newClient, nil
}

// Archive архивирует клиента.
func (s *ClientService) Archive(ctx context.Context, userID string, clientID string) (domain.Client, error) {
	exist, err := s.clientRepo.FindByID(ctx, userID, clientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Client{}, domain.ErrNotFound
		}

		return domain.Client{}, fmt.Errorf("find client: %v", err)
	}

	newClient := domain.Client{
		ID:           exist.ID,
		UserID:       exist.UserID,
		Name:         exist.Name,
		BillableRate: exist.BillableRate,
		Comment:      exist.Comment,
		CreatedAt:    exist.CreatedAt,
		UpdatedAt:    func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:   func() *time.Time { t := s.nowFunc(); return &t }(),
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.clientRepo.Update(ctx, tx, newClient)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newClient.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newClient.UserID,
			TableName:  domain.AuditLogTableNameClients,
			RecordID:   newClient.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newClient); return &v }(),
		})
	})
	if err != nil {
		return domain.Client{}, fmt.Errorf("archive client: %v", err)
	}

	return newClient, nil
}

// Unarchive разархивирует клиента.
func (s *ClientService) Unarchive(ctx context.Context, userID string, clientID string) (domain.Client, error) {
	exist, err := s.clientRepo.FindByID(ctx, userID, clientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Client{}, domain.ErrNotFound
		}

		return domain.Client{}, fmt.Errorf("find client: %v", err)
	}

	newClient := domain.Client{
		ID:           exist.ID,
		UserID:       exist.UserID,
		Name:         exist.Name,
		BillableRate: exist.BillableRate,
		Comment:      exist.Comment,
		CreatedAt:    exist.CreatedAt,
		UpdatedAt:    func() *time.Time { t := s.nowFunc(); return &t }(),
		ArchivedAt:   nil,
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.clientRepo.Update(ctx, tx, newClient)
	}, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{
			ActionAt:   *newClient.UpdatedAt,
			ActionType: domain.AuditLogActionTypeUpdate,
			UserID:     newClient.UserID,
			TableName:  domain.AuditLogTableNameClients,
			RecordID:   newClient.ID,
			OldValue:   func() *string { v := domain.StructToJSON(exist); return &v }(),
			NewValue:   func() *string { v := domain.StructToJSON(newClient); return &v }(),
		})
	})
	if err != nil {
		return domain.Client{}, fmt.Errorf("unarchive client: %v", err)
	}

	return newClient, nil
}
