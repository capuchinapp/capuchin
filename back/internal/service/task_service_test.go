package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"capuchin/internal/domain"
	"capuchin/internal/service/mocks"
)

func TestTaskService_Create(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTaskID    = "01JFEC2Q02JM8NJKGZSAE8AHP5"
		testClientID  = "01JFCN3VVQRQBGJ255PH9YMGY8"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
	)

	type args struct {
		task domain.Task
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TaskService
		args    args
		want    domain.Task
		wantErr error
	}{
		{
			name: "Should create task",
			srvFunc: func(t *testing.T) *TaskService {
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
					Create(ctx, mock.Anything, domain.Task{
						ID:          testTaskID,
						UserID:      testUserID,
						ClientID:    testClientID,
						ClientName:  "Client Name",
						ProjectID:   testProjectID,
						ProjectName: "Project Name",
						CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTasks,
						RecordID:   testTaskID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"2024-12-18T14:15:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null,` +
								`"completedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TaskService{
					taskRepo:     taskRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					idFunc:       func() string { return testTaskID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				task: domain.Task{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want: domain.Task{
				ID:          testTaskID,
				UserID:      testUserID,
				ClientID:    testClientID,
				ClientName:  "Client Name",
				ProjectID:   testProjectID,
				ProjectName: "Project Name",
				CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - project not found",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{}, domain.ErrNotFound)

				return &TaskService{
					projectRepo: projectRepo,
				}
			},
			args: args{
				task: domain.Task{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrProjectNotFound,
		},
		{
			name: "Should fail - project is archived",
			srvFunc: func(t *testing.T) *TaskService {
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

				return &TaskService{
					projectRepo: projectRepo,
				}
			},
			args: args{
				task: domain.Task{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrProjectArchived,
		},
		{
			name: "Should fail - client not found",
			srvFunc: func(t *testing.T) *TaskService {
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

				return &TaskService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				task: domain.Task{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "Should fail - client is archived",
			srvFunc: func(t *testing.T) *TaskService {
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

				return &TaskService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				task: domain.Task{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrClientArchived,
		},
		{
			name: "Should fail - task already exists",
			srvFunc: func(t *testing.T) *TaskService {
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
					Create(ctx, mock.Anything, domain.Task{
						ID:          testTaskID,
						UserID:      testUserID,
						ClientID:    testClientID,
						ClientName:  "Client Name",
						ProjectID:   testProjectID,
						ProjectName: "Project Name",
						CreatedAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTasks,
						RecordID:   testTaskID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"2024-12-18T14:15:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null,` +
								`"completedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TaskService{
					taskRepo:     taskRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, domain.ErrAlreadyExists),
					idFunc:       func() string { return testTaskID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				task: domain.Task{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Create(ctx, tt.args.task)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskService_Update(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTaskID    = "01JFEC2Q02JM8NJKGZSAE8AHP5"
		testClientID  = "01JFCN3VVQRQBGJ255PH9YMGY8"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
	)

	type args struct {
		task domain.Task
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TaskService
		args    args
		want    domain.Task
		wantErr error
	}{
		{
			name: "Should update task",
			srvFunc: func(t *testing.T) *TaskService {
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
						UserID:      testUserID,
						Name:        "Old Name",
						ClientID:    testClientID,
						ClientName:  "Client Name",
						ProjectID:   testProjectID,
						ProjectName: "Project Name",
						Comment:     "Old Comment",
					}, nil)
				taskRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Task{
						ID:          testTaskID,
						UserID:      testUserID,
						Name:        "New Name",
						ClientID:    testClientID,
						ClientName:  "Client Name",
						ProjectID:   testProjectID,
						ProjectName: "Project Name",
						Comment:     "New Comment",
						UpdatedAt:   func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTasks,
						RecordID:   testTaskID,
						OldValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"name":"Old Name",` +
								`"comment":"Old Comment",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null,` +
								`"completedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"name":"New Name",` +
								`"comment":"New Comment",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"archivedAt":null,` +
								`"completedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TaskService{
					taskRepo:     taskRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				task: domain.Task{
					ID:        testTaskID,
					UserID:    testUserID,
					Name:      "New Name",
					ProjectID: testProjectID,
					Comment:   "New Comment",
				},
			},
			want: domain.Task{
				ID:          testTaskID,
				UserID:      testUserID,
				ClientID:    testClientID,
				ClientName:  "Client Name",
				ProjectID:   testProjectID,
				ProjectName: "Project Name",
				Name:        "New Name",
				Comment:     "New Comment",
				UpdatedAt:   func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - task not found",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskService{
					taskRepo: taskRepo,
				}
			},
			args: args{
				task: domain.Task{
					ID:     testTaskID,
					UserID: testUserID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "Should fail - project not found",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						UserID:    testUserID,
						ProjectID: testProjectID,
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{}, domain.ErrNotFound)

				return &TaskService{
					taskRepo:    taskRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				task: domain.Task{
					ID:        testTaskID,
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrProjectNotFound,
		},
		{
			name: "Should fail - project is archived",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						UserID:    testUserID,
						ProjectID: testProjectID,
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:         testProjectID,
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TaskService{
					taskRepo:    taskRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				task: domain.Task{
					ID:        testTaskID,
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrProjectArchived,
		},
		{
			name: "Should fail - client not found",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						UserID:    testUserID,
						ProjectID: testProjectID,
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

				return &TaskService{
					taskRepo:    taskRepo,
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				task: domain.Task{
					ID:        testTaskID,
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "Should fail - client is archived",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						UserID:    testUserID,
						ProjectID: testProjectID,
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

				return &TaskService{
					taskRepo:    taskRepo,
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				task: domain.Task{
					ID:        testTaskID,
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Task{},
			wantErr: domain.ErrClientArchived,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Update(ctx, tt.args.task)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskService_Archive(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTaskID = "01JFEC2Q02JM8NJKGZSAE8AHP5"
	)

	type args struct {
		userID string
		taskID string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TaskService
		args    args
		want    domain.Task
		wantErr error
	}{
		{
			name: "Should archive task",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:     testTaskID,
						UserID: testUserID,
					}, nil)
				taskRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Task{
						ID:         testTaskID,
						UserID:     testUserID,
						UpdatedAt:  func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTasks,
						RecordID:   testTaskID,
						OldValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null,` +
								`"completedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"archivedAt":"2024-12-18T14:15:00Z",` +
								`"completedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TaskService{
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want: domain.Task{
				ID:         testTaskID,
				UserID:     testUserID,
				UpdatedAt:  func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
				ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - task not found",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskService{
					taskRepo: taskRepo,
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want:    domain.Task{},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "Should fail - task is completed",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:          testTaskID,
						UserID:      testUserID,
						CompletedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TaskService{
					taskRepo: taskRepo,
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want:    domain.Task{},
			wantErr: domain.ErrImpossibleAction,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Archive(ctx, tt.args.userID, tt.args.taskID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskService_Unarchive(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTaskID = "01JFEC2Q02JM8NJKGZSAE8AHP5"
	)

	type args struct {
		userID string
		taskID string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TaskService
		args    args
		want    domain.Task
		wantErr error
	}{
		{
			name: "Should unarchive task",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:         testTaskID,
						UserID:     testUserID,
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)
				taskRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Task{
						ID:        testTaskID,
						UserID:    testUserID,
						UpdatedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTasks,
						RecordID:   testTaskID,
						OldValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":"2024-12-18T14:15:00Z",` +
								`"completedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"archivedAt":null,` +
								`"completedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TaskService{
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want: domain.Task{
				ID:        testTaskID,
				UserID:    testUserID,
				UpdatedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - task not found",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskService{
					taskRepo: taskRepo,
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want:    domain.Task{},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Unarchive(ctx, tt.args.userID, tt.args.taskID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskService_Complete(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTaskID = "01JFEC2Q02JM8NJKGZSAE8AHP5"
	)

	type args struct {
		userID string
		taskID string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TaskService
		args    args
		want    domain.Task
		wantErr error
	}{
		{
			name: "Should complete task",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:     testTaskID,
						UserID: testUserID,
					}, nil)
				taskRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Task{
						ID:          testTaskID,
						UserID:      testUserID,
						UpdatedAt:   func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
						CompletedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTasks,
						RecordID:   testTaskID,
						OldValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null,` +
								`"completedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"archivedAt":null,` +
								`"completedAt":"2024-12-18T14:15:00Z"}`

							return &v
						}(),
					}).
					Return(nil)

				return &TaskService{
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want: domain.Task{
				ID:          testTaskID,
				UserID:      testUserID,
				UpdatedAt:   func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
				CompletedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - task not found",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskService{
					taskRepo: taskRepo,
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want:    domain.Task{},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "Should fail - task is archived",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:         testTaskID,
						UserID:     testUserID,
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &TaskService{
					taskRepo: taskRepo,
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want:    domain.Task{},
			wantErr: domain.ErrImpossibleAction,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Complete(ctx, tt.args.userID, tt.args.taskID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskService_Incomplete(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testTaskID = "01JFEC2Q02JM8NJKGZSAE8AHP5"
	)

	type args struct {
		userID string
		taskID string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TaskService
		args    args
		want    domain.Task
		wantErr error
	}{
		{
			name: "Should incomplete task",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{
						ID:          testTaskID,
						UserID:      testUserID,
						CompletedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)
				taskRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Task{
						ID:        testTaskID,
						UserID:    testUserID,
						UpdatedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameTasks,
						RecordID:   testTaskID,
						OldValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null,` +
								`"completedAt":"2024-12-18T14:15:00Z"}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"archivedAt":null,` +
								`"completedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &TaskService{
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want: domain.Task{
				ID:        testTaskID,
				UserID:    testUserID,
				UpdatedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - task not found",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				taskRepo := mocks.NewTaskRepositoryMock(t)
				taskRepo.EXPECT().
					FindByID(ctx, testUserID, testTaskID).
					Return(domain.Task{}, domain.ErrNotFound)

				return &TaskService{
					taskRepo: taskRepo,
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want:    domain.Task{},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Incomplete(ctx, tt.args.userID, tt.args.taskID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskService_Report(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID = "01JMGSY8ZP13VKEFZ1HTMZ997X"
		testTaskID = "01JMGSYCX3ZN2YBK2Q1XRFHXDQ"
	)

	type args struct {
		userID string
		taskID string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *TaskService
		args    args
		want    domain.TaskReport
		wantErr error
	}{
		{
			name: "Should report task",
			srvFunc: func(t *testing.T) *TaskService {
				t.Helper()

				timelogRepo := mocks.NewTimelogRepositoryMock(t)
				timelogRepo.EXPECT().
					TaskReport(ctx, testUserID, testTaskID).
					Return(domain.TaskReport{
						DurationSeconds:  3600,
						BillableAmount:   50000,
						UniqueDatesCount: 1,
					}, nil)

				return &TaskService{
					timelogRepo: timelogRepo,
				}
			},
			args: args{
				userID: testUserID,
				taskID: testTaskID,
			},
			want: domain.TaskReport{
				DurationSeconds:  3600,
				BillableAmount:   50000,
				UniqueDatesCount: 1,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Report(ctx, tt.args.userID, tt.args.taskID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
