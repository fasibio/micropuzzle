package logger

import (
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var zapLogger = zap.NewNop().Sugar()

//Initialize make a new Logger by given level
func Initialize(level string) (*zap.SugaredLogger, error) {

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "@timestamp"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.CallerKey = ""
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = zapLogLevel(level)

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	zapLogger = logger.Sugar()
	return zapLogger, nil
}

// Get return the current logger
func Get() *zap.SugaredLogger {
	return zapLogger
}

func zapLogLevel(level string) zap.AtomicLevel {

	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	// default is info level
	return zap.NewAtomicLevelAt(zap.InfoLevel)
}
