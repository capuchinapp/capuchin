package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"capuchin/internal/domain"
	"capuchin/internal/service/mocks"
)

func TestSessionService_List(t *testing.T) {
	ctx := t.Context()

	const (
		testUserID = "01JFCGS2YRPKT24B3RBZ8PVPBP"
	)

	type args struct {
		userID       string
		outdatedTime time.Time
	}
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *SessionService
		args    args
		want    []domain.Session
		wantErr error
	}{
		{
			name: "Should list sessions",
			srvFunc: func(t *testing.T) *SessionService {
				t.Helper()

				sessionRepo := mocks.NewSessionRepositoryMock(t)
				sessionRepo.EXPECT().
					FindAll(ctx, testUserID).
					Return([]domain.Session{
						{ID: "01JFCN3VVQRQBGJ255PH9YMGY8"},
						{ID: "01JFCN417MEH3ACC91HW4NXEE2"},
					}, nil)
				sessionRepo.EXPECT().
					CancelOld(ctx, time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC)).
					Return(nil)

				return &SessionService{
					sessionRepo: sessionRepo,
				}
			},
			args: args{
				userID:       testUserID,
				outdatedTime: time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
			},
			want: []domain.Session{
				{ID: "01JFCN3VVQRQBGJ255PH9YMGY8"},
				{ID: "01JFCN417MEH3ACC91HW4NXEE2"},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)

			got, err := srv.List(ctx, tt.args.userID, tt.args.outdatedTime)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
