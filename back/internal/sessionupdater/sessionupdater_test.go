package sessionupdater

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"capuchin/internal/sessionupdater/mocks"
)

func TestSessionUpdater_UpdateTime(t *testing.T) {
	su := &SessionUpdater{
		cache:   make(map[string]time.Time),
		nowFunc: func() time.Time { return time.Date(2024, 12, 18, 14, 20, 0, 0, time.UTC) },
	}

	_, ok := su.cache["01JHHXTMRFMJ4DQJVTVCHZSZX3"]
	assert.False(t, ok)

	su.UpdateTime("01JHHXTMRFMJ4DQJVTVCHZSZX3")

	actual, ok := su.cache["01JHHXTMRFMJ4DQJVTVCHZSZX3"]
	assert.True(t, ok)
	assert.Equal(t, time.Date(2024, 12, 18, 14, 20, 0, 0, time.UTC), actual)
}

func TestSessionUpdater_Stop(t *testing.T) {
	su := &SessionUpdater{
		logger:   zap.NewNop(),
		stopChan: make(chan struct{}),
	}

	su.Stop()

	_, ok := <-su.stopChan
	assert.False(t, ok)
}

func TestSessionUpdater_saveLoop(_ *testing.T) {
	testSaveInterval := 3 * time.Second

	su := &SessionUpdater{
		logger: zap.NewNop(),

		saveInterval: testSaveInterval,
		cache:        make(map[string]time.Time),
		stopChan:     make(chan struct{}),
	}

	go su.saveLoop()

	// Ждем немного, чтобы saveLoop успел сработать по таймеру
	time.Sleep(testSaveInterval + time.Second)

	su.Stop()
}

func TestSessionUpdater_save(t *testing.T) {
	tests := []struct {
		name    string
		srvFunc func(t *testing.T) *SessionUpdater
	}{
		{
			name: "Successful save",
			srvFunc: func(t *testing.T) *SessionUpdater {
				t.Helper()

				sessionRepo := mocks.NewSessionRepoMock(t)
				sessionRepo.EXPECT().
					UpdateCheckedAt(mock.Anything, "01JHHXTMRFMJ4DQJVTVCHZSZX3", time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC)).
					Return(nil)

				return &SessionUpdater{
					sessionRepo: sessionRepo,
					logger:      zap.NewNop(),
					cache: map[string]time.Time{
						"01JHHXTMRFMJ4DQJVTVCHZSZX3": time.Date(2024, 12, 18, 14, 15, 0, 0, time.UTC),
					},
					stopChan: make(chan struct{}),
				}
			},
		},
		{
			name: "Empty cache",
			srvFunc: func(t *testing.T) *SessionUpdater {
				t.Helper()

				return &SessionUpdater{
					logger:   zap.NewNop(),
					cache:    map[string]time.Time{},
					stopChan: make(chan struct{}),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.srvFunc(t)
			defer srv.Stop() // на случай если тест упадет

			srv.save()

			srv.saveWG.Wait()

			assert.Empty(t, srv.cache)
		})
	}
}
