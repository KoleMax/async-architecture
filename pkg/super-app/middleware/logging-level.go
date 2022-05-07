package middleware

import (
	"fmt"
	"net/http"

	"github.com/KoleMax/async-architecture/pkg/logging"
	"github.com/gin-gonic/gin"
)

const logLevelHeader = "x-log-level"

type InvalidLogLevel struct {
	Error string `json:"error"`
}

func LogLevelOverride(ctx *gin.Context) {
	if overrideLogLevel := ctx.GetHeader(logLevelHeader); overrideLogLevel != "" {
		newLogger := logging.NewLogger()
		if err := logging.SetLevel(newLogger, overrideLogLevel); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, InvalidLogLevel{
				Error: fmt.Sprintf("Logging set level from '%s' header error: %v", logLevelHeader, err),
			})
			return
		}
		logging.SetFormatter(newLogger)
		logging.CopyGlobalHooks(newLogger)
		logging.Warnf(ctx, "Received '%s', overriding log level in ctx with %s", logLevelHeader, overrideLogLevel)
		c := logging.ToContext(ctx.Request.Context(), newLogger)
		ctx.Request = ctx.Request.WithContext(c)
	}
	ctx.Next()
}
