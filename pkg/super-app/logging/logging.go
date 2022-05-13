// Adapter for gin context
package logging

import (
	"github.com/KoleMax/async-architecture/pkg/logging"
	"github.com/gin-gonic/gin"
)

func Debugf(ctx *gin.Context, format string, args ...interface{}) {
	logging.Debugf(ctx.Request.Context(), format, args...)
}

func Infof(ctx *gin.Context, format string, args ...interface{}) {
	logging.Infof(ctx.Request.Context(), format, args...)
}

func Warnf(ctx *gin.Context, format string, args ...interface{}) {
	logging.Warnf(ctx.Request.Context(), format, args...)
}

func Errorf(ctx *gin.Context, format string, args ...interface{}) {
	logging.Errorf(ctx.Request.Context(), format, args...)
}

func Panicf(ctx *gin.Context, format string, args ...interface{}) {
	logging.Panicf(ctx.Request.Context(), format, args...)
}
