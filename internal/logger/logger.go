package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	level  string
	logger *zap.SugaredLogger
}

func New(level string) *Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return &Logger{level: level, logger: sugar}
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l Logger) Debug(msg string) {
	l.logger.Debug(msg)
}
