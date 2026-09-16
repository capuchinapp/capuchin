package sqlite

import (
	"context"
	"fmt"

	"capuchin/internal/domain"
)

// AuditLogRepository SQLite реализация репозитория лога действий.
type AuditLogRepository struct{}

// NewAuditLogRepository возвращает новый AuditLogRepository.
func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{}
}

// Create создаёт лог действия.
func (*AuditLogRepository) Create(ctx context.Context, db SQLExecutor, auditLog domain.AuditLog) error {
	q := `-- audit_log::Create
		INSERT INTO "audit_log"
		("action_at", "action_type", "user_id", "table_name", "record_id", "old_value", "new_value")
		VALUES
		(:action_at, :action_type, :user_id, :table_name, :record_id, :old_value, :new_value)`

	if _, err := db.NamedExecContext(ctx, q, auditLogModelFromDomain(auditLog)); err != nil {
		return fmt.Errorf("insert audit log: %v", err)
	}

	return nil
}
