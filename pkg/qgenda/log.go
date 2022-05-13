package qgenda

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogConfig() zap.Config {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	return zc
}
