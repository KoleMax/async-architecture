package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func extractToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("API token is required")
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("API token is incorrect")
	}

	return authHeaderParts[1], nil
}

func ParseTokenMiddleware(ctx *gin.Context) {
	tokenString, err := extractToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "Claims", claims))
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		ctx.Next()
		return
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
}
