package middleware

import (
	"fmt"
	"net/http"

	"github.com/KoleMax/async-architecture/pkg/logging"
	"github.com/gin-gonic/gin"
)

type recoveryError struct {
	Error string `json:"error"`
}

// TODO: add stack to logging

func Recover(ctx *gin.Context, err interface{}) {
	logging.Errorf(ctx, "recovered from panic: %v", err)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, recoveryError{
		Error: fmt.Sprintf("recovered from panic: %v", err),
	})
}
