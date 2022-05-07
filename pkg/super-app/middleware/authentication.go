package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const authenticationURL = "http://localhost:8081/api/v1/auth/authenticate"

func makeAuthRequest(ctx *gin.Context) (string, error) {
	resp, err := http.Post(authenticationURL, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res["json"])

	return "", nil
}

func Authenticate(ctx *gin.Context) {
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
