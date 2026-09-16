package sqlite

import (
	"database/sql"

	"capuchin/internal/domain"
)

// auditLogModel представляет лог действия.
type auditLogModel struct {
	ID         int64              `db:"id"`
	ActionAt   int64              `db:"action_at"`
	ActionType auditLogActionType `db:"action_type"`
	UserID     string             `db:"user_id"`
	TableName  auditLogTableName  `db:"table_name"`
	RecordID   string             `db:"record_id"`
	OldValue   sql.NullString     `db:"old_value"`
	NewValue   sql.NullString     `db:"new_value"`
}

type auditLogActionType string

const (
	auditLogActionTypeUnspecified auditLogActionType = "UNSPECIFIED"
	auditLogActionTypeCreate      auditLogActionType = "CREATE"
	auditLogActionTypeUpdate      auditLogActionType = "UPDATE"
	auditLogActionTypeDelete      auditLogActionType = "DELETE"
)

type auditLogTableName string

const (
	auditLogTableNameUnspecified auditLogTableName = "UNSPECIFIED"
	auditLogTableNameClients     auditLogTableName = "CLIENTS"
	auditLogTableNameProjects    auditLogTableName = "PROJECTS"
	auditLogTableNameTasks       auditLogTableName = "TASKS"
	auditLogTableNameFavorites   auditLogTableName = "FAVORITES"
	auditLogTableNameTimelogs    auditLogTableName = "TIMELOGS"
)

// auditLogModelFromDomain преобразует лог действия из домена в sqlite.
func auditLogModelFromDomain(src domain.AuditLog) auditLogModel {
	return auditLogModel{
		ID:         src.ID,
		ActionAt:   timeToUnix(src.ActionAt),
		ActionType: auditLogActionTypeFromDomain(src.ActionType),
		UserID:     src.UserID,
		TableName:  auditLogTableNameFromDomain(src.TableName),
		RecordID:   src.RecordID,
		OldValue:   nullableStringFromDomain(src.OldValue),
		NewValue:   nullableStringFromDomain(src.NewValue),
	}
}

func auditLogActionTypeFromDomain(src domain.AuditLogActionType) auditLogActionType {
	switch src {
	case domain.AuditLogActionTypeCreate:
		return auditLogActionTypeCreate
	case domain.AuditLogActionTypeUpdate:
		return auditLogActionTypeUpdate
	case domain.AuditLogActionTypeDelete:
		return auditLogActionTypeDelete
	default:
		return auditLogActionTypeUnspecified
	}
}

func auditLogTableNameFromDomain(src domain.AuditLogTableName) auditLogTableName {
	switch src {
	case domain.AuditLogTableNameClients:
		return auditLogTableNameClients
	case domain.AuditLogTableNameProjects:
		return auditLogTableNameProjects
	case domain.AuditLogTableNameTasks:
		return auditLogTableNameTasks
	case domain.AuditLogTableNameFavorites:
		return auditLogTableNameFavorites
	case domain.AuditLogTableNameTimelogs:
		return auditLogTableNameTimelogs
	default:
		return auditLogTableNameUnspecified
	}
}
