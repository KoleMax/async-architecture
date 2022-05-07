package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerCtxKey int

const loggerContextKey loggerCtxKey = iota

func FromContext(ctx context.Context) *logrus.Logger {
	logger, ok := ctx.Value(loggerContextKey).(*logrus.Logger)
	if !ok {
		logger = global
	}
	return logger
}

func ToContext(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}
