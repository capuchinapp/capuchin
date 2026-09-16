package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"capuchin/internal/domain"
)

// TransactionManager SQLite реализация менеджера транзакций.
type TransactionManager struct {
	db *sqlx.DB
}

// NewTransactionManager возвращает новый TransactionManager.
func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

// ExecuteInTransaction выполняет две операции в одной транзакции.
// Первая функция — основная операция (например, вставка клиента).
// Вторая функция — операция логирования (например, вставка в audit_log).
// Можно передать несколько операций.
func (m *TransactionManager) ExecuteInTransaction(
	ctx context.Context,
	opFuncs ...func(ctx context.Context, tx SQLExecutor) error,
) error {
	if len(opFuncs) == 0 {
		return nil
	}

	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	for i, opFunc := range opFuncs {
		if err := opFunc(ctx, tx); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("failed to rollback transaction: %v (operation %d failed: %v)", rbErr, i+1, err)
			}

			return fmt.Errorf("operation %d failed: %w", i+1, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// ExampleTransactionManager пример использования TransactionManager.
func ExampleTransactionManager() {
	// Подготовка переменных
	db := &sqlx.DB{}
	ctx := context.Background()
	task := domain.Task{}
	s := struct {
		taskRepo     *TaskRepository
		auditLogRepo *AuditLogRepository
	}{
		taskRepo:     NewTaskRepository(db),
		auditLogRepo: NewAuditLogRepository(),
	}

	// Использование
	txManager := NewTransactionManager(db)

	err := txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx SQLExecutor) error {
		return s.taskRepo.Create(ctx, tx, task)
	}, func(ctx context.Context, tx SQLExecutor) error {
		return s.auditLogRepo.Create(ctx, tx, domain.AuditLog{})
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
