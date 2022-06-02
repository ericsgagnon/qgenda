package app

import (
	"go.uber.org/zap"
)

type LogConfig struct {
	Level       string //zapcore.Level // debug, info, warn, error, dpanic, panic, fatal
	Development bool
}

func NewLogConfig() LogConfig {
	return LogConfig{
		Level:       "warn",
		Development: false,
	}
}

func NewLogger(cfg *LogConfig) *zap.Logger {
	cfgLevel, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil
	}
	// zc := zap.NewProductionConfig()
	// zc.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	// return zc
	zc := zap.NewProductionConfig()

	if cfg != nil {
		zc.Level.SetLevel(cfgLevel.Level())
		zc.Development = cfg.Development
	}
	l, err := zc.Build()
	if err != nil {
		return nil
	}
	return l
}
