package restapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	"capuchin/internal/restapi/mocks"
	"capuchin/internal/validator"
)

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *AuthHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return Created",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				authService := mocks.NewAuthServiceMock(t)
				authService.EXPECT().
					Register(mock.Anything, "0.0.0.0", "ivan@localhost.tld").
					Return(nil)

				return &AuthHandler{
					validate:    validator.New(),
					authService: authService,
				}
			},
			input: `{
				"email": "ivan@localhost.tld"
			}`,
			wantCode: http.StatusCreated,
			wantBody: "Created",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					logger: zap.NewNop(),
				}
			},
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid email",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"email": "ivan@localhost"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "email"`,
		},
		{
			name: "Should return error - user already exists",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				authService := mocks.NewAuthServiceMock(t)
				authService.EXPECT().
					Register(mock.Anything, "0.0.0.0", "ivan@localhost.tld").
					Return(domain.ErrAlreadyExists)

				return &AuthHandler{
					validate:    validator.New(),
					authService: authService,
				}
			},
			input: `{
				"email": "ivan@localhost.tld"
			}`,
			wantCode: http.StatusConflict,
			wantBody: "user with this email already exists",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Post("/", tt.srvFunc(t).Register)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestAuthHandler_Activate(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *AuthHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return Created",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				authService := mocks.NewAuthServiceMock(t)
				authService.EXPECT().
					Activate(mock.Anything, testUserID, "DN3WCdXYCsyPHcSe4Y5xnMTewfaHfD3cCsUvxvtc").
					Return("01JH4GTPGN6TQ86Y2802GPK1MS", nil)

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Create(mock.Anything, "01JH4GTPGN6TQ86Y2802GPK1MS").
					Return()

				return &AuthHandler{
					validate:      validator.New(),
					authService:   authService,
					cookieManager: cookieManager,
				}
			},
			input: `{
				"userId": "01JFYQE72N34P9VZD84AEXC0PC",
				"code": "DN3WCdXYCsyPHcSe4Y5xnMTewfaHfD3cCsUvxvtc"
			}`,
			wantCode: http.StatusCreated,
			wantBody: "Created",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					logger: zap.NewNop(),
				}
			},
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid user id",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"userId": "123",
				"code": "DN3WCdXYCsyPHcSe4Y5xnMTewfaHfD3cCsUvxvtc"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "userId"`,
		},
		{
			name: "Should return error - invalid code",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"userId": "01JFYQE72N34P9VZD84AEXC0PC",
				"code": "DN3WCdXYCsyPHcSe4Y5xnMTewfaHfD3cCsUvx"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "code"`,
		},
		{
			name: "Should return error - code not found",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				authService := mocks.NewAuthServiceMock(t)
				authService.EXPECT().
					Activate(mock.Anything, testUserID, "DN3WCdXYCsyPHcSe4Y5xnMTewfaHfD3cCsUvxvtc").
					Return("", domain.ErrNotFound)

				return &AuthHandler{
					validate:    validator.New(),
					authService: authService,
				}
			},
			input: `{
				"userId": "01JFYQE72N34P9VZD84AEXC0PC",
				"code": "DN3WCdXYCsyPHcSe4Y5xnMTewfaHfD3cCsUvxvtc"
			}`,
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Post("/", tt.srvFunc(t).Activate)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *AuthHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				authService := mocks.NewAuthServiceMock(t)
				authService.EXPECT().
					Login(mock.Anything, "ivan@localhost.tld").
					Return(nil)

				return &AuthHandler{
					validate:    validator.New(),
					authService: authService,
				}
			},
			input: `{
				"email": "ivan@localhost.tld"
			}`,
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					logger: zap.NewNop(),
				}
			},
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid email",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"email": "ivan@localhost"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "email"`,
		},
		{
			name: "Should return error - email not found",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				authService := mocks.NewAuthServiceMock(t)
				authService.EXPECT().
					Login(mock.Anything, "ivan@localhost.tld").
					Return(domain.ErrNotFound)

				return &AuthHandler{
					validate:    validator.New(),
					authService: authService,
				}
			},
			input: `{
				"email": "ivan@localhost.tld"
			}`,
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
		{
			name: "Should return error - generated duplication code",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				authService := mocks.NewAuthServiceMock(t)
				authService.EXPECT().
					Login(mock.Anything, "ivan@localhost.tld").
					Return(domain.ErrAlreadyExists)

				return &AuthHandler{
					validate:    validator.New(),
					authService: authService,
				}
			},
			input: `{
				"email": "ivan@localhost.tld"
			}`,
			wantCode: http.StatusConflict,
			wantBody: "Conflict",
		},
		{
			name: "Should return error - user not activated",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				authService := mocks.NewAuthServiceMock(t)
				authService.EXPECT().
					Login(mock.Anything, "ivan@localhost.tld").
					Return(domain.ErrUserNotActivated)

				return &AuthHandler{
					validate:    validator.New(),
					authService: authService,
				}
			},
			input: `{
				"email": "ivan@localhost.tld"
			}`,
			wantCode: http.StatusConflict,
			wantBody: "user not activated",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Post("/", tt.srvFunc(t).Login)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestAuthHandler_ApplyCode(t *testing.T) {
	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *AuthHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return Created",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				authService := mocks.NewAuthServiceMock(t)
				authService.EXPECT().
					ApplyCode(mock.Anything, "ivan@localhost.tld", "123456").
					Return("01JH4GTPGN6TQ86Y2802GPK1MS", nil)

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Create(mock.Anything, "01JH4GTPGN6TQ86Y2802GPK1MS").
					Return()

				return &AuthHandler{
					validate:      validator.New(),
					authService:   authService,
					cookieManager: cookieManager,
				}
			},
			input: `{
				"email": "ivan@localhost.tld",
				"code": "123456"
			}`,
			wantCode: http.StatusCreated,
			wantBody: "Created",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					logger: zap.NewNop(),
				}
			},
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid email",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"email": "ivan@localhost",
				"code": "123456"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "email"`,
		},
		{
			name: "Should return error - invalid code",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				return &AuthHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"email": "ivan@localhost.tld",
				"code": "123"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "code"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Post("/", tt.srvFunc(t).ApplyCode)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *AuthHandler
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return("01JH4GTPGN6TQ86Y2802GPK1MS")

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					DeleteByCID(mock.Anything, "01JH4GTPGN6TQ86Y2802GPK1MS").
					Return(nil)

				return &AuthHandler{
					validate:      validator.New(),
					cookieManager: cookieManager,
					sessionRepo:   sessionRepo,
				}
			},
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return OK - no cookie",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return("")

				return &AuthHandler{
					validate:      validator.New(),
					cookieManager: cookieManager,
				}
			},
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return error - invalid session id",
			srvFunc: func(t *testing.T) *AuthHandler {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return("123")

				return &AuthHandler{
					validate:      validator.New(),
					cookieManager: cookieManager,
				}
			},
			wantCode: http.StatusBadRequest,
			wantBody: `invalid session id`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Post("/", tt.srvFunc(t).Logout)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}
