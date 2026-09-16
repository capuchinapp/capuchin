package service

import (
	"context"
	"fmt"
	"time"

	"capuchin/internal/domain"
)

// SessionService реализует взаимодействие с сессиями.
type SessionService struct {
	sessionRepo SessionRepository

	nowFunc func() time.Time
}

// NewSessionService возвращает новый SessionService.
func NewSessionService(sessionRepo SessionRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,

		nowFunc: time.Now,
	}
}

// List возвращает сессии пользователя.
func (s *SessionService) List(ctx context.Context, userID string, outdatedTime time.Time) ([]domain.Session, error) {
	if err := s.sessionRepo.CancelOld(ctx, outdatedTime); err != nil {
		return nil, fmt.Errorf("cancel old: %v", err)
	}

	ss, err := s.sessionRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find sessions: %v", err)
	}

	return ss, nil
}
