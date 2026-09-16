package sessionupdater

import (
	"context"
	"time"
)

// SessionRepo представляет репозиторий сессий.
type SessionRepo interface {
	UpdateCheckedAt(ctx context.Context, sessionID string, checkedAt time.Time) error
}
