package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var DefaultLoggerFactory *zap.Logger
var DefaultLogger *zap.SugaredLogger

func StartLogger() error {
	var err error
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.TimeOnly)
	if DefaultLoggerFactory, err = config.Build(); err != nil {
		return err
	}
	DefaultLogger = DefaultLoggerFactory.Sugar()

	return nil
}

func EndLogger() error {
	return DefaultLoggerFactory.Sync()
}
