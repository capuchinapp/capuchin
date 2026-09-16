package cloudbackend

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newLogger создает новый регистратор логов.
func newLogger(logLevel zapcore.Level) (*zap.Logger, error) {
	if logLevel == zapcore.DebugLevel {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("zap.NewDevelopment: %v", err)
		}

		logger.Debug("Debug mode running")

		return logger, nil
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("zap.NewProduction: %v", err)
	}

	return logger, nil
}
