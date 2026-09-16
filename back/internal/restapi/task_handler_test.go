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

func TestTaskHandler_Index(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TaskHandler
		query    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByFilter(mock.Anything, testUserID, domain.TaskFilter{}).
					Return([]domain.Task{
						{
							ID:          "01JFEC2BRSGJKVG1FKNWSNV6HH",
							ClientID:    "01JH749ZAVM5MX8AMCJH7AKBX8",
							ClientName:  "Client 1",
							ProjectID:   "01JHFF8960RHPV35NE8GCZJYWZ",
							ProjectName: "Project 1",
							Name:        "Task 1",
							Comment:     "Comment 1",
							CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
						{
							ID:          "01JH741H5A4Y368PBHBTJFQ723",
							ClientID:    "01JH74AAT06A0P3AR2AGCFGRR5",
							ClientName:  "Client 2",
							ProjectID:   "01JHFF8G4S0RMW5DQNGZ0G3RH1",
							ProjectName: "Project 2",
							Name:        "Task 2",
							CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &TaskHandler{
					validate: validator.New(),
					taskRepo: taskRepo,
				}
			},
			query:    "",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
					"clientName": "Client 1",
					"projectId": "01JHFF8960RHPV35NE8GCZJYWZ",
					"projectName": "Project 1",
					"name": "Task 1",
					"comment": "Comment 1",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null,
					"completedAt": null
				},
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"clientId": "01JH74AAT06A0P3AR2AGCFGRR5",
					"clientName": "Client 2",
					"projectId": "01JHFF8G4S0RMW5DQNGZ0G3RH1",
					"projectName": "Project 2",
					"name": "Task 2",
					"comment": "",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null,
					"completedAt": null
				}
			]`,
		},
		{
			name: "Should return OK - with filter_archived_projects",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByFilter(mock.Anything, testUserID, domain.TaskFilter{
						WithoutArchivedProjects: true,
					}).
					Return([]domain.Task{
						{
							ID:          "01JFEC2BRSGJKVG1FKNWSNV6HH",
							ClientID:    "01JH749ZAVM5MX8AMCJH7AKBX8",
							ClientName:  "Client 1",
							ProjectID:   "01JHFF8960RHPV35NE8GCZJYWZ",
							ProjectName: "Project 1",
							Name:        "Task 1",
							Comment:     "Comment 1",
							CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
						{
							ID:          "01JH741H5A4Y368PBHBTJFQ723",
							ClientID:    "01JH74AAT06A0P3AR2AGCFGRR5",
							ClientName:  "Client 2",
							ProjectID:   "01JHFF8G4S0RMW5DQNGZ0G3RH1",
							ProjectName: "Project 2",
							Name:        "Task 2",
							CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &TaskHandler{
					validate: validator.New(),
					taskRepo: taskRepo,
				}
			},
			query:    "?filter_archived_projects=1",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
					"clientName": "Client 1",
					"projectId": "01JHFF8960RHPV35NE8GCZJYWZ",
					"projectName": "Project 1",
					"name": "Task 1",
					"comment": "Comment 1",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null,
					"completedAt": null
				},
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"clientId": "01JH74AAT06A0P3AR2AGCFGRR5",
					"clientName": "Client 2",
					"projectId": "01JHFF8G4S0RMW5DQNGZ0G3RH1",
					"projectName": "Project 2",
					"name": "Task 2",
					"comment": "",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null,
					"completedAt": null
				}
			]`,
		},
		{
			name: "Should return OK - with project_id",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByFilter(mock.Anything, testUserID, domain.TaskFilter{
						ProjectID: "01JHFF8G4S0RMW5DQNGZ0G3RH1",
					}).
					Return([]domain.Task{
						{
							ID:          "01JH741H5A4Y368PBHBTJFQ723",
							ClientID:    "01JH74AAT06A0P3AR2AGCFGRR5",
							ClientName:  "Client 2",
							ProjectID:   "01JHFF8G4S0RMW5DQNGZ0G3RH1",
							ProjectName: "Project 2",
							Name:        "Task 2",
							CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &TaskHandler{
					validate: validator.New(),
					taskRepo: taskRepo,
				}
			},
			query:    "?project_id=01JHFF8G4S0RMW5DQNGZ0G3RH1",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"clientId": "01JH74AAT06A0P3AR2AGCFGRR5",
					"clientName": "Client 2",
					"projectId": "01JHFF8G4S0RMW5DQNGZ0G3RH1",
					"projectName": "Project 2",
					"name": "Task 2",
					"comment": "",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null,
					"archivedAt": null,
					"completedAt": null
				}
			]`,
		},
		{
			name: "Should return error - invalid filter",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					validate: validator.New(),
				}
			},
			query:    "?filter_archived_projects=123",
			wantCode: http.StatusInternalServerError,
			wantBody: `query parser task index filter: failed to decode: schema: error converting value for "filter_archived_projects"`,
		},
		{
			name: "Should return error - invalid project_id",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					validate: validator.New(),
				}
			},
			query:    "?project_id=123",
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "project_id"`,
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

func TestTaskHandler_Create(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TaskHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return Created",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Create(mock.Anything, domain.Task{
						UserID:    testUserID,
						ProjectID: "01JH74H6YKJNXH6H1CA0XF81DZ",
						Name:      "Task 1",
						Comment:   "Comment 1",
					}).
					Return(domain.Task{
						ID:          "01JH74Z7PNJDHC273QBCZBV1EQ",
						UserID:      testUserID,
						ClientID:    "01JH74Y924XMCA82DGZBVCKD7H",
						ClientName:  "Client 1",
						ProjectID:   "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName: "Project 1",
						Name:        "Task 1",
						Comment:     "Comment 1",
						CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				return &TaskHandler{
					validate:    validator.New(),
					sanitizer:   sanitizer,
					taskService: taskService,
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"name": "Task 1",
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusCreated,
			wantBody: `{
				"id": "01JH74Z7PNJDHC273QBCZBV1EQ",
				"clientId": "01JH74Y924XMCA82DGZBVCKD7H",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"name": "Task 1",
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": null,
				"archivedAt": null,
				"completedAt": null
			}`,
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					logger: zap.NewNop(),
				}
			},
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid projectId",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"projectId": "123",
				"name": "Task 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "projectId"`,
		},
		{
			name: "Should return error - invalid name",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"name": "very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "name"`,
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Create(mock.Anything, domain.Task{
						UserID:    testUserID,
						ProjectID: "01JH74H6YKJNXH6H1CA0XF81DZ",
						Name:      "Task 1",
					}).
					Return(domain.Task{}, domain.ErrProjectNotFound)

				return &TaskHandler{
					validate:    validator.New(),
					sanitizer:   sanitizer,
					taskService: taskService,
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"name": "Task 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "project not found",
		},
		{
			name: "Should return error - task already exists",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Create(mock.Anything, domain.Task{
						UserID:    testUserID,
						ProjectID: "01JH74H6YKJNXH6H1CA0XF81DZ",
						Name:      "Task 1",
					}).
					Return(domain.Task{}, domain.ErrAlreadyExists)

				return &TaskHandler{
					validate:    validator.New(),
					sanitizer:   sanitizer,
					taskService: taskService,
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"name": "Task 1"
			}`,
			wantCode: http.StatusConflict,
			wantBody: "task with this name already exists",
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

func TestTaskHandler_Get(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TaskHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JH74H6YKJNXH6H1CA0XF81DZ").
					Return(domain.Task{
						ID:          "01JFEC2BRSGJKVG1FKNWSNV6HH",
						ClientID:    "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:  "Client 1",
						ProjectID:   "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName: "Project 1",
						Name:        "Task 1",
						Comment:     "Comment 1",
						CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				return &TaskHandler{
					validate: validator.New(),
					taskRepo: taskRepo,
				}
			},
			id:       "01JH74H6YKJNXH6H1CA0XF81DZ",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"name": "Task 1",
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": null,
				"archivedAt": null,
				"completedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid task id",
		},
		{
			name: "Should return error - task not found",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JH74H6YKJNXH6H1CA0XF81DZ").
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskHandler{
					validate: validator.New(),
					taskRepo: taskRepo,
				}
			},
			id:       "01JH74H6YKJNXH6H1CA0XF81DZ",
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

func TestTaskHandler_Update(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TaskHandler
		id       string
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Update(mock.Anything, domain.Task{
						ID:        "01JH74HARHZQ73K9WNTKKK01SQ",
						UserID:    testUserID,
						ProjectID: "01JH74H6YKJNXH6H1CA0XF81DZ",
						Name:      "Task 1 New",
						Comment:   "Comment 1 New",
					}).
					Return(domain.Task{
						ID:          "01JH74HARHZQ73K9WNTKKK01SQ",
						ClientID:    "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:  "Client 1",
						ProjectID:   "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName: "Project 1",
						Name:        "Task 1 New",
						Comment:     "Comment 1 New",
						CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:   func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TaskHandler{
					validate:    validator.New(),
					sanitizer:   sanitizer,
					taskService: taskService,
				}
			},
			id: "01JH74HARHZQ73K9WNTKKK01SQ",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"name": "Task 1 New",
				"comment": "Comment 1 New"
			}`,
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JH74HARHZQ73K9WNTKKK01SQ",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"name": "Task 1 New",
				"comment": "Comment 1 New",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": null,
				"completedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					logger: zap.NewNop(),
				}
			},
			input:    "",
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid task id",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					logger: zap.NewNop(),
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid projectId",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					validate: validator.New(),
				}
			},
			id: "01JH74HARHZQ73K9WNTKKK01SQ",
			input: `{
				"projectId": "123",
				"name": "Task 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "projectId"`,
		},
		{
			name: "Should return error - invalid name",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					validate: validator.New(),
				}
			},
			id: "01JH74HARHZQ73K9WNTKKK01SQ",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"name": "very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name very long name"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "name"`,
		},
		{
			name: "Should return error - task not found",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Update(mock.Anything, domain.Task{
						ID:        "01JH74HARHZQ73K9WNTKKK01SQ",
						UserID:    testUserID,
						ProjectID: "01JH74H6YKJNXH6H1CA0XF81DZ",
						Name:      "Task 1",
					}).
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskHandler{
					validate:    validator.New(),
					sanitizer:   sanitizer,
					taskService: taskService,
				}
			},
			id: "01JH74HARHZQ73K9WNTKKK01SQ",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"name": "Task 1"
			}`,
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Update(mock.Anything, domain.Task{
						ID:        "01JH74HARHZQ73K9WNTKKK01SQ",
						UserID:    testUserID,
						ProjectID: "01JH74H6YKJNXH6H1CA0XF81DZ",
						Name:      "Task 1",
					}).
					Return(domain.Task{}, domain.ErrProjectNotFound)

				return &TaskHandler{
					validate:    validator.New(),
					sanitizer:   sanitizer,
					taskService: taskService,
				}
			},
			id: "01JH74HARHZQ73K9WNTKKK01SQ",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"name": "Task 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "project not found",
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

