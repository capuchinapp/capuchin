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

func TestProjectHandler_Index(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ProjectHandler
		query    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByFilter(mock.Anything, testUserID, domain.ProjectFilter{}).
					Return([]domain.Project{
						{
							ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
							Name:         "Project 1",
							ClientID:     "01JH749ZAVM5MX8AMCJH7AKBX8",
							ClientName:   "Client 1",
							BillableRate: 90000,
							Comment:      "Comment 1",
							CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
						{
							ID:           "01JH741H5A4Y368PBHBTJFQ723",
							Name:         "Project 2",
							ClientID:     "01JH74AAT06A0P3AR2AGCFGRR5",
							ClientName:   "Client 2",
							BillableRate: 50000,
							CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &ProjectHandler{
					validate:    validator.New(),
					projectRepo: projectRepo,
				}
			},
			query:    "",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"name": "Project 1",
					"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
					"clientName": "Client 1",
					"billableRate": 90000,
					"comment": "Comment 1",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null
				},
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"name": "Project 2",
					"clientId": "01JH74AAT06A0P3AR2AGCFGRR5",
					"clientName": "Client 2",
					"billableRate": 50000,
					"comment": "",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null
				}
			]`,
		},
		{
			name: "Should return OK - with filter_archived_clients",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByFilter(mock.Anything, testUserID, domain.ProjectFilter{
						WithoutArchivedClients: true,
					}).
					Return([]domain.Project{
						{
							ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
							Name:         "Project 1",
							ClientID:     "01JH749ZAVM5MX8AMCJH7AKBX8",
							ClientName:   "Client 1",
							BillableRate: 90000,
							Comment:      "Comment 1",
							CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
						{
							ID:           "01JH741H5A4Y368PBHBTJFQ723",
							Name:         "Project 2",
							ClientID:     "01JH74AAT06A0P3AR2AGCFGRR5",
							ClientName:   "Client 2",
							BillableRate: 50000,
							CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &ProjectHandler{
					validate:    validator.New(),
					projectRepo: projectRepo,
				}
			},
			query:    "?filter_archived_clients=1",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"name": "Project 1",
					"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
					"clientName": "Client 1",
					"billableRate": 90000,
					"comment": "Comment 1",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null
				},
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"name": "Project 2",
					"clientId": "01JH74AAT06A0P3AR2AGCFGRR5",
					"clientName": "Client 2",
					"billableRate": 50000,
					"comment": "",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null
				}
			]`,
		},
		{
			name: "Should return OK - with client_id",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByFilter(mock.Anything, testUserID, domain.ProjectFilter{
						ClientID: "01JH749ZAVM5MX8AMCJH7AKBX8",
					}).
					Return([]domain.Project{
						{
							ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
							Name:         "Project 1",
							ClientID:     "01JH749ZAVM5MX8AMCJH7AKBX8",
							ClientName:   "Client 1",
							BillableRate: 90000,
							Comment:      "Comment 1",
							CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &ProjectHandler{
					validate:    validator.New(),
					projectRepo: projectRepo,
				}
			},
			query:    "?client_id=01JH749ZAVM5MX8AMCJH7AKBX8",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"name": "Project 1",
					"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
					"clientName": "Client 1",
					"billableRate": 90000,
					"comment": "Comment 1",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null
				}
			]`,
		},
		{
			name: "Should return error - invalid filter",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					validate: validator.New(),
				}
			},
			query:    "?filter_archived_clients=123",
			wantCode: http.StatusInternalServerError,
			wantBody: `query parser project index filter: failed to decode: schema: error converting value for "filter_archived_clients"`,
		},
		{
			name: "Should return error - invalid client_id",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					validate: validator.New(),
				}
			},
			query:    "?client_id=123",
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "client_id"`,
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

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/"+tt.query, nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestProjectHandler_Create(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ProjectHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return Created",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Create(mock.Anything, domain.Project{
						UserID:       testUserID,
						ClientID:     "01JH74Y924XMCA82DGZBVCKD7H",
						Name:         "Project 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
					}).
					Return(domain.Project{
						ID:           "01JH74Z7PNJDHC273QBCZBV1EQ",
						UserID:       testUserID,
						Name:         "Project 1",
						ClientID:     "01JH74Y924XMCA82DGZBVCKD7H",
						ClientName:   "Client 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				return &ProjectHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					projectService: projectService,
				}
			},
			input: `{
				"clientId": "01JH74Y924XMCA82DGZBVCKD7H",
				"name": "Project 1",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusCreated,
			wantBody: `{
				"id": "01JH74Z7PNJDHC273QBCZBV1EQ",
				"clientId": "01JH74Y924XMCA82DGZBVCKD7H",
				"clientName": "Client 1",
				"name": "Project 1",
				"billableRate": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": null,
				"archivedAt": null
			}`,
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					logger: zap.NewNop(),
				}
			},
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid clientId",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"clientId": "123",
				"name": "Project 1",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "clientId"`,
		},
		{
			name: "Should return error - invalid name",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"clientId": "01JH74Y924XMCA82DGZBVCKD7H",
				"name": "very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "name"`,
		},
		{
			name: "Should return error - invalid billableRate",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"clientId": "01JH74Y924XMCA82DGZBVCKD7H",
				"name": "Project 1",
				"billableRate": -100
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "billableRate"`,
		},
		{
			name: "Should return error - client not found",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Create(mock.Anything, domain.Project{
						UserID:       testUserID,
						ClientID:     "01JH74Y924XMCA82DGZBVCKD7H",
						Name:         "Project 1",
						BillableRate: 90000,
					}).
					Return(domain.Project{}, domain.ErrClientNotFound)

				return &ProjectHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					projectService: projectService,
				}
			},
			input: `{
				"clientId": "01JH74Y924XMCA82DGZBVCKD7H",
				"name": "Project 1",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "client not found",
		},
		{
			name: "Should return error - project already exists",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Create(mock.Anything, domain.Project{
						UserID:       testUserID,
						ClientID:     "01JH74Y924XMCA82DGZBVCKD7H",
						Name:         "Project 1",
						BillableRate: 90000,
					}).
					Return(domain.Project{}, domain.ErrAlreadyExists)

				return &ProjectHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					projectService: projectService,
				}
			},
			input: `{
				"clientId": "01JH74Y924XMCA82DGZBVCKD7H",
				"name": "Project 1",
				"billableRate": 90000
			}`,
			wantCode: http.StatusConflict,
			wantBody: "project with this name already exists",
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

