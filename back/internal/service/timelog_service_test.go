package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"capuchin/internal/domain"
	"capuchin/internal/service/mocks"
	"capuchin/internal/sqlite"
)

func TestTimelogService_Create(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTimelogID = "01JFVEYHENXWAA400DXQH47M2Y"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
		testClientID  = "01JFCN3VVQRQBGJ255PH9YMGY8"
		testTaskID    = "01JFEC2Q02JM8NJKGZSAE8AHP5"

		testTaskName = "Task Name"

		testTime1415 = "14:15:00"
	)

	type args struct {
		timelog domain.Timelog
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TimelogService
		args    args
		want    domain.Timelog
		wantErr error
	}{
		{
			name: "Should create timelog: TimeEnd is nil",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:   testTaskID,
						Name: testTaskName,
					}, nil)

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindRunning(ctx, testUserID).
					Return(domain.Timelog{}, domain.ErrNotFound)
				timelogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.Timelog{
						ID:           testTimelogID,
						UserID:       testUserID,
						ClientID:     testClientID,
						ClientName:   "Client Name",
						ProjectID:    testProjectID,
						ProjectName:  "Project Name",
						TaskID:       func() *string { id := testTaskID; return &id }(),
						TaskName:     func() *string { name := testTaskName; return &name }(),
						Date:         "2024-09-10",
						TimeStart:    "13:18:05",
						BillableRate: 90000,
						Comment:      "Comment",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTimelogs,
						RecordID:   testTimelogID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"date":"2024-09-10",` +
								`"timeStart":"13:18:05",` +
								`"timeEnd":null,` +
								`"durationSeconds":0,` +
								`"billableRate":90000,` +
								`"billableAmount":0,` +
								`"comment":"Comment",` +
								`"createdAt":"2024-12-18T14:15:00Z",` +
								`"updatedAt":null,` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":"Task Name",` +
								`"taskCompletedAt":null}`
							return &v
						}(),
					}).
					Return(nil)

				return &TimelogService{
					timelogRepo:  timelogRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					idFunc:       func() string { return testTimelogID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:       testUserID,
					ProjectID:    testProjectID,
					TaskID:       func() *string { id := testTaskID; return &id }(),
					Date:         "2024-09-10",
					TimeStart:    "13:18:05",
					BillableRate: 90000,
					Comment:      "Comment",
				},
			},
			want: domain.Timelog{
				ID:           testTimelogID,
				UserID:       testUserID,
				ClientID:     testClientID,
				ClientName:   "Client Name",
				ProjectID:    testProjectID,
				ProjectName:  "Project Name",
				TaskID:       func() *string { id := testTaskID; return &id }(),
				TaskName:     func() *string { name := testTaskName; return &name }(),
				Date:         "2024-09-10",
				TimeStart:    "13:18:05",
				BillableRate: 90000,
				Comment:      "Comment",
				CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "Should create timelog: TimeEnd is not nil",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:   testTaskID,
						Name: testTaskName,
					}, nil)

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindRunning(ctx, testUserID).
					Return(domain.Timelog{}, domain.ErrNotFound)
				timelogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.Timelog{
						ID:              testTimelogID,
						UserID:          testUserID,
						ClientID:        testClientID,
						ClientName:      "Client Name",
						ProjectID:       testProjectID,
						ProjectName:     "Project Name",
						TaskID:          func() *string { id := testTaskID; return &id }(),
						TaskName:        func() *string { name := testTaskName; return &name }(),
						Date:            "2024-09-10",
						TimeStart:       "13:18:05",
						TimeEnd:         func() *string { s := testTime1415; return &s }(),
						DurationSeconds: 3415,
						BillableRate:    90000,
						BillableAmount:  85375,
						Comment:         "Comment",
						CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTimelogs,
						RecordID:   testTimelogID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"date":"2024-09-10",` +
								`"timeStart":"13:18:05",` +
								`"timeEnd":"14:15:00",` +
								`"durationSeconds":3415,` +
								`"billableRate":90000,` +
								`"billableAmount":85375,` +
								`"comment":"Comment",` +
								`"createdAt":"2024-12-18T14:15:00Z",` +
								`"updatedAt":null,` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":"Task Name",` +
								`"taskCompletedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TimelogService{
					timelogRepo:  timelogRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					idFunc:       func() string { return testTimelogID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:       testUserID,
					ProjectID:    testProjectID,
					TaskID:       func() *string { id := testTaskID; return &id }(),
					Date:         "2024-09-10",
					TimeStart:    "13:18:05",
					TimeEnd:      func() *string { s := testTime1415; return &s }(),
					BillableRate: 90000,
					Comment:      "Comment",
				},
			},
			want: domain.Timelog{
				ID:              testTimelogID,
				UserID:          testUserID,
				ClientID:        testClientID,
				ClientName:      "Client Name",
				ProjectID:       testProjectID,
				ProjectName:     "Project Name",
				TaskID:          func() *string { id := testTaskID; return &id }(),
				TaskName:        func() *string { name := testTaskName; return &name }(),
				Date:            "2024-09-10",
				TimeStart:       "13:18:05",
				TimeEnd:         func() *string { s := testTime1415; return &s }(),
				DurationSeconds: 3415,
				BillableRate:    90000,
				BillableAmount:  85375,
				Comment:         "Comment",
				CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - project not found",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{}, domain.ErrNotFound)

				return &TimelogService{
					projectRepo: projectRepo,
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrProjectNotFound,
		},
		{
			name: "Should fail - project is archived",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:         testProjectID,
						ClientID:   testClientID,
						Name:       "Project Name",
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TimelogService{
					projectRepo: projectRepo,
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrProjectArchived,
		},
		{
			name: "Should fail - client not found",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{}, domain.ErrNotFound)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				return &TimelogService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "Should fail - client is archived",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:         testClientID,
						Name:       "Client Name",
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				return &TimelogService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrClientArchived,
		},
		{
			name: "Should fail - task not found",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{}, domain.ErrNotFound)

				return &TimelogService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					idFunc:      func() string { return testTimelogID },
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrTaskNotFound,
		},
		{
			name: "Should fail - task is archived",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:         testTaskID,
						Name:       testTaskName,
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TimelogService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					idFunc:      func() string { return testTimelogID },
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrTaskArchived,
		},
		{
			name: "Should fail - task is completed",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:          testTaskID,
						Name:        testTaskName,
						CompletedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TimelogService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					idFunc:      func() string { return testTimelogID },
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrTaskCompleted,
		},
		{
			name: "Should fail - stop running timelog - The timeEnd cannot be less than the timeStart",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:   testTaskID,
						Name: testTaskName,
					}, nil)

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindRunning(ctx, testUserID).
					Return(domain.Timelog{
						ID:        "01JFVGA9HFNWN96ACVK74733RT",
						Date:      "2024-09-10",
						TimeStart: "15:15:00",
					}, nil)

				return &TimelogService{
					timelogRepo: timelogRepo,
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					idFunc:      func() string { return testTimelogID },
					nowFunc:     func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:       testUserID,
					ProjectID:    testProjectID,
					TaskID:       func() *string { id := testTaskID; return &id }(),
					Date:         "2024-09-10",
					TimeStart:    "13:18:05",
					TimeEnd:      func() *string { s := testTime1415; return &s }(),
					BillableRate: 90000,
					Comment:      "Comment",
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrTimeEndLessThanTimeStart,
		},
		{
			name: "Should fail - The timeEnd cannot be less than the timeStart",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:   testTaskID,
						Name: testTaskName,
					}, nil)

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindRunning(ctx, testUserID).
					Return(domain.Timelog{}, domain.ErrNotFound)

				return &TimelogService{
					timelogRepo: timelogRepo,
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					idFunc:      func() string { return testTimelogID },
					nowFunc:     func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:       testUserID,
					ProjectID:    testProjectID,
					TaskID:       func() *string { id := testTaskID; return &id }(),
					Date:         "2024-09-10",
					TimeStart:    "13:18:05",
					TimeEnd:      func() *string { s := "10:15:00"; return &s }(),
					BillableRate: 90000,
					Comment:      "Comment",
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrTimeEndLessThanTimeStart,
		},
		{
			name: "Should fail - timelog already exists",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:   testTaskID,
						Name: testTaskName,
					}, nil)

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindRunning(ctx, testUserID).
					Return(domain.Timelog{}, domain.ErrNotFound)
				timelogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.Timelog{
						ID:           testTimelogID,
						UserID:       testUserID,
						ClientID:     testClientID,
						ClientName:   "Client Name",
						ProjectID:    testProjectID,
						ProjectName:  "Project Name",
						TaskID:       func() *string { id := testTaskID; return &id }(),
						TaskName:     func() *string { name := testTaskName; return &name }(),
						Date:         "2024-09-10",
						TimeStart:    "13:18:05",
						BillableRate: 90000,
						Comment:      "Comment",
						CreatedAt:    time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTimelogs,
						RecordID:   testTimelogID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"date":"2024-09-10",` +
								`"timeStart":"13:18:05",` +
								`"timeEnd":null,` +
								`"durationSeconds":0,` +
								`"billableRate":90000,` +
								`"billableAmount":0,` +
								`"comment":"Comment",` +
								`"createdAt":"2024-12-18T14:15:00Z",` +
								`"updatedAt":null,` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":"Task Name",` +
								`"taskCompletedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TimelogService{
					timelogRepo:  timelogRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, domain.ErrAlreadyExists),
					idFunc:       func() string { return testTimelogID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				timelog: domain.Timelog{
					UserID:       testUserID,
					ProjectID:    testProjectID,
					TaskID:       func() *string { id := testTaskID; return &id }(),
					Date:         "2024-09-10",
					TimeStart:    "13:18:05",
					BillableRate: 90000,
					Comment:      "Comment",
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Create(ctx, tt.args.timelog)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimelogService_Update(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTimelogID = "01JFVEYHENXWAA400DXQH47M2Y"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
		testClientID  = "01JFCN3VVQRQBGJ255PH9YMGY8"
		testTaskID    = "01JFEC2Q02JM8NJKGZSAE8AHP5"

		testTaskName = "Task Name"

		testTime1415 = "14:15:00"

		testNewDate      = "2024-12-19"
		testNewTimeStart = "15:16:00"
		testNewTimeEnd   = "17:16:00"
	)

	type args struct {
		timelog domain.Timelog
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TimelogService
		args    args
		want    domain.Timelog
		wantErr error
	}{
		{
			name: "Should update timelog",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:   testTaskID,
						Name: testTaskName,
					}, nil)

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						ID:           testTimelogID,
						UserID:       testUserID,
						ClientID:     testClientID,
						ClientName:   "Client Name",
						ProjectID:    testProjectID,
						ProjectName:  "Project Name",
						Date:         "2024-12-18",
						TimeStart:    testTime1415,
						TimeEnd:      func() *string { t := testTime1415; return &t }(),
						BillableRate: 1000,
						Comment:      "Old Comment",
					}, nil)
				timelogRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Timelog{
						ID:              testTimelogID,
						UserID:          testUserID,
						ClientID:        testClientID,
						ClientName:      "Client Name",
						ProjectID:       testProjectID,
						ProjectName:     "Project Name",
						TaskID:          func() *string { id := testTaskID; return &id }(),
						TaskName:        func() *string { name := testTaskName; return &name }(),
						Date:            testNewDate,
						TimeStart:       testNewTimeStart,
						TimeEnd:         func() *string { t := testNewTimeEnd; return &t }(),
						DurationSeconds: 7200,
						BillableRate:    2000,
						BillableAmount:  4000,
						Comment:         "New Comment",
						UpdatedAt:       func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTimelogs,
						RecordID:   testTimelogID,
						OldValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"date":"2024-12-18",` +
								`"timeStart":"14:15:00",` +
								`"timeEnd":"14:15:00",` +
								`"durationSeconds":0,` +
								`"billableRate":1000,` +
								`"billableAmount":0,` +
								`"comment":"Old Comment",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"taskId":null,` +
								`"taskName":null,` +
								`"taskCompletedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"date":"2024-12-19",` +
								`"timeStart":"15:16:00",` +
								`"timeEnd":"17:16:00",` +
								`"durationSeconds":7200,` +
								`"billableRate":2000,` +
								`"billableAmount":4000,` +
								`"comment":"New Comment",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":"Task Name",` +
								`"taskCompletedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TimelogService{
					timelogRepo:  timelogRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				timelog: domain.Timelog{
					ID:           testTimelogID,
					UserID:       testUserID,
					ProjectID:    testProjectID,
					TaskID:       func() *string { id := testTaskID; return &id }(),
					Date:         testNewDate,
					TimeStart:    testNewTimeStart,
					TimeEnd:      func() *string { t := testNewTimeEnd; return &t }(),
					BillableRate: 2000,
					Comment:      "New Comment",
				},
			},
			want: domain.Timelog{
				ID:              testTimelogID,
				UserID:          testUserID,
				ClientID:        testClientID,
				ClientName:      "Client Name",
				ProjectID:       testProjectID,
				ProjectName:     "Project Name",
				TaskID:          func() *string { id := testTaskID; return &id }(),
				TaskName:        func() *string { name := testTaskName; return &name }(),
				Date:            testNewDate,
				TimeStart:       testNewTimeStart,
				TimeEnd:         func() *string { t := testNewTimeEnd; return &t }(),
				DurationSeconds: 7200,
				BillableRate:    2000,
				BillableAmount:  4000,
				Comment:         "New Comment",
				UpdatedAt:       func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - timelog not found",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{}, domain.ErrNotFound)

				return &TimelogService{
					timelogRepo: timelogRepo,
				}
			},
			args: args{
				timelog: domain.Timelog{
					ID:      testTimelogID,
					UserID:  testUserID,
					TimeEnd: func() *string { t := testNewTimeEnd; return &t }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "Should fail - cannot edit running timelog",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				return &TimelogService{}
			},
			args: args{
				timelog: domain.Timelog{
					ID:     testTimelogID,
					UserID: testUserID,
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.NewFailedPreconditionError("cannot edit running timelog"),
		},
		{
			name: "Should fail - project not found",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						UserID:    testUserID,
						ProjectID: testProjectID,
						Date:      "2024-12-18",
						TimeStart: testTime1415,
						TimeEnd:   func() *string { t := testTime1415; return &t }(),
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{}, domain.ErrNotFound)

				return &TimelogService{
					timelogRepo: timelogRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				timelog: domain.Timelog{
					ID:        testTimelogID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					Date:      testNewDate,
					TimeStart: testNewTimeStart,
					TimeEnd:   func() *string { t := testNewTimeEnd; return &t }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrProjectNotFound,
		},
		{
			name: "Should fail - project is archived",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						UserID:    testUserID,
						ProjectID: testProjectID,
						Date:      "2024-12-18",
						TimeStart: testTime1415,
						TimeEnd:   func() *string { t := testTime1415; return &t }(),
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:         testProjectID,
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TimelogService{
					timelogRepo: timelogRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				timelog: domain.Timelog{
					ID:        testTimelogID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					Date:      testNewDate,
					TimeStart: testNewTimeStart,
					TimeEnd:   func() *string { t := testNewTimeEnd; return &t }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrProjectArchived,
		},
		{
			name: "Should fail - client not found",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						UserID:    testUserID,
						ProjectID: testProjectID,
						Date:      "2024-12-18",
						TimeStart: testTime1415,
						TimeEnd:   func() *string { t := testTime1415; return &t }(),
					}, nil)

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{}, domain.ErrNotFound)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
					}, nil)

				return &TimelogService{
					timelogRepo: timelogRepo,
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				timelog: domain.Timelog{
					ID:        testTimelogID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					Date:      testNewDate,
					TimeStart: testNewTimeStart,
					TimeEnd:   func() *string { t := testNewTimeEnd; return &t }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "Should fail - client is archived",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						UserID:    testUserID,
						ProjectID: testProjectID,
						Date:      "2024-12-18",
						TimeStart: testTime1415,
						TimeEnd:   func() *string { t := testTime1415; return &t }(),
					}, nil)

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:         testClientID,
						Name:       "Client Name",
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
					}, nil)

				return &TimelogService{
					timelogRepo: timelogRepo,
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				timelog: domain.Timelog{
					ID:        testTimelogID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					Date:      testNewDate,
					TimeStart: testNewTimeStart,
					TimeEnd:   func() *string { t := testNewTimeEnd; return &t }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrClientArchived,
		},
		{
			name: "Should fail - task not found",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						UserID:    testUserID,
						ProjectID: testProjectID,
						Date:      "2024-12-18",
						TimeStart: testTime1415,
						TimeEnd:   func() *string { t := testTime1415; return &t }(),
					}, nil)

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{}, domain.ErrNotFound)

				return &TimelogService{
					timelogRepo: timelogRepo,
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					nowFunc:     func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				timelog: domain.Timelog{
					ID:        testTimelogID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
					Date:      testNewDate,
					TimeStart: testNewTimeStart,
					TimeEnd:   func() *string { t := testNewTimeEnd; return &t }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrTaskNotFound,
		},
		{
			name: "Should fail - task is archived",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						UserID:    testUserID,
						ProjectID: testProjectID,
						Date:      "2024-12-18",
						TimeStart: testTime1415,
						TimeEnd:   func() *string { t := testTime1415; return &t }(),
					}, nil)

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:         testTaskID,
						Name:       testTaskName,
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TimelogService{
					timelogRepo: timelogRepo,
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					nowFunc:     func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				timelog: domain.Timelog{
					ID:        testTimelogID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
					Date:      testNewDate,
					TimeStart: testNewTimeStart,
					TimeEnd:   func() *string { t := testNewTimeEnd; return &t }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrTaskArchived,
		},
		{
			name: "Should fail - task is completed",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						UserID:    testUserID,
						ProjectID: testProjectID,
						Date:      "2024-12-18",
						TimeStart: testTime1415,
						TimeEnd:   func() *string { t := testTime1415; return &t }(),
					}, nil)

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:   testClientID,
						Name: "Client Name",
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:       testProjectID,
						ClientID: testClientID,
						Name:     "Project Name",
					}, nil)

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:          testTaskID,
						Name:        testTaskName,
						CompletedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TimelogService{
					timelogRepo: timelogRepo,
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					nowFunc:     func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				timelog: domain.Timelog{
					ID:        testTimelogID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
					Date:      testNewDate,
					TimeStart: testNewTimeStart,
					TimeEnd:   func() *string { t := testNewTimeEnd; return &t }(),
				},
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrTaskCompleted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Update(ctx, tt.args.timelog)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimelogService_Delete(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTimelogID = "01JFVEYHENXWAA400DXQH47M2Y"
	)

	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TimelogService
	}{
		{
			name: "Should delete timelog",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						ID:     testTimelogID,
						UserID: testUserID,
					}, nil)
				timelogRepo.EXPECT().
					Delete(ctx, mock.Anything, testUserID, testTimelogID).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeDelete,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTimelogs,
						RecordID:   testTimelogID,
						OldValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"date":"",` +
								`"timeStart":"",` +
								`"timeEnd":null,` +
								`"durationSeconds":0,` +
								`"billableRate":0,` +
								`"billableAmount":0,` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"taskId":null,` +
								`"taskName":null,` +
								`"taskCompletedAt":null}`

							return &v
						}(),
						NewValue: nil,
					}).
					Return(nil)

				return &TimelogService{
					timelogRepo:  timelogRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
		},
		{
			name: "Should fail - timelog not found",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{}, domain.ErrNotFound)

				return &TimelogService{
					timelogRepo: timelogRepo,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)
			err := srv.Delete(ctx, testUserID, testTimelogID)
			assert.NoError(t, err)
		})
	}
}

func TestTimelogService_Stop(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTimelogID = "01JFVEYHENXWAA400DXQH47M2Y"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
		testClientID  = "01JFCN3VVQRQBGJ255PH9YMGY8"
		testTaskID    = "01JFEC2Q02JM8NJKGZSAE8AHP5"

		testTaskName = "Task Name"

		testTime1415 = "14:15:00"
	)

	type args struct {
		userID    string
		timelogID string
		date      string
		timeEnd   string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TimelogService
		args    args
		want    domain.Timelog
		wantErr error
	}{
		{
			name: "Should stop timelog",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						ID:              testTimelogID,
						UserID:          testUserID,
						ProjectID:       testProjectID,
						ProjectName:     "Project Name",
						ClientID:        testClientID,
						ClientName:      "Client Name",
						Date:            "2024-09-10",
						TimeStart:       "13:16:18",
						BillableRate:    1000,
						Comment:         "Comment",
						CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						TaskID:          func() *string { s := testTaskID; return &s }(),
						TaskName:        func() *string { s := testTaskName; return &s }(),
						TaskCompletedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)
				timelogRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Timelog{
						ID:              testTimelogID,
						UserID:          testUserID,
						ProjectID:       testProjectID,
						ProjectName:     "Project Name",
						ClientID:        testClientID,
						ClientName:      "Client Name",
						Date:            "2024-09-10",
						TimeStart:       "13:16:18",
						TimeEnd:         func() *string { s := testTime1415; return &s }(),
						DurationSeconds: 3522,
						BillableRate:    1000,
						BillableAmount:  978,
						Comment:         "Comment",
						CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						UpdatedAt:       func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
						TaskID:          func() *string { s := testTaskID; return &s }(),
						TaskName:        func() *string { s := testTaskName; return &s }(),
						TaskCompletedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTimelogs,
						RecordID:   testTimelogID,
						OldValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"date":"2024-09-10",` +
								`"timeStart":"13:16:18",` +
								`"timeEnd":null,` +
								`"durationSeconds":0,` +
								`"billableRate":1000,` +
								`"billableAmount":0,` +
								`"comment":"Comment",` +
								`"createdAt":"2024-12-18T14:15:00Z",` +
								`"updatedAt":null,` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":"Task Name",` +
								`"taskCompletedAt":"2024-12-18T14:15:00Z"}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"date":"2024-09-10",` +
								`"timeStart":"13:16:18",` +
								`"timeEnd":"14:15:00",` +
								`"durationSeconds":3522,` +
								`"billableRate":1000,` +
								`"billableAmount":978,` +
								`"comment":"Comment",` +
								`"createdAt":"2024-12-18T14:15:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":"Task Name",` +
								`"taskCompletedAt":"2024-12-18T14:15:00Z"}`

							return &v
						}(),
					}).
					Return(nil)

				return &TimelogService{
					timelogRepo:  timelogRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID:    testUserID,
				timelogID: testTimelogID,
				date:      "2024-09-10",
				timeEnd:   testTime1415,
			},
			want: domain.Timelog{
				ID:              testTimelogID,
				UserID:          testUserID,
				ProjectID:       testProjectID,
				ProjectName:     "Project Name",
				ClientID:        testClientID,
				ClientName:      "Client Name",
				Date:            "2024-09-10",
				TimeStart:       "13:16:18",
				TimeEnd:         func() *string { s := testTime1415; return &s }(),
				DurationSeconds: 3522,
				BillableRate:    1000,
				BillableAmount:  978,
				Comment:         "Comment",
				CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
				UpdatedAt:       func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
				TaskID:          func() *string { s := testTaskID; return &s }(),
				TaskName:        func() *string { s := testTaskName; return &s }(),
				TaskCompletedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should success - timelog.TimeEnd is not nil",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						ID:        testTimelogID,
						Date:      "2024-09-10",
						TimeStart: "13:16:18",
						TimeEnd:   func() *string { s := testTime1415; return &s }(),
					}, nil)

				return &TimelogService{
					timelogRepo: timelogRepo,
				}
			},
			args: args{
				userID:    testUserID,
				timelogID: testTimelogID,
				date:      "2024-09-10",
				timeEnd:   testTime1415,
			},
			want: domain.Timelog{
				ID:        testTimelogID,
				Date:      "2024-09-10",
				TimeStart: "13:16:18",
				TimeEnd:   func() *string { s := testTime1415; return &s }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - timelog not found",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{}, domain.ErrNotFound)

				return &TimelogService{
					timelogRepo: timelogRepo,
				}
			},
			args: args{
				userID:    testUserID,
				timelogID: testTimelogID,
				date:      "2024-09-10",
				timeEnd:   testTime1415,
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "Should fail - The timeEnd cannot be less than the timeStart",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindByID(ctx, testUserID, testTimelogID).
					Return(domain.Timelog{
						ID:        testTimelogID,
						Date:      "2024-09-10",
						TimeStart: "13:16:18",
					}, nil)

				return &TimelogService{
					timelogRepo: timelogRepo,
					nowFunc:     func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID:    testUserID,
				timelogID: testTimelogID,
				date:      "2024-09-10",
				timeEnd:   "10:15:00",
			},
			want:    domain.Timelog{},
			wantErr: domain.ErrTimeEndLessThanTimeStart,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Stop(ctx, tt.args.userID, tt.args.timelogID, tt.args.date, tt.args.timeEnd)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimelogService_FindLastN(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID = "01JFCGS2YRPKT24B3RBZ8PVPBP"

		testProjectID1 = "01JFVM0MQEBECTTVMZ7EVSTYXS"
		testProjectID2 = "01JFVM0R5P5MR9Z3NQ7XDSC3D9"

		testTaskID1 = "01JFVKZET0KAKXC9JC6FAP17J3"
		testTaskID2 = "01JFVKZJEWB3K0DZWMVTZAMFDX"
		testTaskID3 = "01JKVYWR89D90T1NHYTWBB1FHQ"

		testComment1 = "Comment 1"
		testComment2 = "Comment 2"
	)

	type args struct {
		userID string
		n      int
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TimelogService
		args    args
		want    []domain.Timelog
		wantErr error
	}{
		{
			name: "Should find last n timelogs",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					FindLastN(ctx, testUserID, 20).
					Return([]domain.Timelog{
						{
							ID:        "tl_1",
							ProjectID: testProjectID1,
							TaskID:    func() *string { s := testTaskID1; return &s }(),
							Comment:   testComment1,
						},
						{
							ID:        "tl_2",
							ProjectID: testProjectID2,
							TaskID:    func() *string { s := testTaskID2; return &s }(),
							Comment:   testComment2,
						},
						{
							ID:        "tl_3",
							ProjectID: testProjectID1,
							TaskID:    func() *string { s := testTaskID1; return &s }(),
							Comment:   testComment1,
						},
						{
							ID:        "tl_4",
							ProjectID: testProjectID2,
							TaskID:    func() *string { s := testTaskID2; return &s }(),
							Comment:   testComment2,
						},
						{
							ID:        "tl_5",
							ProjectID: testProjectID1,
							TaskID:    func() *string { s := testTaskID1; return &s }(),
						},
						{
							ID:        "tl_6",
							ProjectID: testProjectID1,
						},
						{
							ID:              "tl_7",
							ProjectID:       testProjectID1,
							TaskID:          func() *string { s := testTaskID3; return &s }(),
							TaskCompletedAt: func() *time.Time { t := time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC); return &t }(),
						},
					}, nil)

				return &TimelogService{
					timelogRepo: timelogRepo,
				}
			},
			args: args{
				userID: testUserID,
				n:      2,
			},
			want: []domain.Timelog{
				{
					ID:        "tl_1",
					ProjectID: testProjectID1,
					TaskID:    func() *string { s := testTaskID1; return &s }(),
					Comment:   testComment1,
				},
				{
					ID:        "tl_2",
					ProjectID: testProjectID2,
					TaskID:    func() *string { s := testTaskID2; return &s }(),
					Comment:   testComment2,
				},
				{
					ID:        "tl_5",
					ProjectID: testProjectID1,
					TaskID:    func() *string { s := testTaskID1; return &s }(),
				},
				{
					ID:        "tl_6",
					ProjectID: testProjectID1,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.FindLastN(ctx, tt.args.userID, tt.args.n)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimelogService_stopRecord(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTimelogID = "01JFVEYHENXWAA400DXQH47M2Y"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
		testTaskID    = "01JFEC2Q02JM8NJKGZSAE8AHP5"

		testTime141520 = "14:15:20"
	)

	type args struct {
		tl      domain.Timelog
		date    string
		timeEnd string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TimelogService
		args    args
		want    domain.Timelog
		wantErr error
	}{
		{
			name: "Should stop record - today",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Timelog{
						ID:              testTimelogID,
						UserID:          testUserID,
						ProjectID:       testProjectID,
						TaskID:          func() *string { s := testTaskID; return &s }(),
						Date:            "2024-09-10",
						TimeStart:       "13:16:18",
						TimeEnd:         func() *string { s := testTime141520; return &s }(),
						DurationSeconds: 3542,
						BillableRate:    90000,
						BillableAmount:  88550,
						Comment:         "Comment",
						UpdatedAt:       func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTimelogs,
						RecordID:   testTimelogID,
						OldValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"date":"2024-09-10",` +
								`"timeStart":"13:16:18",` +
								`"timeEnd":null,` +
								`"durationSeconds":0,` +
								`"billableRate":90000,` +
								`"billableAmount":0,` +
								`"comment":"Comment",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":null,` +
								`"taskCompletedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFVEYHENXWAA400DXQH47M2Y",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"date":"2024-09-10",` +
								`"timeStart":"13:16:18",` +
								`"timeEnd":"14:15:20",` +
								`"durationSeconds":3542,` +
								`"billableRate":90000,` +
								`"billableAmount":88550,` +
								`"comment":"Comment",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":null,` +
								`"taskCompletedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TimelogService{
					timelogRepo:  timelogRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				tl: domain.Timelog{
					ID:           testTimelogID,
					UserID:       testUserID,
					ProjectID:    testProjectID,
					TaskID:       func() *string { s := testTaskID; return &s }(),
					Date:         "2024-09-10",
					TimeStart:    "13:16:18",
					BillableRate: 90000,
					Comment:      "Comment",
				},
				date:    "2024-09-10",
				timeEnd: testTime141520,
			},
			want: domain.Timelog{
				ID:              testTimelogID,
				UserID:          testUserID,
				ProjectID:       testProjectID,
				TaskID:          func() *string { s := testTaskID; return &s }(),
				Date:            "2024-09-10",
				TimeStart:       "13:16:18",
				TimeEnd:         func() *string { s := testTime141520; return &s }(),
				DurationSeconds: 3542,
				BillableRate:    90000,
				BillableAmount:  88550,
				Comment:         "Comment",
				UpdatedAt:       func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should stop record - yesterday",
			srvFunc: func(t *testing.T) *TimelogService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					Create(ctx, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, tx sqlite.SQLExecutor, tl domain.Timelog) error {
						switch tl.Date {
						case "2024-09-11":
							assert.Equal(t, domain.Timelog{
								ID:              "tl_1",
								UserID:          testUserID,
								ProjectID:       testProjectID,
								TaskID:          func() *string { s := testTaskID; return &s }(),
								Date:            "2024-09-11",
								TimeStart:       "00:00:00",
								TimeEnd:         func() *string { s := "23:59:59"; return &s }(),
								DurationSeconds: 86399,
								BillableRate:    90000,
								BillableAmount:  2159975,
								Comment:         "Comment",
								CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
							}, tl)
						case "2024-09-12":
							assert.Equal(t, domain.Timelog{
								ID:              "tl_1",
								UserID:          testUserID,
								ProjectID:       testProjectID,
								TaskID:          func() *string { s := testTaskID; return &s }(),
								Date:            "2024-09-12",
								TimeStart:       "00:00:00",
								TimeEnd:         func() *string { s := testTime141520; return &s }(),
								DurationSeconds: 51320,
								BillableRate:    90000,
								BillableAmount:  1283000,
								Comment:         "Comment",
								CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
							}, tl)
						default:
							assert.Equal(t, domain.Timelog{}, tl)
						}

						return nil
					})
				timelogRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Timelog{
						ID:              "tl_0",
						UserID:          testUserID,
						ProjectID:       testProjectID,
						TaskID:          func() *string { s := testTaskID; return &s }(),
						Date:            "2024-09-10",
						TimeStart:       "13:16:18",
						TimeEnd:         func() *string { s := "23:59:59"; return &s }(),
						DurationSeconds: 38621,
						BillableRate:    90000,
						BillableAmount:  965524,
						Comment:         "Comment",
						UpdatedAt:       func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, tx sqlite.SQLExecutor, auditLog domain.AuditLog) error {
						if auditLog.RecordID == "tl_0" {
							assert.Equal(t, domain.AuditLog{
								ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
								ActionType: domain.AuditLogActionTypeUpdate,
								UserID:     testUserID,
								TableName:  domain.AuditLogTableNameTimelogs,
								RecordID:   "tl_0",
								OldValue: func() *string {
									v := `{"id":"tl_0",` +
										`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
										`"projectName":"",` +
										`"clientId":"",` +
										`"clientName":"",` +
										`"date":"2024-09-10",` +
										`"timeStart":"13:16:18",` +
										`"timeEnd":null,` +
										`"durationSeconds":0,` +
										`"billableRate":90000,` +
										`"billableAmount":0,` +
										`"comment":"Comment",` +
										`"createdAt":"0001-01-01T00:00:00Z",` +
										`"updatedAt":null,` +
										`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
										`"taskName":null,` +
										`"taskCompletedAt":null}`

									return &v
								}(),
								NewValue: func() *string {
									v := `{"id":"tl_0",` +
										`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
										`"projectName":"",` +
										`"clientId":"",` +
										`"clientName":"",` +
										`"date":"2024-09-10",` +
										`"timeStart":"13:16:18",` +
										`"timeEnd":"23:59:59",` +
										`"durationSeconds":38621,` +
										`"billableRate":90000,` +
										`"billableAmount":965524,` +
										`"comment":"Comment",` +
										`"createdAt":"0001-01-01T00:00:00Z",` +
										`"updatedAt":"2024-12-18T14:15:00Z",` +
										`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
										`"taskName":null,` +
										`"taskCompletedAt":null}`

									return &v
								}(),
							}, auditLog)

							return nil
						}

						var newTimelog domain.Timelog
						if err := json.Unmarshal([]byte(*auditLog.NewValue), &newTimelog); err != nil {
							t.Fatalf("failed to unmarshal NewValue: %v", err)
						}

						switch newTimelog.Date {
						case "2024-09-11":
							assert.Equal(t, domain.AuditLog{
								ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
								ActionType: domain.AuditLogActionTypeCreate,
								UserID:     testUserID,
								TableName:  domain.AuditLogTableNameTimelogs,
								RecordID:   "tl_1",
								OldValue:   nil,
								NewValue: func() *string {
									v := `{"id":"tl_1",` +
										`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
										`"projectName":"",` +
										`"clientId":"",` +
										`"clientName":"",` +
										`"date":"2024-09-11",` +
										`"timeStart":"00:00:00",` +
										`"timeEnd":"23:59:59",` +
										`"durationSeconds":86399,` +
										`"billableRate":90000,` +
										`"billableAmount":2159975,` +
										`"comment":"Comment",` +
										`"createdAt":"2024-12-18T14:15:00Z",` +
										`"updatedAt":null,` +
										`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
										`"taskName":null,` +
										`"taskCompletedAt":null}`

									return &v
								}(),
							}, auditLog)
						case "2024-09-12":
							assert.Equal(t, domain.AuditLog{
								ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
								ActionType: domain.AuditLogActionTypeCreate,
								UserID:     testUserID,
								TableName:  domain.AuditLogTableNameTimelogs,
								RecordID:   "tl_1",
								OldValue:   nil,
								NewValue: func() *string {
									v := `{"id":"tl_1",` +
										`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
										`"projectName":"",` +
										`"clientId":"",` +
										`"clientName":"",` +
										`"date":"2024-09-12",` +
										`"timeStart":"00:00:00",` +
										`"timeEnd":"14:15:20",` +
										`"durationSeconds":51320,` +
										`"billableRate":90000,` +
										`"billableAmount":1283000,` +
										`"comment":"Comment",` +
										`"createdAt":"2024-12-18T14:15:00Z",` +
										`"updatedAt":null,` +
										`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
										`"taskName":null,` +
										`"taskCompletedAt":null}`

									return &v
								}(),
							}, auditLog)
						default:
							assert.Equal(t, domain.AuditLog{}, auditLog)
						}

						return nil
					})

				return &TimelogService{
					timelogRepo:  timelogRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					idFunc:       func() string { return "tl_1" },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				tl: domain.Timelog{
					ID:           "tl_0",
					UserID:       testUserID,
					ProjectID:    testProjectID,
					TaskID:       func() *string { s := testTaskID; return &s }(),
					Date:         "2024-09-10",
					TimeStart:    "13:16:18",
					BillableRate: 90000,
					Comment:      "Comment",
				},
				date:    "2024-09-12",
				timeEnd: testTime141520,
			},
			want: domain.Timelog{
				ID:              "tl_1",
				UserID:          testUserID,
				ProjectID:       testProjectID,
				TaskID:          func() *string { s := testTaskID; return &s }(),
				Date:            "2024-09-12",
				TimeStart:       "00:00:00",
				TimeEnd:         func() *string { s := testTime141520; return &s }(),
				DurationSeconds: 51320,
				BillableRate:    90000,
				BillableAmount:  1283000,
				Comment:         "Comment",
				CreatedAt:       time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.stopRecord(ctx, tt.args.tl, tt.args.date, tt.args.timeEnd)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
