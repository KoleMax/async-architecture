package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticateResponse struct {
	AuthAccount AuthAccount `json:"account"`
}

// Authenticate	 godoc
// @Summary      Authenticate
// @Description  Authenticate
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  AuthenticateResponse
// @Router       /api/v1/auth/authenticate [post]
// @Security     BasicAuth
func (s *Service) Authenticate(ctx *gin.Context) {
	email, password, hasAuth := ctx.Request.BasicAuth()
	if !hasAuth {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no auth header"})
		return
	}

	account, err := s.authAccountsRepo.GetByEmail(email)
	// Such detailed errors for debug purpose
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if account == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no such account"})
		return
	}
	if account.Password != password {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, AuthenticateResponse{
		AuthAccount: AuthAccount{
			PublicId: account.PublicId,
			Email:    account.Email,
			Fullname: account.Fullname,
			Position: account.Position,
		},
	})
}
