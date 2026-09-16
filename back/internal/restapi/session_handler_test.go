package restapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"

	"capuchin/internal/domain"
	"capuchin/internal/restapi/mocks"
)

func TestSessionHandler_Index(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *SessionHandler
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *SessionHandler {
				t.Helper()

				sessionService := mocks.NewSessionServiceMock(t)
				sessionService.EXPECT().
					List(mock.Anything, testUserID, time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC).AddDate(0, 0, -1*1)).
					Return([]domain.Session{
						{
							ID:       "01JFCN3VVQRQBGJ255PH9YMGY8",
							CookieID: "01JFYQE72N34P9VZD84AEXC0PC",
						},
						{
							ID:       "01JG356VM30XQ1PY5J251YDYTK",
							CookieID: "01JG3570Y6PKW9GAPFBCG9KA2F",
						},
					}, nil)

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					GetExpiresDays().
					Return(1)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return("01JFYQE72N34P9VZD84AEXC0PC")

				return &SessionHandler{
					sessionService: sessionService,
					cookieManager:  cookieManager,
					nowFunc:        func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFCN3VVQRQBGJ255PH9YMGY8",
					"checkedAt": "0001-01-01T00:00:00Z",
					"isCurrent": true
				},
				{
					"id": "01JG356VM30XQ1PY5J251YDYTK",
					"checkedAt": "0001-01-01T00:00:00Z",
					"isCurrent": false
				}
			]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Get("/", tt.srvFunc(t).Index)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestSessionHandler_Delete(t *testing.T) {
	const (
		testUserID    = "01JFYQE72N34P9VZD84AEXC0PC"
		testSessionID = "01JG356VM30XQ1PY5J251YDYTK"
		testCookieID  = "01JG3570Y6PKW9GAPFBCG9KA2F"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *SessionHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *SessionHandler {
				t.Helper()

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					FindByID(mock.Anything, testUserID, testSessionID).
					Return(domain.Session{
						ID:       testSessionID,
						CookieID: testCookieID,
					}, nil)
				sessionRepo.EXPECT().
					DeleteByCID(mock.Anything, testCookieID).
					Return(nil)

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return("01JG3M23GNVP0N4B1TD4S2AZCN")

				return &SessionHandler{
					sessionRepo:   sessionRepo,
					cookieManager: cookieManager,
				}
			},
			id:       testSessionID,
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *SessionHandler {
				t.Helper()

				return &SessionHandler{}
			},
			id:       "123",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid id",
		},
		{
			name: "Should return error - session not found",
			srvFunc: func(t *testing.T) *SessionHandler {
				t.Helper()

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					FindByID(mock.Anything, testUserID, testSessionID).
					Return(domain.Session{}, domain.ErrNotFound)

				return &SessionHandler{
					sessionRepo: sessionRepo,
				}
			},
			id:       testSessionID,
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return error",
			srvFunc: func(t *testing.T) *SessionHandler {
				t.Helper()

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					FindByID(mock.Anything, testUserID, testSessionID).
					Return(domain.Session{
						ID:       testSessionID,
						CookieID: testCookieID,
					}, nil)

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return(testCookieID)

				return &SessionHandler{
					sessionRepo:   sessionRepo,
					cookieManager: cookieManager,
				}
			},
			id:       testSessionID,
			wantCode: http.StatusForbidden,
			wantBody: "Forbidden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Delete("/:id", tt.srvFunc(t).Delete)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/"+tt.id, nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}
