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

func TestFavoriteHandler_Index(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *FavoriteHandler
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindAll(mock.Anything, testUserID).
					Return([]domain.Favorite{
						{
							ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
							Name:         "Favorite 1",
							ClientID:     "01JH749ZAVM5MX8AMCJH7AKBX8",
							ClientName:   "Client 1",
							ProjectID:    "01JH74AEFJRBTWZEF4KTB64M8V",
							ProjectName:  "Project 1",
							BillableRate: 90000,
							Comment:      "Comment 1",
						},
						{
							ID:           "01JH741H5A4Y368PBHBTJFQ723",
							Name:         "Favorite 2",
							ClientID:     "01JH74AAT06A0P3AR2AGCFGRR5",
							ClientName:   "Client 2",
							ProjectID:    "01JH74AKN2WM4FWXWJ4JBBZRPE",
							ProjectName:  "Project 2",
							TaskID:       func() *string { s := "01JH74AQYEE61Z2ANVRA77NPV6"; return &s }(),
							TaskName:     func() *string { s := "Task One"; return &s }(),
							BillableRate: 50000,
						},
					}, nil)

				return &FavoriteHandler{
					validate:     validator.New(),
					favoriteRepo: favoriteRepo,
				}
			},
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"name": "Favorite 1",
					"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
					"clientName": "Client 1",
					"projectId": "01JH74AEFJRBTWZEF4KTB64M8V",
					"projectName": "Project 1",
					"taskId": null,
					"taskName": null,
					"billableRate": 90000,
					"comment": "Comment 1"
				},
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"name": "Favorite 2",
					"clientId": "01JH74AAT06A0P3AR2AGCFGRR5",
					"clientName": "Client 2",
					"projectId": "01JH74AKN2WM4FWXWJ4JBBZRPE",
					"projectName": "Project 2",
					"taskId": "01JH74AQYEE61Z2ANVRA77NPV6",
					"taskName": "Task One",
					"billableRate": 50000,
					"comment": ""
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

