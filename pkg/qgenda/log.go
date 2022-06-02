package qgenda

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Level zap.AtomicLevel
}

func NewLogConfig() zap.Config {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	return zc
}
