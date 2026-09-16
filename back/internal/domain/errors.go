package domain

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrAlreadyExists    = errors.New("record already exists")
	ErrPermissionDenied = errors.New("permission denied")
	ErrUserNotActivated = errors.New("user not activated")
	ErrImpossibleAction = errors.New("impossible action")

	ErrClientNotFound = NewFailedPreconditionError("client not found")
	ErrClientArchived = NewFailedPreconditionError("client is archived")

	ErrProjectNotFound = NewFailedPreconditionError("project not found")
	ErrProjectArchived = NewFailedPreconditionError("project is archived")

	ErrTaskNotFound  = NewFailedPreconditionError("task not found")
	ErrTaskArchived  = NewFailedPreconditionError("task is archived")
	ErrTaskCompleted = NewFailedPreconditionError("task is completed")

	ErrTimeEndLessThanTimeStart     = NewFailedPreconditionError("timeEnd cannot be less than the timeStart")
	ErrCurrentDateLessThanStartDate = NewFailedPreconditionError("current date cannot be less than the start date")
)

// FailedPreconditionError ошибка не выполнения предварительного условия.
type FailedPreconditionError struct {
	Message string
}

// Error реализует ошибку.
func (e *FailedPreconditionError) Error() string {
	return e.Message
}

// NewFailedPreconditionError создает новую ошибку.
func NewFailedPreconditionError(message string) error {
	return &FailedPreconditionError{Message: message}
}

// IsFailedPreconditionError проверяет, имеет ли ошибка тип FailedPreconditionError.
func IsFailedPreconditionError(err error) bool {
	var e *FailedPreconditionError
	return errors.As(err, &e)
}