func TestProjectHandler_Get(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ProjectHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Project{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						Name:         "Project 1",
						ClientID:     "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:   "Client 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				return &ProjectHandler{
					validate:    validator.New(),
					projectRepo: projectRepo,
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"name": "Project 1",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"billableRate": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": null,
				"archivedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid project id",
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Project{}, domain.ErrNotFound)

				return &ProjectHandler{
					validate:    validator.New(),
					projectRepo: projectRepo,
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

func TestProjectHandler_Update(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ProjectHandler
		id       string
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Update(mock.Anything, domain.Project{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						ClientID:     "01JH74HARHZQ73K9WNTKKK01SQ",
						Name:         "Project 1 New",
						BillableRate: 60000,
						Comment:      "Comment 1 New",
					}).
					Return(domain.Project{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						ClientID:     "01JH74HARHZQ73K9WNTKKK01SQ",
						ClientName:   "Client 1",
						Name:         "Project 1 New",
						BillableRate: 60000,
						Comment:      "Comment 1 New",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &ProjectHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					projectService: projectService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"clientId": "01JH74HARHZQ73K9WNTKKK01SQ",
				"name": "Project 1 New",
				"billableRate": 60000,
				"comment": "Comment 1 New"
			}`,
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"clientId": "01JH74HARHZQ73K9WNTKKK01SQ",
				"clientName": "Client 1",
				"name": "Project 1 New",
				"billableRate": 60000,
				"comment": "Comment 1 New",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					logger: zap.NewNop(),
				}
			},
			input:    "",
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid project id",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					logger: zap.NewNop(),
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid clientId",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"clientId": "123",
				"name": "Project 1",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "clientId"`,
		},
		{
			name: "Should return error - invalid name",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"clientId": "01JH74HARHZQ73K9WNTKKK01SQ",
				"name": "very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "name"`,
		},
		{
			name: "Should return error - invalid billableRate",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"clientId": "01JH74HARHZQ73K9WNTKKK01SQ",
				"name": "Project 1",
				"billableRate": -100
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "billableRate"`,
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Update(mock.Anything, domain.Project{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						ClientID:     "01JH74HARHZQ73K9WNTKKK01SQ",
						Name:         "Project 1",
						BillableRate: 90000,
					}).
					Return(domain.Project{}, domain.ErrNotFound)

				return &ProjectHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					projectService: projectService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"clientId": "01JH74HARHZQ73K9WNTKKK01SQ",
				"name": "Project 1",
				"billableRate": 90000
			}`,
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
		{
			name: "Should return error - client not found",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Update(mock.Anything, domain.Project{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						ClientID:     "01JH74HARHZQ73K9WNTKKK01SQ",
						Name:         "Project 1",
						BillableRate: 90000,
					}).
					Return(domain.Project{}, domain.ErrClientNotFound)

				return &ProjectHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					projectService: projectService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"clientId": "01JH74HARHZQ73K9WNTKKK01SQ",
				"name": "Project 1",
				"billableRate": 90000
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "client not found",
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

func TestProjectHandler_Archive(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ProjectHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Archive(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Project{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						ClientID:     "01JHET28708HXZ3BQQP0KR86CY",
						ClientName:   "Client 1",
						Name:         "Project 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
						ArchivedAt:   func() *time.Time { t := time.Date(2024, 12, 22, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &ProjectHandler{
					validate:       validator.New(),
					projectService: projectService,
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"clientId": "01JHET28708HXZ3BQQP0KR86CY",
				"clientName": "Client 1",
				"name": "Project 1",
				"billableRate": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": "2024-12-22T14:15:00Z"
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid project id",
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Archive(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Project{}, domain.ErrNotFound)

				return &ProjectHandler{
					validate:       validator.New(),
					projectService: projectService,
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

func TestProjectHandler_Unarchive(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *ProjectHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Unarchive(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Project{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						ClientID:     "01JHET28708HXZ3BQQP0KR86CY",
						ClientName:   "Client 1",
						Name:         "Project 1",
						BillableRate: 90000,
						Comment:      "Comment 1",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &ProjectHandler{
					validate:       validator.New(),
					projectService: projectService,
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"clientId": "01JHET28708HXZ3BQQP0KR86CY",
				"clientName": "Client 1",
				"name": "Project 1",
				"billableRate": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				return &ProjectHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid project id",
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *ProjectHandler {
				t.Helper()

				projectService := mocks.NewProjectServiceMock(t)
				projectService.EXPECT().
					Unarchive(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(domain.Project{}, domain.ErrNotFound)

				return &ProjectHandler{
					validate:       validator.New(),
					projectService: projectService,
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
