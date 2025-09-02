package support

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a zap logger configured for the given mode (debug|release|test).
func NewLogger(mode string) (*zap.Logger, error) {
	switch mode {
	case "release":
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		return cfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	case "test":
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		return cfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	default: // debug
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		return cfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	}
}
