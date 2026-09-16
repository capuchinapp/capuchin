package sessionupdater

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"capuchin/internal/domain"
)

const (
	saveInterval = 60 * time.Second // Периодичность сброса данных в БД

	oneSaveTimeoutSec  = 3
	defaultSaveTimeout = 10 * time.Second
)

// SessionUpdater периодически обновляет данные в БД.
type SessionUpdater struct {
	sessionRepo SessionRepo
	logger      *zap.Logger

	saveInterval time.Duration
	saveWG       sync.WaitGroup       // отдельный WaitGroup только для save операций
	cache        map[string]time.Time // string - sessionID
	mu           sync.Mutex
	stopChan     chan struct{}

	nowFunc func() time.Time
}

// New создает новый экземпляр SessionUpdater.
func New(sessionRepo SessionRepo, logger *zap.Logger) *SessionUpdater {
	su := &SessionUpdater{
		sessionRepo: sessionRepo,
		logger:      logger.Named("session_updater"),

		saveInterval: saveInterval,
		cache:        make(map[string]time.Time),
		stopChan:     make(chan struct{}),

		nowFunc: time.Now,
	}

	go su.saveLoop()
	su.logger.Info("Session updater started")

	return su
}

// UpdateTime обновляет время в кэше.
func (su *SessionUpdater) UpdateTime(sessionID string) {
	su.mu.Lock()
	defer su.mu.Unlock()

	su.cache[sessionID] = su.nowFunc()
}

// Stop останавливает фоновую горутину.
func (su *SessionUpdater) Stop() {
	close(su.stopChan)
	su.logger.Info("Session updater stopped")
	su.saveWG.Wait() // ждем завершения всех save операций
}

// saveLoop периодически сохраняет данные в БД.
func (su *SessionUpdater) saveLoop() {
	ticker := time.NewTicker(su.saveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			su.save() // Сохраняем накопленные данные в БД
		case <-su.stopChan:
			su.save() // Сохраняем оставшиеся данные перед выходом
			return
		}
	}
}

// save сохраняет накопленные данные в БД.
func (su *SessionUpdater) save() {
	su.mu.Lock()
	defer su.mu.Unlock()

	if len(su.cache) == 0 {
		return
	}

	// Создаем копию данных для работы
	cacheCopy := su.cache
	su.cache = make(map[string]time.Time)

	su.saveWG.Add(1)
	go func(data map[string]time.Time) {
		defer su.saveWG.Done()

		timeout, err := time.ParseDuration(
			fmt.Sprintf("%ds", oneSaveTimeoutSec*len(data)),
		)
		if err != nil {
			timeout = defaultSaveTimeout
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		for sessionID, checkedAt := range data {
			if err := su.sessionRepo.UpdateCheckedAt(ctx, sessionID, checkedAt); err != nil {
				su.logger.Error("Failed to update session",
					zap.String("part_session_id", domain.PartID(sessionID)),
					zap.Error(err),
				)
			}
		}

		su.logger.Info("Sessions saved", zap.Int("count", len(data)))
	}(cacheCopy)
}
