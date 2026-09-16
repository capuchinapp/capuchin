package restapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"

	"capuchin/internal/domain"
	"capuchin/internal/restapi/mocks"
)

func TestIndexHandler_Index(t *testing.T) {
	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *IndexHandler
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK - isAuth true",
			srvFunc: func(t *testing.T) *IndexHandler {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return("01JFYQE72N34P9VZD84AEXC0PC")

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					FindByCID(mock.Anything, "01JFYQE72N34P9VZD84AEXC0PC").
					Return(domain.Session{
						CookieID: "01JFYQE72N34P9VZD84AEXC0PC",
						UserID:   "01JJ6P2298EPV806CN19ZJTJ1A",
					}, nil)

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindRunning(mock.Anything, "01JJ6P2298EPV806CN19ZJTJ1A").
					Return(domain.Timelog{
						Date:      "2025-01-22",
						TimeStart: "14:39:15",
					}, nil)

				return &IndexHandler{
					appVersion:    "v0.0.0",
					cookieManager: cookieManager,
					sessionRepo:   sessionRepo,
					timelogRepo:   timelogRepo,
				}
			},
			wantCode: http.StatusOK,
			wantBody: `{
				"appVersion": "v0.0.0",
				"isAuth": true,
				"name": "Capuchin API",
				"runningTimelogDatetime": "2025-01-22 14:39:15"
			}`,
		},
		{
			name: "Should return OK - isAuth true and runningTimelog not found",
			srvFunc: func(t *testing.T) *IndexHandler {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return("01JFYQE72N34P9VZD84AEXC0PC")

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					FindByCID(mock.Anything, "01JFYQE72N34P9VZD84AEXC0PC").
					Return(domain.Session{
						CookieID: "01JFYQE72N34P9VZD84AEXC0PC",
						UserID:   "01JJ6P2298EPV806CN19ZJTJ1A",
					}, nil)

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindRunning(mock.Anything, "01JJ6P2298EPV806CN19ZJTJ1A").
					Return(domain.Timelog{}, domain.ErrNotFound)

				return &IndexHandler{
					appVersion:    "v0.0.0",
					cookieManager: cookieManager,
					sessionRepo:   sessionRepo,
					timelogRepo:   timelogRepo,
				}
			},
			wantCode: http.StatusOK,
			wantBody: `{
				"appVersion": "v0.0.0",
				"isAuth": true,
				"name": "Capuchin API",
				"runningTimelogDatetime": null
			}`,
		},
		{
			name: "Should return OK - isAuth false",
			srvFunc: func(t *testing.T) *IndexHandler {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return("01JFYQE72N34P9VZD84AEXC0PC")

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					FindByCID(mock.Anything, "01JFYQE72N34P9VZD84AEXC0PC").
					Return(domain.Session{}, nil)

				return &IndexHandler{
					appVersion:    "v0.0.0",
					cookieManager: cookieManager,
					sessionRepo:   sessionRepo,
				}
			},
			wantCode: http.StatusOK,
			wantBody: `{
				"appVersion": "v0.0.0",
				"isAuth": false,
				"name": "Capuchin API",
				"runningTimelogDatetime": null
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/", tt.srvFunc(t).Index)
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}
