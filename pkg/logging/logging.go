package log

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

var DefaultLogger *Logger = New()

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

// ============================================================================
// ============================================================================
func (l *Logger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}
func Trace(args ...interface{}) {
	DefaultLogger.Trace(args...)
}

// ============================================================================
func (l *Logger) TraceWith(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Trace(msg)
}
func TraceWith(msg string, fields map[string]interface{}) {
	DefaultLogger.TraceWith(msg, fields)
}

// ============================================================================
// ============================================================================
func (l *Logger) Info(msg string) {
	l.logger.Infof(msg)
}
func Info(msg string) {
	DefaultLogger.Info(msg)
}

// ============================================================================
func (l *Logger) InfoWith(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Info(msg)
}
func InfoWith(msg string, fields map[string]interface{}) {
	DefaultLogger.InfoWith(msg, fields)
}

// ============================================================================
// ============================================================================
func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}
func Warn(msg string) {
	DefaultLogger.Warn(msg)
}

// ============================================================================
func (l *Logger) WarnWith(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Warn(msg)
}
func WarnWith(msg string, fields map[string]interface{}) {
	DefaultLogger.WarnWith(msg, fields)
}

// ============================================================================
// ============================================================================
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}
func Error(msg string) {
	DefaultLogger.Error(msg)
}

// ============================================================================
func (l *Logger) ErrorWith(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Error(msg)
}
func ErrorWith(msg string, fields map[string]interface{}) {
	DefaultLogger.ErrorWith(msg, fields)
}

// ============================================================================
// ============================================================================
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal(msg)
}
func Fatal(msg string) {
	DefaultLogger.Fatal(msg)
}

// ============================================================================
func (l *Logger) FatalWith(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Fatal(msg)
}
func FatalWith(msg string, fields map[string]interface{}) {
	DefaultLogger.FatalWith(msg, fields)
}

// ============================================================================
// ============================================================================