func TestTaskHandler_Archive(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TaskHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Archive(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{
						ID:          "01JH74HARHZQ73K9WNTKKK01SQ",
						ClientID:    "01JHET28708HXZ3BQQP0KR86CY",
						ClientName:  "Client 1",
						ProjectID:   "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName: "Project 1",
						Name:        "Task 1",
						Comment:     "Comment 1",
						CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:   func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
						ArchivedAt:  func() *time.Time { t := time.Date(2024, 12, 22, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JH74HARHZQ73K9WNTKKK01SQ",
				"clientId": "01JHET28708HXZ3BQQP0KR86CY",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"name": "Task 1",
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": "2024-12-22T14:15:00Z",
				"completedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid task id",
		},
		{
			name: "Should return error - task not found",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Archive(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
		{
			name: "Should return error - impossible action",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Archive(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{}, domain.ErrImpossibleAction)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			wantCode: http.StatusConflict,
			wantBody: "impossible action",
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

func TestTaskHandler_Unarchive(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TaskHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Unarchive(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{
						ID:          "01JH74HARHZQ73K9WNTKKK01SQ",
						ClientID:    "01JHET28708HXZ3BQQP0KR86CY",
						ClientName:  "Client 1",
						ProjectID:   "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName: "Project 1",
						Name:        "Task 1",
						Comment:     "Comment 1",
						CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:   func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JH74HARHZQ73K9WNTKKK01SQ",
				"clientId": "01JHET28708HXZ3BQQP0KR86CY",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"name": "Task 1",
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": null,
				"completedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid task id",
		},
		{
			name: "Should return error - task not found",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Unarchive(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
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

func TestTaskHandler_Complete(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TaskHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Complete(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{
						ID:          "01JH74HARHZQ73K9WNTKKK01SQ",
						ClientID:    "01JHET28708HXZ3BQQP0KR86CY",
						ClientName:  "Client 1",
						ProjectID:   "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName: "Project 1",
						Name:        "Task 1",
						Comment:     "Comment 1",
						CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:   func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
						CompletedAt: func() *time.Time { t := time.Date(2024, 12, 22, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JH74HARHZQ73K9WNTKKK01SQ",
				"clientId": "01JHET28708HXZ3BQQP0KR86CY",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"name": "Task 1",
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": null,
				"completedAt": "2024-12-22T14:15:00Z"
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid task id",
		},
		{
			name: "Should return error - task not found",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Complete(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
		{
			name: "Should return error - impossible action",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Complete(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{}, domain.ErrImpossibleAction)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			wantCode: http.StatusConflict,
			wantBody: "impossible action",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Post("/:id", tt.srvFunc(t).Complete)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/"+tt.id, nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestTaskHandler_Incomplete(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TaskHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Incomplete(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{
						ID:          "01JH74HARHZQ73K9WNTKKK01SQ",
						ClientID:    "01JHET28708HXZ3BQQP0KR86CY",
						ClientName:  "Client 1",
						ProjectID:   "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName: "Project 1",
						Name:        "Task 1",
						Comment:     "Comment 1",
						CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:   func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JH74HARHZQ73K9WNTKKK01SQ",
				"clientId": "01JHET28708HXZ3BQQP0KR86CY",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"name": "Task 1",
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z",
				"archivedAt": null,
				"completedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid task id",
		},
		{
			name: "Should return error - task not found",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Incomplete(mock.Anything, testUserID, "01JH74HARHZQ73K9WNTKKK01SQ").
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       "01JH74HARHZQ73K9WNTKKK01SQ",
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

			app.Post("/:id", tt.srvFunc(t).Incomplete)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/"+tt.id, nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestTaskHandler_Report(t *testing.T) {
	const (
		testUserID = "01JMGSDTRNFGDE3558DZNGHVDV"
		testTaskID = "01JMGSJ9C8HFY0AGZP7RDWSS8N"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TaskHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Report(mock.Anything, testUserID, testTaskID).
					Return(domain.TaskReport{
						DurationSeconds:  3600,
						BillableAmount:   50000,
						UniqueDatesCount: 2,
					}, nil)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       testTaskID,
			wantCode: http.StatusOK,
			wantBody: `{
				"durationSeconds": 3600,
				"billableAmount": 50000,
				"uniqueDatesCount": 2
			}`,
		},
		{
			name: "Should return OK - report empty",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				taskService := mocks.NewTaskServiceMock(t)
				taskService.EXPECT().
					Report(mock.Anything, testUserID, testTaskID).
					Return(domain.TaskReport{}, nil)

				return &TaskHandler{
					validate:    validator.New(),
					taskService: taskService,
				}
			},
			id:       testTaskID,
			wantCode: http.StatusOK,
			wantBody: `{
				"durationSeconds": 0,
				"billableAmount": 0,
				"uniqueDatesCount": 0
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TaskHandler {
				t.Helper()

				return &TaskHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid task id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Get("/:id", tt.srvFunc(t).Report)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/"+tt.id, nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}
