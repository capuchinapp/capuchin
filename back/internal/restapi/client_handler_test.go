package restapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"capuchin/internal/domain"
	"capuchin/internal/restapi/mocks"
	"capuchin/internal/validator"
)

func TestClientHandler_Index(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ClientHandler
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindAll(mock.Anything, testUserID).
					Return([]domain.Client{
						{
							ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
							Name:         "Client 1",
							BillableRate: 90000,
							Comment:      "Comment 1",
							CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
						{
							ID:           "01JH741H5A4Y368PBHBTJFQ723",
							Name:         "Client 2",
							BillableRate: 50000,
							CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &ClientHandler{
					clientRepo: clientRepo,
				}
			},
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"name": "Client 1",
					"billableRate": 90000,
					"comment": "Comment 1",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null
				},
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"name": "Client 2",
					"billableRate": 50000,
					"comment": "",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null
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

func TestClientHandler_Create(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ClientHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return Created",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				clientService := mocks.NewClientServiceMock(t)
				clientService.EXPECT().
					Create(mock.Anything, domain.Client{
						UserID:       testUserID,
						Name:         "Client 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
					}).
					Return(domain.Client{
						ID:           "01JH74Z7PNJDHC273QBCZBV1EQ",
						UserID:       testUserID,
						Name:         "Client 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				return &ClientHandler{
					validate:      validator.New(),
					sanitizer:     sanitizer,
					clientService: clientService,
				}
			},
			input: `{
				"name": "Client 1",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusCreated,
			wantBody: `{
				"id": "01JH74Z7PNJDHC273QBCZBV1EQ",
				"name": "Client 1",
				"billableRate": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": null,
				"archivedAt": null
			}`,
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
					logger: zap.NewNop(),
				}
			},
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid name",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"name": "very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "name"`,
		},
		{
			name: "Should return error - invalid billableRate",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"name": "Client 1",
				"billableRate": -100
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "billableRate"`,
		},
		{
			name: "Should return error - client already exists",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				clientService := mocks.NewClientServiceMock(t)
				clientService.EXPECT().
					Create(mock.Anything, domain.Client{
						UserID:       testUserID,
						Name:         "Client 1",
						BillableRate: 90000,
					}).
					Return(domain.Client{}, domain.ErrAlreadyExists)

				return &ClientHandler{
					validate:      validator.New(),
					sanitizer:     sanitizer,
					clientService: clientService,
				}
			},
			input: `{
				"name": "Client 1",
				"billableRate": 90000
			}`,
			wantCode: http.StatusConflict,
			wantBody: "client with this name already exists",
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

func TestClientHandler_Get(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ClientHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Client{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						Name:         "Client 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				return &ClientHandler{
					validate:   validator.New(),
					clientRepo: clientRepo,
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"name": "Client 1",
				"billableRate": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": null,
				"archivedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid client id",
		},
		{
			name: "Should return error - client not found",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Client{}, domain.ErrNotFound)

				return &ClientHandler{
					validate:   validator.New(),
					clientRepo: clientRepo,
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

func TestClientHandler_Update(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ClientHandler
		id       string
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				clientService := mocks.NewClientServiceMock(t)
				clientService.EXPECT().
					Update(mock.Anything, domain.Client{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						Name:         "Client 1 New",
						BillableRate: 60000,
						Comment:      "Comment 1 New",
					}).
					Return(domain.Client{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						Name:         "Client 1 New",
						BillableRate: 60000,
						Comment:      "Comment 1 New",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &ClientHandler{
					validate:      validator.New(),
					sanitizer:     sanitizer,
					clientService: clientService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Client 1 New",
				"billableRate": 60000,
				"comment": "Comment 1 New"
			}`,
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"name": "Client 1 New",
				"billableRate": 60000,
				"comment": "Comment 1 New",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
					logger: zap.NewNop(),
				}
			},
			input:    "",
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid client id",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
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
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "name"`,
		},
		{
			name: "Should return error - invalid billableRate",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Client 1",
				"billableRate": -100
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "billableRate"`,
		},
		{
			name: "Should return error - client not found",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				clientService := mocks.NewClientServiceMock(t)
				clientService.EXPECT().
					Update(mock.Anything, domain.Client{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						Name:         "Client 1",
						BillableRate: 90000,
					}).
					Return(domain.Client{}, domain.ErrNotFound)

				return &ClientHandler{
					validate:      validator.New(),
					sanitizer:     sanitizer,
					clientService: clientService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"name": "Client 1",
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

func TestClientHandler_Archive(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ClientHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				clientService := mocks.NewClientServiceMock(t)
				clientService.EXPECT().
					Archive(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Client{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						Name:         "Client 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
						ArchivedAt:   func() *time.Time { t := time.Date(2024, 12, 22, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &ClientHandler{
					validate:      validator.New(),
					clientService: clientService,
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"name": "Client 1",
				"billableRate": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": "2024-12-22T14:15:00Z"
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid client id",
		},
		{
			name: "Should return error - client not found",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				clientService := mocks.NewClientServiceMock(t)
				clientService.EXPECT().
					Archive(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Client{}, domain.ErrNotFound)

				return &ClientHandler{
					validate:      validator.New(),
					clientService: clientService,
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

			app.Post("/:id", tt.srvFunc(t).Archive)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/"+tt.id, nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestClientHandler_Unarchive(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ClientHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				clientService := mocks.NewClientServiceMock(t)
				clientService.EXPECT().
					Unarchive(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Client{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						Name:         "Client 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &ClientHandler{
					validate:      validator.New(),
					clientService: clientService,
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"name": "Client 1",
				"billableRate": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				return &ClientHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid client id",
		},
		{
			name: "Should return error - client not found",
			srvFunc: func(t *testing.T) *ClientHandler {
				t.Helper()

				clientService := mocks.NewClientServiceMock(t)
				clientService.EXPECT().
					Unarchive(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Client{}, domain.ErrNotFound)

				return &ClientHandler{
					validate:      validator.New(),
					clientService: clientService,
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

			app.Post("/:id", tt.srvFunc(t).Unarchive)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/"+tt.id, nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}
