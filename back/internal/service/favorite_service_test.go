package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"capuchin/internal/domain"
	"capuchin/internal/service/mocks"
)

func TestFavoriteService_Create(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID     = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testClientID   = "01JFCN3VVQRQBGJ255PH9YMGY8"
		testProjectID  = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
		testFavoriteID = "01JFEC2BRSGJKVG1FKNWSNV6HH"
		testTaskID     = "01JFEC2Q02JM8NJKGZSAE8AHP5"

		testTaskName = "Task Name"
	)

	type args struct {
		favorite domain.Favorite
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *FavoriteService
		args    args
		want    domain.Favorite
		wantErr error
	}{
		{
			name: "Should create favorite",
			srvFunc: func(t *testing.T) *FavoriteService {
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

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					Create(ctx, mock.Anything, domain.Favorite{
						ID:          testFavoriteID,
						UserID:      testUserID,
						ClientID:    testClientID,
						ClientName:  "Client Name",
						ProjectID:   testProjectID,
						ProjectName: "Project Name",
						TaskID:      func() *string { id := testTaskID; return &id }(),
						TaskName:    func() *string { name := testTaskName; return &name }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameFavorites,
						RecordID:   testFavoriteID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFEC2BRSGJKVG1FKNWSNV6HH",` +
								`"name":"",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":"Task Name",` +
								`"billableRate":0,` +
								`"comment":""}`

							return &v
						}(),
					}).
					Return(nil)

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					idFunc:       func() string { return testFavoriteID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				favorite: domain.Favorite{
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { id := testTaskID; return &id }(),
				},
			},
			want: domain.Favorite{
				ID:          testFavoriteID,
				UserID:      testUserID,
				ClientID:    testClientID,
				ClientName:  "Client Name",
				ProjectID:   testProjectID,
				ProjectName: "Project Name",
				TaskID:      func() *string { id := testTaskID; return &id }(),
				TaskName:    func() *string { name := testTaskName; return &name }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - project not found",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{}, domain.ErrNotFound)

				return &FavoriteService{
					projectRepo: projectRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrProjectNotFound,
		},
		{
			name: "Should fail - project is archived",
			srvFunc: func(t *testing.T) *FavoriteService {
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

				return &FavoriteService{
					projectRepo: projectRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrProjectArchived,
		},
		{
			name: "Should fail - client not found",
			srvFunc: func(t *testing.T) *FavoriteService {
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

				return &FavoriteService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "Should fail - client is archived",
			srvFunc: func(t *testing.T) *FavoriteService {
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

				return &FavoriteService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrClientArchived,
		},
		{
			name: "Should fail - task not found",
			srvFunc: func(t *testing.T) *FavoriteService {
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

				return &FavoriteService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					idFunc:      func() string { return testFavoriteID },
				}
			},
			args: args{
				favorite: domain.Favorite{
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrTaskNotFound,
		},
		{
			name: "Should fail - task is archived",
			srvFunc: func(t *testing.T) *FavoriteService {
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

				return &FavoriteService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					idFunc:      func() string { return testFavoriteID },
				}
			},
			args: args{
				favorite: domain.Favorite{
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrTaskArchived,
		},
		{
			name: "Should fail - task is completed",
			srvFunc: func(t *testing.T) *FavoriteService {
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

				return &FavoriteService{
					clientRepo:  clientRepo,
					projectRepo: projectRepo,
					taskRepo:    taskRepo,
					idFunc:      func() string { return testFavoriteID },
				}
			},
			args: args{
				favorite: domain.Favorite{
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrTaskCompleted,
		},
		{
			name: "Should fail - favorite already exists",
			srvFunc: func(t *testing.T) *FavoriteService {
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

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					Create(ctx, mock.Anything, domain.Favorite{
						ID:          testFavoriteID,
						UserID:      testUserID,
						ClientID:    testClientID,
						ClientName:  "Client Name",
						ProjectID:   testProjectID,
						ProjectName: "Project Name",
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameFavorites,
						RecordID:   testFavoriteID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFEC2BRSGJKVG1FKNWSNV6HH",` +
								`"name":"",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"taskId":null,` +
								`"taskName":null,` +
								`"billableRate":0,` +
								`"comment":""}`

							return &v
						}(),
					}).
					Return(nil)

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, domain.ErrAlreadyExists),
					idFunc:       func() string { return testFavoriteID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				favorite: domain.Favorite{
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Create(ctx, tt.args.favorite)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFavoriteService_Update(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID     = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testClientID   = "01JFCN3VVQRQBGJ255PH9YMGY8"
		testProjectID  = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
		testFavoriteID = "01JFEC2BRSGJKVG1FKNWSNV6HH"
		testTaskID     = "01JFEC2Q02JM8NJKGZSAE8AHP5"

		testTaskName = "Task Name"
	)

	type args struct {
		favorite domain.Favorite
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *FavoriteService
		args    args
		want    domain.Favorite
		wantErr error
	}{
		{
			name: "Should update favorite",
			srvFunc: func(t *testing.T) *FavoriteService {
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

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{
						ID:           testFavoriteID,
						UserID:       testUserID,
						Name:         "Old Name",
						ClientID:     testClientID,
						ClientName:   "Client Name",
						ProjectID:    testProjectID,
						ProjectName:  "Project Name",
						BillableRate: 1000,
						Comment:      "Old Comment",
					}, nil)
				favoriteRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Favorite{
						ID:           testFavoriteID,
						UserID:       testUserID,
						Name:         "New Name",
						ClientID:     testClientID,
						ClientName:   "Client Name",
						ProjectID:    testProjectID,
						ProjectName:  "Project Name",
						TaskID:       func() *string { id := testTaskID; return &id }(),
						TaskName:     func() *string { name := testTaskName; return &name }(),
						BillableRate: 2000,
						Comment:      "New Comment",
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameFavorites,
						RecordID:   testFavoriteID,
						OldValue: func() *string {
							v := `{"id":"01JFEC2BRSGJKVG1FKNWSNV6HH",` +
								`"name":"Old Name",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"taskId":null,` +
								`"taskName":null,` +
								`"billableRate":1000,` +
								`"comment":"Old Comment"}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFEC2BRSGJKVG1FKNWSNV6HH",` +
								`"name":"New Name",` +
								`"projectId":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"projectName":"Project Name",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"taskId":"01JFEC2Q02JM8NJKGZSAE8AHP5",` +
								`"taskName":"Task Name",` +
								`"billableRate":2000,` +
								`"comment":"New Comment"}`

							return &v
						}(),
					}).
					Return(nil)

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					taskRepo:     taskRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				favorite: domain.Favorite{
					ID:           testFavoriteID,
					UserID:       testUserID,
					Name:         "New Name",
					ProjectID:    testProjectID,
					TaskID:       func() *string { id := testTaskID; return &id }(),
					BillableRate: 2000,
					Comment:      "New Comment",
				},
			},
			want: domain.Favorite{
				ID:           testFavoriteID,
				UserID:       testUserID,
				ClientID:     testClientID,
				ClientName:   "Client Name",
				ProjectID:    testProjectID,
				ProjectName:  "Project Name",
				TaskID:       func() *string { id := testTaskID; return &id }(),
				TaskName:     func() *string { name := testTaskName; return &name }(),
				Name:         "New Name",
				BillableRate: 2000,
				Comment:      "New Comment",
			},
			wantErr: nil,
		},
		{
			name: "Should fail - favorite not found",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{}, domain.ErrNotFound)

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					ID:     testFavoriteID,
					UserID: testUserID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "Should fail - project not found",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{
						UserID:    testUserID,
						ProjectID: testProjectID,
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{}, domain.ErrNotFound)

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					projectRepo:  projectRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					ID:        testFavoriteID,
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrProjectNotFound,
		},
		{
			name: "Should fail - project is archived",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{
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

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					projectRepo:  projectRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					ID:        testFavoriteID,
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrProjectArchived,
		},
		{
			name: "Should fail - client not found",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{
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

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					ID:        testFavoriteID,
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "Should fail - client is archived",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{
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

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					ID:        testFavoriteID,
					UserID:    testUserID,
					ProjectID: testProjectID,
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrClientArchived,
		},
		{
			name: "Should fail - task not found",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{
						UserID:    testUserID,
						ProjectID: testProjectID,
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

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					taskRepo:     taskRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					ID:        testFavoriteID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrTaskNotFound,
		},
		{
			name: "Should fail - task is archived",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{
						UserID:    testUserID,
						ProjectID: testProjectID,
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

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					taskRepo:     taskRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					ID:        testFavoriteID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrTaskArchived,
		},
		{
			name: "Should fail - task is completed",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{
						UserID:    testUserID,
						ProjectID: testProjectID,
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

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					taskRepo:     taskRepo,
				}
			},
			args: args{
				favorite: domain.Favorite{
					ID:        testFavoriteID,
					UserID:    testUserID,
					ProjectID: testProjectID,
					TaskID:    func() *string { s := testTaskID; return &s }(),
				},
			},
			want:    domain.Favorite{},
			wantErr: domain.ErrTaskCompleted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Update(ctx, tt.args.favorite)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFavoriteService_Delete(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID     = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testFavoriteID = "01JFEC2BRSGJKVG1FKNWSNV6HH"
	)

	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *FavoriteService
	}{
		{
			name: "Should delete favorite",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{
						ID:     testFavoriteID,
						UserID: testUserID,
					}, nil)
				favoriteRepo.EXPECT().
					Delete(ctx, mock.Anything, testUserID, testFavoriteID).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeDelete,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameFavorites,
						RecordID:   testFavoriteID,
						OldValue: func() *string {
							v := `{"id":"01JFEC2BRSGJKVG1FKNWSNV6HH",` +
								`"name":"",` +
								`"projectId":"",` +
								`"projectName":"",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"taskId":null,` +
								`"taskName":null,` +
								`"billableRate":0,` +
								`"comment":""}`

							return &v
						}(),
						NewValue: nil,
					}).
					Return(nil)

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
		},
		{
			name: "Should fail - favorite not found",
			srvFunc: func(t *testing.T) *FavoriteService {
				t.Helper()

				favoriteRepo := mocks.NewFavoriteRepositoryMock(t)
				favoriteRepo.EXPECT().
					FindByID(ctx, testUserID, testFavoriteID).
					Return(domain.Favorite{}, domain.ErrNotFound)

				return &FavoriteService{
					favoriteRepo: favoriteRepo,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			err := srv.Delete(ctx, testUserID, testFavoriteID)
			assert.NoError(t, err)
		})
	}
}
