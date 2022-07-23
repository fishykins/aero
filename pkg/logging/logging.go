package logging

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func New() *Logger {
	log := logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	log.Formatter.(*logrus.TextFormatter).ForceColors = true
	log.Formatter.(*logrus.TextFormatter).DisableTimestamp = false
	return &Logger{
		logger: log,
	}
}

func (l *Logger) WithLevel(level logrus.Level) *Logger {
	l.logger.SetLevel(level)
	return l
}

func (l *Logger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *Logger) Info(msg string) {
	l.logger.Infof(msg)
}

func (l *Logger) InfoWith(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Logger) WarnWith(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Logger) ErrorWith(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Error(msg)
}

func (l *Logger) Fatal(msg string) {
	l.logger.Fatal(msg)
}

func (l *Logger) FatalWith(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Fatal(msg)
}
