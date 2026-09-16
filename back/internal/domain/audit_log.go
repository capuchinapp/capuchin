package domain

import (
	"time"
)

// AuditLog представляет лог действия.
type AuditLog struct {
	ID         int64
	ActionAt   time.Time
	ActionType AuditLogActionType
	UserID     string
	TableName  AuditLogTableName
	RecordID   string
	OldValue   *string
	NewValue   *string
}

type AuditLogActionType string

const (
	AuditLogActionTypeCreate AuditLogActionType = "CREATE"
	AuditLogActionTypeUpdate AuditLogActionType = "UPDATE"
	AuditLogActionTypeDelete AuditLogActionType = "DELETE"
)

type AuditLogTableName string

const (
	AuditLogTableNameClients   AuditLogTableName = "CLIENTS"
	AuditLogTableNameProjects  AuditLogTableName = "PROJECTS"
	AuditLogTableNameTasks     AuditLogTableName = "TASKS"
	AuditLogTableNameFavorites AuditLogTableName = "FAVORITES"
	AuditLogTableNameTimelogs  AuditLogTableName = "TIMELOGS"
)