func TestFavoriteHandler_Create(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *FavoriteHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return Created",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				favoriteService := mocks.NewFavoriteServiceMock(t)
				favoriteService.EXPECT().
					Create(mock.Anything, domain.Favorite{
						UserID:       testUserID,
						Name:         "Favorite 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						TaskID:       func() *string { s := "01JHFSFMV37B9MVNS3YPKGQXFK"; return &s }(),
						BillableRate: 90000,
						Comment:      "Comment 1",
					}).
					Return(domain.Favorite{
						ID:           "01JH74Z7PNJDHC273QBCZBV1EQ",
						UserID:       testUserID,
						Name:         "Favorite 1",
						ClientID:     "01JH74Y924XMCA82DGZBVCKD7H",
						ClientName:   "Client 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName:  "Project 1",
						TaskID:       func() *string { s := "01JHFSFMV37B9MVNS3YPKGQXFK"; return &s }(),
						TaskName:     func() *string { s := "Task 2"; return &s }(),
						BillableRate: 90000,
						Comment:      "Comment 1",
					}, nil)

				return &FavoriteHandler{
					validate:        validator.New(),
					sanitizer:       sanitizer,
					favoriteService: favoriteService,
				}
			},
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHFSFMV37B9MVNS3YPKGQXFK",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusCreated,
			wantBody: `{
				"id": "01JH74Z7PNJDHC273QBCZBV1EQ",
				"name": "Favorite 1",
				"clientId": "01JH74Y924XMCA82DGZBVCKD7H",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"taskId": "01JHFSFMV37B9MVNS3YPKGQXFK",
				"taskName": "Task 2",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
		},
		{
			name: "Should return Created - without task and comment",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				favoriteService := mocks.NewFavoriteServiceMock(t)
				favoriteService.EXPECT().
					Create(mock.Anything, domain.Favorite{
						UserID:       testUserID,
						Name:         "Favorite 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						BillableRate: 90000,
					}).
					Return(domain.Favorite{
						ID:           "01JH74Z7PNJDHC273QBCZBV1EQ",
						UserID:       testUserID,
						Name:         "Favorite 1",
						ClientID:     "01JH74Y924XMCA82DGZBVCKD7H",
						ClientName:   "Client 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName:  "Project 1",
						BillableRate: 90000,
					}, nil)

				return &FavoriteHandler{
					validate:        validator.New(),
					sanitizer:       sanitizer,
					favoriteService: favoriteService,
				}
			},
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"billableRate": 90000
			}`,
			wantCode: http.StatusCreated,
			wantBody: `{
				"id": "01JH74Z7PNJDHC273QBCZBV1EQ",
				"name": "Favorite 1",
				"clientId": "01JH74Y924XMCA82DGZBVCKD7H",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"taskId": null,
				"taskName": null,
				"billableRate": 90000,
				"comment": ""
			}`,
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					logger: zap.NewNop(),
				}
			},
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid name",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"name": "very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHFSFMV37B9MVNS3YPKGQXFK",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "name"`,
		},
		{
			name: "Should return error - invalid projectId",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"name": "Favorite 1",
				"projectId": "123",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "projectId"`,
		},
		{
			name: "Should return error - invalid taskId",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "123",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "taskId"`,
		},
		{
			name: "Should return error - invalid billableRate",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": null,
				"billableRate": -100
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "billableRate"`,
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				favoriteService := mocks.NewFavoriteServiceMock(t)
				favoriteService.EXPECT().
					Create(mock.Anything, domain.Favorite{
						UserID:       testUserID,
						Name:         "Favorite 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						BillableRate: 90000,
					}).
					Return(domain.Favorite{}, domain.ErrProjectNotFound)

				return &FavoriteHandler{
					validate:        validator.New(),
					sanitizer:       sanitizer,
					favoriteService: favoriteService,
				}
			},
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "project not found",
		},
		{
			name: "Should return error - favorite already exists",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				favoriteService := mocks.NewFavoriteServiceMock(t)
				favoriteService.EXPECT().
					Create(mock.Anything, domain.Favorite{
						UserID:       testUserID,
						Name:         "Favorite 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						BillableRate: 90000,
					}).
					Return(domain.Favorite{}, domain.ErrAlreadyExists)

				return &FavoriteHandler{
					validate:        validator.New(),
					sanitizer:       sanitizer,
					favoriteService: favoriteService,
				}
			},
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"billableRate": 90000
			}`,
			wantCode: http.StatusConflict,
			wantBody: "favorite with this name already exists",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Post("/", tt.srvFunc(t).Create)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestFavoriteHandler_Get(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *FavoriteHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Favorite{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						Name:         "Favorite 1",
						ClientID:     "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:   "Client 1",
						ProjectID:    "01JH74AEFJRBTWZEF4KTB64M8V",
						ProjectName:  "Project 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
					}, nil)

				return &FavoriteHandler{
					validate:     validator.New(),
					favoriteRepo: favoriteRepo,
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"name": "Favorite 1",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"projectId": "01JH74AEFJRBTWZEF4KTB64M8V",
				"projectName": "Project 1",
				"taskId": null,
				"taskName": null,
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid favorite id",
		},
		{
			name: "Should return error - favorite not found",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Favorite{}, domain.ErrNotFound)

				return &FavoriteHandler{
					validate:     validator.New(),
					favoriteRepo: favoriteRepo,
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Get("/:id", tt.srvFunc(t).Get)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/"+tt.id, nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestFavoriteHandler_Update(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *FavoriteHandler
		id       string
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				favoriteService := mocks.NewFavoriteServiceMock(t)
				favoriteService.EXPECT().
					Update(mock.Anything, domain.Favorite{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						Name:         "Favorite 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						TaskID:       func() *string { s := "01JH74HARHZQ73K9WNTKKK01SQ"; return &s }(),
						BillableRate: 90000,
						Comment:      "Comment 1",
					}).
					Return(domain.Favorite{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						Name:         "Favorite 1",
						ClientID:     "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:   "Client 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName:  "Project 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
					}, nil)

				return &FavoriteHandler{
					validate:        validator.New(),
					sanitizer:       sanitizer,
					favoriteService: favoriteService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JH74HARHZQ73K9WNTKKK01SQ",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"name": "Favorite 1",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"taskId": null,
				"taskName": null,
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
		},
		{
			name: "Should return OK - without task and comment",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				favoriteService := mocks.NewFavoriteServiceMock(t)
				favoriteService.EXPECT().
					Update(mock.Anything, domain.Favorite{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						Name:         "Favorite 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						BillableRate: 90000,
					}).
					Return(domain.Favorite{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						Name:         "Favorite 1",
						ClientID:     "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:   "Client 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName:  "Project 1",
						BillableRate: 90000,
					}, nil)

				return &FavoriteHandler{
					validate:        validator.New(),
					sanitizer:       sanitizer,
					favoriteService: favoriteService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": null,
				"billableRate": 90000
			}`,
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"name": "Favorite 1",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"taskId": null,
				"taskName": null,
				"billableRate": 90000,
				"comment": ""
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					logger: zap.NewNop(),
				}
			},
			input:    "",
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid favorite id",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					logger: zap.NewNop(),
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid name",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JH74HARHZQ73K9WNTKKK01SQ",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "name"`,
		},
		{
			name: "Should return error - invalid projectId",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Favorite 1",
				"projectId": "123",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "projectId"`,
		},
		{
			name: "Should return error - invalid taskId",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "123",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "taskId"`,
		},
		{
			name: "Should return error - invalid billableRate",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": null,
				"billableRate": -100
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "billableRate"`,
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				favoriteService := mocks.NewFavoriteServiceMock(t)
				favoriteService.EXPECT().
					Update(mock.Anything, domain.Favorite{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						Name:         "Favorite 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						BillableRate: 90000,
					}).
					Return(domain.Favorite{}, domain.ErrProjectNotFound)

				return &FavoriteHandler{
					validate:        validator.New(),
					sanitizer:       sanitizer,
					favoriteService: favoriteService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "project not found",
		},
		{
			name: "Should return error - favorite not found",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				favoriteService := mocks.NewFavoriteServiceMock(t)
				favoriteService.EXPECT().
					Update(mock.Anything, domain.Favorite{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						Name:         "Favorite 1",
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						BillableRate: 90000,
					}).
					Return(domain.Favorite{}, domain.ErrNotFound)

				return &FavoriteHandler{
					validate:        validator.New(),
					sanitizer:       sanitizer,
					favoriteService: favoriteService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Favorite 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": null,
				"billableRate": 90000
			}`,
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Patch("/:id", tt.srvFunc(t).Update)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPatch, "/"+tt.id, strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestFavoriteHandler_Delete(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *FavoriteHandler
		id       string
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				favoriteService := mocks.NewFavoriteServiceMock(t)
				favoriteService.EXPECT().
					Delete(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(nil)

				return &FavoriteHandler{
					favoriteService: favoriteService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *FavoriteHandler {
				t.Helper()

				return &FavoriteHandler{}
			},
			id:       "123",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid favorite id",
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
