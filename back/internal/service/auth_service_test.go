package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	"capuchin/internal/mailer"
	"capuchin/internal/service/mocks"
	"capuchin/internal/sqlite"
	sqlitemocks "capuchin/internal/sqlite/mocks"
)

func TestAuthService_Register(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID       = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testUserEmail    = "test@capuchin.ru"
		testActivateCode = "123456"
	)

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *AuthService
		args    args
		wantErr error
	}{
		{
			name: "Should register user",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				mailService := mocks.NewMailServiceMock(t)
				mailService.EXPECT().
					Send(ctx, []string{testUserEmail}, "Успешная регистрация", mailer.TemplateRegisterSuccess, map[string]any{
						"userID":       testUserID,
						"activateCode": testActivateCode,
						"timeString":   activateCodeExpiresText,
					}).
					Return(nil)

				userRepo := mocks.NewUserRepositoryMock(t)
				userRepo.EXPECT().
					FindByEmail(ctx, testUserEmail).
					Return(domain.User{}, domain.ErrNotFound)
				userRepo.EXPECT().
					Create(ctx, mock.Anything, domain.User{
						ID:           testUserID,
						Email:        testUserEmail,
						ActivateCode: func() *string { s := testActivateCode; return &s }(),
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				txManager := mocks.NewTransactionManagerMock(t)
				txManager.
					On("ExecuteInTransaction", ctx,
						mock.AnythingOfType("func(context.Context, sqlite.SQLExecutor) error"),
					).
					Run(func(args mock.Arguments) {
						tx := sqlitemocks.NewSQLExecutorMock(t)

						fn1, ok := args.Get(1).(func(context.Context, sqlite.SQLExecutor) error)
						assert.True(t, ok)

						_ = fn1(ctx, tx)
					}).
					Return(nil)

				return &AuthService{
					mailService:    mailService,
					userRepo:       userRepo,
					txManager:      txManager,
					randLetterFunc: func(l int) string { return testActivateCode },
					idFunc:         func() string { return testUserID },
					nowFunc:        func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				email: testUserEmail,
			},
			wantErr: nil,
		},
		{
			name: "Should return ErrAlreadyExists if user already exists for FindByEmail",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				userRepo := mocks.NewUserRepositoryMock(t)
				userRepo.EXPECT().
					FindByEmail(ctx, testUserEmail).
					Return(domain.User{}, nil)

				return &AuthService{
					userRepo: userRepo,
				}
			},
			args: args{
				email: testUserEmail,
			},
			wantErr: domain.ErrAlreadyExists,
		},
		{
			name: "Should return ErrAlreadyExists if user already exists for Create",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				userRepo := mocks.NewUserRepositoryMock(t)
				userRepo.EXPECT().
					FindByEmail(ctx, testUserEmail).
					Return(domain.User{}, domain.ErrNotFound)
				userRepo.EXPECT().
					Create(ctx, mock.Anything, domain.User{
						ID:           testUserID,
						Email:        testUserEmail,
						ActivateCode: func() *string { s := testActivateCode; return &s }(),
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(domain.ErrAlreadyExists)

				txManager := mocks.NewTransactionManagerMock(t)
				txManager.
					On("ExecuteInTransaction", ctx,
						mock.AnythingOfType("func(context.Context, sqlite.SQLExecutor) error"),
					).
					Run(func(args mock.Arguments) {
						tx := sqlitemocks.NewSQLExecutorMock(t)

						fn1, ok := args.Get(1).(func(context.Context, sqlite.SQLExecutor) error)
						assert.True(t, ok)

						_ = fn1(ctx, tx)
					}).
					Return(domain.ErrAlreadyExists)

				return &AuthService{
					userRepo:       userRepo,
					txManager:      txManager,
					randLetterFunc: func(l int) string { return testActivateCode },
					idFunc:         func() string { return testUserID },
					nowFunc:        func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				email: testUserEmail,
			},
			wantErr: domain.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			err := srv.Register(ctx, "127.0.0.1", tt.args.email)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestAuthService_Activate(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID       = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testActivateCode = "123456"
		testCookieID     = "01JFCHR6H7CE0EER7FSTR668V4"
	)

	type args struct {
		userID string
		code   string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *AuthService
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Should return cookieID",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				userRepo := mocks.NewUserRepositoryMock(t)
				userRepo.EXPECT().
					DeleteNotActivated(ctx, time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC).Add(-activateCodeExpires)).
					Return(nil)
				userRepo.EXPECT().
					CancelActivateCode(ctx, domain.User{
						ID:           testUserID,
						ActivateCode: func() *string { s := testActivateCode; return &s }(),
					}).
					Return(nil)
				userRepo.EXPECT().
					FindByID(ctx, testUserID).
					Return(domain.User{
						ID:           testUserID,
						ActivateCode: func() *string { s := testActivateCode; return &s }(),
					}, nil)

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					Create(ctx, domain.Session{
						ID:        testCookieID,
						CookieID:  testCookieID,
						UserID:    testUserID,
						CheckedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				return &AuthService{
					userRepo:    userRepo,
					sessionRepo: sessionRepo,
					logger:      zap.NewNop(),
					nowFunc:     func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
					idFunc:      func() string { return testCookieID },
				}
			},
			args: args{
				userID: testUserID,
				code:   testActivateCode,
			},
			want:    testCookieID,
			wantErr: nil,
		},
		{
			name: "Should return error",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				userRepo := mocks.NewUserRepositoryMock(t)
				userRepo.EXPECT().
					DeleteNotActivated(ctx, time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC).Add(-activateCodeExpires)).
					Return(nil)
				userRepo.EXPECT().
					FindByID(ctx, testUserID).
					Return(domain.User{
						ID:           testUserID,
						ActivateCode: func() *string { s := testActivateCode; return &s }(),
					}, nil)

				return &AuthService{
					userRepo: userRepo,
					logger:   zap.NewNop(),
					nowFunc:  func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID: testUserID,
				code:   "000000",
			},
			want:    "",
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Activate(ctx, tt.args.userID, tt.args.code)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID       = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testUserEmail    = "test@capuchin.ru"
		testActivateCode = "123456"
	)

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *AuthService
		args    args
		wantErr error
	}{
		{
			name: "Should success login",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				mailService := mocks.NewMailServiceMock(t)
				mailService.EXPECT().
					Send(ctx, []string{testUserEmail}, "Код авторизации", mailer.TemplateAuthCode, map[string]any{
						"code":       testActivateCode,
						"timeString": authCodeExpiresText,
					}).
					Return(nil)

				userRepo := mocks.NewUserRepositoryMock(t)
				userRepo.EXPECT().
					FindByEmail(ctx, testUserEmail).
					Return(domain.User{
						ID:           testUserID,
						Email:        testUserEmail,
						ActivateCode: nil,
					}, nil)

				authCodeRepo := mocks.NewAuthCodeRepositoryMock(t)
				authCodeRepo.EXPECT().
					Create(ctx, domain.AuthCode{
						ID:        "01JFCM58B0GDCX7REETVRNCERY",
						UserID:    testUserID,
						UserEmail: testUserEmail,
						Code:      testActivateCode,
						CreatedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				return &AuthService{
					mailService:   mailService,
					userRepo:      userRepo,
					authCodeRepo:  authCodeRepo,
					logger:        zap.NewNop(),
					randDigitFunc: func(l int) string { return testActivateCode },
					nowFunc:       func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
					idFunc:        func() string { return "01JFCM58B0GDCX7REETVRNCERY" },
				}
			},
			args: args{
				email: testUserEmail,
			},
			wantErr: nil,
		},
		{
			name: "Should fail login - email not found",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				userRepo := mocks.NewUserRepositoryMock(t)
				userRepo.EXPECT().
					FindByEmail(ctx, testUserEmail).
					Return(domain.User{}, domain.ErrNotFound)

				return &AuthService{
					userRepo: userRepo,
				}
			},
			args: args{
				email: testUserEmail,
			},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "Should fail login - activate code not nil",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				userRepo := mocks.NewUserRepositoryMock(t)
				userRepo.EXPECT().
					FindByEmail(ctx, testUserEmail).
					Return(domain.User{
						ID:           testUserID,
						Email:        testUserEmail,
						ActivateCode: func() *string { s := testActivateCode; return &s }(),
					}, nil)

				return &AuthService{
					userRepo: userRepo,
					logger:   zap.NewNop(),
				}
			},
			args: args{
				email: testUserEmail,
			},
			wantErr: domain.ErrUserNotActivated,
		},
		{
			name: "Should fail login - auth code already exists",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				userRepo := mocks.NewUserRepositoryMock(t)
				userRepo.EXPECT().
					FindByEmail(ctx, testUserEmail).
					Return(domain.User{
						ID:           testUserID,
						Email:        testUserEmail,
						ActivateCode: nil,
					}, nil)

				authCodeRepo := mocks.NewAuthCodeRepositoryMock(t)
				authCodeRepo.EXPECT().
					Create(ctx, domain.AuthCode{
						ID:        "01JFCM58B0GDCX7REETVRNCERY",
						UserID:    testUserID,
						UserEmail: testUserEmail,
						Code:      testActivateCode,
						CreatedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(domain.ErrAlreadyExists)

				return &AuthService{
					userRepo:      userRepo,
					authCodeRepo:  authCodeRepo,
					logger:        zap.NewNop(),
					randDigitFunc: func(l int) string { return testActivateCode },
					nowFunc:       func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
					idFunc:        func() string { return "01JFCM58B0GDCX7REETVRNCERY" },
				}
			},
			args: args{
				email: testUserEmail,
			},
			wantErr: domain.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			err := srv.Login(ctx, tt.args.email)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestAuthService_ApplyCode(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID       = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testUserEmail    = "test@capuchin.ru"
		testActivateCode = "123456"
		testCookieID     = "01JFCHR6H7CE0EER7FSTR668V4"
	)

	type args struct {
		email string
		code  string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *AuthService
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Should success apply code",
			srvFunc: func(t *testing.T) *AuthService {
				t.Helper()

				authCodeRepo := mocks.NewAuthCodeRepositoryMock(t)
				authCodeRepo.EXPECT().
					Find(ctx, testUserEmail, testActivateCode).
					Return(domain.AuthCode{
						ID:     "01JFCKH3P1D0ZW4KW0KH63Q0G7",
						UserID: testUserID,
					}, nil)
				authCodeRepo.EXPECT().
					Delete(ctx, "01JFCKH3P1D0ZW4KW0KH63Q0G7").
					Return(nil)
				authCodeRepo.EXPECT().
					CancelOld(ctx, time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC).Add(-authCodeExpires)).
					Return(nil)

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					Create(ctx, domain.Session{
						ID:        testCookieID,
						CookieID:  testCookieID,
						UserID:    testUserID,
						CheckedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				return &AuthService{
					authCodeRepo: authCodeRepo,
					sessionRepo:  sessionRepo,
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
					idFunc:       func() string { return testCookieID },
				}
			},
			args: args{
				email: testUserEmail,
				code:  testActivateCode,
			},
			want:    testCookieID,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.ApplyCode(ctx, tt.args.email, tt.args.code)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
