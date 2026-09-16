package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"capuchin/internal/domain"
	"capuchin/internal/mailer"
	"capuchin/internal/sqlite"
)

const (
	authCodeLength      = 6
	authCodeExpires     = 1 * time.Hour
	authCodeExpiresText = "1 час"

	activateCodeLength      = 40 // См. тип поля в файле migrations/00002_alter__users__activate_code.sql
	activateCodeExpires     = 7 * 24 * time.Hour
	activateCodeExpiresText = "7 дней"
)

// AuthService реализует использование аутентификации.
type AuthService struct {
	mailService MailService

	authCodeRepo AuthCodeRepository
	sessionRepo  SessionRepository
	userRepo     UserRepository

	txManager TransactionManager

	logger *zap.Logger

	idFunc         func() string
	nowFunc        func() time.Time
	randDigitFunc  func(l int) string
	randLetterFunc func(l int) string
}

// NewAuthService возвращает новый AuthService.
func NewAuthService(
	mailService MailService,
	authCodeRepo AuthCodeRepository,
	sessionRepo SessionRepository,
	userRepo UserRepository,
	txManager TransactionManager,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		mailService: mailService,

		authCodeRepo: authCodeRepo,
		sessionRepo:  sessionRepo,
		userRepo:     userRepo,

		txManager: txManager,

		logger: logger.Named("auth_service"),

		idFunc:         domain.NewID,
		nowFunc:        time.Now,
		randDigitFunc:  domain.NewDigitCode,
		randLetterFunc: domain.NewString,
	}
}

// Register создаёт пользователя.
func (s *AuthService) Register(ctx context.Context, _ string, email string) error {
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("find user by email: %v", err)
	}
	if err == nil {
		return domain.ErrAlreadyExists
	}

	activateCode := s.randLetterFunc(activateCodeLength)

	user := domain.User{
		ID:           s.idFunc(),
		Email:        email,
		ActivateCode: &activateCode,
		CreatedAt:    s.nowFunc(),
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context, tx sqlite.SQLExecutor) error {
		return s.userRepo.Create(ctx, tx, user)
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("create user: %v", err)
	}

	err = s.mailService.Send(ctx, []string{user.Email}, "Успешная регистрация", mailer.TemplateRegisterSuccess, map[string]any{
		"userID":       user.ID,
		"activateCode": activateCode,
		"timeString":   activateCodeExpiresText,
	})
	if err != nil {
		return fmt.Errorf("send successful registration: %v", err)
	}

	return nil
}

// Activate активирует пользователя.
func (s *AuthService) Activate(ctx context.Context, userID string, code string) (sessionID string, err error) {
	outdatedTime := s.nowFunc().Add(-activateCodeExpires)
	if err := s.userRepo.DeleteNotActivated(ctx, outdatedTime); err != nil {
		return "", fmt.Errorf("cancel old: %v", err)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to find user", zap.Error(err))

		return "", domain.ErrNotFound
	}

	if user.ActivateCode == nil || *user.ActivateCode != code {
		s.logger.Warn("failed to activate user", zap.Error(err))

		return "", domain.ErrNotFound
	}

	cookieID := s.idFunc()

	session := domain.Session{
		ID:        s.idFunc(),
		CookieID:  cookieID,
		UserID:    userID,
		CheckedAt: s.nowFunc(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", fmt.Errorf("create session: %v", err)
	}

	if err := s.userRepo.CancelActivateCode(ctx, user); err != nil {
		s.logger.Error("failed to cancel activate code", zap.Error(err))
	}

	return cookieID, nil
}

// Login отправляет код авторизации на электронную почту пользователя.
func (s *AuthService) Login(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}

		return fmt.Errorf("find user by email: %v", err)
	}

	if user.ActivateCode != nil {
		s.logger.Warn("user not activated", zap.String("user_id", user.ID))

		return domain.ErrUserNotActivated
	}

	code := s.randDigitFunc(authCodeLength)

	authCode := domain.AuthCode{
		ID:        s.idFunc(),
		UserID:    user.ID,
		UserEmail: user.Email,
		Code:      code,
		CreatedAt: s.nowFunc(),
	}

	if err := s.authCodeRepo.Create(ctx, authCode); err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return domain.ErrAlreadyExists
		}

		return fmt.Errorf("create auth code: %v", err)
	}

	err = s.mailService.Send(ctx, []string{user.Email}, "Код авторизации", mailer.TemplateAuthCode, map[string]any{
		"code":       code,
		"timeString": authCodeExpiresText,
	})
	if err != nil {
		return fmt.Errorf("send auth code: %v", err)
	}

	return nil
}

// ApplyCode применяет код авторизации.
func (s *AuthService) ApplyCode(ctx context.Context, email string, code string) (sessionID string, err error) {
	outdatedTime := s.nowFunc().Add(-authCodeExpires)
	if err := s.authCodeRepo.CancelOld(ctx, outdatedTime); err != nil {
		return "", fmt.Errorf("cancel old: %v", err)
	}

	authCode, err := s.authCodeRepo.Find(ctx, email, code)
	if err != nil {
		// Не проверяем domain.ErrNotFound
		return "", fmt.Errorf("find auth code: %v", err)
	}

	cookieID := s.idFunc()

	session := domain.Session{
		ID:        s.idFunc(),
		CookieID:  cookieID,
		UserID:    authCode.UserID,
		CheckedAt: s.nowFunc(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", fmt.Errorf("create session: %v", err)
	}

	if err := s.authCodeRepo.Delete(ctx, authCode.ID); err != nil {
		s.logger.Error("failed to delete auth code", zap.Error(err))
	}

	return cookieID, nil
}
