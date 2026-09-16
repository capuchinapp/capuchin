package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"capuchin/internal/domain"
	"capuchin/internal/service/mocks"
)

func TestProjectService_Create(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testClientID  = "01JFCN3VVQRQBGJ255PH9YMGY8"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
	)

	type args struct {
		project domain.Project
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *ProjectService
		args    args
		want    domain.Project
		wantErr error
	}{
		{
			name: "Should create project",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:        testClientID,
						UserID:    testUserID,
						Name:      "Client Name",
						CreatedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					Create(ctx, mock.Anything, domain.Project{
						ID:         testProjectID,
						UserID:     testUserID,
						ClientID:   testClientID,
						ClientName: "Client Name",
						CreatedAt:  time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameProjects,
						RecordID:   testProjectID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"name":"",` +
								`"billableRate":0,` +
								`"comment":"",` +
								`"createdAt":"2024-12-18T14:15:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &ProjectService{
					projectRepo:  projectRepo,
					clientRepo:   clientRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					idFunc:       func() string { return testProjectID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				project: domain.Project{
					UserID:   testUserID,
					ClientID: testClientID,
				},
			},
			want: domain.Project{
				ID:         testProjectID,
				UserID:     testUserID,
				ClientID:   testClientID,
				ClientName: "Client Name",
				CreatedAt:  time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "Should return error - client not found",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{}, domain.ErrNotFound)

				return &ProjectService{
					clientRepo: clientRepo,
				}
			},
			args: args{
				project: domain.Project{
					UserID:   testUserID,
					ClientID: testClientID,
				},
			},
			want:    domain.Project{},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "Should return error - client is archived",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:         testClientID,
						UserID:     testUserID,
						CreatedAt:  time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &ProjectService{
					clientRepo: clientRepo,
				}
			},
			args: args{
				project: domain.Project{
					UserID:   testUserID,
					ClientID: testClientID,
				},
			},
			want:    domain.Project{},
			wantErr: domain.ErrClientArchived,
		},
		{
			name: "Should return error - project already exists",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:        testClientID,
						UserID:    testUserID,
						CreatedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					Create(ctx, mock.Anything, domain.Project{
						ID:        testProjectID,
						UserID:    testUserID,
						ClientID:  testClientID,
						CreatedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameProjects,
						RecordID:   testProjectID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"",` +
								`"name":"",` +
								`"billableRate":0,` +
								`"comment":"",` +
								`"createdAt":"2024-12-18T14:15:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &ProjectService{
					clientRepo:   clientRepo,
					projectRepo:  projectRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, domain.ErrAlreadyExists),
					idFunc:       func() string { return testProjectID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				project: domain.Project{
					UserID:   testUserID,
					ClientID: testClientID,
				},
			},
			want:    domain.Project{},
			wantErr: domain.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Create(ctx, tt.args.project)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectService_Update(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testClientID  = "01JFCN3VVQRQBGJ255PH9YMGY8"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
	)

	type args struct {
		project domain.Project
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *ProjectService
		args    args
		want    domain.Project
		wantErr error
	}{
		{
			name: "Should update project",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:           testProjectID,
						UserID:       testUserID,
						ClientID:     testClientID,
						Name:         "Old Name",
						BillableRate: 1000,
						Comment:      "Old Comment",
					}, nil)
				projectRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Project{
						ID:           testProjectID,
						UserID:       testUserID,
						ClientID:     testClientID,
						ClientName:   "Client Name",
						Name:         "New Name",
						BillableRate: 2000,
						Comment:      "New Comment",
						UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:        testClientID,
						UserID:    testUserID,
						Name:      "Client Name",
						CreatedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}, nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameProjects,
						RecordID:   testProjectID,
						OldValue: func() *string {
							v := `{"id":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"",` +
								`"name":"Old Name",` +
								`"billableRate":1000,` +
								`"comment":"Old Comment",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"clientId":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"clientName":"Client Name",` +
								`"name":"New Name",` +
								`"billableRate":2000,` +
								`"comment":"New Comment",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"archivedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &ProjectService{
					projectRepo:  projectRepo,
					clientRepo:   clientRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				project: domain.Project{
					ID:           testProjectID,
					UserID:       testUserID,
					ClientID:     testClientID,
					Name:         "New Name",
					BillableRate: 2000,
					Comment:      "New Comment",
				},
			},
			want: domain.Project{
				ID:           testProjectID,
				UserID:       testUserID,
				ClientID:     testClientID,
				ClientName:   "Client Name",
				Name:         "New Name",
				BillableRate: 2000,
				Comment:      "New Comment",
				UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - project not found",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{}, domain.ErrNotFound)

				return &ProjectService{
					projectRepo: projectRepo,
				}
			},
			args: args{
				project: domain.Project{
					ID:           testProjectID,
					UserID:       testUserID,
					ClientID:     testClientID,
					Name:         "New Name",
					BillableRate: 2000,
					Comment:      "New Comment",
				},
			},
			want:    domain.Project{},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "Should return error - client not found",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:           testProjectID,
						UserID:       testUserID,
						ClientID:     testClientID,
						Name:         "Old Name",
						BillableRate: 1000,
						Comment:      "Old Comment",
					}, nil)

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{}, domain.ErrNotFound)

				return &ProjectService{
					projectRepo: projectRepo,
					clientRepo:  clientRepo,
				}
			},
			args: args{
				project: domain.Project{
					ID:       testProjectID,
					UserID:   testUserID,
					ClientID: testClientID,
				},
			},
			want:    domain.Project{},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "Should return error - client is archived",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:           testProjectID,
						UserID:       testUserID,
						ClientID:     testClientID,
						Name:         "Old Name",
						BillableRate: 1000,
						Comment:      "Old Comment",
					}, nil)

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:         testClientID,
						UserID:     testUserID,
						CreatedAt:  time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)

				return &ProjectService{
					projectRepo: projectRepo,
					clientRepo:  clientRepo,
				}
			},
			args: args{
				project: domain.Project{
					ID:       testProjectID,
					UserID:   testUserID,
					ClientID: testClientID,
				},
			},
			want:    domain.Project{},
			wantErr: domain.ErrClientArchived,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Update(ctx, tt.args.project)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectService_Archive(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
	)

	type args struct {
		userID    string
		projectID string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *ProjectService
		args    args
		want    domain.Project
		wantErr error
	}{
		{
			name: "Should archive project",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:     testProjectID,
						UserID: testUserID,
					}, nil)
				projectRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Project{
						ID:         testProjectID,
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
						TableName:  domain.AuditLogTableNameProjects,
						RecordID:   testProjectID,
						OldValue: func() *string {
							v := `{"id":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"billableRate":0,` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"billableRate":0,` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"archivedAt":"2024-12-18T14:15:00Z"}`

							return &v
						}(),
					}).
					Return(nil)

				return &ProjectService{
					projectRepo:  projectRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID:    testUserID,
				projectID: testProjectID,
			},
			want: domain.Project{
				ID:         testProjectID,
				UserID:     testUserID,
				UpdatedAt:  func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
				ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - project not found",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{}, domain.ErrNotFound)

				return &ProjectService{
					projectRepo: projectRepo,
				}
			},
			args: args{
				userID:    testUserID,
				projectID: testProjectID,
			},
			want:    domain.Project{},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Archive(ctx, tt.args.userID, tt.args.projectID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectService_Unarchive(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID    = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testProjectID = "01JFD0MHHR7WMYRZRZFFXT6QBZ"
	)

	type args struct {
		userID    string
		projectID string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *ProjectService
		args    args
		want    domain.Project
		wantErr error
	}{
		{
			name: "Should unarchive project",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{
						ID:         testProjectID,
						UserID:     testUserID,
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)
				projectRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Project{
						ID:        testProjectID,
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
						TableName:  domain.AuditLogTableNameProjects,
						RecordID:   testProjectID,
						OldValue: func() *string {
							v := `{"id":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"billableRate":0,` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":"2024-12-18T14:15:00Z"}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFD0MHHR7WMYRZRZFFXT6QBZ",` +
								`"clientId":"",` +
								`"clientName":"",` +
								`"name":"",` +
								`"billableRate":0,` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":"2024-12-18T14:15:00Z",` +
								`"archivedAt":null}`

							return &v
						}(),
					}).
					Return(nil)

				return &ProjectService{
					projectRepo:  projectRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID:    testUserID,
				projectID: testProjectID,
			},
			want: domain.Project{
				ID:        testProjectID,
				UserID:    testUserID,
				UpdatedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - project not found",
			srvFunc: func(t *testing.T) *ProjectService {
				t.Helper()

				projectRepo := mocks.NewProjectRepositoryMock(t)
				projectRepo.EXPECT().
					FindByID(ctx, testUserID, testProjectID).
					Return(domain.Project{}, domain.ErrNotFound)

				return &ProjectService{
					projectRepo: projectRepo,
				}
			},
			args: args{
				userID:    testUserID,
				projectID: testProjectID,
			},
			want:    domain.Project{},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Unarchive(ctx, tt.args.userID, tt.args.projectID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
