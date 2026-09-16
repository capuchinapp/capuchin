package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"capuchin/internal/domain"
	"capuchin/internal/service/mocks"
)

func TestClientService_Create(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID   = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testClientID = "01JFCN3VVQRQBGJ255PH9YMGY8"
	)

	type args struct {
		client domain.Client
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *ClientService
		args    args
		want    domain.Client
		wantErr error
	}{
		{
			name: "Should create client",
			srvFunc: func(t *testing.T) *ClientService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					Create(ctx, mock.Anything, domain.Client{
						ID:        testClientID,
						UserID:    testUserID,
						CreatedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameClients,
						RecordID:   testClientID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
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

				return &ClientService{
					clientRepo:   clientRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					idFunc:       func() string { return testClientID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				client: domain.Client{
					UserID: testUserID,
				},
			},
			want: domain.Client{
				ID:        testClientID,
				UserID:    testUserID,
				CreatedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "Should return error - userID is empty",
			srvFunc: func(t *testing.T) *ClientService {
				t.Helper()

				return &ClientService{}
			},
			args: args{
				client: domain.Client{
					UserID: "",
				},
			},
			want:    domain.Client{},
			wantErr: ErrUserIDRequired,
		},
		{
			name: "Should return error - client already exists",
			srvFunc: func(t *testing.T) *ClientService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					Create(ctx, mock.Anything, domain.Client{
						ID:        testClientID,
						UserID:    testUserID,
						CreatedAt: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeCreate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameClients,
						RecordID:   testClientID,
						OldValue:   nil,
						NewValue: func() *string {
							v := `{"id":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
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

				return &ClientService{
					clientRepo:   clientRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, domain.ErrAlreadyExists),
					idFunc:       func() string { return testClientID },
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				client: domain.Client{
					UserID: testUserID,
				},
			},
			want:    domain.Client{},
			wantErr: domain.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Create(ctx, tt.args.client)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientService_Update(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID   = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testClientID = "01JFCN3VVQRQBGJ255PH9YMGY8"
	)

	type args struct {
		client domain.Client
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *ClientService
		args    args
		want    domain.Client
		wantErr error
	}{
		{
			name: "Should update client",
			srvFunc: func(t *testing.T) *ClientService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:           testClientID,
						UserID:       testUserID,
						Name:         "Old Name",
						BillableRate: 1000,
						Comment:      "Old Comment",
					}, nil)
				clientRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Client{
						ID:           testClientID,
						UserID:       testUserID,
						Name:         "New Name",
						BillableRate: 2000,
						Comment:      "New Comment",
						UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}).
					Return(nil)

				auditLogRepo := mocks.NewAuditLogRepositoryMock(t)
				auditLogRepo.EXPECT().
					Create(ctx, mock.Anything, domain.AuditLog{
						ActionAt:   time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
						ActionType: domain.AuditLogActionTypeUpdate,
						UserID:     testUserID,
						TableName:  domain.AuditLogTableNameClients,
						RecordID:   testClientID,
						OldValue: func() *string {
							v := `{"id":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"name":"Old Name",` +
								`"billableRate":1000,` +
								`"comment":"Old Comment",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
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

				return &ClientService{
					clientRepo:   clientRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				client: domain.Client{
					ID:           testClientID,
					UserID:       testUserID,
					Name:         "New Name",
					BillableRate: 2000,
					Comment:      "New Comment",
				},
			},
			want: domain.Client{
				ID:           testClientID,
				UserID:       testUserID,
				Name:         "New Name",
				BillableRate: 2000,
				Comment:      "New Comment",
				UpdatedAt:    func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - client not found",
			srvFunc: func(t *testing.T) *ClientService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{}, domain.ErrNotFound)

				return &ClientService{
					clientRepo: clientRepo,
				}
			},
			args: args{
				client: domain.Client{
					ID:           testClientID,
					UserID:       testUserID,
					Name:         "New Name",
					BillableRate: 2000,
					Comment:      "New Comment",
				},
			},
			want:    domain.Client{},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Update(ctx, tt.args.client)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientService_Archive(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID   = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testClientID = "01JFCN3VVQRQBGJ255PH9YMGY8"
	)

	type args struct {
		userID   string
		clientID string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *ClientService
		args    args
		want    domain.Client
		wantErr error
	}{
		{
			name: "Should archive client",
			srvFunc: func(t *testing.T) *ClientService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:     testClientID,
						UserID: testUserID,
					}, nil)
				clientRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Client{
						ID:         testClientID,
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
						TableName:  domain.AuditLogTableNameClients,
						RecordID:   testClientID,
						OldValue: func() *string {
							v := `{"id":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"name":"",` +
								`"billableRate":0,` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":null}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
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

				return &ClientService{
					clientRepo:   clientRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID:   testUserID,
				clientID: testClientID,
			},
			want: domain.Client{
				ID:         testClientID,
				UserID:     testUserID,
				UpdatedAt:  func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
				ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - client not found",
			srvFunc: func(t *testing.T) *ClientService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{}, domain.ErrNotFound)

				return &ClientService{
					clientRepo: clientRepo,
				}
			},
			args: args{
				userID:   testUserID,
				clientID: testClientID,
			},
			want:    domain.Client{},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Archive(ctx, tt.args.userID, tt.args.clientID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientService_Unarchive(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID   = "01JFCGS2YRPKT24B3RBZ8PVPBP"
		testClientID = "01JFCN3VVQRQBGJ255PH9YMGY8"
	)

	type args struct {
		userID   string
		clientID string
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *ClientService
		args    args
		want    domain.Client
		wantErr error
	}{
		{
			name: "Should unarchive client",
			srvFunc: func(t *testing.T) *ClientService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{
						ID:         testClientID,
						UserID:     testUserID,
						ArchivedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
					}, nil)
				clientRepo.EXPECT().
					Update(ctx, mock.Anything, domain.Client{
						ID:        testClientID,
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
						TableName:  domain.AuditLogTableNameClients,
						RecordID:   testClientID,
						OldValue: func() *string {
							v := `{"id":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
								`"name":"",` +
								`"billableRate":0,` +
								`"comment":"",` +
								`"createdAt":"0001-01-01T00:00:00Z",` +
								`"updatedAt":null,` +
								`"archivedAt":"2024-12-18T14:15:00Z"}`

							return &v
						}(),
						NewValue: func() *string {
							v := `{"id":"01JFCN3VVQRQBGJ255PH9YMGY8",` +
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

				return &ClientService{
					clientRepo:   clientRepo,
					auditLogRepo: auditLogRepo,
					txManager:    newTxManager(t, ctx, nil),
					nowFunc:      func() time.Time { return time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC) },
				}
			},
			args: args{
				userID:   testUserID,
				clientID: testClientID,
			},
			want: domain.Client{
				ID:        testClientID,
				UserID:    testUserID,
				UpdatedAt: func() *time.Time { t := time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC); return &t }(),
			},
			wantErr: nil,
		},
		{
			name: "Should fail - client not found",
			srvFunc: func(t *testing.T) *ClientService {
				t.Helper()

				clientRepo := mocks.NewClientRepositoryMock(t)
				clientRepo.EXPECT().
					FindByID(ctx, testUserID, testClientID).
					Return(domain.Client{}, domain.ErrNotFound)

				return &ClientService{
					clientRepo: clientRepo,
				}
			},
			args: args{
				userID:   testUserID,
				clientID: testClientID,
			},
			want:    domain.Client{},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.Unarchive(ctx, tt.args.userID, tt.args.clientID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
