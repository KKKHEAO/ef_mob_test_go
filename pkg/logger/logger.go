package logger

import (
	"ef_mob_test_go/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) (*zap.SugaredLogger, error) {
	if cfg.Server.AppEnv == "production" {
		return newProduction()
	}
	return newDevelopment()
}

func newProduction() (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zLog, err := cfg.Build()
	return zLog.Sugar(), err
}

func newDevelopment() (*zap.SugaredLogger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zLog, err := cfg.Build()
	return zLog.Sugar(), err
}
