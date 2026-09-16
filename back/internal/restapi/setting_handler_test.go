package restapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	"capuchin/internal/restapi/mocks"
	"capuchin/internal/validator"
)

func TestSettingHandler_Index(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *SettingHandler
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *SettingHandler {
				t.Helper()

				settingRepo := mocks.NewSettingRepositoryMock(t)
				settingRepo.EXPECT().
					FindAll(mock.Anything, testUserID).
					Return([]domain.Setting{
						{
							UserID: testUserID,
							Key:    "dateFormat",
							Value:  "yyyy-MM-dd",
						},
						{
							UserID: testUserID,
							Key:    "workingDays",
							Value:  "1,2,3,4,5",
						},
					}, nil)

				return &SettingHandler{
					settingRepo: settingRepo,
				}
			},
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"key": "dateFormat",
					"value": "yyyy-MM-dd"
				},
				{
					"key": "workingDays",
					"value": "1,2,3,4,5"
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

func TestSettingHandler_Update(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *SettingHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *SettingHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				settingRepo := mocks.NewSettingRepositoryMock(t)
				settingRepo.EXPECT().
					InsertOrUpdate(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, setting domain.Setting) error {
						assert.Equal(t, testUserID, setting.UserID)

						switch setting.Key {
						case "dateFormat":
							assert.Equal(t, "DD.MM.YYYY", setting.Value)
						case "workingDays":
							assert.Equal(t, "1,2,3,4,5", setting.Value)
						default:
							t.Errorf("unexpected key: %s", setting.Key)
						}

						return nil
					})
				settingRepo.EXPECT().
					FindAll(mock.Anything, testUserID).
					Return([]domain.Setting{
						{
							UserID: testUserID,
							Key:    "dateFormat",
							Value:  "DD.MM.YYYY",
						},
						{
							UserID: testUserID,
							Key:    "workingDays",
							Value:  "1,2,3,4,5",
						},
					}, nil)

				return &SettingHandler{
					validate:    validator.New(),
					sanitizer:   sanitizer,
					settingRepo: settingRepo,
					logger:      zap.NewNop(),
				}
			},
			input: `{
				"dateFormat": "DD.MM.YYYY",
				"workingDays": "1,2,3,4,5"
			}`,
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"key": "dateFormat",
					"value": "DD.MM.YYYY"
				},
				{
					"key": "workingDays",
					"value": "1,2,3,4,5"
				}
			]`,
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *SettingHandler {
				t.Helper()

				return &SettingHandler{
					logger: zap.NewNop(),
				}
			},
			input:    "",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - validate error",
			srvFunc: func(t *testing.T) *SettingHandler {
				t.Helper()

				return &SettingHandler{
					validate: validator.New(),
				}
			},
			input:    `{"dateFormat": "DD.MM.YY"}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "dateFormat"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Put("/", tt.srvFunc(t).Update)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}
