package logging

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/KoleMax/async-architecture/pkg/logging/syslog"
	"github.com/sirupsen/logrus"
)

const (
	fieldMessage = "message"
	fieldLevel   = "levelname"
	fieldTime    = "asctime"
)

var ErrInvalidLogLevel = errors.New("invalid log level")

var (
	global *logrus.Logger
)

func init() {
	global = logrus.StandardLogger()
}

func GetGlobalLogger() *logrus.Logger {
	return global
}

func NewLogger() *logrus.Logger {
	return logrus.New()
}

func SetLevel(logger *logrus.Logger, level string) error {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logger.SetLevel(l)
	return nil
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   fieldMessage,
			logrus.FieldKeyLevel: fieldLevel,
			logrus.FieldKeyTime:  fieldTime,
		},
	})
}

func CopyGlobalHooks(logger *logrus.Logger) {
	stdLogger := GetGlobalLogger()
	for _, hooks := range stdLogger.Hooks {
		for _, h := range hooks {
			logger.AddHook(h)
		}
	}
}

func SetSyslogHook(logger *logrus.Logger, syslogUrl, level string) error { // TODO: make dynamic syslog level change, try to use mutex in syslog.Writer.
	syslogLevel, err := syslogLevelFromString(level)
	if err != nil {
		return err
	}

	appFullName := strings.Split(os.Args[0], "/") // TODO: realy need this?
	nameTag := appFullName[len(appFullName)-1]

	hook, err := newSyslogHook("udp", syslogUrl, syslogLevel, nameTag)
	if err != nil {
		return err
	}

	logger.AddHook(hook)

	return nil
}

func syslogLevelFromString(level string) (syslog.Priority, error) {
	switch strings.ToLower(level) {
	case "fatal":
		return syslog.LOG_EMERG, nil
	case "error":
		return syslog.LOG_ERR, nil
	case "warn", "warning":
		return syslog.LOG_WARNING, nil
	case "info":
		return syslog.LOG_INFO, nil
	case "debug":
		return syslog.LOG_DEBUG, nil
	}

	return syslog.LOG_EMERG, ErrInvalidLogLevel
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Debugf(format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Infof(format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Warnf(format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Errorf(format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Panicf(format, args...)
}
