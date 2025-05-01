package koolog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(isProd bool, level zapcore.Level) (*zap.Logger, error) {
	var loggerConfig zap.Config
	if isProd {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.Level.SetLevel(level)

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
