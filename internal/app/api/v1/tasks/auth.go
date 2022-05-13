package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/KoleMax/async-architecture/pkg/logging"
	"github.com/gin-gonic/gin"
)

func Authorize(ctx *gin.Context, allowedRoles []string) *AuthAccount {
	logging.Infof(ctx, "authorization")

	authHeader, ok := ctx.Request.Header["Authorization"]
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization token"})
		return nil
	}

	reqURL, _ := url.Parse("http://localhost:8081/api/v1/auth/authenticate")
	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type":  {"application/json; charset=UTF-8"},
			"Authorization": authHeader,
		},
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return nil
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		ctx.AbortWithStatusJSON(resp.StatusCode, gin.H{"error": res["error"]})
		return nil
	}

	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "Account", res["account"]))

	var authAccount AuthAccount

	accountJSON, ok := res["account"]
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no account in auth response"})
		return nil
	}

	jsonBytes, err := json.Marshal(accountJSON)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return nil
	}

	err = json.Unmarshal(jsonBytes, &authAccount)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return nil
	}

	if !isRoleAllowed(authAccount.Position, allowedRoles) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "position isn't allowed"})
		return nil
	}

	return &authAccount
}

func isRoleAllowed(role string, allowedRoles []string) bool {
	for _, aRole := range allowedRoles {
		if role == aRole {
			return true
		}
	}
	return false
}
