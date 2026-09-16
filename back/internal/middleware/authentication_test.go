package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	"capuchin/internal/middleware/mocks"
)

func TestNewAuthentication(t *testing.T) {
	var (
		cid = "01JFYQE72N34P9VZD84AEXC0PC"
		sid = "01JHHXTMRFMJ4DQJVTVCHZSZX3"
		uid = "01JHHXTVAMTYVE7Q3H6VRP7206"
	)

	tests := []struct {
		name     string
		cfgFunc  func(t *testing.T) AuthenticationConfig
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK - not found in cache",
			cfgFunc: func(t *testing.T) AuthenticationConfig {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return(cid)

				sessionRepo := mocks.NewSessionRepoMock(t)
				sessionRepo.EXPECT().
					FindByCID(mock.Anything, cid).
					Return(domain.Session{
						ID:        sid,
						CookieID:  cid,
						UserID:    uid,
						CheckedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				sessionCache := mocks.NewSessionCacheMock(t)
				sessionCache.EXPECT().
					Get(cid).
					Return(domain.Session{}, false)
				sessionCache.EXPECT().
					Add(cid, domain.Session{
						ID:        sid,
						CookieID:  cid,
						UserID:    uid,
						CheckedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(false)

				sessionUpdater := mocks.NewSessionUpdaterMock(t)
				sessionUpdater.EXPECT().
					UpdateTime(sid).
					Return()

				return AuthenticationConfig{
					CookieManager:  cookieManager,
					SessionRepo:    sessionRepo,
					SessionCache:   sessionCache,
					SessionUpdater: sessionUpdater,
					Logger:         zap.NewNop(),
					NowFunc:        func() time.Time { return time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC) },
					Paths:          []string{"/"},
				}
			},
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return OK - finded in cache",
			cfgFunc: func(t *testing.T) AuthenticationConfig {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return(cid)

				sessionCache := mocks.NewSessionCacheMock(t)
				sessionCache.EXPECT().
					Get(cid).
					Return(domain.Session{
						ID:        sid,
						CookieID:  cid,
						UserID:    uid,
						CheckedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, true)

				sessionUpdater := mocks.NewSessionUpdaterMock(t)
				sessionUpdater.EXPECT().
					UpdateTime(sid).
					Return()

				return AuthenticationConfig{
					CookieManager:  cookieManager,
					SessionCache:   sessionCache,
					SessionUpdater: sessionUpdater,
					Logger:         zap.NewNop(),
					NowFunc:        func() time.Time { return time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC) },
					Paths:          []string{"/"},
				}
			},
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return OK - path not in paths",
			cfgFunc: func(t *testing.T) AuthenticationConfig {
				t.Helper()

				return AuthenticationConfig{
					Paths: []string{"/test"},
				}
			},
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return Unauthorized - no cookie",
			cfgFunc: func(t *testing.T) AuthenticationConfig {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return("")

				return AuthenticationConfig{
					CookieManager: cookieManager,
					Paths:         []string{"/"},
				}
			},
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized",
		},
		{
			name: "Should return Unauthorized - session not found",
			cfgFunc: func(t *testing.T) AuthenticationConfig {
				t.Helper()

				cookieManager := mocks.NewCookieManagerMock(t)
				cookieManager.EXPECT().
					Get(mock.Anything).
					Return(cid)

				sessionRepo := mocks.NewSessionRepoMock(t)
				sessionRepo.EXPECT().
					FindByCID(mock.Anything, cid).
					Return(domain.Session{}, domain.ErrNotFound)

				sessionCache := mocks.NewSessionCacheMock(t)
				sessionCache.EXPECT().
					Get(cid).
					Return(domain.Session{}, false)

				return AuthenticationConfig{
					CookieManager: cookieManager,
					SessionRepo:   sessionRepo,
					SessionCache:  sessionCache,
					Paths:         []string{"/"},
				}
			},
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(NewAuthentication(tt.cfgFunc(t)))

			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusOK)
			})

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			t.Logf("StatusCode: %d\n", resp.StatusCode)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			bodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody, string(bodyBytes))
		})
	}
}
