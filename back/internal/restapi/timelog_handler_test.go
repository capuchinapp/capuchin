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

func TestTimelogHandler_Index(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TimelogHandler
		query    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByFilter(mock.Anything, testUserID, domain.TimelogFilter{
						DateFrom: "2024-12-18",
						DateTo:   "2024-12-18",
					}).
					Return([]domain.Timelog{
						{
							ID:              "01JFEC2BRSGJKVG1FKNWSNV6HH",
							ClientID:        "01JH749ZAVM5MX8AMCJH7AKBX8",
							ClientName:      "Client 1",
							ProjectID:       "01JHFF8960RHPV35NE8GCZJYWZ",
							ProjectName:     "Project 1",
							TaskID:          func() *string { id := "01JHFW9PWTTSNGNT99ZK9Q8FAQ"; return &id }(),
							TaskName:        func() *string { s := "Task 1"; return &s }(),
							TaskCompletedAt: func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
							Date:            "2024-12-18",
							TimeStart:       "14:15:00",
							DurationSeconds: 0,
							BillableRate:    90000,
							BillableAmount:  0,
							Comment:         "Comment 1",
							CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
						{
							ID:              "01JH741H5A4Y368PBHBTJFQ723",
							ClientID:        "01JH74AAT06A0P3AR2AGCFGRR5",
							ClientName:      "Client 2",
							ProjectID:       "01JHFF8G4S0RMW5DQNGZ0G3RH1",
							ProjectName:     "Project 2",
							Date:            "2024-12-18",
							TimeStart:       "14:15:00",
							TimeEnd:         func() *string { s := "15:15:00"; return &s }(),
							DurationSeconds: 3600,
							BillableRate:    90000,
							BillableAmount:  90000,
							CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &TimelogHandler{
					validate:    validator.New(),
					timelogRepo: timelogRepo,
				}
			},
			query:    "?date_from=2024-12-18&date_to=2024-12-18",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
					"clientName": "Client 1",
					"projectId": "01JHFF8960RHPV35NE8GCZJYWZ",
					"projectName": "Project 1",
					"taskId": "01JHFW9PWTTSNGNT99ZK9Q8FAQ",
					"taskName": "Task 1",
					"taskCompletedAt": "2024-12-20T14:15:00Z",
					"date": "2024-12-18",
					"timeStart": "14:15:00",
					"timeEnd": null,
					"durationSeconds": 0,
					"billableRate": 90000,
					"billableAmount": 0,
					"comment": "Comment 1",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null
				},
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"clientId": "01JH74AAT06A0P3AR2AGCFGRR5",
					"clientName": "Client 2",
					"projectId": "01JHFF8G4S0RMW5DQNGZ0G3RH1",
					"projectName": "Project 2",
					"taskId": null,
					"taskName": null,
					"taskCompletedAt": null,
					"date": "2024-12-18",
					"timeStart": "14:15:00",
					"timeEnd": "15:15:00",
					"durationSeconds": 3600,
					"billableRate": 90000,
					"billableAmount": 90000,
					"comment": "",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null
				}
			]`,
		},
		{
			name: "Should return OK - with client_id",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByFilter(mock.Anything, testUserID, domain.TimelogFilter{
						DateFrom: "2024-12-18",
						DateTo:   "2024-12-18",
						ClientID: "01JH749ZAVM5MX8AMCJH7AKBX8",
					}).
					Return([]domain.Timelog{
						{
							ID:              "01JFEC2BRSGJKVG1FKNWSNV6HH",
							ClientID:        "01JH749ZAVM5MX8AMCJH7AKBX8",
							ClientName:      "Client 1",
							ProjectID:       "01JHFF8960RHPV35NE8GCZJYWZ",
							ProjectName:     "Project 1",
							TaskID:          func() *string { id := "01JHFW9PWTTSNGNT99ZK9Q8FAQ"; return &id }(),
							TaskName:        func() *string { s := "Task 2"; return &s }(),
							TaskCompletedAt: func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
							Date:            "2024-12-18",
							TimeStart:       "14:15:00",
							DurationSeconds: 0,
							BillableRate:    90000,
							BillableAmount:  0,
							Comment:         "Comment 1",
							CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &TimelogHandler{
					validate:    validator.New(),
					timelogRepo: timelogRepo,
				}
			},
			query:    "?date_from=2024-12-18&date_to=2024-12-18&client_id=01JH749ZAVM5MX8AMCJH7AKBX8",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
					"clientName": "Client 1",
					"projectId": "01JHFF8960RHPV35NE8GCZJYWZ",
					"projectName": "Project 1",
					"taskId": "01JHFW9PWTTSNGNT99ZK9Q8FAQ",
					"taskName": "Task 2",
					"taskCompletedAt": "2024-12-20T14:15:00Z",
					"date": "2024-12-18",
					"timeStart": "14:15:00",
					"timeEnd": null,
					"durationSeconds": 0,
					"billableRate": 90000,
					"billableAmount": 0,
					"comment": "Comment 1",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null
				}
			]`,
		},
		{
			name: "Should return OK - with project_id",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByFilter(mock.Anything, testUserID, domain.TimelogFilter{
						DateFrom:  "2024-12-18",
						DateTo:    "2024-12-18",
						ProjectID: "01JHFF8G4S0RMW5DQNGZ0G3RH1",
					}).
					Return([]domain.Timelog{
						{
							ID:              "01JH741H5A4Y368PBHBTJFQ723",
							ClientID:        "01JH74AAT06A0P3AR2AGCFGRR5",
							ClientName:      "Client 2",
							ProjectID:       "01JHFF8G4S0RMW5DQNGZ0G3RH1",
							ProjectName:     "Project 2",
							Date:            "2024-12-18",
							TimeStart:       "14:15:00",
							TimeEnd:         func() *string { s := "15:16:00"; return &s }(),
							DurationSeconds: 3600,
							BillableRate:    90000,
							BillableAmount:  90000,
							CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &TimelogHandler{
					validate:    validator.New(),
					timelogRepo: timelogRepo,
				}
			},
			query:    "?date_from=2024-12-18&date_to=2024-12-18&project_id=01JHFF8G4S0RMW5DQNGZ0G3RH1",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"clientId": "01JH74AAT06A0P3AR2AGCFGRR5",
					"clientName": "Client 2",
					"projectId": "01JHFF8G4S0RMW5DQNGZ0G3RH1",
					"projectName": "Project 2",
					"taskId": null,
					"taskName": null,
					"taskCompletedAt": null,
					"date": "2024-12-18",
					"timeStart": "14:15:00",
					"timeEnd": "15:16:00",
					"durationSeconds": 3600,
					"billableRate": 90000,
					"billableAmount": 90000,
					"comment": "",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null
				}
			]`,
		},
		{
			name: "Should return error - invalid date_from",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			query:    "?date_from=123&date_to=2024-12-18",
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "date_from"`,
		},
		{
			name: "Should return error - invalid date_to",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			query:    "?date_from=2024-12-18&date_to=123",
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "date_to"`,
		},
		{
			name: "Should return error - invalid client_id",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			query:    "?date_from=2024-12-18&date_to=2024-12-18&client_id=123",
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "client_id"`,
		},
		{
			name: "Should return error - invalid project_id",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			query:    "?date_from=2024-12-18&date_to=2024-12-18&project_id=123",
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

func TestTimelogHandler_Create(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TimelogHandler
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return Created",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					Create(mock.Anything, domain.Timelog{
						UserID:       testUserID,
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						Date:         "2024-12-18",
						TimeStart:    "14:15:00",
						TimeEnd:      func() *string { s := "15:17:00"; return &s }(),
						BillableRate: 90000,
						Comment:      "Comment 1",
						TaskID:       func() *string { id := "01JHHJZEW6TMBS2B56Z0BD59QC"; return &id }(),
					}).
					Return(domain.Timelog{
						ID:              "01JFEC2BRSGJKVG1FKNWSNV6HH",
						ClientID:        "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:      "Client 1",
						ProjectID:       "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName:     "Project 1",
						TaskID:          func() *string { id := "01JHHJZEW6TMBS2B56Z0BD59QC"; return &id }(),
						TaskName:        func() *string { s := "Task 3"; return &s }(),
						Date:            "2024-12-18",
						TimeStart:       "14:15:00",
						TimeEnd:         func() *string { s := "15:17:00"; return &s }(),
						DurationSeconds: 3600,
						BillableRate:    90000,
						BillableAmount:  90000,
						Comment:         "Comment 1",
						CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				return &TimelogHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					timelogService: timelogService,
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHJZEW6TMBS2B56Z0BD59QC",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:17:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusCreated,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"taskId": "01JHHJZEW6TMBS2B56Z0BD59QC",
				"taskName": "Task 3",
				"taskCompletedAt": null,
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:17:00",
				"durationSeconds": 3600,
				"billableRate": 90000,
				"billableAmount": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": null
			}`,
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					logger: zap.NewNop(),
				}
			},
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid projectId",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"projectId": "123",
				"taskId": "01JHHJZEW6TMBS2B56Z0BD59QC",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:15:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "projectId"`,
		},
		{
			name: "Should return error - invalid date",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHJZEW6TMBS2B56Z0BD59QC",
				"date": "123",
				"timeStart": "14:15:00",
				"timeEnd": "15:15:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "date"`,
		},
		{
			name: "Should return error - invalid timeStart",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHJZEW6TMBS2B56Z0BD59QC",
				"date": "2024-12-18",
				"timeStart": "123",
				"timeEnd": "15:15:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "timeStart"`,
		},
		{
			name: "Should return error - invalid timeEnd",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHJZEW6TMBS2B56Z0BD59QC",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "123",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "timeEnd"`,
		},
		{
			name: "Should return error - invalid billableRate",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHJZEW6TMBS2B56Z0BD59QC",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:15:00",
				"billableRate": -1,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "billableRate"`,
		},
		{
			name: "Should return error - invalid taskId",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "123",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:15:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "taskId"`,
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					Create(mock.Anything, domain.Timelog{
						UserID:       testUserID,
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						Date:         "2024-12-18",
						TimeStart:    "14:15:00",
						TimeEnd:      func() *string { s := "15:30:00"; return &s }(),
						BillableRate: 90000,
						Comment:      "Comment 1",
						TaskID:       func() *string { id := "01JHHK2K6K8QNVTD39TYM53WX2"; return &id }(),
					}).
					Return(domain.Timelog{}, domain.ErrProjectNotFound)

				return &TimelogHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					timelogService: timelogService,
				}
			},
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHK2K6K8QNVTD39TYM53WX2",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:30:00",
				"billableRate": 90000,
				"comment": "Comment 1"
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

			app.Post("/", tt.srvFunc(t).Create)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestTimelogHandler_Get(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TimelogHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JH74H6YKJNXH6H1CA0XF81DZ").
					Return(domain.Timelog{
						ID:              "01JH74H6YKJNXH6H1CA0XF81DZ",
						ClientID:        "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:      "Client 1",
						ProjectID:       "01JHFF8960RHPV35NE8GCZJYWZ",
						ProjectName:     "Project 1",
						TaskID:          func() *string { id := "01JHHK079ZB654VMAH2JNRGY94"; return &id }(),
						TaskName:        func() *string { s := "Task 4"; return &s }(),
						TaskCompletedAt: func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
						Date:            "2024-12-18",
						TimeStart:       "14:15:00",
						DurationSeconds: 0,
						BillableRate:    90000,
						BillableAmount:  0,
						Comment:         "Comment 1",
						CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				return &TimelogHandler{
					validate:    validator.New(),
					timelogRepo: timelogRepo,
				}
			},
			id:       "01JH74H6YKJNXH6H1CA0XF81DZ",
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"projectId": "01JHFF8960RHPV35NE8GCZJYWZ",
				"projectName": "Project 1",
				"taskId": "01JHHK079ZB654VMAH2JNRGY94",
				"taskName": "Task 4",
				"taskCompletedAt": "2024-12-20T14:15:00Z",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": null,
				"durationSeconds": 0,
				"billableRate": 90000,
				"billableAmount": 0,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": null
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid timelog id",
		},
		{
			name: "Should return error - timelog not found",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(mock.Anything, testUserID, "01JH74H6YKJNXH6H1CA0XF81DZ").
					Return(domain.Timelog{}, domain.ErrNotFound)

				return &TimelogHandler{
					validate:    validator.New(),
					timelogRepo: timelogRepo,
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

func TestTimelogHandler_Update(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TimelogHandler
		id       string
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					Update(mock.Anything, domain.Timelog{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						Date:         "2024-12-18",
						TimeStart:    "14:15:00",
						TimeEnd:      func() *string { s := "15:18:00"; return &s }(),
						BillableRate: 90000,
						Comment:      "Comment 1 New",
						TaskID:       func() *string { id := "01JHHK0FYG9MCRXNKPQJMPNV42"; return &id }(),
					}).
					Return(domain.Timelog{
						ID:              "01JFEC2BRSGJKVG1FKNWSNV6HH",
						ClientID:        "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:      "Client 1",
						ProjectID:       "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName:     "Project 1",
						TaskID:          func() *string { id := "01JHHK0FYG9MCRXNKPQJMPNV42"; return &id }(),
						TaskName:        func() *string { s := "Task 5"; return &s }(),
						Date:            "2024-12-18",
						TimeStart:       "14:15:00",
						TimeEnd:         func() *string { s := "15:18:00"; return &s }(),
						DurationSeconds: 3600,
						BillableRate:    90000,
						BillableAmount:  90000,
						Comment:         "Comment 1 New",
						CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:       func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TimelogHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					timelogService: timelogService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHK0FYG9MCRXNKPQJMPNV42",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:18:00",
				"billableRate": 90000,
				"comment": "Comment 1 New"
			}`,
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"taskId": "01JHHK0FYG9MCRXNKPQJMPNV42",
				"taskName": "Task 5",
				"taskCompletedAt": null,
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:18:00",
				"durationSeconds": 3600,
				"billableRate": 90000,
				"billableAmount": 90000,
				"comment": "Comment 1 New",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z"
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			input:    "",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid timelog id",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					logger: zap.NewNop(),
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid projectId",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"projectId": "123",
				"taskId": "01JHHK0FYG9MCRXNKPQJMPNV42",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:15:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "projectId"`,
		},
		{
			name: "Should return error - invalid date",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHK0FYG9MCRXNKPQJMPNV42",
				"date": "123",
				"timeStart": "14:15:00",
				"timeEnd": "15:15:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "date"`,
		},
		{
			name: "Should return error - invalid timeStart",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHK0FYG9MCRXNKPQJMPNV42",
				"date": "2024-12-18",
				"timeStart": "123",
				"timeEnd": "15:15:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "timeStart"`,
		},
		{
			name: "Should return error - invalid timeEnd",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHK0FYG9MCRXNKPQJMPNV42",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "123",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "timeEnd"`,
		},
		{
			name: "Should return error - invalid billableRate",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHK0FYG9MCRXNKPQJMPNV42",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:15:00",
				"billableRate": -1,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "billableRate"`,
		},
		{
			name: "Should return error - invalid taskId",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "123",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:15:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "taskId"`,
		},
		{
			name: "Should return error - timelog not found",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					Update(mock.Anything, domain.Timelog{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						Date:         "2024-12-18",
						TimeStart:    "14:15:00",
						TimeEnd:      func() *string { s := "15:19:00"; return &s }(),
						BillableRate: 90000,
						Comment:      "Comment 1",
						TaskID:       func() *string { id := "01JHHK37S7G3JMP744ZXVSATGD"; return &id }(),
					}).
					Return(domain.Timelog{}, domain.ErrNotFound)

				return &TimelogHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					timelogService: timelogService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHK37S7G3JMP744ZXVSATGD",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:19:00",
				"billableRate": 90000,
				"comment": "Comment 1"
			}`,
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
		{
			name: "Should return error - project not found",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				sanitizer := mocks.NewSanitizerMock(t)
				sanitizer.EXPECT().
					Sanitize(mock.Anything).
					RunAndReturn(func(s string) string { return s })

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					Update(mock.Anything, domain.Timelog{
						ID:           "01JFEC2BRSGJKVG1FKNWSNV6HH",
						UserID:       testUserID,
						ProjectID:    "01JH74H6YKJNXH6H1CA0XF81DZ",
						Date:         "2024-12-18",
						TimeStart:    "14:15:00",
						TimeEnd:      func() *string { s := "15:20:00"; return &s }(),
						BillableRate: 90000,
						Comment:      "Comment 1",
						TaskID:       func() *string { id := "01JHHK3CCGHWBSAS77G04FDAA8"; return &id }(),
					}).
					Return(domain.Timelog{}, domain.ErrProjectNotFound)

				return &TimelogHandler{
					validate:       validator.New(),
					sanitizer:      sanitizer,
					timelogService: timelogService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"taskId": "01JHHK3CCGHWBSAS77G04FDAA8",
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:20:00",
				"billableRate": 90000,
				"comment": "Comment 1"
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

func TestTimelogHandler_Delete(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TimelogHandler
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					Delete(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH").
					Return(nil)

				return &TimelogHandler{
					validate:       validator.New(),
					timelogService: timelogService,
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			wantCode: http.StatusOK,
			wantBody: "OK",
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid timelog id",
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

func TestTimelogHandler_Stop(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TimelogHandler
		id       string
		input    string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					Stop(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH", "2024-12-18", "15:21:00").
					Return(domain.Timelog{
						ID:              "01JFEC2BRSGJKVG1FKNWSNV6HH",
						ClientID:        "01JH749ZAVM5MX8AMCJH7AKBX8",
						ClientName:      "Client 1",
						ProjectID:       "01JH74H6YKJNXH6H1CA0XF81DZ",
						ProjectName:     "Project 1",
						TaskID:          func() *string { id := "01JHHK182FRYPHRDRA00A72V0V"; return &id }(),
						TaskName:        func() *string { s := "Task 6"; return &s }(),
						Date:            "2024-12-18",
						TimeStart:       "14:15:00",
						TimeEnd:         func() *string { s := "15:21:00"; return &s }(),
						DurationSeconds: 3600,
						BillableRate:    90000,
						BillableAmount:  90000,
						Comment:         "Comment 1",
						CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:       func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TimelogHandler{
					validate:       validator.New(),
					timelogService: timelogService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"date": "2024-12-18",
				"timeEnd": "15:21:00"
			}`,
			wantCode: http.StatusOK,
			wantBody: `{
				"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
				"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
				"clientName": "Client 1",
				"projectId": "01JH74H6YKJNXH6H1CA0XF81DZ",
				"projectName": "Project 1",
				"taskId": "01JHHK182FRYPHRDRA00A72V0V",
				"taskName": "Task 6",
				"taskCompletedAt": null,
				"date": "2024-12-18",
				"timeStart": "14:15:00",
				"timeEnd": "15:21:00",
				"durationSeconds": 3600,
				"billableRate": 90000,
				"billableAmount": 90000,
				"comment": "Comment 1",
				"createdAt": "2024-12-18T14:15:00Z",
				"updatedAt": "2024-12-20T14:15:00Z"
			}`,
		},
		{
			name: "Should return error - invalid id",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					logger: zap.NewNop(),
				}
			},
			id:       `123`,
			input:    "",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid timelog id",
		},
		{
			name: "Should return error - invalid input",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					logger: zap.NewNop(),
				}
			},
			id:       "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input:    `123`,
			wantCode: http.StatusBadRequest,
			wantBody: "invalid input",
		},
		{
			name: "Should return error - invalid date",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"date": "123",
				"timeEnd": "15:15:00"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "date"`,
		},
		{
			name: "Should return error - invalid timeEnd",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					validate: validator.New(),
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"date": "2024-12-18",
				"timeEnd": "123"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: `invalid field "timeEnd"`,
		},
		{
			name: "Should return error - timelog not found",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					Stop(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH", "2024-12-18", "15:22:00").
					Return(domain.Timelog{}, domain.ErrNotFound)

				return &TimelogHandler{
					validate:       validator.New(),
					timelogService: timelogService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"date": "2024-12-18",
				"timeEnd": "15:22:00"
			}`,
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
		{
			name: "Should return error - current date cannot be less than the start date",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					Stop(mock.Anything, testUserID, "01JFEC2BRSGJKVG1FKNWSNV6HH", "2024-12-18", "15:23:00").
					Return(domain.Timelog{}, domain.ErrCurrentDateLessThanStartDate)

				return &TimelogHandler{
					validate:       validator.New(),
					timelogService: timelogService,
				}
			},
			id: "01JFEC2BRSGJKVG1FKNWSNV6HH",
			input: `{
				"date": "2024-12-18",
				"timeEnd": "15:23:00"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: "current date cannot be less than the start date",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Post("/:id", tt.srvFunc(t).Stop)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/"+tt.id, strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}

func TestTimelogHandler_LastN(t *testing.T) {
	const (
		testUserID = "01JFYQE72N34P9VZD84AEXC0PC"
	)

	tests := []struct {
		name     string
		srvFunc  func(t *testing.T) *TimelogHandler
		n        string
		wantCode int
		wantBody string
	}{
		{
			name: "Should return OK",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				timelogService := mocks.NewTimelogServiceMock(t)
				timelogService.EXPECT().
					FindLastN(mock.Anything, testUserID, 5).
					Return([]domain.Timelog{
						{
							ID:              "01JFEC2BRSGJKVG1FKNWSNV6HH",
							ClientID:        "01JH749ZAVM5MX8AMCJH7AKBX8",
							ClientName:      "Client 1",
							ProjectID:       "01JHFF8960RHPV35NE8GCZJYWZ",
							ProjectName:     "Project 1",
							TaskID:          func() *string { id := "01JHHK1FFYBGGF3QFWV8K1N4ZB"; return &id }(),
							TaskName:        func() *string { s := "Task 7"; return &s }(),
							TaskCompletedAt: func() *time.Time { t := time.Date(2024, 12, 20, 14, 15, 0, 0, time.UTC); return &t }(),
							Date:            "2024-12-18",
							TimeStart:       "14:15:00",
							DurationSeconds: 0,
							BillableRate:    90000,
							BillableAmount:  0,
							Comment:         "Comment 1",
							CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
						{
							ID:              "01JH741H5A4Y368PBHBTJFQ723",
							ClientID:        "01JH74AAT06A0P3AR2AGCFGRR5",
							ClientName:      "Client 2",
							ProjectID:       "01JHFF8G4S0RMW5DQNGZ0G3RH1",
							ProjectName:     "Project 2",
							Date:            "2024-12-18",
							TimeStart:       "14:15:00",
							TimeEnd:         func() *string { s := "15:24:00"; return &s }(),
							DurationSeconds: 3600,
							BillableRate:    90000,
							BillableAmount:  90000,
							CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						},
					}, nil)

				return &TimelogHandler{
					validate:       validator.New(),
					timelogService: timelogService,
				}
			},
			n:        "5",
			wantCode: http.StatusOK,
			wantBody: `[
				{
					"id": "01JFEC2BRSGJKVG1FKNWSNV6HH",
					"clientId": "01JH749ZAVM5MX8AMCJH7AKBX8",
					"clientName": "Client 1",
					"projectId": "01JHFF8960RHPV35NE8GCZJYWZ",
					"projectName": "Project 1",
					"taskId": "01JHHK1FFYBGGF3QFWV8K1N4ZB",
					"taskName": "Task 7",
					"taskCompletedAt": "2024-12-20T14:15:00Z",
					"date": "2024-12-18",
					"timeStart": "14:15:00",
					"timeEnd": null,
					"durationSeconds": 0,
					"billableRate": 90000,
					"billableAmount": 0,
					"comment": "Comment 1",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null
				},
				{
					"id": "01JH741H5A4Y368PBHBTJFQ723",
					"clientId": "01JH74AAT06A0P3AR2AGCFGRR5",
					"clientName": "Client 2",
					"projectId": "01JHFF8G4S0RMW5DQNGZ0G3RH1",
					"projectName": "Project 2",
					"taskId": null,
					"taskName": null,
					"taskCompletedAt": null,
					"date": "2024-12-18",
					"timeStart": "14:15:00",
					"timeEnd": "15:24:00",
					"durationSeconds": 3600,
					"billableRate": 90000,
					"billableAmount": 90000,
					"comment": "",
					"createdAt": "2024-12-18T14:15:00Z",
					"updatedAt": null
				}
			]`,
		},
		{
			name: "Should return error - invalid n",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					logger: zap.NewNop(),
				}
			},
			n:        "string",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid n parameter",
		},
		{
			name: "Should return error - n is zero",
			srvFunc: func(t *testing.T) *TimelogHandler {
				t.Helper()

				return &TimelogHandler{
					logger: zap.NewNop(),
				}
			},
			n:        "0",
			wantCode: http.StatusBadRequest,
			wantBody: "n must be greater than zero",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Use(func(c *fiber.Ctx) error {
				c.Locals(LocalsUserIDKey, testUserID)
				return c.Next()
			})

			app.Get("/:n", tt.srvFunc(t).LastN)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/"+tt.n, nil)

			testRunRequest(t, app, req, tt.wantCode, tt.wantBody)
		})
	}
}
